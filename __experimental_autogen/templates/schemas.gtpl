package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"qdrant-terraform-automation/models"
	"strconv"
)

// {{ pascalize .Name }}Schema mapping representing the {{ .Name }} resource defined in the Terraform configuration
func {{ pascalize .Name }}Schema() map[string]*schema.Schema {
return map[string]*schema.Schema{
{{- range .Properties }}
    "{{ .Name | snakize }}": {
    {{- if eq .Name "id" }}
        Type:     schema.TypeString,
        Computed: true,
    {{- else if eq .GoType "string" }}
        Type: schema.TypeString,
        {{- if .Default }}
            Default: "{{ .Default }}",
        {{- end }}
    {{- else if eq .GoType "boolean" }}
        Type: schema.TypeBool,
        {{- if .Default }}
            Default: {{ .Default }},
        {{- end }}
    {{- else if eq .GoType "int" }}
        Type: schema.TypeInt,
    {{- else if eq .GoType "map[string]interface{}" }}
        Type: schema.TypeMap,
        Elem: &schema.Schema{Type: schema.TypeString},
    {{- else if eq .GoType "[]interface{}" }}
        Type: schema.TypeList,
        Elem: &schema.Resource{
        {{- if .Items }}
            Schema: {{ pascalize .Items.GoType }}Schema(),
        {{- else }}
            Schema: &schema.Schema{Type: schema.TypeString},
        {{- end }}
        },
        ConfigMode: schema.SchemaConfigModeAttr,
    {{- end }}
    {{- if .Required }}
        Required: true,
    {{- else if .ReadOnly }}
        Computed: true,
    {{- else }}
        Optional: true,
    {{- end }}
    },
{{- end }}
}
}

// DataSource{{ pascalize .Name }}Schema mapping representing the resource's respective datasource object defined in Terraform configuration
func DataSource{{ pascalize .Name }}Schema() map[string]*schema.Schema {
return map[string]*schema.Schema{
{{- range .Properties }}
    "{{ .Name | snakize }}": {
    {{- if eq .Name "id" }}
        Type:     schema.TypeString,
        Computed: true,
        Optional: true,
    {{- else if eq .GoType "string" }}
        Type: schema.TypeString,
        {{- if .Default }}
            Default: "{{ .Default }}",
        {{- end }}
    {{- else if eq .GoType "boolean" }}
        Type: schema.TypeBool,
        {{- if .Default }}
            Default: {{ .Default }},
        {{- end }}
    {{- else if eq .GoType "int" }}
        Type: schema.TypeInt,
    {{- else if eq .GoType "map[string]interface{}" }}
        Type: schema.TypeMap,
        Elem: &schema.Schema{Type: schema.TypeString},
    {{- else if eq .GoType "[]interface{}" }}
        Type: schema.TypeList,
        Elem: &schema.Resource{
        {{- if .Items }}
            Schema: {{ pascalize .Items.GoType }}Schema(),
        {{- else }}
            Schema: &schema.Schema{Type: schema.TypeString},
        {{- end }}
        },
        ConfigMode: schema.SchemaConfigModeAttr,
    {{- end }}
    Optional: true,
    },
{{- end }}
"filter": {
Type:     schema.TypeString,
Optional: true,
},
}
}
