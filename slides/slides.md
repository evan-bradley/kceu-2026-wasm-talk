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

<!-- PABLO 

Welcome to Observability Without Borders.

If you want to follow along, this presentation is available as a webpage on the QR code on the bottom left corner which will appear later and at the end as well.

We have added some links to interesting tidbits throughout the slides that you can click to learn more.

-->

---

# About us

<div class="speakers-grid">
  <div class="speaker">
    <img src="/pablo.jpeg" class="speaker-pic" />
    <div class="speaker-name"><a href="https://github.com/mx-psi">Pablo Baeyens</a></div>
    <div class="speaker-org">Datadog</div>
  </div>
  <div class="speaker">
    <img src="/evan.jpg" class="speaker-pic" />
    <div class="speaker-name"><a href="https://github.com/evan-bradley">Evan Bradley</a></div>
    <div class="speaker-org">Dynatrace</div>
  </div>
</div>

<!-- PABLO: 
I am Pablo and this is Evan, we are both maintainers in the Collector SIG.
-->

---

# What we will cover

<div class="icon-grid">
  <carbon-checkmark-outline class="icon" />
  <span>Wasm is ready today for use in production, with caveats.</span>
  <carbon-checkmark class="icon" />
  <span>The Collector already has basic support for compilation to Wasm.</span>
  <carbon-idea class="icon" />
  <span>Upstream support means it's ready for your ideas.</span>
</div>

<!-- 

PABLO: 

Here are some key points to keep in mind as we go through the presentation:

1. Wasm was introduced around 9 years ago now, and has wide browser and server-side support.
   There's still a lot that needs to be done, but it has been used in production and can
   live up to its promises if you're deliberate with where you use it and are willing to use non-standard extensions.
2. With just a few tweaks to make things work, the Collector already has some basic
   compatibility with WebAssembly. We'll be covering more about what does and doesn't
   work today.
3. We've added limited official support for WebAssembly, but as its an advanced use case,
   haven't done any significant development. We'd love your ideas and contributions!

 -->

---

# What is the OpenTelemetry Collector?

<img src="/otel-diagram.svg" style="flex: 1; min-height: 0; max-width: 100%; object-fit: contain; display: block; margin: auto;" />

<!-- PABLO: As most of you know OpenTelemetry is the open standard for telemetry.

The Collector is a tool offered by OpenTelemetry that allows you to build telemetry pipelines to receive, process and export your telemetry from any source to any backend.-->

---

# What is WebAssembly?

<div class="icon-grid">
  <carbon-code class="icon" />
  <span>Portable compilation target available since 2017, 3.0 in 2025.</span>
  <carbon-devices class="icon" />
  <span>Can be run on browsers, embeddable and standalone runtimes.</span>
  <carbon-scale class="icon" />
  <span><a href="https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/wasm">Envoy</a>, <a href="https://istio.io/latest/docs/reference/config/proxy_extensions/wasm-plugin/">Istio</a>, <a href="https://www.openpolicyagent.org/docs/wasm">OPA</a> and <a href="https://github.com/kubernetes-sigs/kube-scheduler-wasm-extension/tree/main">k8s</a> all use it for plugins.</span>
  <carbon-pen-fountain class="icon" />
  <span><a href="https://www.figma.com/blog/webassembly-cut-figmas-load-time-by-3x/">Figma</a>, <a href="https://youtu.be/48ORmla7mak">Adobe</a>  and <a href="https://youtu.be/2En8cj6xlv4">Google</a> use it for thick-client apps and cross-platform code sharing.</span>
</div>


<!-- 

  PABLO:

  WebAssembly has already been in use for large thick-client apps for a long time now.

  1. Some other cloud-native projects like Envoy, Istio, OPA and Kubernetes use it in a limited way to provide filters and plugins.
  2. Figma is written in C++, and switched their C++ to JavaScript compilation target
     from asm.js to WebAssembly and saw a significant gain in document loading speed.
  3. Adobe also has long-standing software written for desktops and has leveraged
     WebAssembly to support running some of their suite in the browser.
  4. Google applications that require heavy processing also offload heavy computations
     to WebAssembly modules to keep their applications performant.

 Source: https://leaddev.com/technical-direction/webassembly-still-waiting-its-moment -->



---

# Why WebAssembly + OTel Collector?

<div class="icon-grid">
  <carbon-devices class="icon" />
  <span>Expands devices the Collector can run on, including user devices.</span>
  <carbon-code class="icon" />
  <span>Compile from Go, Rust, C++ and many other languages.</span>
  <carbon-flash class="icon" />
  <span>Performance for computationally-intensive workloads.</span>
</div>

<!-- 

PABLO:

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

# What is WebAssembly? Wasm and WASI

WASI extends provides standardized interfaces for filesystem, networking...

<div class="comparison-grid">
  <div class="info-box">
    <h3 class="opacity-100">WebAssembly (Wasm)</h3>
    <ul>
      <li v-click="1">Binary format targeted for browsers.</li>
      <li v-click="2">Can only see what the host allows.</li>
      <li v-click="3">Stable (3.0) specification.</li>
      <li v-click="4">Widely supported.</li>
    </ul>
  </div>
  <div class="info-box">
    <h3 class="opacity-100">WASI</h3>
    <ul>
      <li v-click="1">Wasm interfaces for OS interaction.</li>
      <li v-click="2">Standardized but controlled access.</li>
      <li v-click="3">Unstable (WASIp2) specification.</li>
      <li v-click="4">Less widely supported.</li>
    </ul>
  </div>
</div>

<!-- 

  PABLO:

  As we mentioned, WebAssembly can also be run server-side and not just in the browser.
  This means its possible to do certain things (e.g. filesystem access) that aren't
  possible in the browser, which is where The WebAssembly System Interface, or WASI,
  comes in.

  It provides a standard way of providing resources to Wasm applications, but unlike Wasm itself,
  it is unstable and less widely supported.

 -->

---

# What is WebAssembly? WASI previews

<Timeline :items="[
  { year: '~2020', desc: '<div class=tl-card-title>WASIp1</div><ul><li>Single API.</li><li>Limited Go support.</li></ul>' },
  { year: '2024', desc: '<div class=tl-card-title>WASIp2</div><ul><li>Component model.</li><li>Full HTTP support.</li><li>Only TinyGo support.</li></ul>' },
  { year: '<i>2026?</i>', desc: '<div class=tl-card-title>WASIp3</div><ul><li>Async I/O.</li><li>Concurrency support.</li><li>Planned Go support.</li></ul>' },
]" />

<!-- PABLO 

WASI is unstable and has released so far two previews and a release candidate for a third preview.

Go supports WASIp1 but does not natively support WASIp2. There is planned WASIp3 support in Go.
WASIp3 unblocks key features needed for using it for I/O and network-heavy applications.

-->

---
transition: fade
---

# Wasm plugins inside the Collector

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
      <li>Runtime plugins.</li>
      <li><a href="https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/11772">wasmprocessor</a>.</li>
      <li>OTTL custom functions.</li>
      <li><a href="https://github.com/otelwasm/otelwasm">otelwasm project</a>.</li>
    </ul>
  </div>
</div>

<!-- PABLO

There are two main things you may think about when combining Wasm and the Collector.

The first one is to run plugins inside the Collector. There have been prior proposals to 
do this including the wasmprocessor or OTTL custom functions, although there is no official upstream support for this so far.

 -->

---

# Wasm plugins inside the Collector: the vision

<div class="icon-grid">
  <carbon-api class="icon" />
  <span>Dynamically load components distributed as OCI artifacts.</span>
  <carbon-settings class="icon" />
  <span>Don't be constrained by your vendor distro.</span>
  <carbon-security class="icon" />
  <span>Sandboxed execution with controlled access to networking or filesystem.</span>
  <carbon-plug class="icon" />
  <span>Write your Collector components in any* language.</span>
</div>

<!-- PABLO 

There is no support for Wasm plugins today in the upstream Collector.
General support for using Wasm for plugins could look like this in the future: 

1. You could dynamically load components that you can pull from a registry, distributed as OCI artifacts on any distro.
2. You could control on a fine-grained way what your component has access to, allowing you to confidently access a greater array of components.
3. You would be able to write your Collector components on any language you want with a single Component Model.
-->

---

# Wasm plugins inside the Collector: PoC today

<div class="arch-slide">
<div class="arch-details">

* <a href="https://github.com/otelwasm/otelwasm">otelwasm</a> has Wasm components.
* It relies on the <a href="https://github.com/WasmEdge/WasmEdge">WasmEdge</a> to provide HTTP support.
* Limited by WASM features today, e.g. no support for true parallelism.

</div>
<div>

```go
type Stack struct {
	CurrentTraces     ptrace.Traces
	CurrentMetrics    pmetric.Metrics
	CurrentLogs       plog.Logs
	ResultTraces      ptrace.Traces
	ResultMetrics     pmetric.Metrics
	ResultLogs        plog.Logs
	StatusReason      string
	RequestedShutdown atomic.Bool

	OnResultMetricsChange func(pmetric.Metrics)
	OnResultLogsChange    func(plog.Logs)
	OnResultTracesChange  func(ptrace.Traces)

	PluginConfigJSON []byte
}
```

</div>
</div>

<!-- PABLO 

The most interesting project out there combining the Collector and Wasm today is the otelwasm project.

It uses a thin wrapper to pass pdata data to components, and relies on unofficial extensions from WasmEdge, a CNCF Sandbox project, to provide full HTTP support.

We think it is very interesting, and, at the same time it shows that it is challenging to provide Wasm plugin support today that meets upstream's standards of performance and correctness. We look forward to see how WASIp3 developments allow otelwasm to evolve.

-->

---

# Collector running in Wasm

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
      <li>Filtering, sampling and transforming in the browser.</li>
      <li>Run it on a Wasm runtime for sandboxing.</li>
      <li>Run only some parts: <a href="https://ottl.run/">ottl.run</a>.</li>
    </ul>
  </div>
</div>

<!-- PABLO 

An alternative way to combine both is to run a whole Collector on Wasm runtime.

This could be on the browser, where you could leverage its processing capabilities, on your preferred Wasm runtime.

A small example of this can be seen today on the ottl.run website, which uses WebAssembly to run a small part of the Collector.

-->

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

<!-- PABLO 

What is the level of support for this upstream today?

Collector compilation targets are organized by tiers depending on the level of support we provide for them. 

Linux on amd64 belongs to our highest tier, <click>
followed by macOS, Windows and arm architectures <click>,
with a long tail of more niche architectures and since recently js/wasm <click>

This means any Collector release is 'guaranteed to build' on js/wasm, corresponding to browser Wasm runtimes.

There are some platforms that are not officially supported but that some people are using today, like Plan9 or AIX (that will change very soon). WASIp1 is one of these not officially supported platforms, although there are projects like otelwasm that use it today.
-->

---

# Collector running in Wasm: Upstream developments

<div class="icon-grid">
  <carbon-add-alt class="icon" />
  <span><code>js/wasm</code> added as a <a href="https://github.com/open-telemetry/opentelemetry-collector/blob/main/docs/platform-support.md#tier-3---community-support">Tier-3 platform</a> (Feb 2026).</span>
  <carbon-cut class="icon" />
  <span>Custom telemetry provider to strip down binary. (Feb 2026)</span>
  <carbon-chart-bar class="icon" />
  <span>244 of 271 (~90%) Collector components already compile to <code>js/wasm</code>.</span>
</div>

<!-- PABLO 

In particular, in preparation for this talk we focused on the js/wasm support, 
as well as building upon a feature for custom telemetry providers to be able to strip down Wasm binaries.

We can also happily report that about 90% of Collector components already compile to js/wasm.
-->


---

# Collector running in Wasm: OCB manifest

```yaml{all|1-2,6-8|1,3-4,10-12|14-16}
exporters:
  - gomod: go.opentelemetry.io/collector/exporter/otlphttpexporter v0.148.0
  - gomod: github.com/evan-bradley/kceu-2026-wasm-talk/jsexporter v0.0.0
    path: ../jsexporter

processors:
  - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/processor/deltatorateprocessor v0.148.0
  - gomod: github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor v0.148.0

receivers:
  - gomod: github.com/evan-bradley/kceu-2026-wasm-talk/jsreceiver v0.0.0
    path: ../jsreceiver

providers:
  - gomod: go.opentelemetry.io/collector/confmap/provider/httpprovider v1.54.0
  - gomod: go.opentelemetry.io/collector/confmap/provider/httpsprovider v1.54.0
```

<!-- 

  EVAN:

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

# Collector running in Wasm: Creating a Wasm binary

<div class="icon-grid">
  <carbon-terminal class="icon" />
  <span><code>GOOS=js GOARCH=wasm ocb --config manifest.yaml</code></span>
  <carbon-terminal class="icon" />
  <span><code>GOOS=wasip1 GOARCH=wasm ocb --config manifest.yaml</code></span>
</div>

<!-- 

EVAN: 
Cross-architecture compilation with Go is a breeze, so simply specify
the GOOS environment variable for your desired Wasm target and a GOARCH
of `wasm` to compile to a wasm binary.
 -->

---

# Collector running in Wasm: <a href="https://www.datadoghq.com/blog/engineering/agent-go-binaries/">gsa</a> analysis

<img src="/gsa.png" class="h-100 mx-auto" />

<!--
  EVAN

  Since we're talking about running in constrainted environments,
  we used the Go size analyzer tool to examine the binary and
  see if we could understand what's consuming space in the binary.

  We found it's mostly due to the Go runtime, which is the large
  pink box in the lower right.

  Next you have dependencies: all the Collector dependencies and
  its transitive dependencies. These are the dark purple and the
  light green in the upper left.

  Finally, the Go standard library modules are shown on the right,
  which consume a meaningful amount of space, but not as much
  as the other two.

  Ref:
  
  https://github.com/WebAssembly/design/blob/master/BinaryEncoding.md#data-section
  
  https://blog.tangrs.id.au/2022/02/15/notes-on-go-binary-metadata/
-->

---

# Collector running in Wasm: Limitations

<div class="icon-grid">
  <carbon-scale class="icon" />
  <span>Binaries are ≥ 38 MiB: 45% runtime and 55% dependencies</span>
  <carbon-archive class="icon" />
  <span>Our demo Wasm Collector is 66 MiB, but <b>13 MiB</b> compressed</span>
  <carbon-misuse class="icon" />
  <span>Limited TinyGo support: lack of complete stdlib.</span>
  <carbon-settings class="icon" />
  <!-- Source: https://webassembly.org/features/ -->
  <span>Limited Go support: no network (in WASI), concurrency, or Wasm GC.</span>
</div>

<!-- 

EVAN

Overall, binaries come out to 38 MiB uncompressed at a minimum.

The one in this demo is 66 MiB uncompressed. Compression helps: gzip reduces
the size to 13 MiB.

Most of this is in the Go runtime, and the rest of it either in standard library modules
or in dependencies.

TinyGo can help reduce binary sizes: it's an alternative compiler compliant with
the Go language spec targeting restricted platforms like Wasm. However, its
implementation of the Go standard library isn't complete enough to compile
the Collector.

Likewise, Go doesn't support all Wasm features, so some things need to happen
on the Go side before the Collector is able to benefit from them.


-->

---

# Computing on the edge

<div class="edge-continuum">
  <div class="edge-track-label">Far edge</div>
  <div class="edge-track-label">Core</div>
  <div class="edge-track">
    <div class="edge-stage far-edge">
      <carbon-mobile class="stage-icon" />
      <div class="stage-title">User Device</div>
      <div class="stage-subtitle">Browser, mobile app, desktop app</div>
    </div>
    <div class="edge-arrow" aria-hidden="true"></div>
    <div class="edge-stage platform-edge">
      <carbon-edge-node class="stage-icon" />
      <div class="stage-title">Edge Platform</div>
      <div class="stage-subtitle">Edge function</div>
    </div>
    <div class="edge-arrow" aria-hidden="true"></div>
    <div class="edge-stage gateway-edge">
      <carbon-router class="stage-icon" />
      <div class="stage-title">Middleware / Gateway</div>
      <div class="stage-subtitle">Gateway, message broker, regional hub</div>
    </div>
    <div class="edge-arrow" aria-hidden="true"></div>
    <div class="edge-stage core-edge">
      <carbon-data-center class="stage-icon" />
      <div class="stage-title">Core Infrastructure</div>
      <div class="stage-subtitle">Central services, control plane, storage</div>
    </div>
  </div>
</div>

<!-- EVAN -->


---
transition: fade
---

# Computing on the edge

<div class="edge-fanin">
  <div class="edge-track-label">Far edge</div>
  <div class="edge-track-label">Core</div>
  <div class="edge-fanin-grid">
    <div class="fanin-col devices">
      <div class="edge-stage far-edge row-1 fan-arrow arrow-down-strong">
        <carbon-mobile class="stage-icon" />
        <div class="stage-title">Browser app</div>
      </div>
      <div class="edge-stage far-edge row-2 fan-arrow arrow-down-soft">
        <carbon-mobile class="stage-icon" />
        <div class="stage-title">Mobile app</div>
      </div>
      <div class="edge-stage far-edge row-3 fan-arrow arrow-up-soft">
        <carbon-mobile class="stage-icon" />
        <div class="stage-title">Desktop app</div>
      </div>
      <div class="edge-stage far-edge row-4 fan-arrow arrow-up-strong">
        <carbon-mobile class="stage-icon" />
        <div class="stage-title">IoT device</div>
      </div>
    </div>
    <div class="fanin-col platforms">
      <div class="edge-stage platform-edge row-p1 fan-arrow arrow-down-soft">
        <carbon-edge-node class="stage-icon" />
        <div class="stage-title">Edge function</div>
      </div>
      <div class="edge-stage platform-edge row-p2 fan-arrow arrow-flat">
        <carbon-edge-node class="stage-icon" />
        <div class="stage-title">Edge function</div>
      </div>
      <div class="edge-stage platform-edge row-p3 fan-arrow arrow-up-soft">
        <carbon-edge-node class="stage-icon" />
        <div class="stage-title">Edge function</div>
      </div>
    </div>
    <div class="fanin-col gateways">
      <div class="edge-stage gateway-edge row-2 fan-arrow arrow-down-soft">
        <carbon-router class="stage-icon" />
        <div class="stage-title">Gateway</div>
      </div>
      <div class="edge-stage gateway-edge row-3 fan-arrow arrow-up-soft">
        <carbon-router class="stage-icon" />
        <div class="stage-title">Middleware</div>
      </div>
    </div>
    <div class="fanin-col core">
      <div class="edge-stage core-edge row-core">
        <carbon-data-center class="stage-icon" />
        <div class="stage-title">Core Infrastructure</div>
      </div>
    </div>
  </div>
</div>

<!-- EVAN -->

---
transition: fade
---

# Computing on the edge

<div class="edge-fanin">
  <div class="edge-track-label">Far edge</div>
  <div class="edge-track-label">Core</div>
  <div class="edge-fanin-grid">
    <div class="fanin-col devices fanin-dimmed">
      <div class="edge-stage far-edge row-1 fan-arrow arrow-down-strong">
        <carbon-mobile class="stage-icon" />
        <div class="stage-title">Browser app</div>
      </div>
      <div class="edge-stage far-edge row-2 fan-arrow arrow-down-soft">
        <carbon-mobile class="stage-icon" />
        <div class="stage-title">Mobile app</div>
      </div>
      <div class="edge-stage far-edge row-3 fan-arrow arrow-up-soft">
        <carbon-mobile class="stage-icon" />
        <div class="stage-title">Desktop app</div>
      </div>
      <div class="edge-stage far-edge row-4 fan-arrow arrow-up-strong">
        <carbon-mobile class="stage-icon" />
        <div class="stage-title">IoT device</div>
      </div>
    </div>
    <div class="fanin-col platforms fanin-dimmed">
      <div class="edge-stage platform-edge row-p1 fan-arrow arrow-down-soft">
        <carbon-edge-node class="stage-icon" />
        <div class="stage-title">Edge function</div>
      </div>
      <div class="edge-stage platform-edge row-p2 fan-arrow arrow-flat">
        <carbon-edge-node class="stage-icon" />
        <div class="stage-title">Edge function</div>
      </div>
      <div class="edge-stage platform-edge row-p3 fan-arrow arrow-up-soft">
        <carbon-edge-node class="stage-icon" />
        <div class="stage-title">Edge function</div>
      </div>
    </div>
    <div class="fanin-col gateways">
      <div class="edge-stage gateway-edge row-2 fan-arrow arrow-down-soft">
        <carbon-router class="stage-icon" />
        <div class="stage-title">Gateway</div>
      </div>
      <div class="edge-stage gateway-edge row-3 fan-arrow arrow-up-soft">
        <carbon-router class="stage-icon" />
        <div class="stage-title">Middleware</div>
      </div>
    </div>
    <div class="fanin-col core">
      <div class="edge-stage core-edge row-core">
        <carbon-data-center class="stage-icon" />
        <div class="stage-title">Core Infrastructure</div>
      </div>
    </div>
  </div>
</div>

<!-- EVAN -->

---
transition: slide-left
---

# Computing on the edge

<div class="edge-fanin">
  <div class="edge-track-label">Far edge</div>
  <div class="edge-track-label">Core</div>
  <div class="edge-fanin-grid">
    <div class="fanin-col devices">
      <div class="edge-stage far-edge row-1 fan-arrow arrow-down-strong">
        <carbon-mobile class="stage-icon" />
        <div class="stage-title">Browser app</div>
      </div>
      <div class="edge-stage far-edge row-2 fan-arrow arrow-down-soft">
        <carbon-mobile class="stage-icon" />
        <div class="stage-title">Mobile app</div>
      </div>
      <div class="edge-stage far-edge row-3 fan-arrow arrow-up-soft">
        <carbon-mobile class="stage-icon" />
        <div class="stage-title">Desktop app</div>
      </div>
      <div class="edge-stage far-edge row-4 fan-arrow arrow-up-strong">
        <carbon-mobile class="stage-icon" />
        <div class="stage-title">IoT device</div>
      </div>
    </div>
    <div class="fanin-col platforms">
      <div class="edge-stage platform-edge row-p1 fan-arrow arrow-down-soft">
        <carbon-edge-node class="stage-icon" />
        <div class="stage-title">Edge function</div>
      </div>
      <div class="edge-stage platform-edge row-p2 fan-arrow arrow-flat">
        <carbon-edge-node class="stage-icon" />
        <div class="stage-title">Edge function</div>
      </div>
      <div class="edge-stage platform-edge row-p3 fan-arrow arrow-up-soft">
        <carbon-edge-node class="stage-icon" />
        <div class="stage-title">Edge function</div>
      </div>
    </div>
    <div class="fanin-col gateways">
      <div class="edge-stage gateway-edge row-2 fan-arrow arrow-down-soft">
        <carbon-router class="stage-icon" />
        <div class="stage-title">Gateway</div>
      </div>
      <div class="edge-stage gateway-edge row-3 fan-arrow arrow-up-soft">
        <carbon-router class="stage-icon" />
        <div class="stage-title">Middleware</div>
      </div>
    </div>
    <div class="fanin-col core">
      <div class="edge-stage core-edge row-core">
        <carbon-data-center class="stage-icon" />
        <div class="stage-title">Core Infrastructure</div>
      </div>
    </div>
  </div>
</div>

<!-- EVAN -->

---
transition: slide-left
---

# Observability without borders

<div class="icon-grid">
  <carbon-application-web class="icon" />
  <span>Running in a browser (Wasm).</span>
  <carbon-container-software class="icon" />
  <span>Running in Wasm runtimes (WASI).</span>
</div>

<!-- 

EVAN: We've covered how the Collector and Wasm can work together, but
where can you use this?

1. The option we've found works the best out of the box right now,
   and probably the most unexpected one, is running it in the browser.
   This has the most restrictions, but allows you to run Collectors
   on a user's machine.
2. There are also a number of server-side Wasm runtimes available
   that boast low startup times and effective sandboxing without
   container technology like Docker. These can allow you to run
   apps in unusual places like edge functions (e.g. Cloudflare workers)
3. You can also use the Collecotr in-process in one of your applications
   using a Wasm runtime built for that language.


The point of all of this is that you can leverage Wasm to run your
Collectors in the odd nooks and crannies of your infrastructure,
which if used tactfully, may open up new possibilities for
your telemetry pipelines.

Let's cover some of the trade-offs that each of these options provides.

 -->

---

# Observability without borders: browser

<div class="icon-grid">
  <carbon-folder-off class="icon" />
  <span>No FS access.</span>
  <carbon-close-outline class="icon" />
  <span>Can't open ports.</span>
  <carbon-application class="icon" />
  <span>Uses: JS SDK processing supplement, Electron/thick-client apps.</span>
</div>

<!-- 

EVAN: Running in the browser puts the Collector directly on your user's device.
It doesn't get more on the edge than this.

There's no filesystem access at this layer, and you can't open ports,
but you can still make network calls to send or receive data.

You can use this to supplement processing in your JS SDK if you're
working with a JS-based webapp, or possibly for another language's
SDK if you're running your app in Wasm!

We think this will most likely find use in thick-client applications
like Electron apps, where there applications are large and more likely
to want to use the Collector for local processing.

 -->

---

# Observability without borders: Wasm runtime

<div class="icon-grid">
  <carbon-wifi-off class="icon" />
  <span>Limited/no networking currently (Go only supports WASIp1).</span>
  <carbon-folder class="icon" />
  <span>Filesystem access is available if the host grants it.</span>
  <carbon-edge-node-alt class="icon" />
  <span>For use alongside other Wasm applications or in edge functions.</span>
  <carbon-code class="icon" />
  <span>Also can use in an in-process Wasm runtime.</span>
</div>

<!-- 

EVAN

The Collector can be run inside a non-browser Wasm runtime.

It's important to note that as of today, Go only supports compiling
to WASIp1, which doesn't incude networking capabilities. You can still
read from the filesystem or export functions to be called from a Wasm
runtime.

This is going to be most useful if you include the Collector as part of
a Wasm application composed of multiple modules, or for use in edge
function runtimes.

You can also use it for in-process processing as part of another
application. Many languages, including Java and Rust, have libraries
that are Wasm runtimes, and can let you call a Collector pipeline
like you would a function. This is an advanced use case, but could
be used for a custom telemetry pipeline application or to supplement
an SDK like we've shown for the JS SDK.

Further reading: https://go.dev/blog/wasmexport

 -->

---

# Observability without borders: what it's not

<div class="icon-grid">
  <carbon-code class="icon" />
  <span>Doesn't replace OTel SDKs.</span>
  <carbon-container-software class="icon" />
  <span>Unlikely to replace most existing Collector deployments.</span>
  <carbon-floorplan class="icon" />
  <span>Not a working solution, just a blueprint.</span>
</div>

<!-- 

EVAN: Since this is an advanced use case, we want to very clearly call out what
this is NOT.

First, you're not going to replace OTel SDKs with this, and in most cases
should simplify and stick with an SDK and use SDK processors if possible.

You're also not going to likely want to go and switch your Collector
deployment architecture after seeing this presentation. We have tested
ourselves and seen architectures used by users that are battle-tested
in production environments that should be the default recommendations
for most users. Think of this as a way to open possibilities for
maximizing your telemetry pipeline's capabilities.

Finally, nothing we've shown here is a working solution ready for
production right now. While Wasm is currently used in production
environments for large, established applications as we've shown,
what we're showing you today is on the bleeding edge of what's
possible. Again, there is official support for this, so we would
love to get your ideas and contributions for what comes next!

 -->

---

# Looking ahead

<div class="icon-grid">
  <carbon-in-progress class="icon" />
  <span>Go WASIp3 support still under active discussion.</span>
  <carbon-package class="icon" />
  <span>Wider TinyGo stdlib support could allow for smaller binaries.</span>
  <carbon-microphone class="icon" />
  <span>WASI OTel (<a href="https://sched.co/2DY17">WasmCon talk earlier today</a>).</span>
  <carbon-group class="icon" />
  <span>Contributions from YOU in the audience!</span>
</div>

<!-- 

EVAN: Looking ahead, here are some areas where we have seen active development,
or where there needs to be active developments to take this further.

 -->

---

# Demo

<div class="icon-grid">
  <carbon-assembly-reference class="icon" />
  <span>We compiled a basic Collector that communicates with the OTel JS SDK.</span>
  <carbon-mobile class="icon" />
  <span>It runs on any modern browser, so try it on your phone!</span>
  <QrArrow />
</div>

<!-- 

BOTH: To hopefully help demonstrate the cool factor of what's possible
with Wasm, we created a small demo that runs right inside these slides.

 -->

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
          <span class="arch-label-right arch-wasm">Wasm</span>
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

<!-- BOTH -->

---

<WasmDemo />

<!-- BOTH -->

---
layout: center
---

<h1 class="qa-title">Q&A</h1>

<QrArrow />

<img src="/kceu26.svg" class="kceu-logo" />

<!-- BOTH -->
