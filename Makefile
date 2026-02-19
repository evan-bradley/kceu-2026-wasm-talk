otelwasmcol/bin/main.wasm: otelwasmcol/wasm-manifest.yaml
	cd otelwasmcol && go tool go.opentelemetry.io/collector/cmd/builder --skip-compilation --config=wasm-manifest.yaml
	cd otelwasmcol/bin && CGO_ENABLED=0 GOOS=js GOARCH=wasm go build -ldflags '-s -w' -trimpath -o main.wasm

otelwasmcol: otelwasmcol/bin/main.wasm

slides/node_modules: slides/package.json slides/package-lock.json slides/pnpm-lock.yaml
	cd slides && npx pnpm install

slides/dist/index.html: otelwasmcol slides/slides.md slides/style.css slides/vite.config.ts slides/package.json slides/node_modules
	cp otelwasmcol/bin/main.wasm slides/public/otelwasmcol.wasm
	gzip slides/public/otelwasmcol.wasm
	cd slides && npx pnpm build

build: slides/dist/index.html

deploy: build
	cd slides && npx pnpm run deploy

.PHONY: serve
serve: build
	cd slides && npx pnpm dev