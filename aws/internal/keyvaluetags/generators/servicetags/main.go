// +build ignore

package main

import (
	"bytes"
	"go/format"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

const filename = `service_tags_gen.go`

// Representing types such as []*athena.Tag, []*ec2.Tag, ...
var sliceServiceNames = []string{
	"acm",
	"acmpca",
	"appmesh",
	"athena",
	"autoscaling",
	"cloud9",
	"cloudformation",
	"cloudfront",
	"cloudhsmv2",
	"cloudtrail",
	"cloudwatch",
	"cloudwatchevents",
	"codeartifact",
	"codebuild",
	"codedeploy",
	"codepipeline",
	"codestarconnections",
	"configservice",
	"databasemigrationservice",
	"datapipeline",
	"datasync",
	"dax",
	"devicefarm",
	"directconnect",
	"directoryservice",
	"docdb",
	"dynamodb",
	"ec2",
	"ecr",
	"ecs",
	"efs",
	"elasticache",
	"elasticbeanstalk",
	"elasticsearchservice",
	"elb",
	"elbv2",
	"emr",
	"firehose",
	"fms",
	"fsx",
	"gamelift",
	"globalaccelerator",
	"iam",
	"inspector",
	"iot",
	"iotanalytics",
	"iotevents",
	"kinesis",
	"kinesisanalytics",
	"kinesisanalyticsv2",
	"kms",
	"licensemanager",
	"lightsail",
	"mediastore",
	"neptune",
	"networkfirewall",
	"networkmanager",
	"organizations",
	"quicksight",
	"ram",
	"rds",
	"redshift",
	"resourcegroupstaggingapi",
	"route53",
	"route53resolver",
	"s3",
	"s3control",
	"sagemaker",
	"secretsmanager",
	"serverlessapplicationrepository",
	"servicecatalog",
	"servicediscovery",
	"sfn",
	"sns",
	"ssm",
	"ssoadmin",
	"storagegateway",
	"swf",
	"transfer",
	"waf",
	"wafregional",
	"wafv2",
	"workspaces",
	"xray",
}

var mapServiceNames = []string{
	"accessanalyzer",
	"amplify",
	"apigateway",
	"apigatewayv2",
	"appstream",
	"appsync",
	"backup",
	"batch",
	"cloudwatchlogs",
	"codecommit",
	"codestarnotifications",
	"cognitoidentity",
	"cognitoidentityprovider",
	"dataexchange",
	"dlm",
	"eks",
	"glacier",
	"glue",
	"guardduty",
	"greengrass",
	"kafka",
	"kinesisvideo",
	"imagebuilder",
	"lambda",
	"mediaconnect",
	"mediaconvert",
	"medialive",
	"mediapackage",
	"mq",
	"opsworks",
	"qldb",
	"pinpoint",
	"resourcegroups",
	"securityhub",
	"signer",
	"sqs",
	"synthetics",
	"worklink",
}

type TemplateData struct {
	MapServiceNames   []string
	SliceServiceNames []string
}

func main() {
	// Always sort to reduce any potential generation churn
	sort.Strings(mapServiceNames)
	sort.Strings(sliceServiceNames)

	templateData := TemplateData{
		MapServiceNames:   mapServiceNames,
		SliceServiceNames: sliceServiceNames,
	}
	templateFuncMap := template.FuncMap{
		"TagKeyType":                  keyvaluetags.ServiceTagKeyType,
		"TagPackage":                  keyvaluetags.ServiceTagPackage,
		"TagResourceTypeField":        keyvaluetags.ServiceTagResourceTypeField,
		"TagType":                     keyvaluetags.ServiceTagType,
		"TagType2":                    keyvaluetags.ServiceTagType2,
		"TagTypeAdditionalBoolFields": keyvaluetags.ServiceTagTypeAdditionalBoolFields,
		"TagTypeIdentifierField":      keyvaluetags.ServiceTagTypeIdentifierField,
		"TagTypeKeyField":             keyvaluetags.ServiceTagTypeKeyField,
		"TagTypeValueField":           keyvaluetags.ServiceTagTypeValueField,
		"Title":                       strings.Title,
		"ToSnakeCase":                 keyvaluetags.ToSnakeCase,
	}

	tmpl, err := template.New("servicetags").Funcs(templateFuncMap).Parse(templateBody)

	if err != nil {
		log.Fatalf("error parsing template: %s", err)
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, templateData)

	if err != nil {
		log.Fatalf("error executing template: %s", err)
	}

	generatedFileContents, err := format.Source(buffer.Bytes())

	if err != nil {
		log.Fatalf("error formatting generated file: %s", err)
	}

	f, err := os.Create(filename)

	if err != nil {
		log.Fatalf("error creating file (%s): %s", filename, err)
	}

	defer f.Close()

	_, err = f.Write(generatedFileContents)

	if err != nil {
		log.Fatalf("error writing to file (%s): %s", filename, err)
	}
}

var templateBody = `
// Code generated by generators/servicetags/main.go; DO NOT EDIT.

package keyvaluetags

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
{{- range .SliceServiceNames }}
{{- if eq . (. | TagPackage) }}
	"github.com/aws/aws-sdk-go/service/{{ . }}"
{{- end }}
{{- end }}
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// map[string]*string handling
{{- range .MapServiceNames }}

// {{ . | Title }}Tags returns {{ . }} service tags.
func (tags KeyValueTags) {{ . | Title }}Tags() map[string]*string {
	return aws.StringMap(tags.Map())
}

// {{ . | Title }}KeyValueTags creates KeyValueTags from {{ . }} service tags.
func {{ . | Title }}KeyValueTags(tags map[string]*string) KeyValueTags {
	return New(tags)
}
{{- end }}

// []*SERVICE.Tag handling
{{- range .SliceServiceNames }}

{{ if and ( . | TagTypeIdentifierField ) ( . | TagTypeAdditionalBoolFields ) }}
// {{ . | Title }}ListOfMap returns a list of {{ . }} in flattened map.
//
// Compatible with setting Terraform state for strongly typed configuration blocks.
//
// This function strips tag resource identifier and type. Generally, this is
// the desired behavior so the tag schema does not require those attributes.
// Use (keyvaluetags.KeyValueTags).ListOfMap() for full tag information.
func (tags KeyValueTags) {{ . | Title }}ListOfMap() []interface{} {
	var result []interface{}

	for _, key := range tags.Keys() {
		m := map[string]interface{}{
			"key":                   key,
			"value":                 aws.StringValue(tags.KeyValue(key)),
			{{- range . | TagTypeAdditionalBoolFields }}
			"{{ . | ToSnakeCase }}": aws.BoolValue(tags.KeyAdditionalBoolValue(key, "{{ . }}")),
			{{- end }}
		}

		result = append(result, m)
	}

	return result
}
{{- end }}

{{ if eq . "autoscaling" }}
// {{ . | Title }}ListOfStringMap returns a list of {{ . }} tags in flattened map of only string values.
//
// Compatible with setting Terraform state for legacy []map[string]string schema.
// Deprecated: Will be removed in a future major version without replacement.
func (tags KeyValueTags) {{ . | Title }}ListOfStringMap() []interface{} {
	var result []interface{}

	for _, key := range tags.Keys() {
		m := map[string]string{
			"key":                   key,
			"value":                 aws.StringValue(tags.KeyValue(key)),
			{{- range . | TagTypeAdditionalBoolFields }}
			"{{ . | ToSnakeCase }}": strconv.FormatBool(aws.BoolValue(tags.KeyAdditionalBoolValue(key, "{{ . }}"))),
			{{- end }}
		}

		result = append(result, m)
	}

	return result
}
{{- end }}

{{- if . | TagKeyType }}
// {{ . | Title }}TagKeys returns {{ . }} service tag keys.
func (tags KeyValueTags) {{ . | Title }}TagKeys() []*{{ . | TagPackage }}.{{ . | TagKeyType }} {
	result := make([]*{{ . | TagPackage }}.{{ . | TagKeyType }}, 0, len(tags))

	for k := range tags.Map() {
		tagKey := &{{ . | TagPackage }}.{{ . | TagKeyType }}{
			{{ . | TagTypeKeyField }}: aws.String(k),
		}

		result = append(result, tagKey)
	}

	return result
}
{{- end }}

// {{ . | Title }}Tags returns {{ . }} service tags.
func (tags KeyValueTags) {{ . | Title }}Tags() []*{{ . | TagPackage }}.{{ . | TagType }} {
	{{- if or ( . | TagTypeIdentifierField ) ( . | TagTypeAdditionalBoolFields) }}
	var result []*{{ . | TagPackage }}.{{ . | TagType }}

	for _, key := range tags.Keys() {
		tag := &{{ . | TagPackage }}.{{ . | TagType }}{
			{{ . | TagTypeKeyField }}:        aws.String(key),
			{{ . | TagTypeValueField }}:      tags.KeyValue(key),
			{{- if ( . | TagTypeIdentifierField ) }}
			{{ . | TagTypeIdentifierField }}: tags.KeyAdditionalStringValue(key, "{{ . | TagTypeIdentifierField }}"),
			{{- if ( . | TagResourceTypeField ) }}
			{{ . | TagResourceTypeField }}:   tags.KeyAdditionalStringValue(key, "{{ . | TagResourceTypeField }}"),
			{{- end }}
			{{- end }}
			{{- range . | TagTypeAdditionalBoolFields }}
			{{ . }}:                          tags.KeyAdditionalBoolValue(key, "{{ . }}"),
			{{- end }}
		}

		result = append(result, tag)
	}
	{{- else }}
	result := make([]*{{ . | TagPackage }}.{{ . | TagType }}, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &{{ . | TagPackage }}.{{ . | TagType }}{
			{{ . | TagTypeKeyField }}:   aws.String(k),
			{{ . | TagTypeValueField }}: aws.String(v),
		}

		result = append(result, tag)
	}
	{{- end }}

	return result
}

// {{ . | Title }}KeyValueTags creates KeyValueTags from {{ . }} service tags.
{{- if or ( . | TagType2 ) ( . | TagTypeAdditionalBoolFields ) }}
//
// Accepts the following types:
//   - []*{{ . | TagPackage }}.{{ . | TagType }}
{{- if . | TagType2 }}
//   - []*{{ . | TagPackage }}.{{ . | TagType2 }}
{{- end }}
{{- if . | TagTypeAdditionalBoolFields }}
//   - []interface{} (Terraform TypeList configuration block compatible)
//   - *schema.Set (Terraform TypeSet configuration block compatible)
{{- end }}
func {{ . | Title }}KeyValueTags(tags interface{}{{ if . | TagTypeIdentifierField }}, identifier string{{ if . | TagResourceTypeField }}, resourceType string{{ end }}{{ end }}) KeyValueTags {
	switch tags := tags.(type) {
	case []*{{ . | TagPackage }}.{{ . | TagType }}:
		{{- if or ( . | TagTypeIdentifierField ) ( . | TagTypeAdditionalBoolFields) }}
		m := make(map[string]*TagData, len(tags))

		for _, tag := range tags {
			tagData := &TagData{
				Value: tag.{{ . | TagTypeValueField }},
			}

			tagData.AdditionalBoolFields = make(map[string]*bool)
			{{- range . | TagTypeAdditionalBoolFields }}
			tagData.AdditionalBoolFields["{{ . }}"] = tag.{{ . }}
			{{- end }}

			{{- if . | TagTypeIdentifierField }}
			tagData.AdditionalStringFields = make(map[string]*string)
			tagData.AdditionalStringFields["{{ . | TagTypeIdentifierField }}"] = &identifier
			{{- if . | TagResourceTypeField }}
			tagData.AdditionalStringFields["{{ . | TagResourceTypeField }}"] = &resourceType
			{{- end }}
			{{- end }}

			m[aws.StringValue(tag.{{ . | TagTypeKeyField }})] = tagData
		}
		{{- else }}
		m := make(map[string]*string, len(tags))

		for _, tag := range tags {
			m[aws.StringValue(tag.{{ . | TagTypeKeyField }})] = tag.{{ . | TagTypeValueField }}
		}
		{{- end }}

		return New(m)
	case []*{{ . | TagPackage }}.{{ . | TagType2 }}:
		{{- if or ( . | TagTypeIdentifierField ) ( . | TagTypeAdditionalBoolFields) }}
		m := make(map[string]*TagData, len(tags))

		for _, tag := range tags {
			tagData := &TagData{
				Value: tag.{{ . | TagTypeValueField }},
			}

			{{ if . | TagTypeAdditionalBoolFields }}
			tagData.AdditionalBoolFields = make(map[string]*bool)
			{{- range . | TagTypeAdditionalBoolFields }}
			tagData.AdditionalBoolFields["{{ . }}"] = tag.{{ . }}
			{{- end }}
			{{- end }}

			{{- if . | TagTypeIdentifierField }}
			tagData.AdditionalStringFields = make(map[string]*string)
			tagData.AdditionalStringFields["{{ . | TagTypeIdentifierField }}"] = &identifier
			{{- if . | TagResourceTypeField }}
			tagData.AdditionalStringFields["{{ . | TagResourceTypeField }}"] = &resourceType
			{{- end }}
			{{- end }}

			m[aws.StringValue(tag.{{ . | TagTypeKeyField }})] = tagData
		}
		{{- else }}
		m := make(map[string]*string, len(tags))

		for _, tag := range tags {
			m[aws.StringValue(tag.{{ . | TagTypeKeyField }})] = tag.{{ . | TagTypeValueField }}
		}
		{{- end }}

		return New(m)
	{{- if . | TagTypeAdditionalBoolFields }}
	case *schema.Set:
		return {{ . | Title }}KeyValueTags(tags.List(){{ if . | TagTypeIdentifierField }}, identifier{{ if . | TagResourceTypeField }}, resourceType{{ end }}{{ end }})
	case []interface{}:
		result := make(map[string]*TagData)

		for _, tfMapRaw := range tags {
			tfMap, ok := tfMapRaw.(map[string]interface{})

			if !ok {
				continue
			}

			key, ok := tfMap["key"].(string)

			if !ok {
				continue
			}

			tagData := &TagData{}

			if v, ok := tfMap["value"].(string); ok {
				tagData.Value = &v
			}

			{{ if . | TagTypeAdditionalBoolFields }}
			tagData.AdditionalBoolFields = make(map[string]*bool)
			{{- range . | TagTypeAdditionalBoolFields }}
			if v, ok := tfMap["{{ . | ToSnakeCase }}"].(bool); ok {
				tagData.AdditionalBoolFields["{{ . }}"] = &v
			}
			{{- end }}
			{{ if eq . "autoscaling" }}
			// Deprecated: Legacy map handling
			{{- range . | TagTypeAdditionalBoolFields }}
			if v, ok := tfMap["{{ . | ToSnakeCase }}"].(string); ok {
				b, _ := strconv.ParseBool(v)
				tagData.AdditionalBoolFields["{{ . }}"] = &b
			}
			{{- end }}
			{{- end }}
			{{- end }}

			{{ if . | TagTypeIdentifierField }}
			tagData.AdditionalStringFields = make(map[string]*string)
			tagData.AdditionalStringFields["{{ . | TagTypeIdentifierField }}"] = &identifier
			{{- if . | TagResourceTypeField }}
			tagData.AdditionalStringFields["{{ . | TagResourceTypeField }}"] = &resourceType
			{{- end }}
			{{- end }}

			result[key] = tagData
		}

		return New(result)
	{{- end }}
	default:
		return New(nil)
	}
}
{{- else }}
func {{ . | Title }}KeyValueTags(tags []*{{ . | TagPackage }}.{{ . | TagType }}) KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.{{ . | TagTypeKeyField }})] = tag.{{ . | TagTypeValueField }}
	}

	return New(m)
}
{{- end }}
{{- end }}
`
