---
# try also 'default' to start simple
theme: seriph
# random image from a curated Unsplash collection by Anthony
# like them? see https://unsplash.com/collections/94734566/slidev
background: https://cover.sli.dev
# some information about your slides (markdown enabled)
title: Observability Without Borders
class: text-center
# https://sli.dev/features/drawing
drawings:
  persist: false
# slide transition: https://sli.dev/guide/animations.html#slide-transitions
transition: slide-left
# enable MDC Syntax: https://sli.dev/features/mdc
mdc: true
# duration of the presentation
duration: 25min
---

# Observability Without Borders

---

# What is the Collector?

<img src="/otel-diagram.svg" style="max-height: 100%; max-width: 100%; object-fit: contain; display: block; margin: auto;" />

---

# Write once, run everywhere™

---

# Why WebAssembly?

---

# Where can I run my Collector today?

The Collector supports a variety of compilation targets today:

<div class="platforms">
  <div class="tier-group tier1-group">
    <div class="tier1"><code>linux/amd64</code></div>
  </div>
  <div class="tier-group tier2-group">
    <div class="tier2"><code>darwin/arm64</code></div>
    <div class="tier2"><code>linux/arm64</code></div>
    <div class="tier2"><code>windows/amd64</code></div>
  </div>
  <div class="tier-group tier3-group">
    <div class="tier3"><code>darwin/amd64</code></div>
    <div class="tier3"><code>linux/386</code></div>
    <div class="tier3 wasm-special"><code>js/wasm</code></div>
    <div class="tier3"><code>linux/arm</code></div>
    <div class="tier3"><code>linux/ppc64le</code></div>
    <div class="tier3"><code>linux/riscv64</code></div>
    <div class="tier3"><code>linux/s390x</code></div>
    <div class="tier3"><code>windows/386</code></div>
  </div>
  <div class="tier-group unofficial-group">
    <div class="unofficial"><code>aix/ppc64</code></div>
    <div class="unofficial"><code>plan9/amd64</code></div>
    <div class="unofficial"><code>wasip1/wasm</code></div>
    <div class="unofficial"><code>...</code></div>
  </div>
</div>

---

# Challenges

---

# Limitations


---
layout: center
---

# WebAssembly Demo

<div class="flex flex-col items-center justify-center gap-4">
  <button
    @click="sendClick"
    :disabled="!wasmReady"
    class="px-6 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:bg-gray-400 disabled:cursor-not-allowed text-lg"
  >
    {{ wasmReady ? '🖱️ Click me!' : 'Loading...' }}
  </button>
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
import { ref, onMounted } from 'vue'

const wasmReady = ref(false)
const clickCount = ref(0)
const outputEl = ref<HTMLElement | null>(null)
let go: any = null
let mod: any = null
let inst: any = null

// Build an OTLP JSON metrics payload with a single click.count delta data point
function buildClickMetrics(): string {
  const now = Date.now() * 1_000_000 // nanoseconds
  return JSON.stringify({
    resourceMetrics: [{
      scopeMetrics: [{
        metrics: [{
          name: 'click.count',
          sum: {
            dataPoints: [{
              startTimeUnixNano: String(now),
              timeUnixNano: String(now),
              asInt: '1'
            }],
            aggregationTemporality: 1,
            isMonotonic: true
          }
        }]
      }]
    }]
  })
}

function sendClick() {
  if (!wasmReady.value) return
  const payload = buildClickMetrics()
  const encoder = new TextEncoder()
  const bytes = encoder.encode(payload)
  const uint8Array = new Uint8Array(bytes)
  // Send to the Collector's js receiver via CopyBytesToGo
  if (globalThis.__otelReceiveCallback) {
    globalThis.__otelReceiveCallback(uint8Array)
  }
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
