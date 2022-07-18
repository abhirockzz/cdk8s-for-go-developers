package main

import (
	"cdk8s-crd/imports/samplecontrollerk8sio"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

type FooChartProps struct {
	cdk8s.ChartProps
}

func NewFooChart(scope constructs.Construct, id string, props *FooChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	samplecontrollerk8sio.NewFoo(chart, jsii.String("foo1"), &samplecontrollerk8sio.FooProps{Spec: &samplecontrollerk8sio.FooSpec{DeploymentName: jsii.String("foo1-dep"), Replicas: jsii.Number(2)}})

	return chart
}
