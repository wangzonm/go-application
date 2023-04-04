package json_test

import (
	"encoding/json"
	"fmt"
	"testing"
)

// 1.16.15/src/encoding/json/encode.go
// func (me mapEncoder) encode(e *encodeState, v reflect.Value, opts encOpts)
// ...
//  Extract and sort the keys.
//	keys := v.MapKeys()
//	sv := make([]reflectWithString, len(keys))
//	for i, v := range keys {
//		sv[i].v = v
//		if err := sv[i].resolve(); err != nil {
//			e.error(fmt.Errorf("json: encoding error for type %q: %q", v.Type().String(), err.Error()))
//		}
//	}
//	sort.Slice(sv, func(i, j int) bool { return sv[i].s < sv[j].s })
// ...
func TestMarshalKeySort(t *testing.T) {
	m := map[string]interface{}{
		"brand": "company",
		"title": "crud engineer",
		"name":  "alex",
		"city":  "bj",
		"age":   20,
	}
	b, err := json.Marshal(&m)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))
}
