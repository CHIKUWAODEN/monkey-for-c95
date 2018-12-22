package parser

import (
	"fmt"
	"os"
	"strings"
)

var traceLevel int = 0

const traceIndentPlaceholder string = ". "

func identLevel() string {
	return strings.Repeat(traceIndentPlaceholder, traceLevel-1)
}

func tracePrint(fs string) {
	fmt.Printf("%s%s\n", identLevel(), fs)
}

func incIndent() { traceLevel++ }
func decIndent() { traceLevel-- }

func trace(msg string) string {
	incIndent()
	if _, ok := os.LookupEnv("TRACE"); ok {
		tracePrint("BEGIN " + msg)
	}
	return msg
}

func untrace(msg string) {
	if _, ok := os.LookupEnv("TRACE"); ok {
		tracePrint("END " + msg)
	}
	decIndent()
}
