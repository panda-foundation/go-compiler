package expression

import (
	"fmt"

	"github.com/panda-foundation/go-compiler/ast/node"
	"github.com/panda-foundation/go-compiler/ir"
)

type MemberAccess struct {
	Base
	Parent        Expression
	Member        *Identifier
	FullNamespace string
}

func (m *MemberAccess) Type(c *node.Context) ir.Type {
	return m.GenerateIR(c).Type()
}

func (m *MemberAccess) GenerateParentIR(c *node.Context) ir.Value {
	if ident, ok := m.Parent.(*Identifier); ok {
		// resolve here
		parent, _ := c.FindSelector(ident.Name, m.Member.Name)
		return parent
	}
	return m.Parent.GenerateIR(c)
}

func (m *MemberAccess) GenerateIR(c *node.Context) ir.Value {
	if ident, ok := m.Parent.(*Identifier); ok {
		// resolve here
		_, obj := c.FindSelector(ident.Name, m.Member.Name)
		if obj == nil {
			c.Error(m.Position, fmt.Sprintf("%s is undefined", m.Member.Name))
		}
		return obj
	} else if _, ok := m.Parent.(*This); ok {
		this := c.FindObject(node.This)
		fmt.Println(this.Type())
	} /* else {
		// TO-DO
		// generate parent firstly, then check type of parent, then generate ir
	}*/
	return nil
}
