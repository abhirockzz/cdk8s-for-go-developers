package main

import (
	"cdk8s-crd/imports/kafkastrimziio"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

type KafkaChartProps struct {
	cdk8s.ChartProps
}

func NewKafkaChart(scope constructs.Construct, id string, props *KafkaChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	kafkastrimziio.NewKafka(chart, jsii.String("kafka"),
		&kafkastrimziio.KafkaProps{
			Spec: &kafkastrimziio.KafkaSpec{
				Kafka: &kafkastrimziio.KafkaSpecKafka{

					Version:  jsii.String("3.2.0"),
					Replicas: jsii.Number(1),
					Listeners: &[]*kafkastrimziio.KafkaSpecKafkaListeners{
						{
							Name: jsii.String("plain"),
							Port: jsii.Number(9092),
							Type: kafkastrimziio.KafkaSpecKafkaListenersType_INTERNAL,
							Tls:  jsii.Bool(false),
						},
					},
					Config: map[string]interface{}{
						"offsets.topic.replication.factor":         1,
						"transaction.state.log.replication.factor": 1,
						"transaction.state.log.min.isr":            1,
						"default.replication.factor":               1,
						"min.insync.replicas":                      1,
						"inter.broker.protocol.version":            "3.2",
					},
					Storage: &kafkastrimziio.KafkaSpecKafkaStorage{
						Type: kafkastrimziio.KafkaSpecKafkaStorageType_EPHEMERAL,
					},
				},
				Zookeeper: &kafkastrimziio.KafkaSpecZookeeper{
					Replicas: jsii.Number(1),
					Storage: &kafkastrimziio.KafkaSpecZookeeperStorage{
						Type: kafkastrimziio.KafkaSpecZookeeperStorageType_EPHEMERAL,
					},
				},
				EntityOperator: &kafkastrimziio.KafkaSpecEntityOperator{
					TopicOperator: &kafkastrimziio.KafkaSpecEntityOperatorTopicOperator{},
				}}})
	return chart
}
