package parser

type Node interface {
	// node() just to make sure only AST types implement this interface.
	node()
}

// root
type FileNode struct {
	Items []Node
}

// Block represents an object with type, name and properties, this properties can be other blocks or assignments
type BlockNode struct {
	Type string
	Name string
	Body *FileNode
}

// AssignmentNode represents a key=value, the value can be a literal, a magnitude, a reference or a list
type AssignmentNode struct {
	Key   string
	Value Node
}

// A number with a unit
type MagnitudeNode struct {
	Value float64
	Unit  string
}

// References to other values
type ReferenceNode struct {
	Path []string
}

type StringLiteral struct{ Value string }
type NumberLiteral struct{ Value float64 }
type BoolLiteral struct{ Value bool }
type NullLiteral struct{ Value struct{} }

type ListLiteral struct {
	Items []Node
}

func (f FileNode) node()       {}
func (b BlockNode) node()      {}
func (a AssignmentNode) node() {}
func (s StringLiteral) node()  {}
func (n NumberLiteral) node()  {}
func (b BoolLiteral) node()    {}
func (n NullLiteral) node()    {}
func (l ListLiteral) node()    {}
func (m MagnitudeNode) node()  {}
func (r ReferenceNode) node()  {}
