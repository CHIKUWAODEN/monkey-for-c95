package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"bitbucket.org/kandayasu/squirrel-go/ast"
)

type ObjectType string

const (
	FUNCTION_OBJ     = "FUNCTION"
	CLASS_OBJ        = "CLASS"
	INSTANCE_OBJ     = "INSTANCE"
	THIS_OBJ         = "THIS"
	INTEGER_OBJ      = "INTEGER"
	STRING_OBJ       = "STRING"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
	QUOTE_OBJ        = "QUOTE"
	REFERENCE_OBJ    = "REFERENCE"
	MACRO_OBJ        = "MACRO"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type BuiltinFunction func(args ...Object) Object

/*---------------------------------------------------------------------------*/

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

/*---------------------------------------------------------------------------*/

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

/*---------------------------------------------------------------------------*/

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

/*---------------------------------------------------------------------------*/

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

/*---------------------------------------------------------------------------*/

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

/*---------------------------------------------------------------------------*/

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

/*---------------------------------------------------------------------------*/

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

/*---------------------------------------------------------------------------*/

type Class struct {
	Name *ast.Identifier
	Body *ast.BlockStatement
}

func (c *Class) Type() ObjectType { return CLASS_OBJ }
func (c *Class) Inspect() string {
	var out bytes.Buffer

	out.WriteString("class")
	out.WriteString(c.Name.String())
	out.WriteString("{")
	out.WriteString(c.Body.String())
	out.WriteString("}")

	return out.String()
}

/*---------------------------------------------------------------------------*/

type Instance struct {
	Class *Class
	This  *Environment
}

func (i *Instance) Type() ObjectType { return INSTANCE_OBJ }
func (i *Instance) Inspect() string {
	var out bytes.Buffer

	out.WriteString("instanfe of ")
	out.WriteString(i.Class.Name.String())

	return out.String()
}

/*---------------------------------------------------------------------------*/

type This struct {
	Instance *Instance
	Name     string
}

func (t *This) Type() ObjectType { return THIS_OBJ }
func (t *This) Inspect() string {
	return "this is " + t.Instance.Inspect()
}

/*---------------------------------------------------------------------------*/

type Reference struct {
	Env  *Environment
	Name string
}

func (r *Reference) Type() ObjectType { return REFERENCE_OBJ }
func (r *Reference) Inspect() string {
	obj, ok := r.Env.Get(r.Name)
	if ok {
		return obj.Inspect()
	}
	return "<missing reference>"
}

func (r *Reference) Assign(obj Object) Object {
	return r.Env.Set(r.Name, obj)
}

func (r *Reference) Value() Object {
	obj, ok := r.Env.Get(r.Name)
	if ok {
		return obj
	}
	return &Null{}
}

/*---------------------------------------------------------------------------*/

type Macro struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (m *Macro) Type() ObjectType { return MACRO_OBJ }
func (m *Macro) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range m.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(m.Body.String())
	out.WriteString("\n}")

	return out.String()
}

/*---------------------------------------------------------------------------*/

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

/*---------------------------------------------------------------------------*/

type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

/*---------------------------------------------------------------------------*/

// [todo] - Hashkey() の戻り値をキャッシュしておくことで性能と最適化できるかもしれない

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	// 潜在的にコリジョンの可能性を含むが、チェイン法やオープンアドレス法といった方法で回避可能
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

/*---------------------------------------------------------------------------*/

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }

func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type Hashable interface {
	HashKey() HashKey
}

/*---------------------------------------------------------------------------*/

type Quote struct {
	Node ast.Node
}

func (q *Quote) Type() ObjectType { return QUOTE_OBJ }
func (q *Quote) Inspect() string {
	return "QUOTE(" + q.Node.String() + ")"
}

/*---------------------------------------------------------------------------*/
