package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"path"
	"strings"

	"github.com/matt-hoiland/glox/internal/constants/exit"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: generate-ast <output director>")
		os.Exit(exit.Usage)
	}
	outputDir := os.Args[1]
	defineAST(outputDir, "Expr",
		"Assign   : Name *token.Token, Value Expr",
		"Binary   : Left Expr, Operator *token.Token, Right Expr",
		"Grouping : Expression Expr",
		"Literal  : Value loxtype.Type",
		"Unary    : Operator *token.Token, Right Expr",
		"Variable : Name *token.Token",
	)
	defineAST(outputDir, "Stmt",
		"Block      : Statements []Stmt",
		"Expression : Expression Expr",
		"Print      : Expression Expr",
		"Var        : Name *token.Token, Initializer Expr",
	)
}

func defineAST(outputDir, baseName string, productions ...string) {
	w := &bytes.Buffer{}

	fmt.Fprintf(w, "// Code generated by tools/generate-ast. DO NOT EDIT.\n")
	fmt.Fprintln(w, "package ast")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "import (")
	fmt.Fprintln(w, `	"github.com/matt-hoiland/glox/internal/environment"`)
	fmt.Fprintln(w, `	"github.com/matt-hoiland/glox/internal/loxtype"`)
	fmt.Fprintln(w, `	"github.com/matt-hoiland/glox/internal/token"`)
	fmt.Fprintln(w, ")")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "type "+baseName+" interface {")
	fmt.Fprintln(w, "\tAccept(*environment.Environment, "+baseName+"Visitor) (loxtype.Type, error)")
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
	defineVisitor(w, baseName, productions...)
	fmt.Fprintln(w)
	for _, production := range productions {
		typeName := strings.TrimSpace(strings.Split(production, ":")[0])
		fields := strings.TrimSpace(strings.Split(production, ":")[1])
		defineType(w, baseName, typeName, fields)
	}

	data, err := format.Source(w.Bytes())
	if err != nil {
		panic(err)
	}

	path := path.Join(outputDir, strings.ToLower(baseName)+".go")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		panic(err)
	}
	if err := os.WriteFile(path, data, 0755); err != nil {
		panic(err)
	}
}

func defineType(w io.Writer, baseName, typeName, fieldList string) {
	fmt.Fprintf(w, "type %s%s struct {\n", typeName, baseName)
	for field := range strings.SplitSeq(fieldList, ",") {
		fmt.Fprintf(w, "\t%s\n", strings.TrimSpace(field))
	}
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "var _ %s = (*%s%s)(nil)\n", baseName, typeName, baseName)
	fmt.Fprintln(w)
	fmt.Fprintf(w, "func New%s%s(%s) *%s%s {\n", typeName, baseName, fieldList, typeName, baseName)
	fmt.Fprintf(w, "\treturn &%s%s{\n", typeName, baseName)
	for fieldPair := range strings.SplitSeq(fieldList, ",") {
		field := strings.TrimSpace(strings.Split(strings.TrimSpace(fieldPair), " ")[0])
		fmt.Fprintf(w, "\t%s: %s,\n", field, field)
	}
	fmt.Fprintln(w, "\t}")
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "func (e *%s%s) Accept(env *environment.Environment, visitor %sVisitor) (loxtype.Type, error) {\n", typeName, baseName, baseName)
	fmt.Fprintf(w, "\treturn visitor.Visit%s%s(env, e)\n", typeName, baseName)
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
}

func defineVisitor(w io.Writer, baseName string, productions ...string) {
	fmt.Fprintf(w, "type %sVisitor interface {\n", baseName)
	for _, production := range productions {
		typeName := strings.TrimSpace(strings.Split(production, ":")[0])
		fmt.Fprintf(w, "\tVisit%s%s(*environment.Environment, *%s%s) (loxtype.Type, error)\n", typeName, baseName, typeName, baseName)
	}
	fmt.Fprintln(w, `}`)
}
