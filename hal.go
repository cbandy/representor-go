package hypermedia

import (
	"encoding/json"
	"reflect"
	"strings"
)

type HALDocument Resource

// MarshalJSON implements the json.Marshaler interface.
func (h HALDocument) MarshalJSON() ([]byte, error) {
	return json.Marshal(halResource(h))
}

type halLinkField struct {
	IsEmpty func(reflect.Value) bool
	Value   func(reflect.Value) interface{}
	Name    string
}

var halLinkFields []halLinkField

func init() {
	lt := reflect.TypeOf(Link{})
	for i := 0; i < lt.NumField(); i++ {
		lf := lt.Field(i)
		hlf := halLinkField{Name: strings.ToLower(lf.Name)}

		switch lf.Type.Kind() {
		case reflect.String:
			hlf.IsEmpty = func(v reflect.Value) bool { return v.Len() == 0 }
			hlf.Value = func(v reflect.Value) interface{} { return v.String() }
		case reflect.Bool:
			hlf.IsEmpty = func(v reflect.Value) bool { return !v.Bool() }
			hlf.Value = func(v reflect.Value) interface{} { return v.Bool() }
		}

		halLinkFields = append(halLinkFields, hlf)
	}
}

type (
	halLink  Link
	halLinks Links

	halResource  Resource
	halResources Resources
)

// MarshalJSON implements the json.Marshaler interface.
func (h halLink) MarshalJSON() ([]byte, error) {
	o := make(map[string]interface{})
	lv := reflect.ValueOf(Link(h))

	for i, f := range halLinkFields {
		fv := lv.Field(i)
		if !f.IsEmpty(fv) {
			o[f.Name] = f.Value(fv)
		}
	}

	return json.Marshal(o)
}

// MarshalJSON implements the json.Marshaler interface.
func (h halLinks) MarshalJSON() ([]byte, error) {
	o := make(map[Relation]interface{})

	for rel, links := range h {
		switch len(links) {
		case 0:
			continue
		case 1:
			o[rel] = halLink(links[0])
		default:
			a := make([]halLink, 0, len(links))
			for _, link := range links {
				a = append(a, halLink(link))
			}
			o[rel] = a
		}
	}

	return json.Marshal(o)
}

// MarshalJSON implements the json.Marshaler interface.
func (h halResources) MarshalJSON() ([]byte, error) {
	o := make(map[Relation]interface{})

	for rel, resources := range h {
		switch len(resources) {
		case 0:
			continue
		case 1:
			o[rel] = halResource(*resources[0])
		default:
			a := make([]halResource, 0, len(resources))
			for _, resource := range resources {
				a = append(a, halResource(*resource))
			}
			o[rel] = a
		}
	}

	return json.Marshal(o)
}

// MarshalJSON implements the json.Marshaler interface.
func (h halResource) MarshalJSON() ([]byte, error) {
	o := make(map[string]interface{})

	for k, v := range h.Attributes {
		o[k] = v
	}
	if len(h.Links) > 0 {
		o["_links"] = halLinks(h.Links)
	}
	if len(h.Embedded) > 0 {
		o["_embedded"] = halResources(h.Embedded)
	}

	return json.Marshal(o)
}
