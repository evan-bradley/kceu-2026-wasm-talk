build: otelwasmcol/wasm-manifest.yaml
	cd otelwasmcol && go tool go.opentelemetry.io/collector/cmd/builder --skip-compilation --config=wasm-manifest.yaml
	cd otelwasmcol/bin && CGO_ENABLED=0 GOOS=js GOARCH=wasm go build -ldflags '-s -w' -trimpath -o main.wasm

.PHONY: serve
serve:
	python -m http.server 8000