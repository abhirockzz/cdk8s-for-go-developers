package servicesk8saws

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"servicesk8saws.FieldExport",
		reflect.TypeOf((*FieldExport)(nil)).Elem(),
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
			j := jsiiProxy_FieldExport{}
			_jsii_.InitJsiiProxy(&j.Type__cdk8sApiObject)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"servicesk8saws.FieldExportProps",
		reflect.TypeOf((*FieldExportProps)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"servicesk8saws.FieldExportSpec",
		reflect.TypeOf((*FieldExportSpec)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"servicesk8saws.FieldExportSpecFrom",
		reflect.TypeOf((*FieldExportSpecFrom)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"servicesk8saws.FieldExportSpecFromResource",
		reflect.TypeOf((*FieldExportSpecFromResource)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"servicesk8saws.FieldExportSpecTo",
		reflect.TypeOf((*FieldExportSpecTo)(nil)).Elem(),
	)
	_jsii_.RegisterEnum(
		"servicesk8saws.FieldExportSpecToKind",
		reflect.TypeOf((*FieldExportSpecToKind)(nil)).Elem(),
		map[string]interface{}{
			"CONFIGMAP": FieldExportSpecToKind_CONFIGMAP,
			"SECRET": FieldExportSpecToKind_SECRET,
		},
	)
}
