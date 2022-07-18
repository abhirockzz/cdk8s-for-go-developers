package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus22/v2"
)

type WordpressProps struct {
	MySQLImage    *string
	MySQLPassword *string
	MySQLStorage  *float64

	WordpressImage   *string
	WordpressStorage *float64
}

func NewWordpressStack(scope constructs.Construct, id *string, props *WordpressProps) constructs.Construct {

	wordpressConstruct := constructs.NewConstruct(scope, id)

	//-------- Secret --------
	secretName := "mysql-pass"
	password := props.MySQLPassword

	mysqlSecret := cdk8splus22.NewSecret(wordpressConstruct, jsii.String("mysql-secret"),
		&cdk8splus22.SecretProps{
			Metadata: &cdk8s.ApiObjectMetadata{Name: jsii.String(secretName)}})

	secretKey := "password"
	mysqlSecret.AddStringData(jsii.String(secretKey), password)

	//-------- Deployment --------

	dep := cdk8splus22.NewDeployment(wordpressConstruct, jsii.String("mysql-deployment-cdk8splus"), &cdk8splus22.DeploymentProps{})

	containerImage := props.MySQLImage

	mysqlContainer := dep.AddContainer(&cdk8splus22.ContainerProps{
		Name:  jsii.String("mysql-container"),
		Image: containerImage,
		Port:  jsii.Number(3306),
	})

	envValFromSecret := cdk8splus22.EnvValue_FromSecretValue(&cdk8splus22.SecretValue{Key: jsii.String(secretKey), Secret: mysqlSecret}, &cdk8splus22.EnvValueFromSecretOptions{Optional: jsii.Bool(false)})

	mySQLPasswordEnvName := "MYSQL_ROOT_PASSWORD"

	mysqlContainer.Env().AddVariable(jsii.String(mySQLPasswordEnvName), envValFromSecret)

	//-------- Storage - PVC, Volume mount --------

	mysqlPVC := cdk8splus22.NewPersistentVolumeClaim(wordpressConstruct, jsii.String("mysql-pvc"), &cdk8splus22.PersistentVolumeClaimProps{
		AccessModes: &[]cdk8splus22.PersistentVolumeAccessMode{cdk8splus22.PersistentVolumeAccessMode_READ_WRITE_ONCE},
		Storage:     cdk8s.Size_Gibibytes(props.MySQLStorage)})

	mysqlVolumeName := "mysql-persistent-storage"
	mysqlVolume := cdk8splus22.Volume_FromPersistentVolumeClaim(wordpressConstruct, jsii.String("mysql-vol-pvc"), mysqlPVC, &cdk8splus22.PersistentVolumeClaimVolumeOptions{Name: jsii.String(mysqlVolumeName)})

	mySQLVolumeMountPath := "/var/lib/mysql"
	mysqlContainer.Mount(jsii.String(mySQLVolumeMountPath), mysqlVolume, &cdk8splus22.MountOptions{})

	//-------- Service --------

	mySQLServiceName := "mysql-service"
	clusterIPNone := "None"

	cdk8splus22.NewService(wordpressConstruct, jsii.String("mysql-service"), &cdk8splus22.ServiceProps{
		Metadata:  &cdk8s.ApiObjectMetadata{Name: jsii.String(mySQLServiceName)},
		Selector:  dep,
		ClusterIP: jsii.String(clusterIPNone),
		Ports:     &[]*cdk8splus22.ServicePort{{Port: jsii.Number(3306)}},
	})

	//-------- Wordpress --------

	wordPressDeployment := cdk8splus22.NewDeployment(wordpressConstruct, jsii.String("wordpress-deployment-cdk8splus"), &cdk8splus22.DeploymentProps{})

	wordpressContainer := wordPressDeployment.AddContainer(&cdk8splus22.ContainerProps{
		Name:  jsii.String("wordpress-container"),
		Image: props.WordpressImage,
		Port:  jsii.Number(80),
	})

	//envValFromSecret := cdk8splus22.EnvValue_FromSecretValue(&cdk8splus22.SecretValue{Key: jsii.String(secretKey), Secret: mysqlSecret}, &cdk8splus22.EnvValueFromSecretOptions{Optional: jsii.Bool(false)})

	wordpressMySQLPasswordEnvName := "WORDPRESS_DB_PASSWORD"
	wordpressMySQLDBHostEnvName := "WORDPRESS_DB_HOST"
	wordpressMySQLDBHostEnvValue := "mysql-service"

	wordpressContainer.Env().AddVariable(jsii.String(wordpressMySQLPasswordEnvName), envValFromSecret)
	wordpressContainer.Env().AddVariable(jsii.String(wordpressMySQLDBHostEnvName), cdk8splus22.EnvValue_FromValue(jsii.String(wordpressMySQLDBHostEnvValue)))

	wordpressPVC := cdk8splus22.NewPersistentVolumeClaim(wordpressConstruct, jsii.String("wordpress-pvc"), &cdk8splus22.PersistentVolumeClaimProps{
		AccessModes: &[]cdk8splus22.PersistentVolumeAccessMode{cdk8splus22.PersistentVolumeAccessMode_READ_WRITE_ONCE},
		Storage:     cdk8s.Size_Gibibytes(props.WordpressStorage)})

	wordpressVolumeName := "wordpress-persistent-storage"
	wordpressVolume := cdk8splus22.Volume_FromPersistentVolumeClaim(wordpressConstruct, jsii.String("wordpress-vol-pvc"), wordpressPVC, &cdk8splus22.PersistentVolumeClaimVolumeOptions{Name: jsii.String(wordpressVolumeName)})

	wordpressVolumeMountPath := "/var/www/html"
	wordpressContainer.Mount(jsii.String(wordpressVolumeMountPath), wordpressVolume, &cdk8splus22.MountOptions{})

	wordpressServiceName := "wordpress-service"

	wordPressDeployment.ExposeViaService(&cdk8splus22.DeploymentExposeViaServiceOptions{Name: jsii.String(wordpressServiceName), ServiceType: cdk8splus22.ServiceType_LOAD_BALANCER, Ports: &[]*cdk8splus22.ServicePort{{Port: jsii.Number(80)}}})

	return wordpressConstruct
}
