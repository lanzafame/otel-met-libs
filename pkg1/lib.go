package pkg1

import (
	"context"

	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
)

var (
	FooKey = key.New("ex.com/foo")
	BarKey = key.New("ex.com/bar")
)

var meter metric.Meter
var oneMetric metric.Float64Gauge
var commonLabels metric.LabelSet

func InitPkg(m metric.Meter) {
	meter = m
	oneMetric = meter.NewFloat64Gauge("ex.com.one",
		metric.WithKeys(FooKey, BarKey),
		metric.WithDescription("A gauge set to 1.0"),
	)
	commonLabels = meter.Labels(FooKey.Int(1), BarKey.Int(2))
}

func Test(ctx context.Context, value float64) {
	meter.RecordBatch(ctx, commonLabels, oneMetric.Measurement(value))
}
