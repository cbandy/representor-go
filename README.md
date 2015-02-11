
Go library for building and consuming Hypermedia messages. See the [Hypermedia
Project Charter](https://github.com/the-hypermedia-project/charter) for details.

### Proposal

```go

type Attributes map[string]interface{}
type Relation   string

//type Transition struct {
//  Href string
//  Method string
//}

type Link struct {
  Deprecation, Href, HrefLang, Name, Profile, Title, Type string
  //EncType, Method, Render, Target, RequestEncoding  string

  Templated bool
}

type Links map[Relation][]Link

func (Links) Add(Relation, Link)
func (Links) Del(Relation)
func (Links) Get(Relation) Link
func (Links) Set(Relation, Link)

type Resources map[Relation][]*Resource

func (Resources) Add(Relation, *Resource)
func (Resources) Del(Relation)
func (Resources) Get(Relation) *Resource
func (Resources) Set(Relation, *Resource)

type Resource struct {
  Attributes Attributes
  Embedded   Resources
  Links      Links

  // Populated during deserialization
  curies []Link
}

&Resource{
  Attributes: Attributes{"x": a, "y": b},
}

func NewResource() *Resource

// Expands rel using CURIEs
func (Resource) GetLink(Relation) Link
func (Resource) GetLinks(Relation) []Link
func (Resource) GetResource(Relation) *Resource
func (Resource) GetResources(Relation) []*Resource

type HALDocument Resource

func (HALDocument) MarshalJSON()
func (HALDocument) UnmarshalJSON()
func (HALDocument) assignCURIEs(*Resource, curies []Link)

```
