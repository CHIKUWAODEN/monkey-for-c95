package token

/// TokenType represent token type
type TokenType string

/// Token is struct that contain a token
type Token struct {
	Type    TokenType
	Literal string
}

// token
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// identifier + Literal
	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	// operator
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	EQ    = "=="
	NOTEQ = "!="

	// delimiter
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// keyword
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"

	CLASS = "CLASS"
	THIS  = "THIS"
	DOT   = "."

	// macro
	MACRO = "MACRO"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"macro":  MACRO,
	"class":  CLASS,
	"this":   THIS,
}

// LookupIdent : check ident is keyword or identifier
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
