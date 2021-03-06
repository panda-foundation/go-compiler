package ast

import (
	"fmt"

	"github.com/panda-foundation/go-compiler/ir"
)

type MemberAccess struct {
	ExpressionBase
	Parent        Expression
	Member        *Identifier
	FullNamespace string
}

//TO-DO subscripting
func (m *MemberAccess) Type(c *Context, expected ir.Type) ir.Type {
	// parent could be: identifier, member_access, new, subscripting, this, base
	if ident, ok := m.Parent.(*Identifier); ok {
		_, obj, _ := c.FindSelector(ident.Name, m.Member.Name)
		if obj != nil {
			return obj.Type()
		}

	} else if _, ok := m.Parent.(*This); ok {
		return c.Function.Class.MemberType(m.Member.Name)

	} else if _, ok := m.Parent.(*Base); ok {
		return c.Function.Class.MemberType(m.Member.Name)

	} else if n, ok := m.Parent.(*New); ok {
		_, d := c.Program.FindDeclaration(n.Typ)
		if class, ok := d.(*Class); ok {
			return class.MemberType(m.Member.Name)
		}

	} else if memberAccess, ok := m.Parent.(*MemberAccess); ok {
		parentType := memberAccess.Type(c, nil)
		if s, ok := parentType.(*ir.StructType); ok {
			if d, ok := c.Program.Declarations[s.TypeName]; ok {
				if class, ok := d.(*Class); ok {
					return class.MemberType(m.Member.Name)
				} else if _, ok := d.(*Enum); ok {
					return ir.I32
				}
			}
		}
	}

	return nil
}

func (m *MemberAccess) GenerateIR(c *Context, expected ir.Type) ir.Value {
	var v ir.Value
	var p ir.Value
	var isMemberFunction bool
	if ident, ok := m.Parent.(*Identifier); ok {
		p, v, isMemberFunction = c.FindSelector(ident.Name, m.Member.Name)

	} else if _, ok := m.Parent.(*This); ok {
		p = c.FindObject(ClassThis)
		v, isMemberFunction = c.Function.Class.GetMember(c, p, m.Member.Name, true)

	} else if _, ok := m.Parent.(*Base); ok {
		p = c.FindObject(ClassThis)
		v, isMemberFunction = c.Function.Class.Parent.GetMember(c, p, m.Member.Name, true)

	} else if n, ok := m.Parent.(*New); ok {
		qualified, d := c.Program.FindDeclaration(n.Typ)
		if class, ok := d.(*Class); ok {
			p = m.Parent.GenerateIR(c, nil)
			if IsBuiltinClass(qualified) {
				v, isMemberFunction = class.GetMember(c, p, m.Member.Name, false)
			} else {
				p, v, isMemberFunction = class.GetMemberFromCounter(c, p, m.Member.Name)
			}
		}

	} else if memberAccess, ok := m.Parent.(*MemberAccess); ok {
		parentType := memberAccess.Type(c, nil)
		if s, ok := parentType.(*ir.StructType); ok {
			if d, ok := c.Program.Declarations[s.TypeName]; ok {
				if class, ok := d.(*Class); ok {
					p = m.Parent.GenerateIR(c, nil)
					if IsBuiltinClass(s.TypeName) {
						v, isMemberFunction = class.GetMember(c, p, m.Member.Name, false)
					} else {
						p, v, isMemberFunction = class.GetMemberFromCounter(c, p, m.Member.Name)
					}
				} else if enum, ok := d.(*Enum); ok {
					v = enum.GetMember(m.Member.Name)
				}
			}
		}
	}
	if v != nil && m.IsFunction(v) {
		v = c.AutoLoad(v)
		v = ir.NewCall(v)
		if p != nil && isMemberFunction {
			call := v.(*ir.InstCall)
			call.Args = append(call.Args, p)
		}
	}
	if v == nil {
		c.Program.Error(m.Position, fmt.Sprintf("%s undefined", m.Member.Name))
	}
	return v
}

func (m *MemberAccess) IsFunction(v ir.Value) bool {
	if t, ok := v.Type().(*ir.PointerType); ok {
		if _, ok = t.ElemType.(*ir.FuncType); ok {
			return true
		} else if e, ok := t.ElemType.(*ir.PointerType); ok {
			// gep instruction
			_, ok = e.ElemType.(*ir.FuncType)
			return ok
		}
	}
	return false
}

func (m *MemberAccess) IsConstant(p *Program) bool {
	if ident, ok := m.Parent.(*Identifier); ok {
		_, d := p.FindSelector(ident.Name, m.Member.Name)
		if d == nil {
			// could be an enum
			_, e := p.FindSelector("", ident.Name)
			if _, ok := e.(*Enum); ok {
				return true
			}
			return false
		}
		if v, ok := d.(*Variable); ok {
			return v.Const && v.Value.IsConstant(p)
		}
		if _, ok := d.(*Function); ok {
			return true
		}

	} else if memberAccess, ok := m.Parent.(*MemberAccess); ok {
		if identifier, ok := memberAccess.Parent.(*Identifier); ok {
			_, e := p.FindSelector(identifier.Name, memberAccess.Member.Name)
			if _, ok := e.(*Enum); ok {
				return true
			}
		}
	}
	return false
}

func (m *MemberAccess) GenerateConstIR(p *Program, expected ir.Type) ir.Constant {
	if ident, ok := m.Parent.(*Identifier); ok {
		_, d := p.FindSelector(ident.Name, m.Member.Name)
		if d == nil {
			// could be an enum
			_, e := p.FindSelector("", ident.Name)
			if enum, ok := e.(*Enum); ok {
				return enum.GetMember(m.Member.Name)
			}
		}
		if v, ok := d.(*Variable); ok {
			return v.Value.GenerateConstIR(p, expected)
		}
		if f, ok := d.(*Function); ok {
			//TO-DO use function as pointer
			return f.IRFunction
		}
	} else if memberAccess, ok := m.Parent.(*MemberAccess); ok {
		if identifier, ok := memberAccess.Parent.(*Identifier); ok {
			_, e := p.FindSelector(identifier.Name, memberAccess.Member.Name)
			if enum, ok := e.(*Enum); ok {
				return enum.GetMember(m.Member.Name)
			}
		}
	}
	p.Error(m.Position, "invalid constant declaration")
	return nil
}
