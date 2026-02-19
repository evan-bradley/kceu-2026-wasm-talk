otelwasmcol/bin/main.wasm: otelwasmcol/wasm-manifest.yaml
	cd otelwasmcol && go tool go.opentelemetry.io/collector/cmd/builder --skip-compilation --config=wasm-manifest.yaml
	cd otelwasmcol/bin && CGO_ENABLED=0 GOOS=js GOARCH=wasm go build -ldflags '-s -w' -trimpath -o main.wasm

otelwasmcol: otelwasmcol/bin/main.wasm

slides/dist/index.html: otelwasmcol slides/slides.md slides/style.css slides/vite.config.ts slides/package.json
	cp otelwasmcol/bin/main.wasm slides/public/otelwasmcol.wasm
	gzip slides/public/otelwasmcol.wasm
	cd slides && npm run build

build: slides/dist/index.html

deploy: build
	cd slides && npm run deploy

.PHONY: serve
serve: build
	cd slides && npm run dev