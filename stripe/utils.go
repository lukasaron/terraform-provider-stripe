package stripe

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func String(d *schema.ResourceData, key string) string {
	return d.Get(key).(string)
}

//func Int(d *schema.ResourceData, key string) int {
//	return d.Get(key).(int)
//}
//
//func Int64(d *schema.ResourceData, key string) int64 {
//	return d.Get(key).(int64)
//}

func StringSlice(d *schema.ResourceData, key string) []string {
	slice := d.Get(key).([]interface{})
	stringSlice := make([]string, len(slice), len(slice))
	for i := range slice {
		stringSlice[i] = slice[i].(string)
	}
	return stringSlice
}

func Bool(d *schema.ResourceData, key string) bool {
	return d.Get(key).(bool)
}

func Map(d *schema.ResourceData, key string) map[string]interface{} {
	return d.Get(key).(map[string]interface{})
}

func Slice(d *schema.ResourceData, key string) []interface{} {
	return d.Get(key).([]interface{})
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
