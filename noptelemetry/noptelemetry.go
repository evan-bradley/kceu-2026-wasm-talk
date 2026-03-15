package noptelemetry

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/service/telemetry"
)

func NewFactory() telemetry.Factory {
	return telemetry.NewFactory(
		func() component.Config { return struct{}{} },
	)
}
