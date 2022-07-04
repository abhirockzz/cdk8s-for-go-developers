package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus22/v2"
)

type NginxChartProps struct {
	cdk8s.ChartProps
}

func NewNginxChart(scope constructs.Construct, id string, props *NginxChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	dep := cdk8splus22.NewDeployment(chart, jsii.String("deployment"), &cdk8splus22.DeploymentProps{Metadata: &cdk8s.ApiObjectMetadata{Name: jsii.String("nginx-deployment-cdk8s-plus")}})

	dep.AddContainer(&cdk8splus22.ContainerProps{
		Name:  jsii.String("nginx-container"),
		Image: jsii.String("nginx"),
		Port:  jsii.Number(80)})

	dep.ExposeViaService(&cdk8splus22.DeploymentExposeViaServiceOptions{
		Name:        jsii.String("nginx-container-service"),
		ServiceType: cdk8splus22.ServiceType_LOAD_BALANCER,
		Ports:       &[]*cdk8splus22.ServicePort{{Port: jsii.Number(9090), TargetPort: jsii.Number(80)}}})

	return chart
}

func main() {
	app := cdk8s.NewApp(nil)
	NewNginxChart(app, "nginx-cdk8s-plus", nil)
	app.Synth()
}
