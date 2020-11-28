package ast

// Metadata type
type Metadata struct {
	Position int
	Name     string
	Text     string
	Values   map[string]*Literal
}

// Attributes of type
type Attributes struct {
	Modifier *Modifier
	Resolved []*Metadata
	Custom   []*Metadata
}

// Modifier type
type Modifier struct {
	Public bool
	Static bool
	Async  bool
	Inline bool
}

// Equal to compare is two modifiers are same
func (m *Modifier) Equal(modifier *Modifier) bool {
	return m.Public == modifier.Public && m.Static == modifier.Static || m.Async == modifier.Async || m.Inline == modifier.Inline
}

// NewProgram to create new program
func NewProgram(packageName string, parent *Program) *Program {
	return &Program{
		Package: packageName,

		Variables:  make(map[string]*Variable),
		Functions:  make(map[string]*Function),
		Enums:      make(map[string]*Enum),
		Interfaces: make(map[string]*Interface),
		Classes:    make(map[string]*Class),

		Parent:   parent,
		Children: make(map[string]*Program),
	}
}

// Program type
type Program struct {
	Attributes

	Package string

	Variables  map[string]*Variable
	Functions  map[string]*Function
	Enums      map[string]*Enum
	Interfaces map[string]*Interface
	Classes    map[string]*Class

	Parent   *Program
	Children map[string]*Program

	Document *Metadata
}