package stripe

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func String(d *schema.ResourceData, key string) string {
	return ToString(d.Get(key))
}

func ToString(value interface{}) string {
	return value.(string)
}

func Int64(d *schema.ResourceData, key string) int64 {
	return ToInt64(d.Get(key))
}

func ToInt64(value interface{}) int64 {
	return int64(value.(int))
}

func Float64(d *schema.ResourceData, key string) float64 {
	return ToFloat64(d.Get(key))
}

func ToFloat64(value interface{}) float64 {
	return value.(float64)
}

func StringSlice(d *schema.ResourceData, key string) []string {
	return ToStringSlice(d.Get(key))
}

func ToStringSlice(value interface{}) []string {
	slice := value.([]interface{})
	stringSlice := make([]string, len(slice), len(slice))
	for i := range slice {
		stringSlice[i] = ToString(slice[i])
	}
	return stringSlice
}

func Bool(d *schema.ResourceData, key string) bool {
	return ToBool(d.Get(key))
}

func ToBool(value interface{}) bool {
	return value.(bool)
}

func Map(d *schema.ResourceData, key string) map[string]interface{} {
	return ToMap(d.Get(key))
}

func ToMap(value interface{}) map[string]interface{} {
	return value.(map[string]interface{})
}

func ToMapSlice(value interface{}) []map[string]interface{} {
	return value.([]map[string]interface{})
}

func CallSet(err ...error) (d diag.Diagnostics) {
	for _, e := range err {
		if e != nil {
			d = append(d, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  e.Error(),
			})
		}
	}
	return d
}
