<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import {
  MeterProvider,
  PeriodicExportingMetricReader,
} from '@opentelemetry/sdk-metrics'
import { WasmCollectorExporter } from '../src/WasmCollectorExporter'
import { Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  LineElement,
  PointElement,
  LinearScale,
  CategoryScale,
  Filler,
  Tooltip,
  Legend,
} from 'chart.js'

ChartJS.register(LineElement, PointElement, LinearScale, CategoryScale, Filler, Tooltip, Legend)

const clickMetricName = 'click.rate'
const wasmReady = ref(false)
const totalClicks = ref(0)
const chartLabels = reactive<string[]>([])
const chartValues = reactive<number[]>([])
const logsEl = ref<HTMLElement | null>(null)
let startTime: number | null = null

const chartData = computed(() => ({
  labels: [...chartLabels],
  datasets: [
    {
      label: clickMetricName,
      data: [...chartValues],
      borderColor: '#8be9fd',
      backgroundColor: 'rgba(139, 233, 253, 0.15)',
      pointBackgroundColor: '#50fa7b',
      pointBorderColor: '#50fa7b',
      pointRadius: 4,
      tension: 0.3,
      fill: true,
    },
  ],
}))

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  devicePixelRatio: window.devicePixelRatio * 2,
  animation: { duration: 300 },
  scales: {
    x: {
      title: { display: true, text: 'Time (s)', color: '#ccc', font: { size: 11 } },
      ticks: { color: '#aaa', font: { size: 10 } },
      grid: { color: '#333' },
    },
    y: {
      title: { display: true, text: 'Clicks/s', color: '#ccc', font: { size: 11 } },
      ticks: { color: '#aaa', font: { size: 10 }, stepSize: 1 },
      grid: { color: '#333' },
      beginAtZero: true,
    },
  },
  plugins: {
    legend: { display: false },
  },
}

// Shared reactive store for collector stderr logs
if (!(globalThis as any).__collectorLogs) {
  (globalThis as any).__collectorLogs = reactive<string[]>([])
}
const collectorLogs = (globalThis as any).__collectorLogs as string[]

// Auto-scroll logs to bottom when new logs arrive
watch(() => collectorLogs.length, async () => {
  await nextTick()
  if (logsEl.value) {
    logsEl.value.scrollTop = logsEl.value.scrollHeight
  }
})

// Singleton OTel SDK and WASM state, shared across all mounts of this component
const g = globalThis as any
if (!g.__wasmDemoState) {
  const exporter = new WasmCollectorExporter()
  const reader = new PeriodicExportingMetricReader({
    exporter,
    exportIntervalMillis: 60_000, // long interval; we flush manually on click
  })
  const meterProvider = new MeterProvider({ readers: [reader] })
  const meter = meterProvider.getMeter('wasm-demo')
  const clickCounter = meter.createCounter('click.count', {
    description: 'Number of button clicks',
  })
  g.__wasmDemoState = {
    reader,
    clickCounter,
    wasmInitStarted: false,
    go: null,
    mod: null,
    inst: null,
    exportListeners: new Set<(obj: any) => void>(),
  }
}
const demoState = g.__wasmDemoState

async function sendClick() {
  if (!wasmReady.value) return
  demoState.clickCounter.add(1)
  await demoState.reader.forceFlush()
}

function addChartPoint(value: number) {
  const maxPoints = 8
  
  if (startTime === null) startTime = Date.now()
  const elapsed = ((Date.now() - startTime) / 1000).toFixed(1)
  chartLabels.push(elapsed)
  chartValues.push(value)
  if (chartLabels.length > maxPoints) {
    chartLabels.splice(0, chartLabels.length - maxPoints)
    chartValues.splice(0, chartValues.length - maxPoints)
  }
}

function updateChartData(obj: any) {
  // Walk the OTLP JSON structure to find click.rate and use the value directly
  for (const rm of obj?.resourceMetrics ?? []) {
    for (const sm of rm?.scopeMetrics ?? []) {
      for (const metric of sm?.metrics ?? []) {
        if (metric.name !== clickMetricName) continue
        const dataPoints = metric.sum?.dataPoints ?? metric.gauge?.dataPoints ?? []
        for (const dp of dataPoints) {
          const value = Number(dp.asInt ?? dp.asDouble ?? 0)
          totalClicks.value += value
          addChartPoint(value)
        }
      }
    }
  }
}

// Register this instance's handler and clean up on unmount
demoState.exportListeners.add(updateChartData)
onUnmounted(() => {
  demoState.exportListeners.delete(updateChartData)
})

// The Wasm file is pre-compressed when uploaded to Cloudflare,
// and as a result must be manually decompressed here.
async function decompressGzip(response: Response) {
  if (response.headers.get("server") !== "cloudflare") {
    return response
  }

  const compressedBlob = await response.blob();
  const decompressedStream = compressedBlob.stream().pipeThrough(new DecompressionStream('gzip'));
  return new Response(decompressedStream, {
    headers: { 'Content-Type': 'application/wasm' }
  });
}

onMounted(async () => {
  // If WASM is already running from another mount, just sync the ready state
  if (demoState.wasmInitStarted) {
    if (demoState.inst) wasmReady.value = true
    return
  }
  demoState.wasmInitStarted = true

  const script = document.createElement('script')
  script.src = import.meta.env.BASE_URL + 'wasm_exec.js'
  script.onload = async () => {
    if (!WebAssembly.instantiateStreaming) {
      WebAssembly.instantiateStreaming = async (resp: any, importObject: any) => {
        const source = await (await resp).arrayBuffer()
        return await WebAssembly.instantiate(source, importObject)
      }
    }

    // @ts-ignore
    demoState.go = new Go()

    // Intercept stderr (fd 2) from the Go WASM runtime to capture collector logs
    const origWriteSync = globalThis.fs.writeSync.bind(globalThis.fs)
    let stderrBuf = ''
    const stderrDecoder = new TextDecoder('utf-8')
    globalThis.fs.writeSync = (fd: number, buf: Uint8Array) => {
      if (fd === 2) {
        stderrBuf += stderrDecoder.decode(buf, { stream: true })
        const nl = stderrBuf.lastIndexOf('\n')
        if (nl !== -1) {
          const lines = stderrBuf.substring(0, nl).split('\n')
          for (const line of lines) {
            if (line.trim()) collectorLogs.push(line)
          }
          stderrBuf = stderrBuf.substring(nl + 1)
        }
      }
      return origWriteSync(fd, buf)
    }

    // Set up a callback for the JS exporter to send telemetry data
    // from the Go Wasm runtime to JavaScript via js.CopyBytesToJS.
    // Broadcasts to all mounted component instances.
    const decoder = new TextDecoder('utf-8')
    globalThis.__otelExportCallback = (uint8Array: Uint8Array) => {
      const json = decoder.decode(uint8Array)
      try {
        const obj = JSON.parse(json)
        for (const listener of demoState.exportListeners) {
          listener(obj)
        }
      } catch (e) {
        console.error('Failed to parse exported telemetry:', e)
      }
    }

    try {
      const result = await WebAssembly.instantiateStreaming(
        fetch(import.meta.env.BASE_URL + "otelwasmcol.wasm.gz")
        .then(decompressGzip),
        demoState.go.importObject
      )
      demoState.mod = result.module
      demoState.inst = result.instance
      wasmReady.value = true
      runWasm()
    } catch (err) {
      console.error(err)
    }
  }
  document.head.appendChild(script)
})

async function runWasm() {
  if (!demoState.go || !demoState.inst || !demoState.mod) return

  totalClicks.value = 0
  chartLabels.length = 0
  chartValues.length = 0
  startTime = null
  const configUrl = `${window.location.origin}/otelcol-config.yaml`
  demoState.go.argv = ["otelwasmcol.wasm", `--config=${configUrl}`]
  await demoState.go.run(demoState.inst)
  demoState.inst = await WebAssembly.instantiate(demoState.mod, demoState.go.importObject)
}
</script>

<template>
  <div class="demo-layout">
    <div class="demo-header">
      <h1>Demo</h1>
      <button
        @click="sendClick"
        :disabled="!wasmReady"
        class="demo-button"
      >
        {{ wasmReady ? '🖱️ Click me!' : 'Loading...' }}
      </button>
    </div>

    <div class="collector-logs-container">
      <div class="collector-logs-header">Collector logs</div>
      <div class="collector-logs" ref="logsEl">
        <div v-for="(line, i) in collectorLogs" :key="i" class="log-line">{{ line }}</div>
        <div v-if="collectorLogs.length === 0" class="log-line dim">Waiting for collector output…</div>
      </div>
    </div>

    <div class="chart-container">
      <Line :data="chartData" :options="chartOptions" />
    </div>
  </div>
</template>

<style scoped>
.demo-layout {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  width: 100%;
}
.demo-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 90%;
  margin: 0 auto;
}
.demo-header h1 {
  margin: 0;
}
.demo-button {
  padding: 12px 24px;
  background: #5D5DFF;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 1.1rem;
  cursor: pointer;
  flex-shrink: 0;
}
.demo-button:hover {
  background: #4a4aed;
}
.demo-button:disabled {
  background: #9ca3af;
  cursor: not-allowed;
}
.chart-container {
  width: 90%;
  height: 200px;
  background: linear-gradient(135deg, rgba(54, 56, 85, 0.5), rgba(26, 28, 44, 0.45));
  border: 1.5px solid rgba(141, 141, 255, 0.15);
  border-radius: 0.75rem;
  padding: 12px;
}
.collector-logs-container {
  width: 90%;
  margin: 0 auto;
}
.collector-logs-header {
  background: linear-gradient(135deg, rgba(93, 93, 255, 0.2), rgba(141, 141, 255, 0.1));
  color: #50fa7b;
  padding: 6px 14px;
  font-size: 0.8rem;
  font-family: 'Fira Code', 'Cascadia Code', 'JetBrains Mono', monospace;
  border-radius: 0.75rem 0.75rem 0 0;
  border-bottom: 1px solid rgba(141, 141, 255, 0.15);
}
.collector-logs {
  width: 100%;
  max-height: 200px;
  overflow-y: auto;
  background: linear-gradient(135deg, rgba(54, 56, 85, 0.5), rgba(26, 28, 44, 0.45));
  color: #d4d4d4;
  padding: 10px 14px;
  border-radius: 0 0 0.75rem 0.75rem;
  font-family: 'Fira Code', 'Cascadia Code', 'JetBrains Mono', monospace;
  font-size: 0.65rem;
  line-height: 1.5;
  text-align: left;
}
.log-line {
  white-space: pre-wrap;
  word-break: break-all;
}
.log-line.dim {
  color: #666;
  font-style: italic;
}

/* Light theme overrides */
:global(html:not(.dark) .collector-logs-header) {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.15), rgba(139, 92, 246, 0.08));
  color: #1a1a1a;
  border-bottom-color: rgba(99, 102, 241, 0.15);
}
:global(html:not(.dark) .collector-logs) {
  background: #1e1e2e;
}
:global(html:not(.dark) .chart-container) {
  background: #1e1e2e;
  border-color: rgba(99, 102, 241, 0.2);
}
</style>
