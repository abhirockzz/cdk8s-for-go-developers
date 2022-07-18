package main

import (
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

type MyChartProps struct {
	cdk8s.ChartProps
}

func main() {
	app := cdk8s.NewApp(nil)
	//NewFooChart(app, "FooApp", nil)
	NewKafkaChart(app, "KafkaApp", nil)
	app.Synth()
}
