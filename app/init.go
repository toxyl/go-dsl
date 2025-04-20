package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"
)

type initTemplateParam struct {
	Index int
	Name  string
	Type  string
	Unit  string
	Desc  string
	Min   any
	Max   any
	Def   any
}

type initTemplateFunc struct {
	OrgName string
	Name    string
	Desc    string
	Params  []initTemplateParam
	Returns []initTemplateParam
}

type initTemplateVar struct {
	OrgName string
	Name    string
	Type    string
	Unit    string
	Desc    string
	Min     any
	Max     any
	Def     any
}

type initTemplate struct {
	Package      string
	Imports      []string
	ID           string
	Name         string
	Description  string
	Version      string
	Extension    string
	VarRegistry  []initTemplateVar
	FuncRegistry []initTemplateFunc
}

func (data *initTemplate) generateVarRegistrations(variables []metaVar) (requiredImports []string) {
	fnCheckNil := func(v any) string {
		if v == nil {
			return "nil"
		}
		return fmt.Sprintf("%v", v)
	}
	data.VarRegistry = []initTemplateVar{}
	for _, v := range variables {
		switch v.typ {
		case "*image.RGBA", "*image.NRGBA", "*image.RGBA64", "*image.NRGBA64":
			requiredImports = append(requiredImports, "image")
		case "color.RGBA", "color.RGBA64", "color.NRGBA", "color.NRGBA64":
			requiredImports = append(requiredImports, "image/color")
		}
		data.VarRegistry = append(data.VarRegistry, initTemplateVar{
			OrgName: v.orgName,
			Name:    v.name,
			Type:    v.typ,
			Unit:    fnCheckNil(v.unit),
			Desc:    v.desc,
			Min:     fnCheckNil(v.min),
			Max:     fnCheckNil(v.max),
			Def:     fnCheckNil(v.def),
		})
	}
	return
}

func (data *initTemplate) generateFuncRegistrations(functions []metaFunc) (requiredImports []string) {
	data.FuncRegistry = []initTemplateFunc{}
	for _, fn := range functions {
		tmplData := initTemplateFunc{
			OrgName: fn.orgName,
			Name:    fn.name,
			Desc:    fn.desc,
			Params:  []initTemplateParam{},
			Returns: []initTemplateParam{},
		}
		for i, param := range fn.params {
			switch param.typ {
			case "*image.RGBA", "*image.NRGBA", "*image.RGBA64", "*image.NRGBA64":
				requiredImports = append(requiredImports, "image")
			case "color.RGBA", "color.RGBA64", "color.NRGBA", "color.NRGBA64":
				requiredImports = append(requiredImports, "image/color")
			}

			tmplData.Params = append(tmplData.Params, initTemplateParam{
				Index: i,
				Name:  param.name,
				Type:  param.typ,
				Unit:  param.unit,
				Desc:  param.desc,
				Min:   param.min,
				Max:   param.max,
				Def:   param.def,
			})
		}
		for i, ret := range fn.returns {
			switch ret.typ {
			case "*image.RGBA", "*image.NRGBA", "*image.RGBA64", "*image.NRGBA64":
				requiredImports = append(requiredImports, "image")
			case "color.RGBA", "color.RGBA64", "color.NRGBA", "color.NRGBA64":
				requiredImports = append(requiredImports, "image/color")
			}

			tmplData.Returns = append(tmplData.Returns, initTemplateParam{
				Index: i,
				Name:  ret.name,
				Type:  ret.typ,
				Unit:  ret.unit,
				Desc:  ret.desc,
				Min:   ret.min,
				Max:   ret.max,
				Def:   ret.def,
			})
		}
		data.FuncRegistry = append(data.FuncRegistry, tmplData)
	}
	return
}

//go:embed init.tmpl
var tmplInit string

func getUniqueStrings(slice []string) []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, item := range slice {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func genInitCode(id string, name string, desc string, version string, extension string, pkg string, fns []metaFunc, vars []metaVar) string {
	tmpl, err := template.New("init").Parse(tmplInit)
	if err != nil {
		panic(err)
	}

	data := initTemplate{
		Package:     pkg,
		Imports:     []string{},
		ID:          id,
		Name:        name,
		Description: desc,
		Version:     version,
		Extension:   extension,
	}
	imports := []string{}
	imports = append(imports, data.generateVarRegistrations(vars)...)
	imports = append(imports, data.generateFuncRegistrations(fns)...)
	data.Imports = getUniqueStrings(imports)

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		panic(err)
	}

	return buf.String()
}
