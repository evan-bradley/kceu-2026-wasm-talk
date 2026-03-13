---
theme: apple-basic
background: https://cover.sli.dev
title: Observability Without Borders
class: text-center topic-both
transition: slide-left
mdc: true
layout: intro
---

# Observability Without Borders

<p class="intro-subtitle">The OpenTelemetry Collector in a WebAssembly World</p>

<div class="intro-meta">
  <p class="intro-speakers">Pablo Baeyens <span class="intro-org">(Datadog)</span> · Evan Bradley <span class="intro-org">(Dynatrace)</span></p>
  <p class="intro-conference">Observability Day Europe 2026</p>
</div>

<QrArrow />

<img src="/kceu26.svg" class="kceu-logo" />

---

# About us

<div class="speakers-grid">
  <div class="speaker">
    <img src="/pablo.jpeg" class="speaker-pic" />
    <div class="speaker-name">Pablo Baeyens</div>
    <div class="speaker-org">Datadog</div>
  </div>
  <div class="speaker">
    <img src="/evan.jpg" class="speaker-pic" />
    <div class="speaker-name">Evan Bradley</div>
    <div class="speaker-org">Dynatrace</div>
  </div>
</div>

---
class: topic-otel
---

# What is the Collector?

<img src="/otel-diagram.svg" style="flex: 1; min-height: 0; max-width: 100%; object-fit: contain; display: block; margin: auto;" />

---
class: topic-wasm
---

# Write once, run everywhere™

<div class="timeline-wrapper">
<div class="timeline">
  <div class="timeline-track"></div>
  <div class="timeline-items">
    <div class="tl-item">
      <div class="tl-dot"></div>
      <div class="tl-year">1972</div>
      <div class="tl-desc">C specification prioritizes ease of writing compilers</div>
    </div>
    <div class="tl-item">
      <div class="tl-dot"></div>
      <div class="tl-year">1995</div>
      <div class="tl-desc">Java promises "write once, run everywhere"</div>
    </div>
    <div class="tl-item">
      <div class="tl-dot"></div>
      <div class="tl-year">2007</div>
      <div class="tl-desc">HTML5 paves the way to replace browser plugins</div>
    </div>
    <div class="tl-item">
      <div class="tl-dot"></div>
      <div class="tl-year">2017</div>
      <div class="tl-desc">WebAssembly MVP declared ready</div>
    </div>
    <div class="tl-item highlight">
      <div class="tl-dot"></div>
      <div class="tl-year">2026</div>
      <div class="tl-desc">An upstream OTel Collector runs in a browser</div>
    </div>
  </div>
</div>
</div>

---
class: topic-wasm
---

# Why WebAssembly?

<div class="icon-grid">
  <carbon-devices class="icon" />
  <span>Most likely format to run on a user device</span>
  <carbon-code class="icon" />
  <span>Compile from Go, Rust, C++ and many other languages</span>
  <carbon-flash class="icon" />
  <span>Performance for computationally-intensive workloads</span>
</div>

---
class: topic-wasm
---

# WASM and WASI

<div class="comparison-grid">
  <div class="info-box">
    <h3 class="opacity-100">WebAssembly (Wasm)</h3>
    <ul>
      <li v-click="1">Binary format targeted for browsers</li>
      <li v-click="2">Can only see what the host allows</li>
      <li v-click="3">Stable (3.0) specification</li>
      <li v-click="4">Widely supported</li>
    </ul>
  </div>
  <div class="info-box">
    <h3 class="opacity-100">WASI</h3>
    <ul>
      <li v-click="1">WASM interfaces for OS interaction</li>
      <li v-click="2">Standardized but controlled access</li>
      <li v-click="3">Unstable (WASIp2) specification</li>
      <li v-click="4">Less widely supported</li>
    </ul>
  </div>
</div>

<!-- TODO: Add examples of WebAssembly usage -->

---
transition: fade
class: topic-both
---

# WASM in the Collector

<div class="arch-slide">
  <div class="arch-diagram">
    <div class="arch-box-root arch-block collector-block">
      <span class="arch-label-left">Collector</span>
      <span class="arch-label-right arch-native">Native</span>
      <div class="arch-block runtime-block">
        <span class="arch-label-left">WASM runtime</span>
        <span class="arch-label-right arch-native">Native (Go library)</span>
        <div class="arch-block plugin-block">
          <span class="arch-label-left">Plugin</span>
          <span class="arch-label-right arch-wasm">WASM</span>
        </div>
      </div>
    </div>
  </div>
  <div class="arch-details">
    <ul>
      <li>Runtime plugins</li>
      <li><a href="https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/11772">wasmprocessor</a></li>
      <li>OTTL custom functions</li>
      <li><a href="https://github.com/otelwasm/otelwasm">otelwasm project</a></li>
    </ul>
  </div>
</div>

---
class: topic-both
---

# Collector in WASM

<div class="arch-slide">
  <div class="arch-diagram">
    <div class="arch-box-root arch-block runtime-block">
      <span class="arch-label-left">WASM runtime</span>
      <span class="arch-label-right arch-native">Native</span>
      <div class="arch-block collector-block">
        <span class="arch-label-left">Collector</span>
        <span class="arch-label-right arch-wasm">WASM</span>
      </div>
    </div>
  </div>
  <div class="arch-details">
    <ul>
      <li><a href="https://ottl.run/">ottl.run</a></li>
      <li>Filtering, sampling and transforming in the browser</li>
      <li>Run it on your WASM runtime for sandboxing</li>
    </ul>
  </div>
</div>

---
class: topic-otel
---

# Where can I run my Collector today?

The Collector supports a variety of compilation targets today:

<div class="platforms">
  <div v-click="1" class="tier-group tier1-group">
    <div class="tier1"><code>linux/amd64</code></div>
  </div>
  <div v-click="2" class="tier-group tier2-group">
    <div class="tier2"><code>darwin/arm64</code></div>
    <div class="tier2"><code>linux/arm64</code></div>
    <div class="tier2"><code>windows/amd64</code></div>
  </div>
  <div v-click="3" class="tier-group tier3-group">
    <div class="tier3"><code>darwin/amd64</code></div>
    <div class="tier3"><code>linux/386</code></div>
    <div class="tier3 wasm-special"><code>js/wasm</code></div>
    <div class="tier3"><code>linux/arm</code></div>
    <div class="tier3"><code>linux/ppc64le</code></div>
    <div class="tier3"><code>linux/riscv64</code></div>
    <div class="tier3"><code>linux/s390x</code></div>
    <div class="tier3"><code>windows/386</code></div>
  </div>
  <div v-click="4" class="tier-group unofficial-group">
    <div class="unofficial"><code>aix/ppc64</code></div>
    <div class="unofficial"><code>plan9/amd64</code></div>
    <div class="unofficial"><code>wasip1/wasm</code></div>
    <div class="unofficial"><code>...</code></div>
  </div>
</div>
---
class: topic-both
---

# Collector in WASM: Upstream developments

<div class="icon-grid">
  <carbon-add-alt class="icon" />
  <span><code>js/wasm</code> added as a Tier-3 platform (Feb 2026)</span>
  <carbon-cut class="icon" />
  <span>Allow stripping down telemetry provider to remove unneeded functionality in the browser. (Feb 2026)</span>
  <carbon-chart-bar class="icon" />
  <span>244 of 271 (~90%) Collector components already compile to <code>js/wasm</code></span>
</div>

---
class: topic-wasm
---

# WASI previews

<div class="timeline-wrapper">
<div class="timeline">
  <div class="timeline-track"></div>
  <div class="timeline-items">
    <div class="tl-item">
      <div class="tl-dot"></div>
      <div class="tl-year">~2020</div>
      <div class="tl-desc">WASIp1
      <ul>
      <li>Single API</li>
      <li>Limited Go support</li>
      </ul>
      </div>
    </div>
    <div class="tl-item">
      <div class="tl-dot"></div>
      <div class="tl-year">2024</div>
      <div class="tl-desc">WASIp2
      <ul>
      <li>Component model</li>
      <li>HTTP support</li>
      <li>Only TinyGo support</li>
      </ul>
      </div>
    </div>
    <div class="tl-item">
      <div class="tl-dot"></div>
      <div class="tl-year"><i>EOY 2026</i></div>
      <div class="tl-desc">WASIp3
      <ul>
      <li>Async I/O</li>
      <li>Concurrency support</li>
      <li>Planned Go support</li>
      </ul>
      </div>
    </div>
  </div>
</div>
</div>

---
class: topic-both
---

# Challenges

<div class="icon-grid">
  <carbon-warning-alt class="icon" />
  <span>Different operating environment: no equivalent APIs for many syscalls.</span>
  <carbon-wifi-off class="icon" />
  <span>No networking in WASI.</span>
  <carbon-locked class="icon" />
  <span>Go only supports WASIp1; we need the WebAssembly Component Model for networking and filesystem access.</span>
</div>

<div class="icon-grid">
  <carbon-time class="icon" />
  <span>TODO: Mention concurrency requirements</span>
  <carbon-settings class="icon" />
  <span>TODO: Mention Collector runtime components</span>
</div>


---
class: topic-both
---

# Collector in WASM: Limitations

<div class="icon-grid">
  <carbon-scale class="icon" />
  <span>A Collector has at least a 10 MiB uncompressed binary size.</span>
  <carbon-misuse class="icon" />
  <span>TinyGo stdlib doesn't reimplement enough Go stdlib network packages (e.g. net/http/httputil).</span>
</div>

<div class="icon-grid">
  <carbon-chart-treemap class="icon" />
  <span>TODO: Mention gsa output; link to Datadog blogpost</span>
</div>

---
class: topic-both
---

# Collector in WASM: Creating a Wasm binary

<div class="icon-grid">
  <carbon-tool-box class="icon" />
  <span>Using OCB</span>
  <carbon-list class="icon" />
  <span>Check if your component is supported</span>
  <carbon-terminal class="icon" />
  <span><code>GOOS=js GOARCH=wasm go build .</code></span>
  <carbon-terminal class="icon" />
  <span><code>GOOS=wasip1 GOARCH=wasm go build .</code></span>
  <carbon-settings class="icon" />
  <span>Load configuration via confmap providers (HTTP or inline YAML)</span>
</div>

---
class: topic-both
---

# Observability without borders

<div class="icon-grid">
  <carbon-application-web class="icon" />
  <span>Running in a browser</span>
  <carbon-container-software class="icon" />
  <span>Running in a Wasm runtime</span>
  <carbon-plug class="icon" />
  <span>Running in a language plugin</span>
</div>

---
class: topic-both
---

# Observability without borders: browser

<div class="icon-grid">
  <carbon-folder-off class="icon" />
  <span>No FS access</span>
  <carbon-close-outline class="icon" />
  <span>Can't open ports</span>
  <carbon-application class="icon" />
  <span>Uses: SDK, thick-client apps, heavy in-browser apps and electron apps</span>
</div>

---
class: topic-both
---

# Observability without borders: Wasm runtime

<div class="icon-grid">
  <carbon-wifi-off class="icon" />
  <span>Limited/no networking currently (Go only supports WASIp1)</span>
  <carbon-folder class="icon" />
  <span>Filesystem access available if the host grants it</span>
  <carbon-partnership class="icon" />
  <span>Realistically should only be used alongside other WASM applications</span>
</div>

---
class: topic-both
---

# Looking ahead

<div class="icon-grid">
  <carbon-in-progress class="icon" />
  <span>Go WASIp3 support still under active discussion.</span>
  <carbon-flow class="icon" />
  <span>Concurrency support in the WebAssembly Component Model.</span>
  <carbon-package class="icon" />
  <span>TinyGo compatibility opens up the possibility of smaller binaries and WASIp2.</span>
  <carbon-microphone class="icon" />
  <span>WASI OTel (talk happening at WasmCon).</span>
  <carbon-group class="icon" />
  <span>Contributions from YOU in the audience!</span>
</div>

---
class: topic-both
---

# WebAssembly demo

<div class="flex flex-col items-center justify-center gap-4">
  <div ref="outputEl" class="wasm-output">
    <table class="otel-table">
      <thead>
        <tr><th>Metric</th><th>Value</th></tr>
      </thead>
      <tbody>
        <tr>
          <td>click.count</td>
          <td class="num">{{ clickCount }}</td>
        </tr>
      </tbody>
    </table>
  </div>
    <button
    @click="sendClick"
    :disabled="!wasmReady"
    class="px-6 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:bg-gray-400 disabled:cursor-not-allowed text-lg"
  >
    {{ wasmReady ? '🖱️ Click me!' : 'Loading...' }}
  </button>
</div>

<style>
.wasm-output {
  width: 90%;
  max-height: 380px;
  overflow-y: auto;
  background: #1e1e1e;
  color: #d4d4d4;
  padding: 12px;
  border-radius: 8px;
  font-size: 0.85rem;
  line-height: 1.4;
  text-align: left;
}
.otel-table {
  width: 100%;
  border-collapse: collapse;
}
.otel-table th {
  text-align: left;
  padding: 6px 12px;
  border-bottom: 2px solid #555;
  color: #8be9fd;
  font-weight: 600;
}
.otel-table td {
  padding: 4px 12px;
  border-bottom: 1px solid #333;
}
.otel-table td.num {
  text-align: right;
  font-variant-numeric: tabular-nums;
}
.otel-table tr:hover td {
  background: #2a2a2a;
}
</style>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { metrics } from '@opentelemetry/api'
import {
  MeterProvider,
  PeriodicExportingMetricReader,
} from '@opentelemetry/sdk-metrics'
import { WasmCollectorExporter } from './src/WasmCollectorExporter'

const wasmReady = ref(false)
const clickCount = ref(0)
const outputEl = ref<HTMLElement | null>(null)
let go: any = null
let mod: any = null
let inst: any = null

// Shared reactive store for collector stderr logs (accessible across slides)
if (!(globalThis as any).__collectorLogs) {
  (globalThis as any).__collectorLogs = reactive<string[]>([])
}
const collectorLogs = (globalThis as any).__collectorLogs as string[]

// Set up the OTel metrics SDK with our custom exporter
const exporter = new WasmCollectorExporter()
const reader = new PeriodicExportingMetricReader({
  exporter,
  exportIntervalMillis: 60_000, // long interval; we flush manually on click
})
const meterProvider = new MeterProvider({ readers: [reader] })
metrics.setGlobalMeterProvider(meterProvider)
const meter = metrics.getMeter('wasm-demo')
const clickCounter = meter.createCounter('click.count', {
  description: 'Number of button clicks',
})

async function sendClick() {
  if (!wasmReady.value) return
  clickCounter.add(1)
  await reader.forceFlush()
}

function updateTable(obj: any) {
  // Walk the OTLP JSON structure to find click.count and accumulate
  for (const rm of obj?.resourceMetrics ?? []) {
    for (const sm of rm?.scopeMetrics ?? []) {
      for (const metric of sm?.metrics ?? []) {
        if (metric.name !== 'click.count') continue
        const dataPoints = metric.sum?.dataPoints ?? metric.gauge?.dataPoints ?? []
        for (const dp of dataPoints) {
          const value = Number(dp.asInt ?? dp.asDouble ?? 0)
          clickCount.value += value
        }
      }
    }
  }
}

// The Wasm file is pre-compressed when uploaded to Cloudflare,
// and as a result must be manually decompressed here.
async function decompressGzip(response) {
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
  const script = document.createElement('script')
  script.src = '/wasm_exec.js'
  script.onload = async () => {
    if (!WebAssembly.instantiateStreaming) {
      WebAssembly.instantiateStreaming = async (resp: any, importObject: any) => {
        const source = await (await resp).arrayBuffer()
        return await WebAssembly.instantiate(source, importObject)
      }
    }

    // @ts-ignore
    go = new Go()

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
    const decoder = new TextDecoder('utf-8')
    globalThis.__otelExportCallback = (uint8Array: Uint8Array) => {
      const json = decoder.decode(uint8Array)
      try {
        const obj = JSON.parse(json)
        updateTable(obj)
      } catch (e) {
        console.error('Failed to parse exported telemetry:', e)
      }
    }

    try {
      const result = await WebAssembly.instantiateStreaming(
        fetch("/otelwasmcol.wasm.gz")
        .then(decompressGzip),
        go.importObject
      )
      mod = result.module
      inst = result.instance
      wasmReady.value = true
      runWasm()
    } catch (err) {
      console.error(err)
    }
  }
  document.head.appendChild(script)
})

async function runWasm() {
  if (!go || !inst || !mod) return

  clickCount.value = 0
  const configUrl = `${window.location.origin}/otelcol-config.yaml`
  go.argv = ["otelwasmcol.wasm", `--config=${configUrl}`]
  await go.run(inst)
  inst = await WebAssembly.instantiate(mod, go.importObject)
}
</script>

---
class: topic-both
---

# WebAssembly demo (2)

<div class="collector-logs-container">
  <div class="collector-logs-header">Collector logs</div>
  <div class="collector-logs" ref="logsEl">
    <div v-for="(line, i) in logs" :key="i" class="log-line">{{ line }}</div>
    <div v-if="logs.length === 0" class="log-line dim">Waiting for collector output…</div>
  </div>
</div>

<style>
.collector-logs-container {
  width: 90%;
  margin: 0 auto;
}
.collector-logs-header {
  background: #2d2d2d;
  color: #50fa7b;
  padding: 6px 14px;
  font-size: 0.8rem;
  font-family: 'Fira Code', 'Cascadia Code', 'JetBrains Mono', monospace;
  border-radius: 8px 8px 0 0;
  border-bottom: 1px solid #444;
}
.collector-logs {
  width: 100%;
  max-height: 360px;
  overflow-y: auto;
  background: #1e1e1e;
  color: #d4d4d4;
  padding: 10px 14px;
  border-radius: 0 0 8px 8px;
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
</style>

<script setup lang="ts">
import { reactive, ref, watch, nextTick } from 'vue'

// Access or create the shared reactive log store
if (!(globalThis as any).__collectorLogs) {
  (globalThis as any).__collectorLogs = reactive<string[]>([])
}
const logs = (globalThis as any).__collectorLogs as string[]
const logsEl = ref<HTMLElement | null>(null)

// Auto-scroll to bottom when new logs arrive
watch(() => logs.length, async () => {
  await nextTick()
  if (logsEl.value) {
    logsEl.value.scrollTop = logsEl.value.scrollHeight
  }
})
</script>

---
layout: center
class: topic-both
---

<h1 class="qa-title">Q&A</h1>

<QrArrow />

<img src="/kceu26.svg" class="kceu-logo" />
