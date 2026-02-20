package jsreceiver

import (
	"context"
	"syscall/js"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
)

type jsReceiver struct {
	set         receiver.Settings
	consumer    consumer.Metrics
	unmarshaler pmetric.JSONUnmarshaler
	callback    js.Func
}

func newJSReceiver(set receiver.Settings, consumer consumer.Metrics) *jsReceiver {
	return &jsReceiver{
		set:      set,
		consumer: consumer,
	}
}

func (r *jsReceiver) Start(_ context.Context, _ component.Host) error {
	// Register a JS function that the browser can call to send metrics into the Collector.
	// JS calls: globalThis.__otelReceiveCallback(uint8Array)
	r.callback = js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 1 {
			return nil
		}

		// Copy bytes from JS Uint8Array to Go using js.CopyBytesToGo
		jsArray := args[0]
		length := jsArray.Get("length").Int()
		buf := make([]byte, length)
		js.CopyBytesToGo(buf, jsArray)

		// Unmarshal the OTLP JSON metrics
		metrics, err := r.unmarshaler.UnmarshalMetrics(buf)
		if err != nil {
			r.set.Logger.Error("failed to unmarshal metrics from JS: " + err.Error())
			return nil
		}

		// Pass to the consumer pipeline (run in a goroutine to avoid blocking JS event loop)
		go func() {
			if err := r.consumer.ConsumeMetrics(context.Background(), metrics); err != nil {
				r.set.Logger.Error("failed to consume metrics: " + err.Error())
			}
		}()

		return nil
	})

	js.Global().Set("__otelReceiveCallback", r.callback)
	return nil
}

func (r *jsReceiver) Shutdown(_ context.Context) error {
	js.Global().Delete("__otelReceiveCallback")
	r.callback.Release()
	return nil
}
