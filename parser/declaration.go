package parser

import (
	"fmt"

	"github.com/panda-foundation/go-compiler/ast/declaration"
	"github.com/panda-foundation/go-compiler/ast/expression"
	"github.com/panda-foundation/go-compiler/token"
)

func (p *Parser) parseVariable(modifier *declaration.Modifier, attributes []*declaration.Attribute) *declaration.Variable {
	d := &declaration.Variable{}
	d.Modifier = modifier
	d.Custom = attributes
	d.Token = p.token
	p.next()
	d.Name = p.parseIdentifier()
	d.Type = p.parseType()

	if p.token == token.Assign {
		p.next()
		d.Value = p.parseExpression()
	}
	p.expect(token.Semi)
	return d
}

func (p *Parser) parseFunction(modifier *declaration.Modifier, attributes []*declaration.Attribute, class *declaration.Class) *declaration.Function {
	d := &declaration.Function{}
	d.Class = class
	d.Modifier = modifier
	d.Custom = attributes
	p.next()
	tilde := false
	if p.token == token.Complement {
		if class == nil {
			p.error(p.position, "'~' is not allow outside class as function name")
		}
		tilde = true
		p.next()
	}
	d.Name = p.parseIdentifier()
	if tilde {
		if d.Name.Name != class.Name.Name {
			p.error(p.position, "invalid destructor name")
		}
		d.Name.Name = "~" + d.Name.Name
	}
	if p.token == token.Less {
		d.TypeParameters = p.parseTypeParameters()
	}
	d.Parameters = p.parseParameters()
	if p.token != token.Semi && p.token != token.LeftBrace {
		d.ReturnType = p.parseType()
	}
	if p.token == token.LeftBrace {
		d.Body = p.parseCompoundStatement()
	} else if p.token == token.Semi {
		p.next()
	}
	return d
}

func (p *Parser) parseEnum(modifier *declaration.Modifier, attributes []*declaration.Attribute) *declaration.Enum {
	d := &declaration.Enum{
		Members: make(map[string]*declaration.Variable),
	}
	d.Modifier = modifier
	d.Custom = attributes
	p.next()
	d.Name = p.parseIdentifier()
	p.expect(token.LeftBrace)
	for p.token != token.RightBrace {
		v := &declaration.Variable{}
		v.Name = p.parseIdentifier()
		if p.token == token.Assign {
			p.next()
			v.Value = p.parseExpression()
		}
		if p.token != token.Comma {
			break
		}
		p.next()
	}
	p.expect(token.RightBrace)
	return d
}

func (p *Parser) parseInterface(modifier *declaration.Modifier, attributes []*declaration.Attribute) *declaration.Interface {
	d := &declaration.Interface{
		Functions:  make(map[string]*declaration.Function),
		Interfaces: make(map[string]*declaration.Interface),
	}
	d.Modifier = modifier
	d.Custom = attributes
	p.next()
	d.Name = p.parseIdentifier()
	if p.token == token.Less {
		d.TypeParameters = p.parseTypeParameters()
	}
	if p.token == token.Colon {
		d.Parents = p.parseIneritanceTypes()
	}
	p.expect(token.LeftBrace)
	for p.token != token.RightBrace {
		attr := p.parseAttributes()
		switch p.token {
		case token.Function:
			f := p.parseFunction(nil, attr, nil)
			name := f.Name.Name
			if _, ok := d.Functions[name]; ok {
				p.error(f.Name.Position, fmt.Sprintf("function %s redeclared", name))
			}
			d.Functions[name] = f

		case token.Interface:
			i := p.parseInterface(nil, attr)
			name := i.Name.Name
			if _, ok := d.Interfaces[name]; ok {
				p.error(i.Name.Position, fmt.Sprintf("interface %s redeclared", name))
			}
			d.Interfaces[name] = i

		default:
			p.expectedError(p.position, "declaration")
		}
	}
	p.expect(token.RightBrace)
	return d
}

func (p *Parser) parseClass(modifier *declaration.Modifier, attributes []*declaration.Attribute) *declaration.Class {
	d := &declaration.Class{
		Variables:  make(map[string]*declaration.Variable),
		Functions:  make(map[string]*declaration.Function),
		Enums:      make(map[string]*declaration.Enum),
		Interfaces: make(map[string]*declaration.Interface),
		Classes:    make(map[string]*declaration.Class),
	}
	d.Modifier = modifier
	d.Custom = attributes
	p.next()
	d.Name = p.parseIdentifier()
	if p.token == token.Less {
		d.TypeParameters = p.parseTypeParameters()
	}
	if p.token == token.Colon {
		d.Parents = p.parseIneritanceTypes()
	}
	p.expect(token.LeftBrace)
	for p.token != token.RightBrace {
		attr := p.parseAttributes()
		modifier := p.parseModifier()
		switch p.token {
		case token.Const, token.Var:
			v := p.parseVariable(modifier, attr)
			name := v.Name.Name
			if _, ok := d.Variables[name]; ok {
				p.error(v.Name.Position, fmt.Sprintf("variable %s redeclared", name))
			}
			d.Variables[name] = v

		case token.Function:
			f := p.parseFunction(modifier, attr, d)
			name := f.Name.Name
			if _, ok := d.Functions[name]; ok {
				p.error(f.Name.Position, fmt.Sprintf("function %s redeclared", name))
			}
			d.Functions[name] = f

		case token.Enum:
			e := p.parseEnum(modifier, attr)
			name := e.Name.Name
			if _, ok := d.Enums[name]; ok {
				p.error(e.Name.Position, fmt.Sprintf("enum %s redeclared", name))
			}
			d.Enums[name] = e

		case token.Interface:
			i := p.parseInterface(modifier, attr)
			name := i.Name.Name
			if _, ok := d.Interfaces[name]; ok {
				p.error(i.Name.Position, fmt.Sprintf("interface %s redeclared", name))
			}
			d.Interfaces[name] = i

		case token.Class:
			c := p.parseClass(modifier, attr)
			name := c.Name.Name
			if _, ok := d.Classes[name]; ok {
				p.error(c.Name.Position, fmt.Sprintf("class %s redeclared", name))
			}
			d.Classes[name] = c

		default:
			p.expectedError(p.position, "declaration")
		}
	}
	p.expect(token.RightBrace)
	return d
}

func (p *Parser) parseModifier() *declaration.Modifier {
	m := &declaration.Modifier{}
	if p.token == token.Public {
		m.Public = true
		p.next()
	}
	if p.token == token.Static {
		m.Static = true
		p.next()
	}
	return m
}

func (p *Parser) parseAttributes() []*declaration.Attribute {
	if p.token != token.META {
		return nil
	}
	var attr []*declaration.Attribute
	for p.token == token.META {
		p.next()
		if p.token != token.IDENT {
			p.expect(token.IDENT)
		}
		m := &declaration.Attribute{Position: p.position}
		m.Name = p.literal
		p.next()

		if p.token == token.STRING {
			m.Text = p.literal
			p.next()
		} else if p.token == token.LeftParen {
			p.next()
			if p.token == token.STRING {
				m.Text = p.literal
				p.next()
			} else {
				m.Values = make(map[string]*expression.Literal)
				for {
					if p.token == token.IDENT {
						name := p.literal
						p.next()
						p.expect(token.Assign)
						switch p.token {
						case token.INT, token.FLOAT, token.CHAR, token.STRING, token.BOOL:
							if _, ok := m.Values[name]; ok {
								p.error(p.position, "duplicated attribute "+name)
							}
							m.Values[name] = &expression.Literal{
								Type:  p.token,
								Value: p.literal,
							}
							m.Values[name].Position = p.position
						default:
							p.expectedError(p.position, "basic literal (bool, char, int, float, string)")
						}
						p.next()
						if p.token == token.RightParen {
							break
						}
						p.expect(token.Comma)
					} else {
						p.expect(token.IDENT)
					}
				}
			}
			p.expect(token.RightParen)
		}
		attr = append(attr, m)
	}
	return attr
}
