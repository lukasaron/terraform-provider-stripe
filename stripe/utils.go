package stripe

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stripe/stripe-go/v76"
)

func ExtractString(d *schema.ResourceData, key string) string {
	return ToString(d.Get(key))
}

func ToString(value interface{}) string {
	switch value.(type) {
	case string:
		return value.(string)
	case *string:
		return *(value.(*string))
	default:
		return ""
	}
}

func ExtractInt(d *schema.ResourceData, key string) int {
	return ToInt(d.Get(key))
}

func ExtractInt64(d *schema.ResourceData, key string) int64 {
	return ToInt64(d.Get(key))
}

func ToInt(value interface{}) int {
	switch value.(type) {
	case int:
		return value.(int)
	case *int:
		return *(value.(*int))
	case int64:
		return int(value.(int64))
	case *int64:
		return int(*(value.(*int64)))
	default:
		return 0
	}
}

func ToInt64(value interface{}) int64 {
	switch value.(type) {
	case int:
		return int64(value.(int))
	case *int:
		return int64(*(value.(*int)))
	case int64:
		return value.(int64)
	case *int64:
		return *(value.(*int64))
	default:
		return 0
	}
}

func NonZeroInt64(value interface{}) *int64 {
	valueInt64 := ToInt64(value)
	if valueInt64 == 0 {
		return nil
	}

	// amount is -1 for free price
	if valueInt64 < 0 {
		valueInt64 = 0
	}

	return &valueInt64
}

func NonZeroFloat64(value interface{}) *float64 {
	valueFloat64 := ToFloat64(value)
	if valueFloat64 == 0 {
		return nil
	}

	// amount is -1 for free price
	if valueFloat64 < 0 {
		valueFloat64 = 0
	}

	return &valueFloat64
}

func NonZeroString(value interface{}) *string {
	valueString := ToString(value)
	if len(valueString) == 0 {
		return nil
	}
	return &valueString
}

func ExtractFloat64(d *schema.ResourceData, key string) float64 {
	return ToFloat64(d.Get(key))
}

func ToFloat64(value interface{}) float64 {
	switch value.(type) {
	case float32:
		return float64(value.(float32))
	case *float32:
		return float64(*(value.(*float32)))
	case float64:
		return value.(float64)
	case *float64:
		return *(value.(*float64))
	default:
		return 0
	}
}

func ToSlice(value interface{}) []interface{} {
	switch value.(type) {
	case []interface{}:
		return value.([]interface{})
	default:
		return []interface{}{}
	}
}

func ExtractStringSlice(d *schema.ResourceData, key string) []string {
	return ToStringSlice(d.Get(key))
}

func ToStringSlice(value interface{}) []string {
	slice := ToSlice(value)
	stringSlice := make([]string, len(slice), len(slice))
	for i := range slice {
		stringSlice[i] = ToString(slice[i])
	}
	return stringSlice
}

func ExtractBool(d *schema.ResourceData, key string) bool {
	return ToBool(d.Get(key))
}

func ToBool(value interface{}) bool {
	switch value.(type) {
	case bool:
		return value.(bool)
	case *bool:
		return *(value.(*bool))
	default:
		return false
	}
}

func ExtractMap(d *schema.ResourceData, key string) map[string]interface{} {
	return ToMap(d.Get(key))
}

func ToMap(value interface{}) map[string]interface{} {
	switch value.(type) {
	case map[string]interface{}:
		return value.(map[string]interface{})
	case []interface{}:
		sl := value.([]interface{})
		if len(sl) > 0 {
			return sl[0].(map[string]interface{})
		}
	}
	return map[string]interface{}{}
}

func ExtractMapSlice(d *schema.ResourceData, key string) []map[string]interface{} {
	return ToMapSlice(d.Get(key))
}

func ToMapSlice(value interface{}) []map[string]interface{} {
	slice := ToSlice(value)
	mapSlice := make([]map[string]interface{}, len(slice), len(slice))
	for i := range slice {
		mapSlice[i] = ToMap(slice[i])
	}
	return mapSlice
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

type MetadataAdder interface {
	AddMetadata(key, value string)
}

func UpdateMetadata(d *schema.ResourceData, adder MetadataAdder) {
	oldMeta, newMeta := d.GetChange("metadata")
	oldMetaMap := ToMap(oldMeta)
	newMetaMap := ToMap(newMeta)
	for k := range newMetaMap {
		if _, set := oldMetaMap[k]; set {
			delete(oldMetaMap, k)
		}
	}

	for k, v := range newMetaMap {
		adder.AddMetadata(k, ToString(v))
	}
	for k := range oldMetaMap {
		// when meta is empty string it's going be removed
		adder.AddMetadata(k, "")
	}
}

func isRateLimitErr(e error) bool {
	err, ok := e.(*stripe.Error)
	return ok && err.HTTPStatusCode == 429
}

func retryWithBackOff(call func() error) error {
	var backOff time.Duration = 1
	for {
		err := call()
		switch {
		case err == nil:
			return nil
		case isRateLimitErr(err):
			time.Sleep(time.Second * backOff)
			backOff *= 2
		default:
			return err
		}
	}
}
