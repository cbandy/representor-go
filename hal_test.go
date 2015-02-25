package hypermedia

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

var HALDocumentJSONTests = []struct {
	json []byte
	doc  HALDocument
}{
	{[]byte(`{}`), HALDocument{}},
	{[]byte(`{"a":"xyz","b":true,"c":123.4}`),
		HALDocument{Attributes: Attributes{
			"a": "xyz",
			"b": true,
			"c": 123.4,
		}}},
	{[]byte(`{"_links":{"a":{"href":"xyz"},"b":[{"name":"j:k"},{"templated":true}]}}`),
		HALDocument{Links: Links{
			"a": {{"href": "xyz"}},
			"b": {{"name": "j:k"}, {"templated": true}},
		}}},

	{[]byte(`{"_links":{"curies":[{"href":"any","name":"pre","templated":true}]}}`),
		HALDocument{Links: Links{
			"curies": {{"name": "pre", "href": "any", "templated": true}},
		}}},

	{[]byte(`{"_embedded":{"a":{"x":"y"},"b":[{"c":"d"},{}]}}`),
		HALDocument{Embedded: Resources{
			"a": {&Resource{Attributes: Attributes{"x": "y"}}},
			"b": {&Resource{Attributes: Attributes{"c": "d"}}, &Resource{}},
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

func TestHALDocumentMarshalJSONExcludesEmptyTransitions(t *testing.T) {
	for _, tt := range []struct {
		json []byte
		doc  HALDocument
	}{
		{[]byte(`{"_links":{"a":{"href":"xyz"}}}`),
			HALDocument{Links: Links{
				"a": {{"href": "xyz"}},
				"n": {},
			}}},
		{[]byte(`{"_embedded":{"a":{"x":"y"}}}`),
			HALDocument{Embedded: Resources{
				"a": {&Resource{Attributes: Attributes{"x": "y"}}},
				"n": {},
			}}},
	} {
		result, err := json.Marshal(tt.doc)

		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(result, tt.json) {
			t.Errorf("Expected '%s', got '%s'", tt.json, result)
		}
	}
}

func TestHALDocumentUnmarshalJSON(t *testing.T) {
	var emptyMaps func(*Resource)
	emptyMaps = func(r *Resource) {
		if r.Attributes == nil {
			r.Attributes = make(Attributes)
		}
		if r.Embedded == nil {
			r.Embedded = make(Resources)
		} else {
			for _, s := range r.Embedded {
				for _, r := range s {
					emptyMaps(r)
				}
			}
		}
		if r.Links == nil {
			r.Links = make(Links)
		}
	}

	for _, tt := range HALDocumentJSONTests {
		var result HALDocument
		err := json.Unmarshal(tt.json, &result)

		if err != nil {
			t.Fatal(err)
		}

		// reflect.DeepEqual distinguishes between empty map and nil map.
		// We expect empty.
		doc := tt.doc
		emptyMaps((*Resource)(&doc))

		if !reflect.DeepEqual(result, doc) {
			t.Errorf("Expected %#v, got %#v", doc, result)

			// Print the content of Resource pointers
			for k, v := range doc.Embedded {
				for i, r := range v {
					if !reflect.DeepEqual(result.Embedded[k][i], r) {
						t.Logf("%q,%v: expected %#v, got %#v", k, i, r, result.Embedded[k][i])
					}
				}
			}
		}
	}
}

func TestHALDocumentUnmarshalJSONErrors(t *testing.T) {
	for _, tt := range []string{
		// Its root object MUST be a Resource Object.
		``, `null`, `[]`, `123`,

		// It is an object whose property names are link relation types (as
		// defined by [RFC5988]) and values are either a Link Object or an array
		// of Link Objects.
		`{"_links":null}`, `{"_links":[]}`, `{"_links":123}`,
		`{"_links":{"a":null}}`, `{"_links":{"a":123}}`,

		// It is an object whose property names are link relation types (as
		// defined by [RFC5988]) and values are either a Resource Object or an
		// array of Resource Objects.
		`{"_embedded":null}`, `{"_embedded":[]}`, `{"_embedded":123}`,
		`{"_embedded":{"a":null}}`, `{"_embedded":{"a":123}}`,
	} {
		err := json.Unmarshal([]byte(tt), new(HALDocument))

		if err == nil {
			t.Errorf("Expected error for invalid hal+json: '%s'", tt)
		}
	}
}
