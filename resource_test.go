package hypermedia

import (
	"reflect"
	"testing"
)

func TestLinksAdd(t *testing.T) {
	links := make(Links)
	links.Add("a", Link{Href: "b"})

	expected := Links{"a": []Link{{Href: "b"}}}
	if !reflect.DeepEqual(links, expected) {
		t.Errorf("Expected %+v, got %+v", expected, links)
	}
}

func TestLinksDel(t *testing.T) {
	links := Links{"a": []Link{{Href: "b"}}}
	links.Del("a")

	expected := Links{}
	if !reflect.DeepEqual(links, expected) {
		t.Errorf("Expected %+v, got %+v", expected, links)
	}
}

func TestLinksSet(t *testing.T) {
	links := Links{"a": []Link{{Href: "b"}}}
	links.Set("a", Link{Href: "c"})

	expected := Links{"a": []Link{{Href: "c"}}}
	if !reflect.DeepEqual(links, expected) {
		t.Errorf("Expected %+v, got %+v", expected, links)
	}
}

func TestLinksGetEmpty(t *testing.T) {
	link := make(Links).Get("a")

	if !reflect.DeepEqual(link, Link{}) {
		t.Errorf("Expected empty Link, got %+v", link)
	}
}

func TestLinksGet(t *testing.T) {
	link := Links{"a": []Link{{Href: "b"}}}.Get("a")

	expected := Link{Href: "b"}
	if link != expected {
		t.Errorf("Expected %+v, got %+v", expected, link)
	}
}

func TestResourcesAdd(t *testing.T) {
	resource := new(Resource)
	resources := make(Resources)
	resources.Add("a", resource)

	expected := Resources{"a": []*Resource{resource}}
	if !reflect.DeepEqual(resources, expected) {
		t.Errorf("Expected %+v, got %+v", expected, resources)
	}
}

func TestResourcesDel(t *testing.T) {
	resources := Resources{"a": []*Resource{new(Resource)}}
	resources.Del("a")

	expected := Resources{}
	if !reflect.DeepEqual(resources, expected) {
		t.Errorf("Expected %+v, got %+v", expected, resources)
	}
}

func TestResourcesSet(t *testing.T) {
	a := new(Resource)
	b := new(Resource)
	resources := Resources{"a": []*Resource{a}}
	resources.Set("a", b)

	expected := Resources{"a": []*Resource{b}}
	if !reflect.DeepEqual(resources, expected) {
		t.Errorf("Expected %+v, got %+v", expected, resources)
	}
}

func TestResourcesGetEmpty(t *testing.T) {
	resource := make(Resources).Get("a")

	if resource != nil {
		t.Errorf("Expected nil, got %+v", resource)
	}
}

func TestResourcesGet(t *testing.T) {
	a := new(Resource)
	resource := Resources{"a": []*Resource{a}}.Get("a")

	if resource != a {
		t.Errorf("Expected %+v, got %+v", a, resource)
	}
}

func TestNewResource(t *testing.T) {
	// No panics
	resource := NewResource()
	resource.Attributes["a"] = "b"
	resource.Embedded.Add("c", new(Resource))
	resource.Links.Add("d", Link{})
}
