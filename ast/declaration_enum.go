package ast

import (
	"fmt"
	"strconv"

	"github.com/panda-foundation/go-compiler/ir"
	"github.com/panda-foundation/go-compiler/token"
)

type Enum struct {
	DeclarationBase
	Members []*Variable

	IRMembers map[string]*ir.Global
}

func (e *Enum) AddVariable(m *Variable) error {
	for _, v := range e.Members {
		if v.Name.Name == m.Name.Name {
			return fmt.Errorf("%s redeclared", m.Name.Name)
		}
	}
	e.Members = append(e.Members, m)
	return nil
}

func (e *Enum) GenerateIR(c *Context) {
	var index int64 = 0
	for _, v := range e.Members {
		if v.Value == nil {
			c.Program.Module.NewGlobalDef(v.Qualified(c.Module.Namespace), ir.NewInt(ir.I32, index))
			index++
		} else {
			if literal, ok := v.Value.(*Literal); ok {
				if literal.Typ == token.INT {
					if i, _ := strconv.Atoi(literal.Value); int64(i) >= index {
						index = int64(i)
						c.Program.Module.NewGlobalDef(v.Qualified(c.Module.Namespace), ir.NewInt(ir.I32, index))
						index++
					} else {
						c.Error(v.Position, fmt.Sprintf("enum value here should be greater than %d.", i-1))
					}
				} else {
					c.Error(v.Position, "enum value must be integer.")
				}
			} else {
				c.Error(v.Position, "enum value must be integer.")
			}
		}
	}
}