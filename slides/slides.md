---
# try also 'default' to start simple
theme: seriph
# random image from a curated Unsplash collection by Anthony
# like them? see https://unsplash.com/collections/94734566/slidev
background: https://cover.sli.dev
# some information about your slides (markdown enabled)
title: Welcome to Slidev
info: |
  ## Slidev Starter Template
  Presentation slides for developers.

  Learn more at [Sli.dev](https://sli.dev)
# apply UnoCSS classes to the current slide
class: text-center
# https://sli.dev/features/drawing
drawings:
  persist: false
# slide transition: https://sli.dev/guide/animations.html#slide-transitions
transition: slide-left
# enable MDC Syntax: https://sli.dev/features/mdc
mdc: true
# duration of the presentation
duration: 35min
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
</div>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const wasmReady = ref(false)
let go: any = null
let mod: any = null
let inst: any = null

onMounted(async () => {
  // Load wasm_exec.js
  const script = document.createElement('script')
  script.src = '/wasm_exec.js'
  script.onload = async () => {
    // Polyfill for browsers that don't support instantiateStreaming
    if (!WebAssembly.instantiateStreaming) {
      WebAssembly.instantiateStreaming = async (resp: any, importObject: any) => {
        const source = await (await resp).arrayBuffer()
        return await WebAssembly.instantiate(source, importObject)
      }
    }

    // @ts-ignore
    go = new Go()

    try {
      const result = await WebAssembly.instantiateStreaming(
        fetch("/otelwasmcol/bin/main.wasm"),
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

  console.clear()
  const configUrl = `${window.location.origin}/github-receiver-config.yaml`
  go.argv = ["main.wasm", `--config=${configUrl}`]
  await go.run(inst)
  inst = await WebAssembly.instantiate(mod, go.importObject)
}
</script>
