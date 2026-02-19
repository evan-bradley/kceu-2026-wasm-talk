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
    @click="runWasm"
    :disabled="!wasmReady"
    class="px-6 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:bg-gray-400 disabled:cursor-not-allowed"
  >
    {{ wasmReady ? 'Run' : 'Loading...' }}
  </button>
  <pre
    ref="outputEl"
    class="wasm-output"
  ><code>{{ outputText }}</code></pre>
</div>

<style>
.wasm-output {
  width: 90%;
  max-height: 340px;
  overflow-y: auto;
  background: #1e1e1e;
  color: #d4d4d4;
  padding: 12px;
  border-radius: 8px;
  font-size: 0.75rem;
  line-height: 1.4;
  text-align: left;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'

const wasmReady = ref(false)
const outputText = ref('')
const outputEl = ref<HTMLPreElement | null>(null)
let go: any = null
let mod: any = null
let inst: any = null

function appendOutput(text: string) {
  outputText.value += text
  nextTick(() => {
    if (outputEl.value) {
      outputEl.value.scrollTop = outputEl.value.scrollHeight
    }
  })
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

    // Intercept fs.writeSync to capture wasm output
    const decoder = new TextDecoder('utf-8')
    let outputBuf = ''
    const origWriteSync = globalThis.fs.writeSync
    globalThis.fs.writeSync = (fd: number, buf: Uint8Array) => {
      outputBuf += decoder.decode(buf)
      const nl = outputBuf.lastIndexOf('\n')
      if (nl !== -1) {
        appendOutput(outputBuf.substring(0, nl + 1))
        outputBuf = outputBuf.substring(nl + 1)
      }
      return buf.length
    }

    const origWrite = globalThis.fs.write
    globalThis.fs.write = (fd: number, buf: Uint8Array, offset: number, length: number, position: any, callback: Function) => {
      if (offset !== 0 || length !== buf.length || position !== null) {
        callback(new Error('not implemented'))
        return
      }
      const n = globalThis.fs.writeSync(fd, buf)
      callback(null, n)
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
    } catch (err) {
      console.error(err)
    }
  }
  document.head.appendChild(script)
})

async function runWasm() {
  if (!go || !inst || !mod) return

  outputText.value = ''
  const configUrl = `${window.location.origin}/github-receiver-config.yaml`
  const ghReceiver = `${window.location.origin}/gh_org.yaml`
  const btExtension = `${window.location.origin}/gh_pat.yaml`
  go.argv = ["otelwasmcol.wasm", `--config=${configUrl}`, `--config=${ghReceiver}`, `--config=${btExtension}`]
  await go.run(inst)
  inst = await WebAssembly.instantiate(mod, go.importObject)
}
</script>
