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

<Timeline :items="[
  { year: '1972', desc: 'C specification prioritizes ease of writing compilers' },
  { year: '1995', desc: 'Java promises &quot;write once, run everywhere&quot;' },
  { year: '2007', desc: 'HTML5 paves the way to replace browser plugins' },
  { year: '2017', desc: 'WebAssembly MVP declared ready' },
  { year: '2026', desc: 'An upstream OTel Collector runs in a browser', highlight: true },
]" />

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

# Wasm and WASI

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
      <li v-click="1">Wasm interfaces for OS interaction</li>
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

# Wasm in the Collector

<div class="arch-slide">
  <div class="arch-diagram">
    <div class="arch-box-root arch-block collector-block">
      <span class="arch-label-left">Collector</span>
      <span class="arch-label-right arch-native">Native</span>
      <div class="arch-block runtime-block">
        <span class="arch-label-left">Wasm runtime</span>
        <span class="arch-label-right arch-native">Native (Go library)</span>
        <div class="arch-block plugin-block">
          <span class="arch-label-left">Plugin</span>
          <span class="arch-label-right arch-wasm">Wasm</span>
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

# Collector in Wasm

<div class="arch-slide">
  <div class="arch-diagram">
    <div class="arch-box-root arch-block runtime-block">
      <span class="arch-label-left">Wasm runtime</span>
      <span class="arch-label-right arch-native">Native</span>
      <div class="arch-block collector-block">
        <span class="arch-label-left">Collector</span>
        <span class="arch-label-right arch-wasm">Wasm</span>
      </div>
    </div>
  </div>
  <div class="arch-details">
    <ul>
      <li><a href="https://ottl.run/">ottl.run</a></li>
      <li>Filtering, sampling and transforming in the browser</li>
      <li>Run it on your Wasm runtime for sandboxing</li>
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

# Collector in Wasm: Upstream developments

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

<Timeline :items="[
  { year: '~2020', desc: '<b>WASIp1</b><hr/><ul><li>Single API</li><li>Limited Go support</li></ul>' },
  { year: '2024', desc: '<b>WASIp2</b><hr/><ul><li>Component model</li><li>HTTP support</li><li>Only TinyGo support</li></ul>' },
  { year: '<i>EOY 2026</i>', desc: '<b>WASIp3</b><hr/><ul><li>Async I/O</li><li>Concurrency support</li><li>Planned Go support</li></ul>' },
]" />

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

# Collector in Wasm: Limitations

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

# Collector in Wasm: Creating a Wasm binary

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
  <span>Realistically should only be used alongside other Wasm applications</span>
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

<WasmDemo />

# Questions?

---
layout: center
class: topic-both
---

<h1 class="qa-title">Q&A</h1>

<QrArrow />

<img src="/kceu26.svg" class="kceu-logo" />