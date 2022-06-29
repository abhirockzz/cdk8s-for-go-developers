package main

import (
	"getting-started-with-cdk8s-go/imports/k8s"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
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

	selector := &k8s.LabelSelector{MatchLabels: &map[string]*string{"app": jsii.String("hello-nginx")}}

	labels := &k8s.ObjectMeta{Labels: &map[string]*string{"app": jsii.String("hello-nginx")}}

	nginxContainer := &k8s.Container{Name: jsii.String("nginx-container"), Image: jsii.String("nginx"), Ports: &[]*k8s.ContainerPort{{ContainerPort: jsii.Number(80)}}}

	deployment := k8s.NewKubeDeployment(chart, jsii.String("deployment"),
		&k8s.KubeDeploymentProps{
			Spec: &k8s.DeploymentSpec{
				Replicas: jsii.Number(1),
				Selector: selector,
				Template: &k8s.PodTemplateSpec{
					Metadata: labels,
					Spec: &k8s.PodSpec{
						Containers: &[]*k8s.Container{nginxContainer}}}}})

	service := k8s.NewKubeService(chart, jsii.String("service"), &k8s.KubeServiceProps{
		Spec: &k8s.ServiceSpec{
			Type:     jsii.String("LoadBalancer"),
			Ports:    &[]*k8s.ServicePort{{Port: jsii.Number(9090), TargetPort: k8s.IntOrString_FromNumber(jsii.Number(80))}},
			Selector: &map[string]*string{"app": jsii.String("hello-nginx")}}})

	deployment.AddDependency(service)

	/*cdk8s.NewInclude(chart, jsii.String("existing service"), &cdk8s.IncludeProps{Url: jsii.String("service.yaml")})

	cdk8s.NewHelm(chart, jsii.String("bitnami nginx helm chart"), &cdk8s.HelmProps{
	Chart:  jsii.String("bitnami/nginx"),
	Values: &map[string]interface{}{"service.type": "ClusterIP"}})*/

	return chart
}

func main() {
	app := cdk8s.NewApp(nil)
	NewNginxChart(app, "nginx", nil)
	app.Synth()
}
