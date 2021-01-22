// Package types declares the data types of LLVM IR.
package ir

import (
	"fmt"
	"strings"
)

// === [ Types ] ===
// Type is an LLVM IR type.
//
// A Type has one of the following underlying types.
//
//    *types.VoidType       // https://godoc.org/github.com/llir/llvm/ir/types#VoidType
//    *types.FuncType       // https://godoc.org/github.com/llir/llvm/ir/types#FuncType
//    *types.IntType        // https://godoc.org/github.com/llir/llvm/ir/types#IntType
//    *types.FloatType      // https://godoc.org/github.com/llir/llvm/ir/types#FloatType
//    *types.PointerType    // https://godoc.org/github.com/llir/llvm/ir/types#PointerType
//    *types.VectorType     // https://godoc.org/github.com/llir/llvm/ir/types#VectorType
//    *types.LabelType      // https://godoc.org/github.com/llir/llvm/ir/types#LabelType
//    *types.TokenType      // https://godoc.org/github.com/llir/llvm/ir/types#TokenType
//    *types.MetadataType   // https://godoc.org/github.com/llir/llvm/ir/types#MetadataType
//    *types.ArrayType      // https://godoc.org/github.com/llir/llvm/ir/types#ArrayType
//    *types.StructType     // https://godoc.org/github.com/llir/llvm/ir/types#StructType

type Type interface {
	fmt.Stringer
	LLString() string
	Equal(u Type) bool
	Name() string
	SetName(name string)
}

// Convenience types.
var (
	// Basic types.
	Void  = &VoidType{}  // void
	Label = &LabelType{} // label
	Token = &TokenType{} // token

	// Integer types.
	I1   = &IntType{BitSize: 1}                  // i1
	I8   = &IntType{BitSize: 8}                  // i8
	I16  = &IntType{BitSize: 16}                 // i16
	I32  = &IntType{BitSize: 32}                 // i32
	I64  = &IntType{BitSize: 64}                 // i64
	UI8  = &IntType{BitSize: 8, Unsigned: true}  // i8
	UI16 = &IntType{BitSize: 16, Unsigned: true} // i16
	UI32 = &IntType{BitSize: 32, Unsigned: true} // i32
	UI64 = &IntType{BitSize: 64, Unsigned: true} // i64

	// Floating-point types.
	Float32 = &FloatType{Kind: FloatKindFloat}  // float
	Float64 = &FloatType{Kind: FloatKindDouble} // double

	// Integer pointer types.
	I8Ptr  = &PointerType{ElemType: I8}  // i8*
	I16Ptr = &PointerType{ElemType: I16} // i16*
	I32Ptr = &PointerType{ElemType: I32} // i32*
	I64Ptr = &PointerType{ElemType: I64} // i64*
)

func IsVoid(t Type) bool {
	_, ok := t.(*VoidType)
	return ok
}

// IsFunc reports whether the given type is a function type.
func IsFunc(t Type) bool {
	_, ok := t.(*FuncType)
	return ok
}

// IsInt reports whether the given type is an integer type.
func IsInt(t Type) bool {
	_, ok := t.(*IntType)
	return ok
}

// IsFloat reports whether the given type is a floating-point type.
func IsFloat(t Type) bool {
	_, ok := t.(*FloatType)
	return ok
}

func IsNumber(t Type) bool {
	return IsInt(t) || IsFloat(t)
}

// IsPointer reports whether the given type is a pointer type.
func IsPointer(t Type) bool {
	_, ok := t.(*PointerType)
	return ok
}

// IsVector reports whether the given type is a vector type.
func IsVector(t Type) bool {
	_, ok := t.(*VectorType)
	return ok
}

// IsLabel reports whether the given type is a label type.
func IsLabel(t Type) bool {
	_, ok := t.(*LabelType)
	return ok
}

// IsToken reports whether the given type is a token type.
func IsToken(t Type) bool {
	_, ok := t.(*TokenType)
	return ok
}

// IsArray reports whether the given type is an array type.
func IsArray(t Type) bool {
	_, ok := t.(*ArrayType)
	return ok
}

// IsStruct reports whether the given type is a struct type.
func IsStruct(t Type) bool {
	_, ok := t.(*StructType)
	return ok
}

// Equal reports whether t and u are of equal type.
func Equal(t, u Type) bool {
	return t.Equal(u)
}

// --- [ Void types ] ----------------------------------------------------------

// VoidType is an LLVM IR void type.
type VoidType struct {
	// Type name; or empty if not present.
	TypeName string
}

// Equal reports whether t and u are of equal type.
func (t *VoidType) Equal(u Type) bool {
	if _, ok := u.(*VoidType); ok {
		return true
	}
	return false
}

// String returns the string representation of the void type.
func (t *VoidType) String() string {
	if len(t.TypeName) > 0 {
		return TypeName(t.TypeName)
	}
	return t.LLString()
}

// LLString returns the LLVM syntax representation of the definition of the
func (t *VoidType) LLString() string {
	return "void"
}

// Name returns the type name of the type.
func (t *VoidType) Name() string {
	return t.TypeName
}

// SetName sets the type name of the type.
func (t *VoidType) SetName(name string) {
	t.TypeName = name
}

// --- [ Function types ] ------------------------------------------------------

// FuncType is an LLVM IR function type.
type FuncType struct {
	// Type name; or empty if not present.
	TypeName string
	// Return type.
	RetType Type
	// Function parameters.
	Params []Type
	// Variable number of function arguments.
	Variadic bool
}

// NewFunc returns a new function type based on the given return type and
// function parameter types.
func NewFuncType(retType Type, params ...Type) *FuncType {
	return &FuncType{
		RetType: retType,
		Params:  params,
	}
}

// Equal reports whether t and u are of equal type.
func (t *FuncType) Equal(u Type) bool {
	if u, ok := u.(*FuncType); ok {
		if !t.RetType.Equal(u.RetType) {
			return false
		}
		if len(t.Params) != len(u.Params) {
			return false
		}
		for i := range t.Params {
			if !t.Params[i].Equal(u.Params[i]) {
				return false
			}
		}
		return t.Variadic == u.Variadic
	}
	return false
}

// String returns the string representation of the function type.
func (t *FuncType) String() string {
	if len(t.TypeName) > 0 {
		return TypeName(t.TypeName)
	}
	return t.LLString()
}

// LLString returns the LLVM syntax representation of the definition of the
// type.
//
// RetType=Type '(' Params ')'
func (t *FuncType) LLString() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s (", t.RetType)
	for i, param := range t.Params {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(param.String())
	}
	if t.Variadic {
		if len(t.Params) > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString("...")
	}
	buf.WriteString(")")
	return buf.String()
}

// Name returns the type name of the type.
func (t *FuncType) Name() string {
	return t.TypeName
}

// SetName sets the type name of the type.
func (t *FuncType) SetName(name string) {
	t.TypeName = name
}

// --- [ Integer types ] -------------------------------------------------------

// IntType is an LLVM IR integer type.
type IntType struct {
	// Type name; or empty if not present.
	TypeName string
	// Integer size in number of bits.
	BitSize uint64
	// If int is unsigned
	Unsigned bool
}

// NewIntType returns a new integer type based on the given integer bit size.
func NewIntType(bitSize uint64) *IntType {
	return &IntType{
		BitSize: bitSize,
	}
}

// Equal reports whether t and u are of equal type.
func (t *IntType) Equal(u Type) bool {
	if u, ok := u.(*IntType); ok {
		return t.BitSize == u.BitSize
	}
	return false
}

// String returns the string representation of the integer type.
func (t *IntType) String() string {
	if len(t.TypeName) > 0 {
		return TypeName(t.TypeName)
	}
	return t.LLString()
}

// LLString returns the LLVM syntax representation of the definition of the
// type.
//
// int_type_tok
func (t *IntType) LLString() string {
	return fmt.Sprintf("i%d", t.BitSize)
}

// Name returns the type name of the type.
func (t *IntType) Name() string {
	return t.TypeName
}

// SetName sets the type name of the type.
func (t *IntType) SetName(name string) {
	t.TypeName = name
}

// --- [ Floating-point types ] ------------------------------------------------

// FloatType is an LLVM IR floating-point type.
type FloatType struct {
	// Type name; or empty if not present.
	TypeName string
	// Floating-point kind.
	Kind FloatKind
}

// Equal reports whether t and u are of equal type.
func (t *FloatType) Equal(u Type) bool {
	if u, ok := u.(*FloatType); ok {
		return t.Kind == u.Kind
	}
	return false
}

// String returns the string representation of the floating-point type.
func (t *FloatType) String() string {
	if len(t.TypeName) > 0 {
		return TypeName(t.TypeName)
	}
	return t.LLString()
}

// LLString returns the LLVM syntax representation of the definition of the
// type.
//
// FloatKind
func (t *FloatType) LLString() string {
	return t.Kind.String()
}

// Name returns the type name of the type.
func (t *FloatType) Name() string {
	return t.TypeName
}

// SetName sets the type name of the type.
func (t *FloatType) SetName(name string) {
	t.TypeName = name
}

// FloatKind represents the set of floating-point kinds.
type FloatKind uint8

// Floating-point kinds.
const (
	// 32-bit floating-point type (IEEE 754 single precision).
	FloatKindFloat FloatKind = iota // float
	// 64-bit floating-point type (IEEE 754 double precision).
	FloatKindDouble // double
)

func (i FloatKind) String() string {
	if i == FloatKindFloat {
		return "float"
	}
	return "double"
}

// --- [ Pointer types ] -------------------------------------------------------

// PointerType is an LLVM IR pointer type.
type PointerType struct {
	// Type name; or empty if not present.
	TypeName string
	// Element type.
	ElemType Type
}

// NewPointerType returns a new pointer type based on the given element type.
func NewPointerType(elemType Type) *PointerType {
	return &PointerType{
		ElemType: elemType,
	}
}

// Equal reports whether t and u are of equal type.
func (t *PointerType) Equal(u Type) bool {
	// HACK: to prevent infinite loops (e.g. struct foo containing field of type
	// pointer to foo).
	return t.String() == u.String()
}

// String returns the string representation of the pointer type.
func (t *PointerType) String() string {
	if len(t.TypeName) > 0 {
		return TypeName(t.TypeName)
	}
	return t.LLString()
}

// LLString returns the LLVM syntax representation of the definition of the
// type.
//
func (t *PointerType) LLString() string {
	buf := &strings.Builder{}
	buf.WriteString(t.ElemType.String())
	buf.WriteString("*")
	return buf.String()
}

// SetName sets the type name of the type.
func (t *PointerType) SetName(name string) {
	t.TypeName = name
}

// Name returns the type name of the type.
func (t *PointerType) Name() string {
	return t.TypeName
}

// --- [ Vector types ] --------------------------------------------------------

// VectorType is an LLVM IR vector type.
type VectorType struct {
	// Type name; or empty if not present.
	TypeName string
	// Scalable vector type.
	Scalable bool
	// Vector length.
	Len uint64
	// Element type.
	ElemType Type
}

// NewVectorType returns a new vector type based on the given vector length and
// element type.
func NewVectorType(len uint64, elemType Type) *VectorType {
	return &VectorType{
		Len:      len,
		ElemType: elemType,
	}
}

// Equal reports whether t and u are of equal type.
func (t *VectorType) Equal(u Type) bool {
	if u, ok := u.(*VectorType); ok {
		if t.Scalable != u.Scalable {
			return false
		}
		if t.Len != u.Len {
			return false
		}
		return t.ElemType.Equal(u.ElemType)
	}
	return false
}

// String returns the string representation of the vector type.
func (t *VectorType) String() string {
	if len(t.TypeName) > 0 {
		return TypeName(t.TypeName)
	}
	return t.LLString()
}

// LLString returns the LLVM syntax representation of the definition of the
// type.
//
// scalable: '<' 'vscale' 'x' Len=UintLit 'x' Elem=Type '>'
// non-scalable: '<' Len=UintLit 'x' Elem=Type '>'
func (t *VectorType) LLString() string {
	if t.Scalable {
		// '<' 'vscale' 'x' Len=UintLit 'x' Elem=Type '>'
		return fmt.Sprintf("<vscale x %d x %s>", t.Len, t.ElemType)
	}
	// '<' Len=UintLit 'x' Elem=Type '>'
	return fmt.Sprintf("<%d x %s>", t.Len, t.ElemType)
}

// Name returns the type name of the type.
func (t *VectorType) Name() string {
	return t.TypeName
}

// SetName sets the type name of the type.
func (t *VectorType) SetName(name string) {
	t.TypeName = name
}

// --- [ Label types ] ---------------------------------------------------------

// LabelType is an LLVM IR label type, which is used for basic block values.
type LabelType struct {
	// Type name; or empty if not present.
	TypeName string
}

// Equal reports whether t and u are of equal type.
func (t *LabelType) Equal(u Type) bool {
	if _, ok := u.(*LabelType); ok {
		return true
	}
	return false
}

// String returns the string representation of the label type.
func (t *LabelType) String() string {
	if len(t.TypeName) > 0 {
		return TypeName(t.TypeName)
	}
	return t.LLString()
}

// LLString returns the LLVM syntax representation of the definition of the
// type.
//
// 'label'
func (t *LabelType) LLString() string {
	return "label"
}

// Name returns the type name of the type.
func (t *LabelType) Name() string {
	return t.TypeName
}

// SetName sets the type name of the type.
func (t *LabelType) SetName(name string) {
	t.TypeName = name
}

// --- [ Token types ] ---------------------------------------------------------

// TokenType is an LLVM IR token type.
type TokenType struct {
	// Type name; or empty if not present.
	TypeName string
}

// Equal reports whether t and u are of equal type.
func (t *TokenType) Equal(u Type) bool {
	if _, ok := u.(*TokenType); ok {
		return true
	}
	return false
}

// String returns the string representation of the token type.
func (t *TokenType) String() string {
	if len(t.TypeName) > 0 {
		return TypeName(t.TypeName)
	}
	return t.LLString()
}

// LLString returns the LLVM syntax representation of the definition of the
// type.
//
// 'token'
func (t *TokenType) LLString() string {
	return "token"
}

// Name returns the type name of the type.
func (t *TokenType) Name() string {
	return t.TypeName
}

// SetName sets the type name of the type.
func (t *TokenType) SetName(name string) {
	t.TypeName = name
}

// --- [ Array types ] ---------------------------------------------------------

// ArrayType is an LLVM IR array type.
type ArrayType struct {
	// Type name; or empty if not present.
	TypeName string
	// Array length.
	Len uint64
	// Element type.
	ElemType Type
}

// NewArrayType returns a new array type based on the given array length and element
// type.
func NewArrayType(len uint64, elemType Type) *ArrayType {
	return &ArrayType{
		Len:      len,
		ElemType: elemType,
	}
}

// Equal reports whether t and u are of equal type.
func (t *ArrayType) Equal(u Type) bool {
	if u, ok := u.(*ArrayType); ok {
		if t.Len != u.Len {
			return false
		}
		return t.ElemType.Equal(u.ElemType)
	}
	return false
}

// String returns the string representation of the array type.
func (t *ArrayType) String() string {
	if len(t.TypeName) > 0 {
		return TypeName(t.TypeName)
	}
	return t.LLString()
}

// LLString returns the LLVM syntax representation of the definition of the
// type.
//
// '[' Len=UintLit 'x' Elem=Type ']'
func (t *ArrayType) LLString() string {
	return fmt.Sprintf("[%d x %s]", t.Len, t.ElemType)
}

// Name returns the type name of the type.
func (t *ArrayType) Name() string {
	return t.TypeName
}

// SetName sets the type name of the type.
func (t *ArrayType) SetName(name string) {
	t.TypeName = name
}

// --- [ Structure types ] -----------------------------------------------------

// StructType is an LLVM IR structure type. Identified (named) struct types are
// uniqued by type names, not by structural identity.
type StructType struct {
	// Type name; or empty if not present.
	TypeName string
	// Packed memory layout.
	Packed bool
	// Struct fields.
	Fields []Type
	// Opaque struct type.
	Opaque bool
}

// NewStructType returns a new struct type based on the given field types.
func NewStructType(fields ...Type) *StructType {
	return &StructType{
		Fields: fields,
	}
}

// Equal reports whether t and u are of equal type.
func (t *StructType) Equal(u Type) bool {
	if u, ok := u.(*StructType); ok {
		if len(t.TypeName) > 0 || len(u.TypeName) > 0 {
			// Identified struct types are uniqued by type names, not by structural
			// identity.
			//
			// t or u is an identified struct type.
			return t.TypeName == u.TypeName
		}
		// Literal struct types are uniqued by structural identity.
		if t.Packed != u.Packed {
			return false
		}
		if len(t.Fields) != len(u.Fields) {
			return false
		}
		for i := range t.Fields {
			if !t.Fields[i].Equal(u.Fields[i]) {
				return false
			}
		}
		return true
	}
	return false
}

// String returns the string representation of the structure type.
func (t *StructType) String() string {
	if len(t.TypeName) > 0 {
		return TypeName(t.TypeName)
	}
	return t.LLString()
}

// LLString returns the LLVM syntax representation of the definition of the
// type.
//
// Opaque struct type.
//
//    'opaque'
//
// Struct type.
//
//    '{' Fields=(Type separator ',')+? '}'
//
// Packed struct type.
//
//    '<' '{' Fields=(Type separator ',')+? '}' '>'   -> PackedStructType
func (t *StructType) LLString() string {
	if t.Opaque {
		return "opaque"
	}
	if len(t.Fields) == 0 {
		if t.Packed {
			return "<{}>"
		}
		return "{}"
	}
	buf := &strings.Builder{}
	if t.Packed {
		buf.WriteString("<")
	}
	buf.WriteString("{ ")
	for i, field := range t.Fields {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(field.String())
	}
	buf.WriteString(" }")
	if t.Packed {
		buf.WriteString(">")
	}
	return buf.String()
}

// Name returns the type name of the type.
func (t *StructType) Name() string {
	return t.TypeName
}

// SetName sets the type name of the type.
func (t *StructType) SetName(name string) {
	t.TypeName = name
}
