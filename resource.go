package hypermedia

type Attributes map[string]interface{}
type Relation string

type Link struct {
	Deprecation, Href, HrefLang, Name, Profile, Title, Type string

	Templated bool
}

type Links map[Relation][]Link

func (ls Links) Add(rel Relation, link Link) { ls[rel] = append(ls[rel], link) }
func (ls Links) Del(rel Relation)            { delete(ls, rel) }
func (ls Links) Set(rel Relation, link Link) { ls[rel] = []Link{link} }
func (ls Links) Get(rel Relation) (link Link) {
	if len(ls[rel]) > 0 {
		link = ls[rel][0]
	}
	return
}

type Resources map[Relation][]*Resource

func (rs Resources) Add(rel Relation, res *Resource) { rs[rel] = append(rs[rel], res) }
func (rs Resources) Del(rel Relation)                { delete(rs, rel) }
func (rs Resources) Set(rel Relation, res *Resource) { rs[rel] = []*Resource{res} }
func (rs Resources) Get(rel Relation) (res *Resource) {
	if len(rs[rel]) > 0 {
		res = rs[rel][0]
	}
	return
}

type Resource struct {
	Attributes Attributes
	Embedded   Resources
	Links      Links
}

func NewResource() *Resource {
	return &Resource{
		Attributes: make(Attributes),
		Embedded:   make(Resources),
		Links:      make(Links),
	}
}
