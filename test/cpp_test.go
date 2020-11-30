package test

import (
	"testing"

	"github.com/panda-foundation/go-compiler/parser"
)

func TestVector(t *testing.T) {
	p := parser.NewParser([]string{"cpp"})

	p.ParseFile("./vector.pd")
}