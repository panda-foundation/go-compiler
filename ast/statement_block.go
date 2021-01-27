package ast

type Block struct {
	StatementBase
	Statements []Statement
}

func (b *Block) GenerateIR(c *Context) {
	for _, stmt := range b.Statements {
		ctx := c
		if _, ok := stmt.(*Block); ok {
			ctx = c.NewContext()
		}
		stmt.GenerateIR(ctx)
	}
}
