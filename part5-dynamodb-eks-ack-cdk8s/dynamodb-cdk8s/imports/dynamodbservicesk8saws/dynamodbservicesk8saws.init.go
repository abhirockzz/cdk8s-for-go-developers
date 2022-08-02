package dynamodbservicesk8saws

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"dynamodbservicesk8saws.Table",
		reflect.TypeOf((*Table)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addDependency", GoMethod: "AddDependency"},
			_jsii_.MemberMethod{JsiiMethod: "addJsonPatch", GoMethod: "AddJsonPatch"},
			_jsii_.MemberProperty{JsiiProperty: "apiGroup", GoGetter: "ApiGroup"},
			_jsii_.MemberProperty{JsiiProperty: "apiVersion", GoGetter: "ApiVersion"},
			_jsii_.MemberProperty{JsiiProperty: "chart", GoGetter: "Chart"},
			_jsii_.MemberProperty{JsiiProperty: "kind", GoGetter: "Kind"},
			_jsii_.MemberProperty{JsiiProperty: "metadata", GoGetter: "Metadata"},
			_jsii_.MemberProperty{JsiiProperty: "name", GoGetter: "Name"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "toJson", GoMethod: "ToJson"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
		},
		func() interface{} {
			j := jsiiProxy_Table{}
			_jsii_.InitJsiiProxy(&j.Type__cdk8sApiObject)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"dynamodbservicesk8saws.TableProps",
		reflect.TypeOf((*TableProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"dynamodbservicesk8saws.TableSpec",
		reflect.TypeOf((*TableSpec)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"dynamodbservicesk8saws.TableSpecAttributeDefinitions",
		reflect.TypeOf((*TableSpecAttributeDefinitions)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"dynamodbservicesk8saws.TableSpecGlobalSecondaryIndexes",
		reflect.TypeOf((*TableSpecGlobalSecondaryIndexes)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"dynamodbservicesk8saws.TableSpecGlobalSecondaryIndexesKeySchema",
		reflect.TypeOf((*TableSpecGlobalSecondaryIndexesKeySchema)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"dynamodbservicesk8saws.TableSpecGlobalSecondaryIndexesProjection",
		reflect.TypeOf((*TableSpecGlobalSecondaryIndexesProjection)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"dynamodbservicesk8saws.TableSpecGlobalSecondaryIndexesProvisionedThroughput",
		reflect.TypeOf((*TableSpecGlobalSecondaryIndexesProvisionedThroughput)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"dynamodbservicesk8saws.TableSpecKeySchema",
		reflect.TypeOf((*TableSpecKeySchema)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"dynamodbservicesk8saws.TableSpecLocalSecondaryIndexes",
		reflect.TypeOf((*TableSpecLocalSecondaryIndexes)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"dynamodbservicesk8saws.TableSpecLocalSecondaryIndexesKeySchema",
		reflect.TypeOf((*TableSpecLocalSecondaryIndexesKeySchema)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"dynamodbservicesk8saws.TableSpecLocalSecondaryIndexesProjection",
		reflect.TypeOf((*TableSpecLocalSecondaryIndexesProjection)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"dynamodbservicesk8saws.TableSpecProvisionedThroughput",
		reflect.TypeOf((*TableSpecProvisionedThroughput)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"dynamodbservicesk8saws.TableSpecSseSpecification",
		reflect.TypeOf((*TableSpecSseSpecification)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"dynamodbservicesk8saws.TableSpecStreamSpecification",
		reflect.TypeOf((*TableSpecStreamSpecification)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"dynamodbservicesk8saws.TableSpecTags",
		reflect.TypeOf((*TableSpecTags)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"dynamodbservicesk8saws.TableSpecTimeToLive",
		reflect.TypeOf((*TableSpecTimeToLive)(nil)).Elem(),
	)
}
