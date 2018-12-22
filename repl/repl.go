package repl

import (
	"bufio"
	"fmt"
	"io"

	"bitbucket.org/kandayasu/squirrel-go/object"

	"bitbucket.org/kandayasu/squirrel-go/evaluator"
	"bitbucket.org/kandayasu/squirrel-go/lexer"
	"bitbucket.org/kandayasu/squirrel-go/parser"
)

// PROMPT : prompt character
const PROMPT = ">> "

// Start : start REPL
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluator.DefineMacros(program, macroEnv)
		expanded := evaluator.ExpandMacros(program, macroEnv)

		evaluated := evaluator.Eval(expanded, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

const AA = `
　　　　　　　　　,, -──- ､._
　　　　　　．-"´　　　　 　 　 ＼．
　　　　　：/ 　　　_ノ 　　　ヽ､_　ヽ.：
　　 　　：/ 　 　oﾟ(（●）) (（●）)ﾟoヽ：
　　　　：|　　　　 　　（__人__）　　　 |：
　　　　：l　　　　　　　 )　 (　　　 　 l：
　　　　： ､　　　　　　  'ー'　　 　 /：
　　　　　：, -‐ (_).　　　　　　　 ／
　　　　　：ｌ_ｊ_ｊ_ｊ と)丶─‐┬．''´
　　　　　　　　：ヽ　　　：i　|：
　 　 　 　 　 　：/　　：⊂ノ|：
`

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, AA)
	io.WriteString(out, "Woops! we ran int some monkey buisness here\n")
	io.WriteString(out, " parser erros:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"˜\n")
	}
}
