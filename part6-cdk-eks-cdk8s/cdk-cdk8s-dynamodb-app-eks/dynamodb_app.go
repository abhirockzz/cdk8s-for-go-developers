package main

import (
	"log"
	"os"
	"strconv"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus22/v2"
)

type AppChartProps struct {
	cdk8s.ChartProps
	serviceAccountName string
	image              *string
	appPort            int
	lbPort             int
	tableName          *string
	region             string
}

func NewAppChartProps(image, tableName *string) AppChartProps {

	serviceAccountName := os.Getenv("SERVICE_ACCOUNT_NAME")
	if serviceAccountName == "" {
		log.Fatal("missing env variable SERVICE_ACCOUNT_NAME")
	}

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		log.Fatal("missing env variable APP_PORT")
	}
	appPortNum, _ := strconv.Atoi(appPort)

	lbPort := os.Getenv("LB_PORT")
	if lbPort == "" {
		log.Println("missing env variable LB_PORT. setting it to APP_PORT", appPort)
		lbPort = appPort
	}
	lbPortNum, _ := strconv.Atoi(lbPort)

	region := os.Getenv("AWS_REGION")
	if region == "" {
		log.Fatal("missing env variable AWS_REGION")
	}

	return AppChartProps{
		serviceAccountName: serviceAccountName,
		image:              image,
		appPort:            appPortNum,
		lbPort:             lbPortNum,
		tableName:          tableName,
		region:             region}

}

func NewDynamoDBAppChart(scope constructs.Construct, id string, props *AppChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	dep := cdk8splus22.NewDeployment(chart, jsii.String("dynamodb-app-deployment"), &cdk8splus22.DeploymentProps{
		Metadata: &cdk8s.ApiObjectMetadata{
			Name: jsii.String("dynamodb-app")},
		ServiceAccount: cdk8splus22.ServiceAccount_FromServiceAccountName(
			chart,
			jsii.String("aws-irsa"),
			jsii.String(props.serviceAccountName))})

	container := dep.AddContainer(
		&cdk8splus22.ContainerProps{
			Name:  jsii.String("dynamodb-app-container"),
			Image: props.image,
			Port:  jsii.Number(float64(props.appPort))})

	container.Env().AddVariable(jsii.String("TABLE_NAME"), cdk8splus22.EnvValue_FromValue(props.tableName))

	container.Env().AddVariable(jsii.String("AWS_REGION"), cdk8splus22.EnvValue_FromValue(&props.region))

	dep.ExposeViaService(
		&cdk8splus22.DeploymentExposeViaServiceOptions{
			Name:        jsii.String("dynamodb-app-service"),
			ServiceType: cdk8splus22.ServiceType_LOAD_BALANCER,
			Ports: &[]*cdk8splus22.ServicePort{
				{Protocol: cdk8splus22.Protocol_TCP,
					Port:       jsii.Number(float64(props.lbPort)),
					TargetPort: jsii.Number(float64(props.appPort))}}})

	return chart
}
