package hypermedia

import (
	"bytes"
	"encoding/json"
	"testing"
)

var HALDocumentJSONTests = []struct {
	json []byte
	doc  HALDocument
}{
	{[]byte(`{}`), HALDocument{}},
	{[]byte(`{"a":"xyz","b":true,"c":12345}`),
		HALDocument{Attributes: Attributes{
			"a": "xyz",
			"b": true,
			"c": 12345,
		}}},
	{[]byte(`{"_links":{"a":{"href":"xyz"},"b":[{"name":"j:k"},{"templated":true}]}}`),
		HALDocument{Links: Links{
			"a": {{Href: "xyz"}},
			"b": {{Name: "j:k"}, {Templated: true}},
			"n": {},
		}}},

	//// TODO sort link fields
	//{[]byte(`{"_links":{"curies":[{"name":"pre","href":"any","templated":true}]}}`),
	//	HALDocument{Links: Links{
	//		"curies": {{Name: "pre", Href: "any", Templated: true}},
	//	}}},

	{[]byte(`{"_embedded":{"a":{"x":"y"},"b":[{"c":"d"},{}]}}`),
		HALDocument{Embedded: Resources{
			"a": {&Resource{Attributes: Attributes{"x": "y"}}},
			"b": {&Resource{Attributes: Attributes{"c": "d"}}, &Resource{}},
			"n": {},
		}}},
}

func TestHALDocumentMarshalJSON(t *testing.T) {
	for _, tt := range HALDocumentJSONTests {
		result, err := json.Marshal(tt.doc)

		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(result, tt.json) {
			t.Errorf("Expected '%s', got '%s'", tt.json, result)
		}
	}
}
