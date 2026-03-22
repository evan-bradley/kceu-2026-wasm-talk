//go:build js && wasm

package jsreceiver

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

var componentType = component.MustNewType("js")

// NewFactory creates a factory for the JS receiver.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		componentType,
		createDefaultConfig,
		receiver.WithMetrics(createMetrics, component.StabilityLevelDevelopment),
	)
}

func createDefaultConfig() component.Config {
	return &Config{}
}

func createMetrics(_ context.Context, set receiver.Settings, cfg component.Config, consumer consumer.Metrics) (receiver.Metrics, error) {
	return newJSReceiver(set, consumer), nil
}
