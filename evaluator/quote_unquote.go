package evaluator

import (
	"fmt"

	"github.com/CHIKUWAODEN/monkey-for-c95/ast"
	"github.com/CHIKUWAODEN/monkey-for-c95/object"
	"github.com/CHIKUWAODEN/monkey-for-c95/token"
)

func quote(node ast.Node, env *object.Environment) object.Object {
	node = evalUnquoteCalls(node, env)
	return &object.Quote{Node: node}
}

func evalUnquoteCalls(quoted ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(quoted, func(node ast.Node) ast.Node {
		if !isUnquotedCall(node) {
			return node
		}

		call, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		if len(call.Arguments) != 1 {
			return node
		}

		// 呼び出し
		unquoted := Eval(call.Arguments[0], env)
		return convertObjectToASTNode(unquoted)
	})
}

func isUnquotedCall(node ast.Node) bool {
	CallExpression, ok := node.(*ast.CallExpression)
	if !ok {
		return false
	}
	return CallExpression.Function.TokenLiteral() == "unquote"
}

func convertObjectToASTNode(obj object.Object) ast.Node {

	switch obj := obj.(type) {

	case *object.Integer:
		t := token.Token{
			Type:    token.INT,
			Literal: fmt.Sprintf("%d", obj.Value),
		}
		return &ast.IntegerLiteral{Token: t, Value: obj.Value}

	case *object.Boolean:
		var t token.Token
		if obj.Value {
			t = token.Token{Type: token.TRUE, Literal: "true"}
		} else {
			t = token.Token{Type: token.FALSE, Literal: "false"}
		}
		return &ast.Boolean{Token: t, Value: obj.Value}

	// Reference が来た場合は Value() の値でフォールバックしてやればよい
	// Reference が指示するものが Integer や Boolean である場合、これまでどおり振る舞う
	case *object.Reference:
		return convertObjectToASTNode(obj.Value())

	case *object.Quote:
		return obj.Node

	// 処理できないオブジェクト、たとえば現状だと String などはここに来るようになっている
	default:
		return nil
	}
}
