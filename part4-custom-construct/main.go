package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

type MyChartProps struct {
	cdk8s.ChartProps
}

func NewMyChart(scope constructs.Construct, id string, props *MyChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	NewWordpressStack(chart, jsii.String("wordpress-stack"), &WordpressProps{
		MySQLImage:       jsii.String("mariadb"),
		MySQLPassword:    jsii.String("Password123"),
		MySQLStorage:     jsii.Number(3),
		WordpressImage:   jsii.String("wordpress:4.8-apache"),
		WordpressStorage: jsii.Number(2)})

	return chart
}

func main() {
	app := cdk8s.NewApp(nil)
	NewMyChart(app, "wordpress-custom-stack", nil)
	app.Synth()
}
