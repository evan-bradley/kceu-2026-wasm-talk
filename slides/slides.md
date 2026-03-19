---
theme: apple-basic
background: https://cover.sli.dev
title: Observability Without Borders
class: text-center
transition: slide-left
mdc: true
layout: intro
fonts:
  sans: 'Inter'
  mono: 'Fira Code'
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

<!-- Mention that we are both maintainers in the Collector SIG -->

---

# What is the Collector?

<img src="/otel-diagram.svg" style="flex: 1; min-height: 0; max-width: 100%; object-fit: contain; display: block; margin: auto;" />

<!-- Show of hands.

After that, mention pitch succintly: the Collector allows you to build telemetry pipelines to receiver, process and export your telemetry from any source to any backend.-->

---

# What we will cover

<div class="icon-grid">
  <carbon-checkmark-outline class="icon" />
  <span>WASM is ready today for use in production</span>
  <carbon-checkmark class="icon" />
  <span>The Collector already largely supports compilation to WASM</span>
  <carbon-idea class="icon" />
  <span>Upstream support means it's ready for your ideas</span>
</div>

<!-- 

Some key points to keep in mind as we go through the presentation:

1. Wasm was introduced around 9 years ago now, and has wide browser and server-side support.
   There's still a lot that needs to be done, but it has been used in production and can
   live up to its promises if you're deliberate with where you use it.
2. With just a few tweaks to make things work, the Collector already has some basic
   compatibility with WebAssembly. We'll be covering more about what does and doesn't
   work today.
3. We've added limited official support for WebAssembly, but as its an advanced use case,
   haven't done any significant development. We'd love your ideas and contributions!

 -->

---

# Write once, run everywhere™

<Timeline :items="[
  { year: '1972', desc: 'C specification prioritizes ease of writing compilers' },
  { year: '1995', desc: 'Java promises &quot;write once, run everywhere&quot;' },
  { year: '2007', desc: 'HTML5 paves the way to replace browser plugins' },
  { year: '2017', desc: 'WebAssembly MVP declared ready' },
  { year: '2026', desc: 'An upstream OTel Collector runs in a browser', highlight: true },
]" />

<!-- 

Writing a program a single time and being able to run it on any machine has long been a
goal. Starting in 1972 ... [slide content]

This list is by no means exhaustive, there are multiple items we had to leave out.
But hopefully this offers a bit of context behind the goals of WebAssembly.

 -->

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

<!-- 

WebAssembly has wide support, both in terms of runtime implementations
and in programming language support. All major browsers have supported
WebAssembly for years now, and if you want to run some software on
a device with a browser, chances are it will work. There are also
server-side WebAssembly runtimes that can run wherever you might think.
One of the more compelling use-cases is probably in edge functions
that run close to users.

A number of languages can compile to WebAssembly, meaning if you can
operate within some of it's constraints, you should be able to easy
port existing code or work in your favorite language.

WebAssembly is performant, meaning it's useful as a target for
computationally-intensive workloads. Many companies have used it
in their thick client web apps with success as we'll discuss later.

 -->

---

# WASM and WASI

WASI extends provides standardized interfaces for filesystem, networking...

<div class="comparison-grid">
  <div class="info-box">
    <h3 class="opacity-100">WebAssembly (WASM)</h3>
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

<!-- 

  As we mentioned, WebAssembly can also be run server-side and not just in the browser.
  This means its possible to do certain things (e.g. filesystem access) that aren't
  possible in the browser, which is where The WebAssembly System Interface, or WASI,
  comes in.

 -->

---

# WASM in production today

<div class="icon-grid">
  <carbon-pen-fountain class="icon" />
  <span><a href="https://www.figma.com/blog/webassembly-cut-figmas-load-time-by-3x/">Figma</a> — WASM cut load times by 3× for all document sizes</span>
  <carbon-image class="icon" />
  <span><a href="https://youtu.be/48ORmla7mak">Adobe</a> — Acrobat, Photoshop, and Lightroom run in the browser and leverage WASM</span>
  <carbon-logo-google class="icon" />
  <span><a href="https://youtu.be/2En8cj6xlv4">Google</a> — Earth, Sheets, Photos and Meet, use WASM for cross-platform code sharing</span>
</div>


<!-- 

  WebAssembly has already been in use for large thick-client apps for a long time now.

  1. Figma is written in C++, and switched their C++ to JavaScript compilation target
     from asm.js to WebAssembly and saw a significant gain in document loading speed.
  2. Adobe also has long-standing software written for desktops and has leveraged
     WebAssembly to support running some of their suite in the browser.
  3. Google applications that require heavy processing also offload heavy computations
     to WebAssembly modules to keep their applications performant.

 -->

<!-- Source: https://leaddev.com/technical-direction/webassembly-still-waiting-its-moment -->

---

# WASI previews

<Timeline :items="[
  { year: '~2020', desc: '<div class=tl-card-title>WASIp1</div><ul><li>Single API</li><li>Limited Go support</li></ul>' },
  { year: '2024', desc: '<div class=tl-card-title>WASIp2</div><ul><li>Component model</li><li>HTTP support</li><li>Only TinyGo support</li></ul>' },
  { year: '<i>2026?</i>', desc: '<div class=tl-card-title>WASIp3</div><ul><li>Async I/O</li><li>Concurrency support</li><li>Planned Go support</li></ul>' },
]" />

---
transition: fade
---

# WASM plugins inside the Collector

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

# WASM plugins inside the Collector: the vision

<div class="icon-grid">
  <carbon-api class="icon" />
  <span>Dynamically load components distributed as OCI artifacts</span>
  <carbon-settings class="icon" />
  <span>Don't be constrained by your vendor distro</span>
  <carbon-security class="icon" />
  <span>Sandboxed execution with controlled access to networking or filesystem</span>
  <carbon-plug class="icon" />
  <span>Write your Collector components in any* language</span>
</div>


---

# Collector running in WASM

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

# Collector running in WASM: Upstream developments

<div class="icon-grid">
  <carbon-add-alt class="icon" />
  <span><code>js/wasm</code> added as a Tier-3 platform (Feb 2026)</span>
  <carbon-cut class="icon" />
  <span>Custom telemetry provider to strip down binary. (Feb 2026)</span>
  <carbon-chart-bar class="icon" />
  <span>244 of 271 (~90%) Collector components already compile to <code>js/wasm</code></span>
</div>


---

# Collector running in WASM: Limitations

<div class="icon-grid">
  <carbon-scale class="icon" />
  <span>A Collector WASM binary is ≥ 38 MiB uncompressed:
  <ul>
  <li>~45%: Go runtime and other necessary data</li>
  <li>~39%: Third-party libraries</li>
  <li>~16%: Go stdlib</li>
  </ul>
  </span>

   <carbon-settings class="icon" />
  <!-- Source: https://webassembly.org/features/ -->
  <span>Limited Go support: no network, parallelism, components or WASM GC</span>
  <carbon-misuse class="icon" />
  <span>Limited TinyGo support: lack of complete stdlib.</span>
</div>

---

# Collector running in WASM: <a href="https://www.datadoghq.com/blog/engineering/agent-go-binaries/">gsa</a> analysis

<img src="/gsa.png" class="h-100 mx-auto" />

<!--
  Ref:
  https://github.com/WebAssembly/design/blob/master/BinaryEncoding.md#data-section
  https://blog.tangrs.id.au/2022/02/15/notes-on-go-binary-metadata/
-->

---

# Collector running in WASM: OCB manifest

```yaml{all|1-2,6-8|1,3-4,10-12|14-16|18-19}
exporters:
  - gomod: go.opentelemetry.io/collector/exporter/otlphttpexporter v0.147.0
  - gomod: github.com/evan-bradley/kceu-2026-wasm-talk/jsexporter v0.0.0
    path: ../jsexporter

processors:
  - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/processor/deltatorateprocessor v0.147.0
  - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor v0.147.0

receivers:
  - gomod: github.com/evan-bradley/kceu-2026-wasm-talk/jsreceiver v0.0.0
    path: ../jsreceiver

providers:
  - gomod: go.opentelemetry.io/collector/confmap/provider/httpprovider v1.53.0
  - gomod: go.opentelemetry.io/collector/confmap/provider/httpsprovider v1.53.0

conf_resolver:
  default_uri_scheme: http
```

<!-- 

  Here's the build manifest used by the OpenTelemetry Collector Builder, which we call OCB.
  This is a slightly cut-down version of the manifest we use to build the Collector you'll
  see later in the slides.

  1. First: as you can see, many of the upstream components you know and love are supported.
  2. However, when running the Collector in Wasm, there's a good chance you may want or need
     custom components. We have these two components to communicate data to and from the
     JavaScript runtime the Collector is running alongside.
  3. It's also worth noting that you'll need to pay close attention with how to configure
     your Collectors. Since in the browser there is no filesystem access, we cut out those
     providers and only use two that get config using an HTTP request. You could write a
     custom one too, if you had another way you wanted to grab the config.
  4. One important note if you use this, our environment variable substitution syntax
     inside Collector configs is customizable, but defaults to reading environment
     variables. You will want to change this for Collectors running in the browser.

 -->

---

# Collector running in WASM: Creating a WASM binary

<div class="icon-grid">
  <carbon-terminal class="icon" />
  <span><code>GOOS=js GOARCH=wasm go build .</code></span>
  <carbon-terminal class="icon" />
  <span><code>GOOS=wasip1 GOARCH=wasm go build .</code></span>
</div>

<!-- 
Cross-architecture compilation with Go is a breeze, so simply specify
the GOOS environment variable for your desired Wasm target and a GOARCH
of `wasm` to compile to a wasm binary.
 -->

---

# Observability without borders

<div class="icon-grid">
  <carbon-application-web class="icon" />
  <span>Running in a browser</span>
  <carbon-container-software class="icon" />
  <span>Running in WASM</span>
  <carbon-plug class="icon" />
  <span>Running as a language plugin</span>
</div>

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

# Observability without borders: WASM runtime

<div class="icon-grid">
  <carbon-wifi-off class="icon" />
  <span>Limited/no networking currently (Go only supports WASIp1)</span>
  <carbon-folder class="icon" />
  <span>Filesystem access available if the host grants it</span>
  <carbon-partnership class="icon" />
  <span>Realistically should only be used alongside other WASM applications</span>
</div>

---

# Looking ahead

<div class="icon-grid">
  <carbon-in-progress class="icon" />
  <span>Go WASIp3 support still under active discussion.</span>
  <carbon-flow class="icon" />
  <span>Concurrency support in the Component Model and WASI 1.0.</span>
  <carbon-package class="icon" />
  <span>TinyGo compatibility opens up the possibility of smaller binaries and WASIp2.</span>
  <carbon-microphone class="icon" />
  <span>WASI OTel (talk happening at WASMCon).</span>
  <carbon-group class="icon" />
  <span>Contributions from YOU in the audience!</span>
</div>

---

# Demo

<div class="icon-grid">
  <carbon-assembly-reference class="icon" />
  <span>We compiled a basic Collector that communicates with the OTel JS SDK</span>
  <carbon-mobile class="icon" />
  <span>It runs on any modern browser, so try it on your phone!</span>
  <QrArrow />
</div>

---

# Demo

<div class="arch-slide">
  <div class="arch-diagram">
    <div class="arch-box-root arch-block browser-block">
      <span class="arch-label-left">Browser</span>
      <div class="demo-grid">
        <div class="demo-sdk-group arch-block demo-node-ui">
          <div class="demo-node demo-node-ui">👆 Button</div>
          <div class="demo-vert-arrow"></div>
          <div class="demo-node demo-node-ui">OTel JS SDK</div>
        </div>
        <div class="demo-pipe-arrow">
          <div class="demo-pipe-label">OTLP metrics</div>
        </div>
        <div class="demo-collector-wrapper arch-block wasm-col-block">
          <span class="arch-label-left">Collector</span>
          <span class="arch-label-right arch-wasm">WASM</span>
          <div class="demo-collector-inner">
            <div class="demo-subcomp">JS Receiver</div>
            <div class="demo-inner-arrow"></div>
            <div class="demo-subcomp">JS Exporter</div>
          </div>
        </div>
        <div class="demo-pipe-arrow">
          <div class="demo-pipe-label">OTLP metrics</div>
        </div>
        <div class="demo-node demo-node-ui">📊 Chart</div>
      </div>
    </div>
  </div>
</div>

---

<WasmDemo />

---
layout: center
---

<h1 class="qa-title">Q&A</h1>

<QrArrow />

<img src="/kceu26.svg" class="kceu-logo" />
