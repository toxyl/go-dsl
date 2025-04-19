package main

import (
	"go/ast"
	"log"
	"strconv"
	"strings"
)

type metaVar struct {
	orgName string
	name    string
	typ     string
	unit    string
	min     any
	max     any
	def     any
	desc    string
}

type metaFunc struct {
	orgName string
	name    string
	desc    string
	params  []metaParam
	returns []metaParam
}

type metaParam struct {
	name string
	typ  string
	min  any
	max  any
	def  any
	unit string
	desc string
}

func extractFunctionMeta(node *ast.File, functions []metaFunc) []metaFunc {
	for _, decl := range node.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if fn.Doc == nil {
				continue
			}
			types := map[string]string{}
			if fn.Type != nil && fn.Type.Params != nil && fn.Type.Params.List != nil {
				for _, param := range fn.Type.Params.List {
					for _, name := range param.Names {
						switch param.Type.(type) {
						case *ast.StarExpr:
							types[name.Name] = "*" + param.Type.(*ast.StarExpr).X.(*ast.SelectorExpr).X.(*ast.Ident).Name + "." + param.Type.(*ast.StarExpr).X.(*ast.SelectorExpr).Sel.Name
						case *ast.Ident:
							types[name.Name] = param.Type.(*ast.Ident).Name
						}
					}
				}
			}
			if fn.Type != nil && fn.Type.Results != nil && fn.Type.Results.List != nil {
				for _, result := range fn.Type.Results.List {
					for _, name := range result.Names {
						types[name.Name] = result.Type.(*ast.Ident).Name
					}
				}
			}

			meta := metaFunc{
				orgName: fn.Name.Name,
			}

			// parse doc comments
			for _, comment := range fn.Doc.List {
				text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))

				if strings.HasPrefix(text, "@") {
					parts := strings.SplitN(text[1:], ":", 2)
					if len(parts) < 2 {
						continue
					}

					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					switch key {
					case "Name":
						meta.name = value
					case "Desc":
						meta.desc = value
					case "Param":
						meta.params = append(meta.params, parseParam(value, types))
					case "Returns":
						meta.returns = append(meta.returns, parseParam(value, types))
					}
				}
			}

			if meta.name != "" {
				functions = append(functions, meta)
			}
		}
	}
	return functions
}

func extractVariableMeta(node *ast.File, variables []metaVar) []metaVar {
	var err error
	for _, decl := range node.Decls {
		if vspec, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range vspec.Specs {
				if vspec, ok := spec.(*ast.ValueSpec); ok {
					if vspec.Doc != nil {
						meta := metaVar{
							orgName: vspec.Names[0].Name,
						}
						switch vspec.Values[0].(type) {
						case *ast.BasicLit:
							meta.def = parseValue(vspec.Values[0].(*ast.BasicLit).Value)
							meta.typ = strings.ToLower(vspec.Values[0].(*ast.BasicLit).Kind.String())
						case *ast.Ident:
							meta.typ = "bool"
							meta.def, err = strconv.ParseBool(vspec.Values[0].(*ast.Ident).Name)
							if err != nil {
								log.Fatal(err)
							}
						}

						switch meta.def.(type) {
						case int:
							meta.typ = "int"
						case int8:
							meta.typ = "int8"
						case int16:
							meta.typ = "int16"
						case int32:
							meta.typ = "int32"
						case int64:
							meta.typ = "int64"
						case uint:
							meta.typ = "uint"
						case uint8:
							meta.typ = "uint8"
						case uint16:
							meta.typ = "uint16"
						case uint32:
							meta.typ = "uint32"
						case uint64:
							meta.typ = "uint64"
						case float32:
							meta.typ = "float32"
						case float64:
							meta.typ = "float64"
						case bool:
							meta.typ = "bool"
						case string:
							meta.typ = "string"
						case []byte:
							meta.typ = "[]byte"
						default:
							meta.typ = "any"
						}

						for _, comment := range vspec.Doc.List {
							text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
							if strings.HasPrefix(text, "@") {
								parts := strings.SplitN(text[1:], ":", 2)
								if len(parts) < 2 {
									continue
								}
								key := strings.TrimSpace(parts[0])
								value := strings.TrimSpace(parts[1])
								switch key {
								case "Name":
									meta.name = value
								case "Desc":
									meta.desc = value
								case "Range":
									if meta.typ == "bool" || meta.typ == "string" {
										continue
									}
									el := strings.Split(value, "..")
									if len(el) == 2 {
										meta.min = parseValue(el[0])
										meta.max = parseValue(el[1])
									}
								case "Unit":
									if meta.typ == "bool" || meta.typ == "string" || value == "" {
										continue
									}
									meta.unit = value
								}
							}
						}
						if meta.name != "" {
							variables = append(variables, meta)
						}
					}
				}
			}
		}
	}
	return variables
}
