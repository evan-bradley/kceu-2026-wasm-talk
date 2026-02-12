build: otelwasmcol/wasm-manifest.yaml
	cd otelwasmcol && go tool go.opentelemetry.io/collector/cmd/builder --skip-compilation --config=wasm-manifest.yaml
	cd otelwasmcol/bin && CGO_ENABLED=0 GOOS=js GOARCH=wasm go build -ldflags '-s -w' -trimpath -o main.wasm
	cp otelwasmcol/bin/main.wasm slides/public/otelwasmcol.wasm

.PHONY: serve
serve: build
	cd slides && npm run dev