package main

import (
	"log"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecrassets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"

	"github.com/aws/constructs-go/constructs/v10"
)

type CdkStackProps struct {
	awscdk.StackProps
}

var eksClusterName string
var kubectlRoleARN string

func init() {
	eksClusterName = os.Getenv("EKS_CLUSTER_NAME")
	if eksClusterName == "" {
		log.Fatal("missing env variable EKS_CLUSTER_NAME")
	}

	kubectlRoleARN = os.Getenv("KUBECTL_ROLE_ARN")
	if kubectlRoleARN == "" {
		log.Fatal("missing env variable KUBECTL_ROLE_ARN")
	}
}

func main() {
	app := awscdk.NewApp(nil)

	NewDynamoDBAppStack(app, "DynamoDBEKSAppStack", &CdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

var eksCluster awseks.Cluster

const dynamoDBPartitionKey = "shortcode"
const appDirectory = "../../part5-dynamodb-eks-ack-cdk8s/dynamodb-app"
const tableName = "urls"

func NewDynamoDBAppStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {

	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	table := awsdynamodb.NewTable(stack, jsii.String("dynamodb-table"),
		&awsdynamodb.TableProps{
			TableName: jsii.String(tableName),
			PartitionKey: &awsdynamodb.Attribute{
				Name: jsii.String(dynamoDBPartitionKey),
				Type: awsdynamodb.AttributeType_STRING,
			},
			BillingMode:   awsdynamodb.BillingMode_PAY_PER_REQUEST,
			RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		})

	eksCluster := awseks.Cluster_FromClusterAttributes(stack, jsii.String("existing cluster"),
		&awseks.ClusterAttributes{
			ClusterName:    jsii.String(eksClusterName),
			KubectlRoleArn: jsii.String(kubectlRoleARN)})

	appDockerImage := awsecrassets.NewDockerImageAsset(stack, jsii.String("app-image"),
		&awsecrassets.DockerImageAssetProps{
			Directory: jsii.String(appDirectory)})

	app := cdk8s.NewApp(nil)
	appProps := NewAppChartProps(appDockerImage.ImageUri(), table.TableName())

	eksCluster.AddCdk8sChart(jsii.String("dynamodbapp-chart"), NewDynamoDBAppChart(app, "dynamodb-cdk8s", &appProps), nil)

	return stack
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	return nil
}
