package hypermedia

import (
	"encoding/json"
	"reflect"
)

type HALDocument Resource

// MarshalJSON implements the json.Marshaler interface. The result is
// appropriate for the "application/hal+json" media type.
func (h HALDocument) MarshalJSON() ([]byte, error) {
	return json.Marshal(halResource(h))
}

// UnmarshalJSON implements the json.Unmarshaler interface. The input should
// be of the "application/hal+json" media type.
func (h *HALDocument) UnmarshalJSON(data []byte) (err error) {
	return json.Unmarshal(data, (*halResource)(h))
}

type (
	halLinks Links

	halResource  Resource
	halResources Resources
)

// MarshalJSON implements the json.Marshaler interface.
func (h halLinks) MarshalJSON() ([]byte, error) {
	o := make(map[Relation]interface{})

	for rel, links := range h {
		switch {
		case len(links) == 0:
			continue
		case len(links) == 1 && rel != "curies":
			o[rel] = links[0]
		default:
			o[rel] = links
		}
	}

	return json.Marshal(o)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (h *halLinks) UnmarshalJSON(data []byte) (err error) {
	var o map[Relation]json.RawMessage
	if err = json.Unmarshal(data, &o); err != nil {
		return
	}
	if o == nil {
		return &json.UnmarshalTypeError{Value: "null", Type: reflect.TypeOf(h)}
	}

	for rel, data := range o {
		var l []Link
		if err = json.Unmarshal(data, &l); err != nil {
			var l Link
			if err = json.Unmarshal(data, &l); err != nil {
				return
			}

			(*h)[rel] = []Link{l}
		} else if l == nil {
			return &json.UnmarshalTypeError{Value: "null", Type: reflect.TypeOf(l)}
		} else {
			(*h)[rel] = l
		}
	}

	return
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

// UnmarshalJSON implements the json.Unmarshaler interface.
func (h *halResources) UnmarshalJSON(data []byte) (err error) {
	var o map[Relation]json.RawMessage
	if err = json.Unmarshal(data, &o); err != nil {
		return
	}
	if o == nil {
		return &json.UnmarshalTypeError{Value: "null", Type: reflect.TypeOf(h)}
	}

	for rel, data := range o {
		var r []*halResource
		if err = json.Unmarshal(data, &r); err != nil {
			var r *halResource
			if err = json.Unmarshal(data, &r); err != nil {
				return
			}

			(*h)[rel] = []*Resource{(*Resource)(r)}
		} else if r == nil {
			return &json.UnmarshalTypeError{Value: "null", Type: reflect.TypeOf(r)}
		} else {
			for _, r := range r {
				(*h)[rel] = append((*h)[rel], (*Resource)(r))
			}
		}
	}

	return
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

// UnmarshalJSON implements the json.Unmarshaler interface.
func (h *halResource) UnmarshalJSON(data []byte) (err error) {
	var o map[string]json.RawMessage
	if err = json.Unmarshal(data, &o); err != nil {
		return
	}
	if o == nil {
		return &json.UnmarshalTypeError{Value: "null", Type: reflect.TypeOf(h)}
	}

	r := NewResource()

	if data, ok := o["_links"]; ok {
		delete(o, "_links")
		if err = json.Unmarshal(data, (*halLinks)(&r.Links)); err != nil {
			return
		}
	}

	if data, ok := o["_embedded"]; ok {
		delete(o, "_embedded")
		if err = json.Unmarshal(data, (*halResources)(&r.Embedded)); err != nil {
			return
		}
	}

	for k, data := range o {
		var v interface{}
		json.Unmarshal(data, &v)
		r.Attributes[k] = v
	}

	*h = halResource(*r)

	return
}
