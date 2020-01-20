package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/lanzafame/otel-met-libs/pkg1"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporter/metric/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/batcher/defaultkeys"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

var (
	LemonsKey = key.New("ex.com/lemons")
)

func initMeter() *push.Controller {
	selector := simple.NewWithExactMeasure()
	exporter, err := prometheus.NewExporter(prometheus.Options{
		DefaultHistogramBuckets: []float64{0., 10., 15., 20.},
	})

	if err != nil {
		log.Panicf("failed to initialize metric stdout exporter %v", err)
	}
	batcher := defaultkeys.New(selector, sdkmetric.NewDefaultLabelEncoder(), false)
	pusher := push.New(batcher, exporter, time.Second)
	pusher.Start()

	go func() {
		_ = http.ListenAndServe(":2222", exporter)
	}()

	global.SetMeterProvider(pusher)
	return pusher
}

func main() {
	defer initMeter().Stop()

	meter := global.MeterProvider().Meter("ex.com/basic")

	measureTwo := meter.NewFloat64Measure("ex.com.two", metric.WithKeys(key.New("A")))
	measureThree := meter.NewFloat64Counter("ex.com.three")

	commonLabels := meter.Labels(LemonsKey.Int(10), key.String("A", "1"), key.String("B", "2"), key.String("C", "3"))
	notSoCommonLabels := meter.Labels(LemonsKey.Int(13))

	ctx := context.Background()

	pkg1.InitPkg(meter)

	for i := 0; i < 100; i++ {
		pkg1.Test(ctx, float64(i))
	}

	meter.RecordBatch(
		ctx,
		commonLabels,
		measureTwo.Measurement(2.0),
		measureThree.Measurement(12.0),
	)

	meter.RecordBatch(
		ctx,
		notSoCommonLabels,
		measureTwo.Measurement(2.0),
		measureThree.Measurement(22.0),
	)

	meter.RecordBatch(
		ctx,
		commonLabels,
		measureTwo.Measurement(12.0),
		measureThree.Measurement(13.0),
	)

	time.Sleep(60 * time.Second)
}
