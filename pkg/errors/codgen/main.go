package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

var (
	typeNames = flag.String("type", "", "comma-separated list of type names; must be set")
	output    = flag.String("output", "", "output file name; default srcdir/<type>_string.go")
)

var tpl = `
// Code generated by codegen. DO NOT EDIT.

package code

func init() {
{{ range .ErrCodes }}
	register({{ .Name }}, {{ .HttpCode }}, "{{ .Desc }}")
{{ end }}
}
`

type errGenerate struct {
	ErrCodes []errCode
}

type errCode struct {
	Name     string
	HttpCode string
	Desc     string
}

// Parses command line flags and processes each specified type.
func main() {
	flag.Parse()
	if len(*typeNames) == 0 {
		fmt.Fprintf(os.Stderr, "codegen: -type flag must be set\n")
		os.Exit(2)
	}

	// 解析类型名
	types := strings.Split(*typeNames, ",")

	// 处理每个类型
	for _, typeName := range types {
		// 为每个类型设置不同的输出文件名
		outputFile := fmt.Sprintf("%s_generated.go", strings.ToLower(typeName))
		generate(typeName, outputFile)
	}
}

// generate processes a single type and generates its error code file.
// It reads the source file, collects error codes, and writes the generated code.
func generate(typeName string, outputFile string) {
	// 构建输入文件路径
	inputFile := fmt.Sprintf("pkg\\errors\\code\\%s.go", typeName)

	// 解析文件
	f := ReadFile(inputFile)
	if f == nil {
		return
	}

	// 生成代码
	buf := new(bytes.Buffer)
	allErrs := make([]errCode, 0)

	for _, decl := range f.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			errs, err := processGenDecl(genDecl)
			if err != nil {
				panic(err)
			}
			allErrs = append(allErrs, errs...)
		}
	}

	// 执行模板
	tmpl, err := template.New("err").Parse(strings.TrimSpace(tpl))
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, errGenerate{ErrCodes: allErrs}); err != nil {
		panic(err)
	}

	// 写入文件
	outputPath := filepath.Join(filepath.Dir(inputFile), outputFile)
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		panic(err)
	}
	fmt.Printf("Generated %s\n", outputPath)
}

// ReadFile parses a Go source file and returns its AST.
// Returns nil if the file cannot be parsed.
func ReadFile(path string) *ast.File {
	f, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "codegen: %v\n", err)
		return nil
	}
	return f
}

// processGenDecl extracts error codes from a generic declaration.
// It handles both doc comments and line comments.
func processGenDecl(decl *ast.GenDecl) ([]errCode, error) {
	errs := make([]errCode, 0)
	seen := make(map[string]bool)

	for _, spec := range decl.Specs {
		v := spec.(*ast.ValueSpec)
		for _, name := range v.Names {
			if seen[name.Name] {
				continue
			}
			seen[name.Name] = true

			var comment string
			if v.Doc != nil && v.Doc.Text() != "" {
				comment = v.Doc.Text()
			} else if c := v.Comment; c != nil && len(c.List) > 0 {
				comment = c.Text()
			}
			httpCode, desc := parseComment(comment)
			tmp := errCode{Name: name.Name, HttpCode: httpCode, Desc: desc}
			errs = append(errs, tmp)
		}
	}

	return errs, nil
}

// parseComment extracts HTTP status code and description from a comment.
// Returns default values (500, "Internal server error") if the format is invalid.
func parseComment(comment string) (httpCode string, desc string) {
	reg := regexp.MustCompile(`\w\s*-\s*(\d{3})\s*:\s*([A-Z].*?)(?:\.|$)`)
	if !reg.MatchString(comment) {
		return "500", "Internal server error"
	}
	groups := reg.FindStringSubmatch(comment)
	if len(groups) != 3 {
		return "500", "Internal server error"
	}
	return groups[1], groups[2]
}
