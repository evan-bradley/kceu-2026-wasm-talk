//go:build js && wasm

package jsexporter

import (
	"context"
	"syscall/js"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type jsExporter struct {
	metricsMarshaler pmetric.JSONMarshaler
	tracesMarshaler  ptrace.JSONMarshaler
	logsMarshaler    plog.JSONMarshaler
}

func newJSExporter() *jsExporter {
	return &jsExporter{}
}

func (e *jsExporter) pushMetrics(_ context.Context, md pmetric.Metrics) error {
	data, err := e.metricsMarshaler.MarshalMetrics(md)
	if err != nil {
		return err
	}
	return copyToJS(data)
}

func (e *jsExporter) pushTraces(_ context.Context, td ptrace.Traces) error {
	data, err := e.tracesMarshaler.MarshalTraces(td)
	if err != nil {
		return err
	}
	return copyToJS(data)
}

func (e *jsExporter) pushLogs(_ context.Context, ld plog.Logs) error {
	data, err := e.logsMarshaler.MarshalLogs(ld)
	if err != nil {
		return err
	}
	return copyToJS(data)
}

// copyToJS copies the serialized JSON bytes to the JavaScript runtime
// using js.CopyBytesToJS and invokes a global callback.
func copyToJS(data []byte) error {
	uint8Array := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(uint8Array, data)

	callback := js.Global().Get("__otelExportCallback")
	if !callback.IsUndefined() && !callback.IsNull() {
		callback.Invoke(uint8Array)
	}
	return nil
}
