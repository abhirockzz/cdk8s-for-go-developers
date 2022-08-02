package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"

	"github.com/aws/constructs-go/constructs/v10"
)

type CdkStackProps struct {
	awscdk.StackProps
}

func main() {
	app := awscdk.NewApp(nil)

	NewNginxOnEKSStack(app, "NginxOnEKSStack", &CdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func NewNginxOnEKSStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {

	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	//vpc
	vpc := awsec2.NewVpc(stack, jsii.String("demo-vpc"), nil)

	//--------EKS---------

	eksSecurityGroup := awsec2.NewSecurityGroup(stack, jsii.String("eks-demo-sg"),
		&awsec2.SecurityGroupProps{
			Vpc:               vpc,
			SecurityGroupName: jsii.String("eks-demo-sg"),
			AllowAllOutbound:  jsii.Bool(true)})

	eksCluster := awseks.NewCluster(stack, jsii.String("demo-eks"),
		&awseks.ClusterProps{
			ClusterName:   jsii.String("demo-eks-cluster"),
			Version:       awseks.KubernetesVersion_V1_21(),
			Vpc:           vpc,
			SecurityGroup: eksSecurityGroup,
			VpcSubnets: &[]*awsec2.SubnetSelection{
				{Subnets: vpc.PrivateSubnets()}},
			DefaultCapacity:         jsii.Number(2),
			DefaultCapacityInstance: awsec2.InstanceType_Of(awsec2.InstanceClass_BURSTABLE3, awsec2.InstanceSize_SMALL), DefaultCapacityType: awseks.DefaultCapacityType_NODEGROUP,
			OutputConfigCommand: jsii.Bool(true),
			EndpointAccess:      awseks.EndpointAccess_PUBLIC()})

	deployNginxUsingCDK(eksCluster)
	deployNginxUsingCDK8s(eksCluster)

	return stack
}

func deployNginxUsingCDK(eksCluster awseks.Cluster) {

	appLabel := map[string]*string{
		"app": jsii.String("nginx-eks-cdk"),
	}

	deployment := map[string]interface{}{
		"apiVersion": jsii.String("apps/v1"),
		"kind":       jsii.String("Deployment"),
		"metadata": map[string]*string{
			"name": jsii.String("nginx-deployment-cdk"),
		},
		"spec": map[string]interface{}{
			"replicas": jsii.Number(1),
			"selector": map[string]map[string]*string{
				"matchLabels": appLabel,
			},
			"template": map[string]interface{}{
				"metadata": map[string]map[string]*string{
					"labels": appLabel,
				},
				"spec": map[string][]map[string]interface{}{
					"containers": {
						{
							"name":  jsii.String("nginx"),
							"image": jsii.String("nginx"),
							"ports": []map[string]*float64{
								{
									"containerPort": jsii.Number(80),
								},
							},
						},
					},
				},
			},
		},
	}

	service := map[string]interface{}{
		"apiVersion": jsii.String("v1"),
		"kind":       jsii.String("Service"),
		"metadata": map[string]*string{
			"name": jsii.String("nginx-service-cdk"),
		},
		"spec": map[string]interface{}{
			"type": jsii.String("LoadBalancer"),
			"ports": []map[string]*float64{
				{
					"port":       jsii.Number(9090),
					"targetPort": jsii.Number(80),
				},
			},
			"selector": appLabel,
		},
	}

	eksCluster.AddManifest(jsii.String("app-deployment"), &service, &deployment)

}

func deployNginxUsingCDK8s(eksCluster awseks.Cluster) {

	app := cdk8s.NewApp(nil)
	eksCluster.AddCdk8sChart(jsii.String("nginx-eks-chart"), NewNginxChart(app, "nginx-cdk8s", nil), nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	return nil

}
