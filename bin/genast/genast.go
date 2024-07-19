package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

type AST struct {
	// upper level type name
	TypeName string
	// lower cased type name for filename and method names
	LowerTypeName string
	// types that implement the upper level type
	Types []Types
	// any imports needed for the generated file
	Imports []string
	// visitor returns a type, or just has side effects
	VisitorHasType bool
}

type Field struct {
	Name string
	Type string
}

type Types struct {
	Name   string
	Fields []Field
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <output dir>\n", os.Args[0])
		os.Exit(64)
	}

	exprTypes := []Types{
		{
			"Assign",
			[]Field{
				{"Name", "*token.Token"},
				{"Value", "Expr"},
			},
		},
		{
			"Binary",
			[]Field{
				{"Left", "Expr"},
				{"Operator", "*token.Token"},
				{"Right", "Expr"},
			},
		},
		{
			"Call",
			[]Field{
				{"Callee", "Expr"},
				{"Paren", "*token.Token"},
				{"Arguments", "[]Expr"},
			},
		},
		{
			"Grouping",
			[]Field{
				{"Expression", "Expr"},
			},
		},
		{
			"Literal",
			[]Field{
				{"Value", "any"},
			},
		},
		{
			"Logical",
			[]Field{
				{"Left", "Expr"},
				{"Operator", "*token.Token"},
				{"Right", "Expr"},
			},
		},
		{
			"Unary",
			[]Field{
				{"Operator", "*token.Token"},
				{"Right", "Expr"},
			},
		},
		{
			"ExprVar",
			[]Field{
				{"Name", "*token.Token"},
			},
		},
	}

	exprAst := AST{
		TypeName:       "Expr",
		LowerTypeName:  "expr",
		Types:          exprTypes,
		Imports:        []string{"github.com/mkeesey/craftinginterpreters/pkg/token"},
		VisitorHasType: true,
	}
	defineAst(os.Args[1], exprAst)

	stmtTypes := []Types{
		{
			"Block",
			[]Field{
				{"Statements", "[]Stmt"},
			},
		},
		{
			"Class",
			[]Field{
				{"Name", "*token.Token"},
				{"Methods", "[]*Function"},
			},
		},
		{
			"Expression",
			[]Field{
				{"Expression", "Expr"},
			},
		},
		{
			"Function",
			[]Field{
				{"Name", "*token.Token"},
				{"Params", "[]*token.Token"},
				{"Body", "[]Stmt"},
			},
		},
		{
			"If",
			[]Field{
				{"Condition", "Expr"},
				{"ThenBranch", "Stmt"},
				{"ElseBranch", "Stmt"},
			},
		},
		{
			"Print",
			[]Field{
				{"Expression", "Expr"},
			},
		},
		{
			"Return",
			[]Field{
				{"Keyword", "*token.Token"},
				{"Value", "Expr"},
			},
		},
		{
			"StmtVar",
			[]Field{
				{"Name", "*token.Token"},
				{"Initializer", "Expr"},
			},
		},
		{
			"While",
			[]Field{
				{"Condition", "Expr"},
				{"Body", "Stmt"},
			},
		},
	}

	stmtAst := AST{
		TypeName:       "Stmt",
		LowerTypeName:  "stmt",
		Types:          stmtTypes,
		Imports:        []string{"github.com/mkeesey/craftinginterpreters/pkg/token"},
		VisitorHasType: false, // only side effects
	}
	defineAst(os.Args[1], stmtAst)
}

func defineAst(outputDir string, ast AST) {
	path := filepath.Join(outputDir, ast.LowerTypeName+".go")
	file, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create file: %s\n", path)
		os.Exit(1)
	}
	defer file.Close()

	tmpl := template.Must(template.New("ast").Parse(tmplBody))
	err = tmpl.Execute(file, ast)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to file: %s\n%s", path, err)
		os.Exit(1)
	}
}

var (
	tmplBody = `// Code generated by genast.go; DO NOT EDIT.
package ast

import (
	"fmt"
{{- if .Imports }}
{{ range .Imports }}
	"{{ . }}"
{{- end }}
{{- end }}
)

type {{ .TypeName }}Visitor{{ if .VisitorHasType }}[T any]{{end}} interface {
{{- range .Types }}
	Visit{{ .Name }}(*{{ .Name }}){{ if $.VisitorHasType }} T{{end}}
{{- end }}
}

func Visit{{ .TypeName }}{{ if .VisitorHasType }}[T any]{{end}}({{ .LowerTypeName }} {{ .TypeName }}, visitor {{ .TypeName }}
{{- if $.VisitorHasType -}}
Visitor[T]) T {
{{- else -}}
Visitor) {
{{- end}}
	switch n := {{ .LowerTypeName }}.(type) {
{{- range .Types }}
	case *{{ .Name }}:
{{- if $.VisitorHasType }}
		return visitor.Visit{{ .Name }}(n)
{{- else }}
		visitor.Visit{{ .Name }}(n)
{{- end}}
{{- end }}
	default:
		panic(fmt.Sprintf("Unknown {{ .TypeName }} type %T", {{ .LowerTypeName }}))
	}
}

type {{ .TypeName }} interface {
	{{ .LowerTypeName }}()
}

{{ range .Types -}}
type {{ .Name }} struct {
{{- range .Fields }}
	{{ .Name }} {{ .Type }}
{{- end }}
}

func (b *{{ .Name }}) {{ $.LowerTypeName }}() {}

{{ end }}
`
)
