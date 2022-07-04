package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus22/v2"
)

type MyChartProps struct {
	cdk8s.ChartProps
}

func NewMySQLChart(scope constructs.Construct, id string, props *MyChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	secretName := "mysql-pass"
	password := "Password123"

	mysqlSecret := cdk8splus22.NewSecret(chart, jsii.String("mysql-secret"),
		&cdk8splus22.SecretProps{
			Metadata: &cdk8s.ApiObjectMetadata{Name: jsii.String(secretName)}})

	secretKey := "password"
	mysqlSecret.AddStringData(jsii.String(secretKey), jsii.String(password))

	dep := cdk8splus22.NewDeployment(chart, jsii.String("mysql-deployment-cdk8splus"), &cdk8splus22.DeploymentProps{})

	containerImage := "mysql"

	mysqlContainer := dep.AddContainer(&cdk8splus22.ContainerProps{
		Name:  jsii.String("mysql-container"),
		Image: jsii.String(containerImage),
		Port:  jsii.Number(3306),
	})

	envValFromSecret := cdk8splus22.EnvValue_FromSecretValue(&cdk8splus22.SecretValue{Key: jsii.String(secretKey), Secret: mysqlSecret}, &cdk8splus22.EnvValueFromSecretOptions{Optional: jsii.Bool(false)})

	mySQLPasswordEnvName := "MYSQL_ROOT_PASSWORD"

	mysqlContainer.Env().AddVariable(jsii.String(mySQLPasswordEnvName), envValFromSecret)

	mysqlPVC := cdk8splus22.NewPersistentVolumeClaim(chart, jsii.String("mysql-pvc"), &cdk8splus22.PersistentVolumeClaimProps{
		AccessModes: &[]cdk8splus22.PersistentVolumeAccessMode{cdk8splus22.PersistentVolumeAccessMode_READ_WRITE_ONCE},
		Storage:     cdk8s.Size_Gibibytes(jsii.Number(2))})

	mysqlVolumeName := "mysql-persistent-storage"
	mysqlVolume := cdk8splus22.Volume_FromPersistentVolumeClaim(chart, jsii.String("mysql-vol-pvc"), mysqlPVC, &cdk8splus22.PersistentVolumeClaimVolumeOptions{Name: jsii.String(mysqlVolumeName)})

	mySQLVolumeMountPath := "/var/lib/mysql"
	mysqlContainer.Mount(jsii.String(mySQLVolumeMountPath), mysqlVolume, &cdk8splus22.MountOptions{})

	mySQLServiceName := "mysql-service"
	clusterIPNone := "None"

	cdk8splus22.NewService(chart, jsii.String("mysql-service"), &cdk8splus22.ServiceProps{
		Metadata:  &cdk8s.ApiObjectMetadata{Name: jsii.String(mySQLServiceName)},
		Selector:  dep,
		ClusterIP: jsii.String(clusterIPNone),
		Ports:     &[]*cdk8splus22.ServicePort{{Port: jsii.Number(3306)}},
	})

	return chart
}

func NewWordpressChart(scope constructs.Construct, id string, props *MyChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	secretName := "mysql-pass"
	mysqlSecret := cdk8splus22.Secret_FromSecretName(chart, jsii.String("existing-secret"), jsii.String(secretName))

	secretKey := "password"

	dep := cdk8splus22.NewDeployment(chart, jsii.String("wordpress-deployment-cdk8splus"), &cdk8splus22.DeploymentProps{})

	containerImage := "wordpress:4.8-apache"

	wordpressContainer := dep.AddContainer(&cdk8splus22.ContainerProps{
		Name:  jsii.String("wordpress-container"),
		Image: jsii.String(containerImage),
		Port:  jsii.Number(80),
	})

	envValFromSecret := cdk8splus22.EnvValue_FromSecretValue(&cdk8splus22.SecretValue{Key: jsii.String(secretKey), Secret: mysqlSecret}, &cdk8splus22.EnvValueFromSecretOptions{Optional: jsii.Bool(false)})

	wordpressMySQLPasswordEnvName := "WORDPRESS_DB_PASSWORD"
	wordpressMySQLDBHostEnvName := "WORDPRESS_DB_HOST"
	wordpressMySQLDBHostEnvValue := "mysql-service"

	wordpressContainer.Env().AddVariable(jsii.String(wordpressMySQLPasswordEnvName), envValFromSecret)
	wordpressContainer.Env().AddVariable(jsii.String(wordpressMySQLDBHostEnvName), cdk8splus22.EnvValue_FromValue(jsii.String(wordpressMySQLDBHostEnvValue)))

	wordpressPVC := cdk8splus22.NewPersistentVolumeClaim(chart, jsii.String("wordpress-pvc"), &cdk8splus22.PersistentVolumeClaimProps{
		AccessModes: &[]cdk8splus22.PersistentVolumeAccessMode{cdk8splus22.PersistentVolumeAccessMode_READ_WRITE_ONCE},
		Storage:     cdk8s.Size_Gibibytes(jsii.Number(2))})

	wordpressVolumeName := "wordpress-persistent-storage"
	wordpressVolume := cdk8splus22.Volume_FromPersistentVolumeClaim(chart, jsii.String("wordpress-vol-pvc"), wordpressPVC, &cdk8splus22.PersistentVolumeClaimVolumeOptions{Name: jsii.String(wordpressVolumeName)})

	wordpressVolumeMountPath := "/var/www/html"
	wordpressContainer.Mount(jsii.String(wordpressVolumeMountPath), wordpressVolume, &cdk8splus22.MountOptions{})

	wordpressServiceName := "wordpress-service"

	dep.ExposeViaService(&cdk8splus22.DeploymentExposeViaServiceOptions{Name: jsii.String(wordpressServiceName), ServiceType: cdk8splus22.ServiceType_LOAD_BALANCER, Ports: &[]*cdk8splus22.ServicePort{{Port: jsii.Number(80)}}})

	return chart
}

func main() {
	app := cdk8s.NewApp(nil)

	mySQLChart := NewMySQLChart(app, "mysql", nil)
	wordpressChart := NewWordpressChart(app, "wordpress", nil)

	wordpressChart.AddDependency(mySQLChart)

	app.Synth()
}
