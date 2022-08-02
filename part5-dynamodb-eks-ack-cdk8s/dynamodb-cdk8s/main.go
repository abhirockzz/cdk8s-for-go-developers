package main

import (
	"log"
	"os"

	ddbcrd "example.com/dynamodb-cdk8s/imports/dynamodbservicesk8saws"
	"example.com/dynamodb-cdk8s/imports/servicesk8saws"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus22/v2"
)

type MyChartProps struct {
	cdk8s.ChartProps
}

var tableName string
var serviceAccountName string
var image string

func init() {
	tableName = os.Getenv("TABLE_NAME")
	if tableName == "" {
		log.Fatal("missing environment variable TABLE_NAME")
	}

	serviceAccountName = os.Getenv("SERVICE_ACCOUNT")
	if serviceAccountName == "" {
		log.Fatal("missing environment variable SERVICE_ACCOUNT")
	}

	image = os.Getenv("DOCKER_IMAGE")
	if image == "" {
		log.Fatal("missing environment variable DOCKER_IMAGE")
	}
}

const primaryKeyName = "shortcode"
const billingMode = "PAY_PER_REQUEST"
const hashKeyType = "HASH"

const appPort = 8080
const lbPort = 9090

const fieldExportNameForTable = "export-dynamodb-tablename"
const fieldExportNameForRegion = "export-dynamodb-region"
const configMapName = "export-dynamodb-urls-info"

var cfgMap cdk8splus22.ConfigMap

var fieldExportForTable servicesk8saws.FieldExport
var fieldExportForRegion servicesk8saws.FieldExport

func NewDynamoDBChart(scope constructs.Construct, id string, props *MyChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	table := ddbcrd.NewTable(chart, jsii.String("dynamodb-ack-cdk8s-table"), &ddbcrd.TableProps{
		Spec: &ddbcrd.TableSpec{
			AttributeDefinitions: &[]*ddbcrd.TableSpecAttributeDefinitions{
				{AttributeName: jsii.String(primaryKeyName), AttributeType: jsii.String("S")}},
			BillingMode: jsii.String(billingMode),
			TableName:   jsii.String(tableName),
			KeySchema: &[]*ddbcrd.TableSpecKeySchema{
				{AttributeName: jsii.String(primaryKeyName),
					KeyType: jsii.String(hashKeyType)}}}})

	cfgMap = cdk8splus22.NewConfigMap(chart, jsii.String("config-map"),
		&cdk8splus22.ConfigMapProps{
			Metadata: &cdk8s.ApiObjectMetadata{
				Name: jsii.String(configMapName)}})

	fieldExportForTable = servicesk8saws.NewFieldExport(chart, jsii.String("fexp-table"), &servicesk8saws.FieldExportProps{
		Metadata: &cdk8s.ApiObjectMetadata{Name: jsii.String(fieldExportNameForTable)},
		Spec: &servicesk8saws.FieldExportSpec{
			From: &servicesk8saws.FieldExportSpecFrom{Path: jsii.String(".spec.tableName"),
				Resource: &servicesk8saws.FieldExportSpecFromResource{
					Group: jsii.String("dynamodb.services.k8s.aws"),
					Kind:  jsii.String("Table"),
					Name:  table.Name()}},
			To: &servicesk8saws.FieldExportSpecTo{
				Name: cfgMap.Name(),
				Kind: servicesk8saws.FieldExportSpecToKind_CONFIGMAP}}})

	fieldExportForRegion = servicesk8saws.NewFieldExport(chart, jsii.String("fexp-region"), &servicesk8saws.FieldExportProps{
		Metadata: &cdk8s.ApiObjectMetadata{Name: jsii.String(fieldExportNameForRegion)},
		Spec: &servicesk8saws.FieldExportSpec{
			From: &servicesk8saws.FieldExportSpecFrom{
				Path: jsii.String(".status.ackResourceMetadata.region"),
				Resource: &servicesk8saws.FieldExportSpecFromResource{
					Group: jsii.String("dynamodb.services.k8s.aws"),
					Kind:  jsii.String("Table"),
					Name:  table.Name()}},
			To: &servicesk8saws.FieldExportSpecTo{
				Name: cfgMap.Name(),
				Kind: servicesk8saws.FieldExportSpecToKind_CONFIGMAP}}})

	return chart
}

func NewDeploymentChart(scope constructs.Construct, id string, props *MyChartProps) cdk8s.Chart {
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
			jsii.String(serviceAccountName))})

	container := dep.AddContainer(
		&cdk8splus22.ContainerProps{
			Name:  jsii.String("dynamodb-app-container"),
			Image: jsii.String(image),
			Port:  jsii.Number(appPort)})

	container.Env().AddVariable(jsii.String("TABLE_NAME"), cdk8splus22.EnvValue_FromConfigMap(
		cfgMap, jsii.String("default."+*fieldExportForTable.Name()),
		&cdk8splus22.EnvValueFromConfigMapOptions{Optional: jsii.Bool(false)}))

	container.Env().AddVariable(jsii.String("AWS_REGION"), cdk8splus22.EnvValue_FromConfigMap(
		cfgMap, jsii.String("default."+*fieldExportForRegion.Name()),
		&cdk8splus22.EnvValueFromConfigMapOptions{Optional: jsii.Bool(false)}))

	dep.ExposeViaService(
		&cdk8splus22.DeploymentExposeViaServiceOptions{
			Name:        jsii.String("dynamodb-app-service"),
			ServiceType: cdk8splus22.ServiceType_LOAD_BALANCER,
			Ports: &[]*cdk8splus22.ServicePort{
				{Protocol: cdk8splus22.Protocol_TCP,
					Port:       jsii.Number(lbPort),
					TargetPort: jsii.Number(appPort)}}})

	return chart
}

func main() {
	app := cdk8s.NewApp(nil)

	dynamodDB := NewDynamoDBChart(app, "dynamodb", nil)
	deployment := NewDeploymentChart(app, "deployment", nil)

	deployment.AddDependency(dynamodDB)

	app.Synth()
}
