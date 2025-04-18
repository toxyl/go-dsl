package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/toxyl/math"

	"github.com/toxyl/flo"
)

func createTestLanguage() {
	pos := 0
	isEnabled := false
	list := make([]any, 21)

	theme := dsl.defaultColorTheme()
	theme.EditorBackground = "#000000"
	theme.EditorForeground = "#AFAFAF"
	theme.VariableAssignments = "#D77C97"
	dsl.initDSL("test-script", "Test Script", "Testing", "0.0.0", "test", theme)
	dsl.vars.register(
		"pos", "int", "index", "The position of something in a list",
		0, 10, 20,
		func() any { return pos },
		func(a any) { pos = a.(int) },
	)
	dsl.vars.register(
		"on", "bool", "", "Whether or not the feature is enabled",
		nil, nil, true,
		func() any { return isEnabled },
		func(a any) { isEnabled = a.(bool) },
	)
	dsl.vars.register(
		"item", "any", "", "The item at the current position in the list",
		nil, nil, true,
		func() any { return list[pos] },
		func(a any) { list[pos] = a },
	)
	dsl.funcs.register(
		"add",
		"Adds two numbers together",
		[]dslParamMeta{
			{name: "x", typ: "int", def: 0, unit: "", desc: "The first number to add"},
			{name: "y", typ: "int", def: 0, unit: "", desc: "The second number to add"},
		},
		[]dslParamMeta{
			{name: "result", typ: "int", def: 0, unit: "", desc: "The sum of the two numbers"},
		},
		func(args ...any) (any, error) {
			return args[0].(int) + args[1].(int), nil
		},
	)
	dsl.funcs.register(
		"mul",
		"Multiplies two numbers together",
		[]dslParamMeta{
			{name: "x", typ: "int", def: 0, unit: "", desc: "The first number to multiply"},
			{name: "y", typ: "int", def: 0, unit: "", desc: "The second number to multiply"},
		},
		[]dslParamMeta{
			{name: "result", typ: "int", def: 0, unit: "", desc: "The product of the two numbers"},
		},
		func(args ...any) (any, error) {
			return args[0].(int) * args[1].(int), nil
		},
	)
	dsl.funcs.register(
		"concat",
		"Concatenates two strings together",
		[]dslParamMeta{
			{name: "x", typ: "string", def: "", unit: "", desc: "The first string to concatenate"},
			{name: "y", typ: "string", def: "", unit: "", desc: "The second string to concatenate"},
		},
		[]dslParamMeta{
			{name: "result", typ: "string", def: "", unit: "", desc: "The concatenated string"},
		},
		func(args ...any) (any, error) {
			values := make([]string, len(args))
			for i, arg := range args {
				values[i] = fmt.Sprintf("%v", arg)
			}
			return strings.Join(values, ""), nil
		},
	)
	dsl.funcs.register(
		"test-function-1", "This is a test function",
		[]dslParamMeta{
			{name: "x", typ: "int", min: 0, max: 10, def: 0, unit: "px", desc: "Position on the x axis"},
			{name: "y", typ: "int", min: 0, max: 10, def: 0, unit: "px", desc: "Position on the y axis"},
			{name: "str", typ: "string", def: "hi", desc: "String to print"},
		},
		[]dslParamMeta{
			{name: "z", typ: "int", min: 0, max: 20, def: 0, unit: "px", desc: "Position on the z axis"},
		},
		func(a ...any) (any, error) {
			x := a[0].(int)
			y := a[1].(int)
			str := a[2].(string)
			z := x + y
			fmt.Println(str, x, y, z)
			return z, nil
		},
	)
	dsl.funcs.register(
		"test-function-2", "This is a test function",
		[]dslParamMeta{
			{name: "lat", typ: "float64", min: -90, max: 90, def: 0, unit: "°", desc: "Latitude"},
			{name: "lon", typ: "float64", min: -180, max: 180, def: 0, unit: "°", desc: "Longitude"},
		},
		[]dslParamMeta{
			{name: "z", typ: "bool", def: false, desc: "Is the point in the ocean?"},
			{name: "err", typ: "error", desc: "Something bad happened"},
		},
		func(a ...any) (any, error) {
			lat, lon := a[0].(float64), a[1].(float64)
			z := lat+lon > 0
			return z, nil
		},
	)
	dsl.funcs.register("img-nrgba64", "This is a function to process an NRGBA64 image",
		[]dslParamMeta{{name: "img", typ: "*image.NRGBA64", desc: "The image to process"}},
		[]dslParamMeta{{name: "res", typ: "*image.NRGBA64", def: false, desc: "The image with a pixel colored red"}},
		func(a ...any) (any, error) {
			img := a[0].(*image.NRGBA64)
			file, err := os.Create("../../LANGUAGE-NRGBA64.png")
			if err != nil {
				return nil, err
			}
			defer file.Close()
			png.Encode(file, img)
			return img, nil
		},
	)
	dsl.funcs.register("img-rgba64", "This is another function to process an RGBA64 image",
		[]dslParamMeta{{name: "img", typ: "*image.RGBA64", desc: "The image to process"}},
		[]dslParamMeta{{name: "res", typ: "*image.RGBA64", def: false, desc: "The image with a pixel colored red"}},
		func(a ...any) (any, error) {
			img := a[0].(*image.RGBA64)
			file, err := os.Create("../../LANGUAGE-RGBA64.png")
			if err != nil {
				return nil, err
			}
			defer file.Close()
			png.Encode(file, img)
			return img, nil
		},
	)
	dsl.funcs.register("img-nrgba", "This is a function to process an NRGBA image",
		[]dslParamMeta{{name: "img", typ: "*image.NRGBA64", desc: "The image to process"}},
		[]dslParamMeta{{name: "res", typ: "*image.NRGBA", def: false, desc: "The image with a pixel colored red"}},
		func(a ...any) (any, error) {
			img := a[0].(*image.NRGBA64)
			file, err := os.Create("../../LANGUAGE-NRGBA64-NRGBA.png")
			if err != nil {
				return nil, err
			}
			defer file.Close()
			png.Encode(file, img)
			return img, nil
		},
	)
	dsl.funcs.register("img-rgba", "This is another function to process an RGBA image",
		[]dslParamMeta{{name: "img", typ: "*image.RGBA64", desc: "The image to process"}},
		[]dslParamMeta{{name: "res", typ: "*image.RGBA", def: false, desc: "The image with a pixel colored red"}},
		func(a ...any) (any, error) {
			img := a[0].(*image.RGBA64)
			file, err := os.Create("../../LANGUAGE-RGBA64-RGBA.png")
			if err != nil {
				return nil, err
			}
			defer file.Close()
			png.Encode(file, img)
			return img, nil
		},
	)
	dsl.storeState()
}

func testResult(t *testing.T, name string, want any, wantErr bool, got any, err error) {
	if (err != nil) != wantErr {
		t.Errorf("\x1b[31m[FAIL] %s: error = %v, wantErr %v\x1b[0m", name, err, wantErr)
		return
	}
	if !wantErr {
		if !dsl.deepEqual(got, want) {
			t.Errorf("\x1b[31m[FAIL] %s: got %v (%T), want %v (%T)\x1b[0m", name, got, got, want, want)
			return
		}
	}
	t.Logf("\x1b[32m[PASS] %s\x1b[0m", name)
}

func TestBasicExpressions(t *testing.T) {
	t.Run("Basic Expressions", func(t *testing.T) {
		type TestCase struct {
			name    string
			script  string
			args    []any
			want    *dslResult
			wantErr bool
		}

		c := func(name string, script string, args []any, want *dslResult, wantErr bool) TestCase {
			return TestCase{name, script, args, want, wantErr}
		}

		tests := []TestCase{
			c("just an int", `42`, []any{}, &dslResult{int64(42), nil}, false),
			c("just a float", `42.1`, []any{}, &dslResult{42.1, nil}, false),
			c("just a string", `"hello"`, []any{}, &dslResult{"hello", nil}, false),
			c("just a bool", `true`, []any{}, &dslResult{true, nil}, false),
			c("basic argument usage", `$1`, []any{42}, &dslResult{42, nil}, false),
			c("multiple arguments", `add($1 $2)`, []any{5, 3}, &dslResult{8, nil}, false),
			c("argument in variable assignment", `x: $1 y: $2 add(x y)`, []any{10, 20}, &dslResult{30, nil}, false),
			c("argument out of range", `$3`, []any{1, 2}, &dslResult{nil, fmt.Errorf("argument $3 out of range")}, true),
			c("mixed types", `concat($1 $2)`, []any{"hello", 42}, &dslResult{"hello42", nil}, false),
			c("nested argument usage", `add($1 mul($2 $3))`, []any{1, 2, 3}, &dslResult{7, nil}, false),
		}

		createTestLanguage()
		for _, tt := range tests {
			dsl.restoreState()
			t.Run(tt.name, func(t *testing.T) {
				got, err := dsl.run(tt.script, false, tt.args...)
				testResult(t, tt.name, tt.want, tt.wantErr, got, err)
			})
		}
	})
}

func TestTokenizer(t *testing.T) {
	t.Run("Tokenizer", func(t *testing.T) {
		type (
			fields struct {
				source string
			}
			TestCase struct {
				name    string
				fields  fields
				want    []*dslToken
				wantErr bool
			}
		)
		c := func(name, source string, wantErr bool, want ...*dslToken) TestCase {
			return TestCase{
				name: name,
				fields: fields{
					source: source,
				},
				want:    want,
				wantErr: wantErr,
			}
		}
		tkn := dsl.newToken

		tests := []TestCase{
			c("example 1", `l:+(test-function-1(1 2 "hi, this will be printed") *(50 2)) +(l +(1 gx))`, false,
				tkn(`l:`, dsl.tokens.assign), tkn(`+(`, dsl.tokens.callStart), tkn(`test-function-1(`, dsl.tokens.callStart), tkn(`1`, dsl.tokens.integer), tkn(`2`, dsl.tokens.integer), tkn(`hi, this will be printed`, dsl.tokens.str), tkn(`)`, dsl.tokens.callEnd), tkn(`*(`, dsl.tokens.callStart), tkn(`50`, dsl.tokens.integer), tkn(`2`, dsl.tokens.integer), tkn(`)`, dsl.tokens.callEnd), tkn(`)`, dsl.tokens.callEnd), tkn(`;`, dsl.tokens.terminator),
				tkn(`+(`, dsl.tokens.callStart), tkn(`l`, dsl.tokens.varRef), tkn(`+(`, dsl.tokens.callStart), tkn(`1`, dsl.tokens.integer), tkn(`gx`, dsl.tokens.varRef), tkn(`)`, dsl.tokens.callEnd), tkn(`)`, dsl.tokens.callEnd),
			),
			c("example 2", `func(); other(x= 4.5 do = false) x: c(5 "hello world" 1 # I can comment inline using \# to escape the hash sign #) `+"\n"+`y: "you can also escape \" in strings"; yetAnother("with a string and a newline`+"\n"+`this time" y)`, false,
				tkn(`func(`, dsl.tokens.callStart), tkn(`)`, dsl.tokens.callEnd), tkn(`;`, dsl.tokens.terminator),
				tkn(`other(`, dsl.tokens.callStart), tkn(`x=`, dsl.tokens.namedArg), tkn(`4.5`, dsl.tokens.float), tkn(`do=`, dsl.tokens.namedArg), tkn(`false`, dsl.tokens.boolean), tkn(`)`, dsl.tokens.callEnd), tkn(`;`, dsl.tokens.terminator),
				tkn(`x:`, dsl.tokens.assign), tkn(`c(`, dsl.tokens.callStart), tkn(`5`, dsl.tokens.integer), tkn(`hello world`, dsl.tokens.str), tkn(`1`, dsl.tokens.integer), tkn(`I can comment inline using # to escape the hash sign`, dsl.tokens.comment), tkn(`)`, dsl.tokens.callEnd), tkn(`;`, dsl.tokens.terminator),
				tkn(`y:`, dsl.tokens.assign), tkn(`you can also escape " in strings`, dsl.tokens.str), tkn(`;`, dsl.tokens.terminator),
				tkn(`yetAnother(`, dsl.tokens.callStart), tkn(`with a string and a newline`+"\n"+`this time`, dsl.tokens.str), tkn(`y`, dsl.tokens.varRef), tkn(`)`, dsl.tokens.callEnd),
			),
			c("example 3", `a: func1(x=123 y=50) b: func2(x=123 y =50) c: func3(x= 123 y=50) d: func4(x = 123 y=50) e: func5(x=123 y= 50) f: func6(x=123 y =50) g: func7(x=func1(1 2) y=func2(3 4))`, false,
				tkn(`a:`, dsl.tokens.assign), tkn(`func1(`, dsl.tokens.callStart), tkn(`x=`, dsl.tokens.namedArg), tkn(`123`, dsl.tokens.integer), tkn(`y=`, dsl.tokens.namedArg), tkn(`50`, dsl.tokens.integer), tkn(`)`, dsl.tokens.callEnd), tkn(`;`, dsl.tokens.terminator),
				tkn(`b:`, dsl.tokens.assign), tkn(`func2(`, dsl.tokens.callStart), tkn(`x=`, dsl.tokens.namedArg), tkn(`123`, dsl.tokens.integer), tkn(`y=`, dsl.tokens.namedArg), tkn(`50`, dsl.tokens.integer), tkn(`)`, dsl.tokens.callEnd), tkn(`;`, dsl.tokens.terminator),
				tkn(`c:`, dsl.tokens.assign), tkn(`func3(`, dsl.tokens.callStart), tkn(`x=`, dsl.tokens.namedArg), tkn(`123`, dsl.tokens.integer), tkn(`y=`, dsl.tokens.namedArg), tkn(`50`, dsl.tokens.integer), tkn(`)`, dsl.tokens.callEnd), tkn(`;`, dsl.tokens.terminator),
				tkn(`d:`, dsl.tokens.assign), tkn(`func4(`, dsl.tokens.callStart), tkn(`x=`, dsl.tokens.namedArg), tkn(`123`, dsl.tokens.integer), tkn(`y=`, dsl.tokens.namedArg), tkn(`50`, dsl.tokens.integer), tkn(`)`, dsl.tokens.callEnd), tkn(`;`, dsl.tokens.terminator),
				tkn(`e:`, dsl.tokens.assign), tkn(`func5(`, dsl.tokens.callStart), tkn(`x=`, dsl.tokens.namedArg), tkn(`123`, dsl.tokens.integer), tkn(`y=`, dsl.tokens.namedArg), tkn(`50`, dsl.tokens.integer), tkn(`)`, dsl.tokens.callEnd), tkn(`;`, dsl.tokens.terminator),
				tkn(`f:`, dsl.tokens.assign), tkn(`func6(`, dsl.tokens.callStart), tkn(`x=`, dsl.tokens.namedArg), tkn(`123`, dsl.tokens.integer), tkn(`y=`, dsl.tokens.namedArg), tkn(`50`, dsl.tokens.integer), tkn(`)`, dsl.tokens.callEnd), tkn(`;`, dsl.tokens.terminator),
				tkn(`g:`, dsl.tokens.assign), tkn(`func7(`, dsl.tokens.callStart), tkn(`x=`, dsl.tokens.namedArg), tkn(`func1(`, dsl.tokens.callStart), tkn(`1`, dsl.tokens.integer), tkn(`2`, dsl.tokens.integer), tkn(`)`, dsl.tokens.callEnd), tkn(`y=`, dsl.tokens.namedArg), tkn(`func2(`, dsl.tokens.callStart), tkn(`3`, dsl.tokens.integer), tkn(`4`, dsl.tokens.integer), tkn(`)`, dsl.tokens.callEnd), tkn(`)`, dsl.tokens.callEnd),
			),
			c("example 4", `this is not a valid statement I believe(`, true),
			c("example 5", `musthaveclosingbracketkn(`, true),
			c("example 6", `mustassignsomething:`, true),
			c("example 7", `func1(1 2 3`, true),
			c("simple function call 1", `func1(1 2 3)`, false,
				tkn(`func1(`, dsl.tokens.callStart), tkn(`1`, dsl.tokens.integer), tkn(`2`, dsl.tokens.integer), tkn(`3`, dsl.tokens.integer), tkn(`)`, dsl.tokens.callEnd),
			),
			c("simple function call 2", `test-function-1(1 2 "hello \" mean\"world!\"")`, false,
				tkn(`test-function-1(`, dsl.tokens.callStart), tkn(`1`, dsl.tokens.integer), tkn(`2`, dsl.tokens.integer), tkn(`hello " mean"world!"`, dsl.tokens.str), tkn(`)`, dsl.tokens.callEnd),
			),
			c("missing closing parenthesis", `test-function-1(1 2`, true),
			c("empty imput", ``, false),
			c("just an int", `42`, false, tkn("42", dsl.tokens.integer)),
			c("just a float", `42.0`, false, tkn("42.0", dsl.tokens.float)),
			c("just a string", `"hi"`, false, tkn("hi", dsl.tokens.str)),
			c("just a bool", `true`, false, tkn("true", dsl.tokens.boolean)),
			c("empty var name", `: hi`, true),
			c("named arguments", `test-function-1(x=1 y=2 str="hello")`, false,
				tkn("test-function-1(", dsl.tokens.callStart), tkn("x=", dsl.tokens.namedArg), tkn("1", dsl.tokens.integer), tkn("y=", dsl.tokens.namedArg), tkn("2", dsl.tokens.integer), tkn("str=", dsl.tokens.namedArg), tkn("hello", dsl.tokens.str), tkn(")", dsl.tokens.callEnd),
			),
			c("argument references 1", `func($1 $2) x: $3 y: $4`, false,
				tkn(`func(`, dsl.tokens.callStart), tkn(`$1`, dsl.tokens.argRef), tkn(`$2`, dsl.tokens.argRef), tkn(`)`, dsl.tokens.callEnd), tkn(`;`, dsl.tokens.terminator),
				tkn(`x:`, dsl.tokens.assign), tkn(`$3`, dsl.tokens.argRef), tkn(`;`, dsl.tokens.terminator),
				tkn(`y:`, dsl.tokens.assign), tkn(`$4`, dsl.tokens.argRef),
			),
			c("argument references 2", `x: $3 y: $4 func($1 $2)`, false,
				tkn(`x:`, dsl.tokens.assign), tkn(`$3`, dsl.tokens.argRef), tkn(`;`, dsl.tokens.terminator),
				tkn(`y:`, dsl.tokens.assign), tkn(`$4`, dsl.tokens.argRef), tkn(`;`, dsl.tokens.terminator),
				tkn(`func(`, dsl.tokens.callStart), tkn(`$1`, dsl.tokens.argRef), tkn(`$2`, dsl.tokens.argRef), tkn(`)`, dsl.tokens.callEnd),
			),
			c("invalid argument reference", `func($)`, true),
			c("nested function calls", `test-function-2(test-function-1(1 10 "Hello \" World") 0)`, false,
				tkn("test-function-2(", dsl.tokens.callStart), tkn("test-function-1(", dsl.tokens.callStart), tkn("1", dsl.tokens.integer), tkn("10", dsl.tokens.integer), tkn(`Hello " World`, dsl.tokens.str), tkn(")", dsl.tokens.callEnd), tkn("0", dsl.tokens.integer), tkn(")", dsl.tokens.callEnd),
			),
			c("single var assign", `a: 100`, false,
				tkn("a:", dsl.tokens.assign), tkn("100", dsl.tokens.integer),
			),
			c("multiple func calls", `add(1 2) sub(5 3) mul(4 2) div(10 2) pow(2 3) sqrtkn(16) sin(0) cos(0) tan(0) log(100)`, false,
				tkn(`add(`, dsl.tokens.callStart), tkn(`1`, dsl.tokens.integer), tkn(`2`, dsl.tokens.integer), tkn(`)`, dsl.tokens.callEnd), tkn(`;`, dsl.tokens.terminator),
				tkn(`sub(`, dsl.tokens.callStart), tkn(`5`, dsl.tokens.integer), tkn(`3`, dsl.tokens.integer), tkn(`)`, dsl.tokens.callEnd), tkn(`;`, dsl.tokens.terminator),
				tkn(`mul(`, dsl.tokens.callStart), tkn(`4`, dsl.tokens.integer), tkn(`2`, dsl.tokens.integer), tkn(`)`, dsl.tokens.callEnd), tkn(`;`, dsl.tokens.terminator),
				tkn(`div(`, dsl.tokens.callStart), tkn(`10`, dsl.tokens.integer), tkn(`2`, dsl.tokens.integer), tkn(`)`, dsl.tokens.callEnd), tkn(`;`, dsl.tokens.terminator),
				tkn(`pow(`, dsl.tokens.callStart), tkn(`2`, dsl.tokens.integer), tkn(`3`, dsl.tokens.integer), tkn(`)`, dsl.tokens.callEnd), tkn(`;`, dsl.tokens.terminator),
				tkn(`sqrtkn(`, dsl.tokens.callStart), tkn(`16`, dsl.tokens.integer), tkn(`)`, dsl.tokens.callEnd), tkn(`;`, dsl.tokens.terminator),
				tkn(`sin(`, dsl.tokens.callStart), tkn(`0`, dsl.tokens.integer), tkn(`)`, dsl.tokens.callEnd), tkn(`;`, dsl.tokens.terminator),
				tkn(`cos(`, dsl.tokens.callStart), tkn(`0`, dsl.tokens.integer), tkn(`)`, dsl.tokens.callEnd), tkn(`;`, dsl.tokens.terminator),
				tkn(`tan(`, dsl.tokens.callStart), tkn(`0`, dsl.tokens.integer), tkn(`)`, dsl.tokens.callEnd), tkn(`;`, dsl.tokens.terminator),
				tkn(`log(`, dsl.tokens.callStart), tkn(`100`, dsl.tokens.integer), tkn(`)`, dsl.tokens.callEnd),
			),
			c("multiple func calls with args (complex)", `x:10 y:20 z:add(x y) result:mul(z pow(2 3)) sqrt(result) sin(cos(tan(0))) log(pow(10 2))`, false,
				tkn(`x:`, dsl.tokens.assign), tkn(`10`, dsl.tokens.integer), tkn(`;`, dsl.tokens.terminator),
				tkn(`y:`, dsl.tokens.assign), tkn(`20`, dsl.tokens.integer), tkn(`;`, dsl.tokens.terminator),
				tkn(`z:`, dsl.tokens.assign), tkn(`add(`, dsl.tokens.callStart), tkn(`x`, dsl.tokens.varRef), tkn(`y`, dsl.tokens.varRef), tkn(`)`, dsl.tokens.callEnd), tkn(`;`, dsl.tokens.terminator),
				tkn(`result:`, dsl.tokens.assign), tkn(`mul(`, dsl.tokens.callStart), tkn(`z`, dsl.tokens.varRef), tkn(`pow(`, dsl.tokens.callStart), tkn(`2`, dsl.tokens.integer), tkn(`3`, dsl.tokens.integer), tkn(`)`, dsl.tokens.callEnd), tkn(`)`, dsl.tokens.callEnd), tkn(`;`, dsl.tokens.terminator),
				tkn(`sqrt(`, dsl.tokens.callStart), tkn(`result`, dsl.tokens.varRef), tkn(`)`, dsl.tokens.callEnd), tkn(`;`, dsl.tokens.terminator),
				tkn(`sin(`, dsl.tokens.callStart), tkn(`cos(`, dsl.tokens.callStart), tkn(`tan(`, dsl.tokens.callStart), tkn(`0`, dsl.tokens.integer), tkn(`)`, dsl.tokens.callEnd), tkn(`)`, dsl.tokens.callEnd), tkn(`)`, dsl.tokens.callEnd), tkn(`;`, dsl.tokens.terminator),
				tkn(`log(`, dsl.tokens.callStart), tkn(`pow(`, dsl.tokens.callStart), tkn(`10`, dsl.tokens.integer), tkn(`2`, dsl.tokens.integer), tkn(`)`, dsl.tokens.callEnd), tkn(`)`, dsl.tokens.callEnd),
			),
			c("assigns with func call", `x: 25 y: 25 add(x y)`, false,
				tkn(`x:`, dsl.tokens.assign), tkn(`25`, dsl.tokens.integer), tkn(`;`, dsl.tokens.terminator),
				tkn(`y:`, dsl.tokens.assign), tkn(`25`, dsl.tokens.integer), tkn(`;`, dsl.tokens.terminator),
				tkn(`add(`, dsl.tokens.callStart), tkn(`x`, dsl.tokens.varRef), tkn(`y`, dsl.tokens.varRef), tkn(`)`, dsl.tokens.callEnd),
			),
			c("named arg with func call", `add(x=1 y=subtract(1 2))`, false,
				tkn("add(", dsl.tokens.callStart), tkn("x=", dsl.tokens.namedArg), tkn("1", dsl.tokens.integer), tkn("y=", dsl.tokens.namedArg), tkn("subtract(", dsl.tokens.callStart), tkn("1", dsl.tokens.integer), tkn("2", dsl.tokens.integer), tkn(")", dsl.tokens.callEnd), tkn(")", dsl.tokens.callEnd),
			),
			c("named arg with empty string", `users(search="")`, false,
				tkn("users(", dsl.tokens.callStart), tkn("search=", dsl.tokens.namedArg), tkn("", dsl.tokens.str), tkn(")", dsl.tokens.callEnd),
			),
		}

		createTestLanguage()
		for _, tt := range tests {
			dsl.restoreState()
			t.Run(tt.name, func(t *testing.T) {
				dsl.load(tt.fields.source)

				if err := dsl.tokenizer.tokenize(); err != nil {
					testResult(t, tt.name, tt.want, tt.wantErr, nil, err)
				} else if err := dsl.tokenizer.lex(); err != nil {
					testResult(t, tt.name, tt.want, tt.wantErr, nil, err)
				} else {
					t.Logf("\x1b[33mSCRIPT: %s\x1b[0m", strings.ReplaceAll(dsl.tokenizer.String(), "\n", " "))
					gotTypes := dsl.tokenizer.getTypes()
					dsl.tokenizer.tokens = tt.want
					wantTypes := dsl.tokenizer.getTypes()
					testResult(t, tt.name, wantTypes, tt.wantErr, gotTypes, err)
				}
			})
		}
	})
}

func TestTypeConversions(t *testing.T) {
	maxInt8 := int8(127)
	minInt8 := int8(-128)
	maxInt16 := int16(32767)
	minInt16 := int16(-32768)
	maxUint8 := uint8(255)
	maxUint16 := uint16(65535)

	t.Run("String Conversions", func(t *testing.T) {
		tests := []struct {
			name       string
			value      any
			targetType string
			want       any
			wantErr    bool
		}{
			{name: "bool false to string", value: false, targetType: "string", want: "false"},
			{name: "bool true to string", value: true, targetType: "string", want: "true"},
			{name: "chan to string", value: make(chan int), targetType: "string", wantErr: true},
			{name: "empty string to int", value: "", targetType: "int", wantErr: true},
			{name: "float64 to string", value: 3.14, targetType: "string", want: "3.14"},
			{name: "int to string", value: 42, targetType: "string", want: "42"},
			{name: "int16 to string", value: int16(42), targetType: "string", want: "42"},
			{name: "int32 to string", value: int32(42), targetType: "string", want: "42"},
			{name: "int64 to string", value: int64(42), targetType: "string", want: "42"},
			{name: "int8 to string", value: int8(42), targetType: "string", want: "42"},
			{name: "map to string", value: make(map[string]int), targetType: "string", wantErr: true},
			{name: "slice to string", value: []int{1, 2, 3}, targetType: "string", wantErr: true},
			{name: "string 0 to bool", value: "0", targetType: "bool", want: false},
			{name: "string 1 to bool", value: "1", targetType: "bool", want: true},
			{name: "string decimal large to int", value: "777", targetType: "int", want: 777},
			{name: "string decimal point to uint", value: "123.45", targetType: "uint", want: uint(123)},
			{name: "string decimal to int", value: "1010", targetType: "int", want: 1010},
			{name: "string decimal to uint", value: "255", targetType: "uint", want: uint(255)},
			{name: "string empty to bool", value: "", targetType: "bool", wantErr: true},
			{name: "string empty to uint", value: "", targetType: "uint", wantErr: true},
			{name: "string f to bool", value: "f", targetType: "bool", want: false},
			{name: "string false to bool", value: "false", targetType: "bool", want: false},
			{name: "string FALSE to bool", value: "FALSE", targetType: "bool", want: false},
			{name: "string hex to int", value: "0xff", targetType: "int", wantErr: true},
			{name: "string hex to uint", value: "0xff", targetType: "uint", wantErr: true},
			{name: "string hex uppercase to uint", value: "0xFF", targetType: "uint", wantErr: true},
			{name: "string hex with 0x prefix to uint", value: "0x1234", targetType: "uint", wantErr: true},
			{name: "string invalid chars to uint", value: "123abc", targetType: "uint", wantErr: true},
			{name: "string invalid float to float64", value: "not_a_float", targetType: "float64", wantErr: true},
			{name: "string invalid number to float64", value: "abc123", targetType: "float64", wantErr: true},
			{name: "string invalid number to int", value: "abc123", targetType: "int", wantErr: true},
			{name: "string invalid number to uint", value: "abc123", targetType: "uint", wantErr: true},
			{name: "string invalid to bool", value: "invalid", targetType: "bool", wantErr: true},
			{name: "string invalid to int", value: "not a number", targetType: "int", wantErr: true},
			{name: "string negative to uint", value: "-123", targetType: "uint", want: uint(0)},
			{name: "string octal to uint", value: "0777", targetType: "uint", want: uint(777)},
			{name: "string pointer to string pointer", value: new(string), targetType: "*string", wantErr: true},
			{name: "string scientific notation to float64", value: "1.23e-4", targetType: "float64", want: 1.23e-4},
			{name: "string t to bool", value: "t", targetType: "bool", want: true},
			{name: "string to error interface", value: "error", targetType: "error", wantErr: true},
			{name: "string to float64", value: "3.14", targetType: "float64", want: 3.14},
			{name: "string to int", value: "42", targetType: "int", want: 42},
			{name: "string to invalid type", value: "hello", targetType: "invalid", wantErr: true},
			{name: "string to negative float64", value: "-3.14", targetType: "float64", want: -3.14},
			{name: "string to negative int", value: "-42", targetType: "int", want: -42},
			{name: "string to string", value: "42", targetType: "string", want: "42"},
			{name: "string true to bool", value: "true", targetType: "bool", want: true},
			{name: "string TRUE to bool", value: "TRUE", targetType: "bool", want: true},
			{name: "string whitespace to bool", value: "   ", targetType: "bool", wantErr: true},
			{name: "string whitespace to uint", value: "   ", targetType: "uint", wantErr: true},
			{name: "string with invalid chars to int", value: "42abc", targetType: "int", wantErr: true},
			{name: "string with leading spaces to int", value: "   42", targetType: "int", want: 42},
			{name: "string with only spaces to int", value: "   ", targetType: "int", wantErr: true},
			{name: "string with plus to float64", value: "+123.45", targetType: "float64", want: 123.45},
			{name: "string with plus to int", value: "+123", targetType: "int", want: 123},
			{name: "string with scientific notation to int", value: "1.23e2", targetType: "int", want: 123},
			{name: "string with sign", value: "-123", targetType: "int", want: -123},
			{name: "string with spaces to int", value: "  42  ", targetType: "int", want: 42},
			{name: "string with trailing dot to float64", value: "123.", targetType: "float64", want: 123.0},
			{name: "string with trailing spaces to int", value: "42   ", targetType: "int", want: 42},
			{name: "struct to string", value: struct{}{}, targetType: "string", wantErr: true},
			{name: "uint to string", value: uint(42), targetType: "string", want: "42"},
			{name: "uint16 to string", value: uint16(42), targetType: "string", want: "42"},
			{name: "uint32 to string", value: uint32(42), targetType: "string", want: "42"},
			{name: "uint64 to string", value: uint64(42), targetType: "string", want: "42"},
			{name: "uint8 to string", value: uint8(42), targetType: "string", want: "42"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := dsl.cast(tt.value, tt.targetType)
				testResult(t, tt.name, tt.want, tt.wantErr, got, err)
			})
		}
	})

	t.Run("Bool Conversions", func(t *testing.T) {
		tests := []struct {
			name       string
			value      any
			targetType string
			want       any
			wantErr    bool
		}{
			{name: "bool false to float64", value: false, targetType: "float64", want: 0.0},
			{name: "bool false to int", value: false, targetType: "int", want: 0},
			{name: "bool false to uint", value: false, targetType: "uint", want: uint(0)},
			{name: "bool to bool false", value: false, targetType: "bool", want: false},
			{name: "bool to bool type false", value: false, targetType: "bool", want: false},
			{name: "bool to bool type true", value: true, targetType: "bool", want: true},
			{name: "bool to bool", value: true, targetType: "bool", want: true},
			{name: "bool to invalid type", value: true, targetType: "invalid", wantErr: true},
			{name: "bool true to float64", value: true, targetType: "float64", want: 1.0},
			{name: "bool true to int", value: true, targetType: "int", want: 1},
			{name: "bool true to uint", value: true, targetType: "uint", want: uint(1)},
			{name: "float64 to bool", value: 3.14, targetType: "bool", want: true},
			{name: "int to bool", value: 42, targetType: "bool", want: true},
			{name: "int16 to bool", value: int16(42), targetType: "bool", want: true},
			{name: "int32 to bool", value: int32(42), targetType: "bool", want: true},
			{name: "int64 to bool", value: int64(42), targetType: "bool", want: true},
			{name: "int8 to bool", value: int8(42), targetType: "bool", want: true},
			{name: "uint to bool", value: uint(42), targetType: "bool", want: true},
			{name: "uint16 to bool", value: uint16(42), targetType: "bool", want: true},
			{name: "uint32 to bool", value: uint32(42), targetType: "bool", want: true},
			{name: "uint64 to bool", value: uint64(42), targetType: "bool", want: true},
			{name: "uint8 to bool", value: uint8(42), targetType: "bool", want: true},
			{name: "unsupported type to bool", value: struct{}{}, targetType: "bool", wantErr: true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := dsl.cast(tt.value, tt.targetType)
				testResult(t, tt.name, tt.want, tt.wantErr, got, err)
			})
		}
	})

	t.Run("Float64 Conversions", func(t *testing.T) {
		tests := []struct {
			name       string
			value      any
			targetType string
			want       any
			wantErr    bool
		}{
			{name: "float32 inf to float64", value: float32(math.Inf(1)), targetType: "float64", want: math.Inf(1)},
			{name: "float32 nan to float64", value: math.NaN[float32](), targetType: "float64", want: math.NaN[float64]()},
			{name: "float32 neg inf to float64", value: float32(math.Inf(-1)), targetType: "float64", want: math.Inf(-1)},
			{name: "float32 to float64", value: float32(3.14), targetType: "float64", want: float64(float32(3.14))},
			{name: "float64 infinity to float32", value: math.Inf(1), targetType: "float32", want: float32(math.Inf(1))},
			{name: "float64 infinity to int", value: math.Inf(1), targetType: "int", want: math.MaxInt64},
			{name: "float64 infinity to int16", value: math.Inf(1), targetType: "int16", want: int16(math.MaxInt16)},
			{name: "float64 infinity to int32", value: math.Inf(1), targetType: "int32", want: int32(math.MaxInt32)},
			{name: "float64 infinity to int64", value: math.Inf(1), targetType: "int64", want: int64(math.MaxInt64)},
			{name: "float64 infinity to int8", value: math.Inf(1), targetType: "int8", want: int8(math.MaxInt8)},
			{name: "float64 infinity to uint", value: math.Inf(1), targetType: "uint", want: uint(math.MaxUint64)},
			{name: "float64 infinity to uint16", value: math.Inf(1), targetType: "uint16", want: uint16(math.MaxUint16)},
			{name: "float64 infinity to uint32", value: math.Inf(1), targetType: "uint32", want: uint32(math.MaxUint32)},
			{name: "float64 infinity to uint64", value: math.Inf(1), targetType: "uint64", want: uint64(math.MaxUint64)},
			{name: "float64 infinity to uint8", value: math.Inf(1), targetType: "uint8", want: uint8(math.MaxUint8)},
			{name: "float64 large to uint", value: float64(9223372036854775808), targetType: "uint", want: uint(9223372036854775808)},
			{name: "float64 max to float32", value: math.MaxFloat64, targetType: "float32", want: float32(math.MaxFloat32)},
			{name: "float64 max to int", value: math.MaxFloat64, targetType: "int", want: int(math.MaxInt)},
			{name: "float64 max to int16", value: math.MaxFloat64, targetType: "int16", want: int16(math.MaxInt16)},
			{name: "float64 max to int32", value: math.MaxFloat64, targetType: "int32", want: int32(math.MaxInt32)},
			{name: "float64 max to int64", value: math.MaxFloat64, targetType: "int64", want: int64(math.MaxInt64)},
			{name: "float64 max to int8", value: math.MaxFloat64, targetType: "int8", want: int8(math.MaxInt8)},
			{name: "float64 maxfloat32 to float32", value: float64(math.MaxFloat32), targetType: "float32", want: float32(math.MaxFloat32)},
			{name: "float64 min to int", value: -math.MaxFloat64, targetType: "int", want: int(math.MinInt)},
			{name: "float64 min to int16", value: -math.MaxFloat64, targetType: "int16", want: int16(math.MinInt16)},
			{name: "float64 min to int32", value: -math.MaxFloat64, targetType: "int32", want: int32(math.MinInt32)},
			{name: "float64 min to int64", value: -math.MaxFloat64, targetType: "int64", want: int64(math.MinInt64)},
			{name: "float64 min to int8", value: -math.MaxFloat64, targetType: "int8", want: int8(math.MinInt8)},
			{name: "float64 NaN to float32", value: math.NaN[float64](), targetType: "float32", want: math.NaN[float32]()},
			{name: "float64 NaN to int", value: math.NaN[float64](), targetType: "int", want: int(0)},
			{name: "float64 NaN to int16", value: math.NaN[float64](), targetType: "int16", want: int16(0)},
			{name: "float64 NaN to int32", value: math.NaN[float64](), targetType: "int32", want: int32(0)},
			{name: "float64 NaN to int64", value: math.NaN[float64](), targetType: "int64", want: int64(0)},
			{name: "float64 NaN to int8", value: math.NaN[float64](), targetType: "int8", want: int8(0)},
			{name: "float64 nan to uint", value: math.NaN[float64](), targetType: "uint", want: uint(0)},
			{name: "float64 NaN to uint", value: math.NaN[float64](), targetType: "uint", want: uint(0)},
			{name: "float64 NaN to uint16", value: math.NaN[float64](), targetType: "uint16", want: uint16(0)},
			{name: "float64 NaN to uint32", value: math.NaN[float64](), targetType: "uint32", want: uint32(0)},
			{name: "float64 NaN to uint64", value: math.NaN[float64](), targetType: "uint64", want: uint64(0)},
			{name: "float64 NaN to uint8", value: math.NaN[float64](), targetType: "uint8", want: uint8(0)},
			{name: "float64 negative infinity to float32", value: math.Inf(-1), targetType: "float32", want: float32(math.Inf(-1))},
			{name: "float64 negative infinity to int", value: math.Inf(-1), targetType: "int", want: math.MinInt64},
			{name: "float64 negative infinity to int16", value: math.Inf(-1), targetType: "int16", want: int16(math.MinInt16)},
			{name: "float64 negative infinity to int32", value: math.Inf(-1), targetType: "int32", want: int32(math.MinInt32)},
			{name: "float64 negative infinity to int64", value: math.Inf(-1), targetType: "int64", want: int64(math.MinInt64)},
			{name: "float64 negative infinity to int8", value: math.Inf(-1), targetType: "int8", want: int8(math.MinInt8)},
			{name: "float64 negative infinity to uint", value: math.Inf(-1), targetType: "uint", want: uint(0)},
			{name: "float64 negative infinity to uint16", value: math.Inf(-1), targetType: "uint16", want: uint16(0)},
			{name: "float64 negative infinity to uint32", value: math.Inf(-1), targetType: "uint32", want: uint32(0)},
			{name: "float64 negative infinity to uint64", value: math.Inf(-1), targetType: "uint64", want: uint64(0)},
			{name: "float64 negative infinity to uint8", value: math.Inf(-1), targetType: "uint8", want: uint8(0)},
			{name: "float64 negative max to float32", value: -math.MaxFloat64, targetType: "float32", want: float32(-math.MaxFloat32)},
			{name: "float64 negative maxfloat32 to float32", value: float64(-math.MaxFloat32), targetType: "float32", want: float32(-math.MaxFloat32)},
			{name: "float64 negative to uint", value: float64(-42), targetType: "uint", want: uint(0)},
			{name: "float64 overflow", value: math.MaxFloat64 + 1, targetType: "float64", want: math.MaxFloat64},
			{name: "float64 pointer to float64 pointer", value: new(float64), targetType: "*float64", wantErr: true},
			{name: "float64 slightly above maxfloat32 to float32", value: float64(math.MaxFloat32) * 1.1, targetType: "float32", want: float32(math.MaxFloat32)},
			{name: "float64 slightly below negative maxfloat32 to float32", value: float64(-math.MaxFloat32) * 1.1, targetType: "float32", want: float32(-math.MaxFloat32)},
			{name: "float64 smallest to float32", value: math.SmallestNonzeroFloat64, targetType: "float32", want: float32(0)},
			{name: "float64 to float32", value: 3.14, targetType: "float32", want: float32(3.14)},
			{name: "float64 to float64", value: 3.14, targetType: "float64", want: 3.14},
			{name: "float64 to int", value: 3.14, targetType: "int", want: 3},
			{name: "float64 to int", value: 42.7, targetType: "int", want: 42},
			{name: "float64 to int16", value: 3.14, targetType: "int16", want: int16(3)},
			{name: "float64 to int32", value: 3.14, targetType: "int32", want: int32(3)},
			{name: "float64 to int64", value: 3.14, targetType: "int64", want: int64(3)},
			{name: "float64 to int8", value: 3.14, targetType: "int8", want: int8(3)},
			{name: "float64 to uint", value: 3.14, targetType: "uint", want: uint(3)},
			{name: "float64 to uint16", value: 3.14, targetType: "uint16", want: uint16(3)},
			{name: "float64 to uint32", value: 3.14, targetType: "uint32", want: uint32(3)},
			{name: "float64 to uint64", value: 3.14, targetType: "uint64", want: uint64(3)},
			{name: "float64 to uint8", value: 3.14, targetType: "uint8", want: uint8(3)},
			{name: "float64 underflow", value: -math.MaxFloat64 - 1, targetType: "float64", want: -math.MaxFloat64},
			{name: "float64 very large to float32", value: math.MaxFloat64, targetType: "float32", want: float32(math.MaxFloat32)},
			{name: "float64 very small to float32", value: math.SmallestNonzeroFloat64, targetType: "float32", want: float32(0)},
			{name: "float64 zero to uint", value: float64(0), targetType: "uint", want: uint(0)},
			{name: "int to float64", value: 42, targetType: "float64", want: 42.0},
			{name: "int to float64", value: 42, targetType: "float64", want: float64(42)},
			{name: "int16 to float64", value: int16(42), targetType: "float64", want: float64(42)},
			{name: "int32 to float64", value: int32(42), targetType: "float64", want: float64(42)},
			{name: "int64 to float64", value: int64(42), targetType: "float64", want: float64(42)},
			{name: "int8 to float64", value: int8(42), targetType: "float64", want: float64(42)},
			{name: "NaN to float64", value: math.NaN[float64](), targetType: "float64", want: math.NaN[float64]()},
			{name: "uint to float64", value: uint(42), targetType: "float64", want: float64(42)},
			{name: "uint16 to float64", value: uint16(42), targetType: "float64", want: float64(42)},
			{name: "uint32 to float64", value: uint32(42), targetType: "float64", want: float64(42)},
			{name: "uint64 to float64", value: uint64(42), targetType: "float64", want: float64(42)},
			{name: "uint8 to float64", value: uint8(42), targetType: "float64", want: float64(42)},
			{name: "very small float to float64", value: math.SmallestNonzeroFloat64, targetType: "float64", want: math.SmallestNonzeroFloat64},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := dsl.cast(tt.value, tt.targetType)
				testResult(t, tt.name, tt.want, tt.wantErr, got, err)
			})
		}
	})

	t.Run("Float32 Conversions", func(t *testing.T) {
		tests := []struct {
			name       string
			value      any
			targetType string
			want       any
			wantErr    bool
		}{
			{name: "float32 infinity to uint", value: float32(math.Inf(1)), targetType: "uint", want: uint(math.MaxUint64)},
			{name: "float32 nan to uint", value: math.NaN[float32](), targetType: "uint", want: uint(0)},
			{name: "float32 negative infinity to uint", value: float32(math.Inf(-1)), targetType: "uint", want: uint(0)},
			{name: "float32 negative to int", value: float32(-42.9), targetType: "int", want: -42},
			{name: "float32 negative to uint", value: float32(-42.9), targetType: "uint", want: uint(0)},
			{name: "float32 to float32", value: float32(3.14), targetType: "float32", want: float32(3.14)},
			{name: "float32 to int", value: float32(42.9), targetType: "int", want: 42},
			{name: "float32 to uint", value: float32(42.9), targetType: "uint", want: uint(42)},
			{name: "int to float32", value: 42, targetType: "float32", want: float32(42)},
			{name: "int16 to float32", value: int16(42), targetType: "float32", want: float32(42)},
			{name: "int32 to float32", value: int32(42), targetType: "float32", want: float32(42)},
			{name: "int64 to float32", value: int64(42), targetType: "float32", want: float32(42)},
			{name: "int8 to float32", value: int8(42), targetType: "float32", want: float32(42)},
			{name: "NaN to float32", value: math.NaN[float64](), targetType: "float32", want: math.NaN[float32]()},
			{name: "uint to float32", value: uint(42), targetType: "float32", want: float32(42)},
			{name: "uint16 to float32", value: uint16(42), targetType: "float32", want: float32(42)},
			{name: "uint32 to float32", value: uint32(42), targetType: "float32", want: float32(42)},
			{name: "uint64 to float32", value: uint64(42), targetType: "float32", want: float32(42)},
			{name: "uint8 to float32", value: uint8(42), targetType: "float32", want: float32(42)},
			{name: "very small float to float32", value: math.SmallestNonzeroFloat64, targetType: "float32", want: float32(0)},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := dsl.cast(tt.value, tt.targetType)
				testResult(t, tt.name, tt.want, tt.wantErr, got, err)
			})
		}
	})

	t.Run("Int Conversions", func(t *testing.T) {
		tests := []struct {
			name       string
			value      any
			targetType string
			want       any
			wantErr    bool
		}{
			{name: "-Inf to int", value: math.Inf(-1), targetType: "int", want: math.MinInt64},
			{name: "+Inf to int", value: math.Inf(1), targetType: "int", want: math.MaxInt64},
			{name: "array to array", value: [3]int{1, 2, 3}, targetType: "[3]int", wantErr: true},
			{name: "array to slice", value: [3]int{1, 2, 3}, targetType: "[]int", wantErr: true},
			{name: "channel type", value: make(chan int), targetType: "int", wantErr: true},
			{name: "empty interface", value: any(nil), targetType: "int", wantErr: true},
			{name: "float with decimal to int", value: 42.7, targetType: "int", want: 42},
			{name: "float with exactly .5 to int", value: 42.5, targetType: "int", want: 42},
			{name: "float with large decimal to int", value: 42.999999, targetType: "int", want: 42},
			{name: "float with small decimal to int", value: 42.000001, targetType: "int", want: 42},
			{name: "int pointer to int pointer", value: new(int), targetType: "*int", wantErr: true},
			{name: "int to error interface", value: 42, targetType: "error", wantErr: true},
			{name: "int to int", value: 42, targetType: "int", want: 42},
			{name: "int to int16", value: 42, targetType: "int16", want: int16(42)},
			{name: "int to int32", value: 42, targetType: "int32", want: int32(42)},
			{name: "int to int64", value: 42, targetType: "int64", want: int64(42)},
			{name: "int to int8", value: 42, targetType: "int8", want: int8(42)},
			{name: "int to invalid type", value: 42, targetType: "invalid", wantErr: true},
			{name: "int to uint", value: 42, targetType: "uint", want: uint(42)},
			{name: "int to uint16", value: 42, targetType: "uint16", want: uint16(42)},
			{name: "int to uint32", value: 42, targetType: "uint32", want: uint32(42)},
			{name: "int to uint64", value: 42, targetType: "uint64", want: uint64(42)},
			{name: "int to uint8", value: 42, targetType: "uint8", want: uint8(42)},
			{name: "int16 to int", value: int16(42), targetType: "int", want: int(42)},
			{name: "int32 to int", value: int32(42), targetType: "int", want: int(42)},
			{name: "int8 to int", value: int8(42), targetType: "int", want: int(42)},
			{name: "large int to uint8 overflow", value: 256, targetType: "uint8", want: uint8(math.MaxUint8)},
			{name: "NaN to int", value: math.NaN[float64](), targetType: "int", want: 0},
			{name: "negative int to uint", value: -1, targetType: "uint", want: uint(0)},
			{name: "negative int to uint16", value: -1, targetType: "uint16", want: uint16(0)},
			{name: "negative int to uint32", value: -1, targetType: "uint32", want: uint32(0)},
			{name: "negative int to uint64", value: -1, targetType: "uint64", want: uint64(0)},
			{name: "negative int to uint8", value: -1, targetType: "uint8", want: uint8(0)},
			{name: "nil value", value: nil, targetType: "int", wantErr: true},
			{name: "slice to array", value: []int{1, 2, 3}, targetType: "[3]int", wantErr: true},
			{name: "slice to slice", value: []int{1, 2, 3}, targetType: "[]int", wantErr: true},
			{name: "uint to int", value: uint(42), targetType: "int", want: int(42)},
			{name: "uint16 to int", value: uint16(42), targetType: "int", want: int(42)},
			{name: "uint32 to int", value: uint32(42), targetType: "int", want: int(42)},
			{name: "uint64 to int", value: uint64(42), targetType: "int", want: int(42)},
			{name: "uint8 to int", value: uint8(42), targetType: "int", want: int(42)},
			{name: "unsupported type", value: struct{}{}, targetType: "int", wantErr: true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := dsl.cast(tt.value, tt.targetType)
				testResult(t, tt.name, tt.want, tt.wantErr, got, err)
			})
		}
	})

	t.Run("Int8 Conversions", func(t *testing.T) {
		tests := []struct {
			name       string
			value      any
			targetType string
			want       any
			wantErr    bool
		}{
			{name: "int16 to int8", value: int16(42), targetType: "int8", want: int8(42)},
			{name: "int32 to int8", value: int32(42), targetType: "int8", want: int8(42)},
			{name: "int64 to int8", value: int64(42), targetType: "int8", want: int8(42)},
			{name: "int8 overflow", value: math.MaxInt16, targetType: "int8", want: int8(math.MaxInt8)},
			{name: "int8 to int16", value: int8(42), targetType: "int16", want: int16(42)},
			{name: "int8 to int32", value: int8(42), targetType: "int32", want: int32(42)},
			{name: "int8 to int64", value: int8(42), targetType: "int64", want: int64(42)},
			{name: "int8 to int8", value: int8(42), targetType: "int8", want: int8(42)},
			{name: "int8 to uint16", value: int8(42), targetType: "uint16", want: uint16(42)},
			{name: "int8 to uint32", value: int8(42), targetType: "uint32", want: uint32(42)},
			{name: "int8 to uint64", value: int8(42), targetType: "uint64", want: uint64(42)},
			{name: "int8 to uint8", value: int8(42), targetType: "uint8", want: uint8(42)},
			{name: "int8 underflow", value: math.MinInt16, targetType: "int8", want: int8(math.MinInt8)},
			{name: "max int8 to int16", value: maxInt8, targetType: "int16", want: int16(maxInt8)},
			{name: "min int8 to int16", value: minInt8, targetType: "int16", want: int16(minInt8)},
			{name: "uint16 to int8", value: uint16(42), targetType: "int8", want: int8(42)},
			{name: "uint32 to int8", value: uint32(42), targetType: "int8", want: int8(42)},
			{name: "uint64 to int8", value: uint64(42), targetType: "int8", want: int8(42)},
			{name: "uint8 to int8", value: uint8(42), targetType: "int8", want: int8(42)},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := dsl.cast(tt.value, tt.targetType)
				testResult(t, tt.name, tt.want, tt.wantErr, got, err)
			})
		}
	})

	t.Run("Int16 Conversions", func(t *testing.T) {
		tests := []struct {
			name       string
			value      any
			targetType string
			want       any
			wantErr    bool
		}{
			{name: "int16 overflow", value: math.MaxInt32, targetType: "int16", want: int16(math.MaxInt16)},
			{name: "int16 to int16", value: int16(42), targetType: "int16", want: int16(42)},
			{name: "int16 to int32", value: int16(42), targetType: "int32", want: int32(42)},
			{name: "int16 to int64", value: int16(42), targetType: "int64", want: int64(42)},
			{name: "int16 to uint16", value: int16(42), targetType: "uint16", want: uint16(42)},
			{name: "int16 to uint32", value: int16(42), targetType: "uint32", want: uint32(42)},
			{name: "int16 to uint64", value: int16(42), targetType: "uint64", want: uint64(42)},
			{name: "int16 to uint8", value: int16(42), targetType: "uint8", want: uint8(42)},
			{name: "int16 underflow", value: math.MinInt32, targetType: "int16", want: int16(math.MinInt16)},
			{name: "int32 to int16", value: int32(42), targetType: "int16", want: int16(42)},
			{name: "int64 to int16", value: int64(42), targetType: "int16", want: int16(42)},
			{name: "max int16 to int32", value: maxInt16, targetType: "int32", want: int32(maxInt16)},
			{name: "max uint8 to int16", value: maxUint8, targetType: "int16", want: int16(maxUint8)},
			{name: "min int16 to int32", value: minInt16, targetType: "int32", want: int32(minInt16)},
			{name: "uint16 to int16", value: uint16(42), targetType: "int16", want: int16(42)},
			{name: "uint32 to int16", value: uint32(42), targetType: "int16", want: int16(42)},
			{name: "uint64 large to int16", value: uint64(math.MaxInt16 + 1), targetType: "int16", want: int16(math.MaxInt16)},
			{name: "uint64 to int16", value: uint64(42), targetType: "int16", want: int16(42)},
			{name: "uint8 to int16", value: uint8(42), targetType: "int16", want: int16(42)},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := dsl.cast(tt.value, tt.targetType)
				testResult(t, tt.name, tt.want, tt.wantErr, got, err)
			})
		}
	})

	t.Run("Int32 Conversions", func(t *testing.T) {
		tests := []struct {
			name       string
			value      any
			targetType string
			want       any
			wantErr    bool
		}{
			{name: "int32 overflow", value: math.MaxInt64, targetType: "int32", want: int32(math.MaxInt32)},
			{name: "int32 to int32", value: int32(42), targetType: "int32", want: int32(42)},
			{name: "int32 to int64", value: int32(42), targetType: "int64", want: int64(42)},
			{name: "int32 to uint16", value: int32(42), targetType: "uint16", want: uint16(42)},
			{name: "int32 to uint32", value: int32(42), targetType: "uint32", want: uint32(42)},
			{name: "int32 to uint64", value: int32(42), targetType: "uint64", want: uint64(42)},
			{name: "int32 to uint8", value: int32(42), targetType: "uint8", want: uint8(42)},
			{name: "int32 underflow", value: math.MinInt64, targetType: "int32", want: int32(math.MinInt32)},
			{name: "int64 to int32", value: int64(42), targetType: "int32", want: int32(42)},
			{name: "max uint16 to int32", value: maxUint16, targetType: "int32", want: int32(maxUint16)},
			{name: "uint32 to int32", value: uint32(42), targetType: "int32", want: int32(42)},
			{name: "uint64 large to int32", value: uint64(math.MaxInt32 + 1), targetType: "int32", want: int32(math.MaxInt32)},
			{name: "uint8 to int32", value: uint8(42), targetType: "int32", want: int32(42)},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := dsl.cast(tt.value, tt.targetType)
				testResult(t, tt.name, tt.want, tt.wantErr, got, err)
			})
		}
	})

	t.Run("Int64 Conversions", func(t *testing.T) {
		tests := []struct {
			name       string
			value      any
			targetType string
			want       any
			wantErr    bool
		}{
			{name: "int64 to int64", value: int64(42), targetType: "int64", want: int64(42)},
			{name: "int64 to uint16", value: int64(42), targetType: "uint16", want: uint16(42)},
			{name: "int64 to uint32", value: int64(42), targetType: "uint32", want: uint32(42)},
			{name: "int64 to uint64", value: int64(42), targetType: "uint64", want: uint64(42)},
			{name: "int64 to uint8", value: int64(42), targetType: "uint8", want: uint8(42)},
			{name: "uint16 to int64", value: uint16(42), targetType: "int64", want: int64(42)},
			{name: "uint32 to int64", value: uint32(42), targetType: "int64", want: int64(42)},
			{name: "uint64 max to int64", value: uint64(math.MaxUint64), targetType: "int64", want: int64(math.MaxInt64)},
			{name: "uint64 maxint64 plus one to int64", value: uint64(math.MaxInt64 + 1), targetType: "int64", want: int64(math.MaxInt64)},
			{name: "uint64 maxint64 to int64", value: uint64(math.MaxInt64), targetType: "int64", want: int64(math.MaxInt64)},
			{name: "uint8 to int64", value: uint8(42), targetType: "int64", want: int64(42)},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := dsl.cast(tt.value, tt.targetType)
				testResult(t, tt.name, tt.want, tt.wantErr, got, err)
			})
		}
	})

	t.Run("Uint Conversions", func(t *testing.T) {
		tests := []struct {
			name       string
			value      any
			targetType string
			want       any
			wantErr    bool
		}{
			{name: "-Inf to uint", value: math.Inf(-1), targetType: "uint", want: uint(0)},
			{name: "+Inf to uint", value: math.Inf(1), targetType: "uint", want: uint(math.MaxUint64)},
			{name: "int16 to uint", value: int16(42), targetType: "uint", want: uint(42)},
			{name: "int32 to uint", value: int32(42), targetType: "uint", want: uint(42)},
			{name: "int64 to uint", value: int64(42), targetType: "uint", want: uint(42)},
			{name: "int8 to uint", value: int8(42), targetType: "uint", want: uint(42)},
			{name: "NaN to uint", value: math.NaN[float64](), targetType: "uint", want: uint(0)},
			{name: "negative to uint", value: -1, targetType: "uint", want: uint(0)},
			{name: "uint to int16", value: uint(42), targetType: "int16", want: int16(42)},
			{name: "uint to int32", value: uint(42), targetType: "int32", want: int32(42)},
			{name: "uint to int64", value: uint(42), targetType: "int64", want: int64(42)},
			{name: "uint to int8", value: uint(42), targetType: "int8", want: int8(42)},
			{name: "uint to uint", value: uint(42), targetType: "uint", want: uint(42)},
			{name: "uint to uint16", value: uint(42), targetType: "uint16", want: uint16(42)},
			{name: "uint to uint32", value: uint(42), targetType: "uint32", want: uint32(42)},
			{name: "uint to uint64", value: uint(42), targetType: "uint64", want: uint64(42)},
			{name: "uint to uint8", value: uint(42), targetType: "uint8", want: uint8(42)},
			{name: "uint underflow", value: -1, targetType: "uint", want: uint(0)},
			{name: "uint32 to uint", value: uint32(42), targetType: "uint", want: uint(42)},
			{name: "uint64 to uint", value: uint64(42), targetType: "uint", want: uint(42)},
			{name: "uint8 to uint", value: uint8(42), targetType: "uint", want: uint(42)},
			{name: "zero to uint", value: 0, targetType: "uint", want: uint(0)},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := dsl.cast(tt.value, tt.targetType)
				testResult(t, tt.name, tt.want, tt.wantErr, got, err)
			})
		}
	})

	t.Run("Uint8 Conversions", func(t *testing.T) {
		tests := []struct {
			name       string
			value      any
			targetType string
			want       any
			wantErr    bool
		}{
			{name: "negative to uint8", value: -1, targetType: "uint8", want: uint8(0)},
			{name: "uint32 to uint8", value: uint32(42), targetType: "uint8", want: uint8(42)},
			{name: "uint64 to uint8", value: uint64(42), targetType: "uint8", want: uint8(42)},
			{name: "uint8 overflow", value: math.MaxUint16, targetType: "uint8", want: uint8(math.MaxUint8)},
			{name: "uint8 to uint16", value: uint8(42), targetType: "uint16", want: uint16(42)},
			{name: "uint8 to uint32", value: uint8(42), targetType: "uint32", want: uint32(42)},
			{name: "uint8 to uint64", value: uint8(42), targetType: "uint64", want: uint64(42)},
			{name: "uint8 to uint8", value: uint8(42), targetType: "uint8", want: uint8(42)},
			{name: "uint8 underflow", value: -1, targetType: "uint8", want: uint8(0)},
			{name: "zero to uint8", value: 0, targetType: "uint8", want: uint8(0)},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := dsl.cast(tt.value, tt.targetType)
				testResult(t, tt.name, tt.want, tt.wantErr, got, err)
			})
		}
	})

	t.Run("Uint16 Conversions", func(t *testing.T) {
		tests := []struct {
			name       string
			value      any
			targetType string
			want       any
			wantErr    bool
		}{
			{name: "negative to uint16", value: -1, targetType: "uint16", want: uint16(0)},
			{name: "uint16 overflow", value: math.MaxUint32, targetType: "uint16", want: uint16(math.MaxUint16)},
			{name: "uint16 to uint16", value: uint16(42), targetType: "uint16", want: uint16(42)},
			{name: "uint16 to uint32", value: uint16(42), targetType: "uint32", want: uint32(42)},
			{name: "uint16 to uint64", value: uint16(42), targetType: "uint64", want: uint64(42)},
			{name: "uint16 underflow", value: -1, targetType: "uint16", want: uint16(0)},
			{name: "uint32 to uint16", value: uint32(42), targetType: "uint16", want: uint16(42)},
			{name: "uint64 to uint16", value: uint64(42), targetType: "uint16", want: uint16(42)},
			{name: "zero to uint16", value: 0, targetType: "uint16", want: uint16(0)},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := dsl.cast(tt.value, tt.targetType)
				testResult(t, tt.name, tt.want, tt.wantErr, got, err)
			})
		}
	})

	t.Run("Uint32 Conversions", func(t *testing.T) {
		tests := []struct {
			name       string
			value      any
			targetType string
			want       any
			wantErr    bool
		}{
			{name: "negative to uint32", value: -1, targetType: "uint32", want: uint32(0)},
			{name: "uint32 to uint32", value: uint32(42), targetType: "uint32", want: uint32(42)},
			{name: "uint32 to uint64", value: uint32(42), targetType: "uint64", want: uint64(42)},
			{name: "uint32 underflow", value: -1, targetType: "uint32", want: uint32(0)},
			{name: "uint64 max to uint32", value: uint64(math.MaxUint64), targetType: "uint32", want: uint32(math.MaxUint32)},
			{name: "uint64 to uint32", value: uint64(42), targetType: "uint32", want: uint32(42)},
			{name: "zero to uint32", value: 0, targetType: "uint32", want: uint32(0)},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := dsl.cast(tt.value, tt.targetType)
				testResult(t, tt.name, tt.want, tt.wantErr, got, err)
			})
		}
	})

	t.Run("Uint64 Conversions", func(t *testing.T) {
		tests := []struct {
			name       string
			value      any
			targetType string
			want       any
			wantErr    bool
		}{
			{name: "negative to uint64", value: -1, targetType: "uint64", want: uint64(0)},
			{name: "uint64 to uint64", value: uint64(42), targetType: "uint64", want: uint64(42)},
			{name: "uint64 underflow", value: -1, targetType: "uint64", want: uint64(0)},
			{name: "zero to uint64", value: 0, targetType: "uint64", want: uint64(0)},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := dsl.cast(tt.value, tt.targetType)
				testResult(t, tt.name, tt.want, tt.wantErr, got, err)
			})
		}
	})

	t.Run("Special Cases", func(t *testing.T) {
		tests := []struct {
			name       string
			value      any
			targetType string
			want       any
			wantErr    bool
		}{
			{name: "channel type", value: make(chan int), targetType: "int", wantErr: true},
			{name: "complex128 to complex128", value: complex128(1 + 2i), targetType: "complex128", wantErr: true},
			{name: "complex128 to complex64", value: complex128(1 + 2i), targetType: "complex64", wantErr: true},
			{name: "complex64 to complex128", value: complex64(1 + 2i), targetType: "complex128", wantErr: true},
			{name: "complex64 to complex64", value: complex64(1 + 2i), targetType: "complex64", wantErr: true},
			{name: "custom type to error interface", value: struct{}{}, targetType: "error", wantErr: true},
			{name: "empty interface", value: any(nil), targetType: "int", wantErr: true},
			{name: "function to function", value: func() {}, targetType: "func()", wantErr: true},
			{name: "function to interface", value: func() {}, targetType: "interface{}", wantErr: true},
			{name: "map to interface", value: map[string]int{"a": 1}, targetType: "interface{}", wantErr: true},
			{name: "map to map", value: map[string]int{"a": 1}, targetType: "map[string]int", wantErr: true},
			{name: "nil value", value: nil, targetType: "int", wantErr: true},
			{name: "unsupported target type", value: 42, targetType: "map", wantErr: true},
			{name: "unsupported type", value: struct{}{}, targetType: "int", wantErr: true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := dsl.cast(tt.value, tt.targetType)
				testResult(t, tt.name, tt.want, tt.wantErr, got, err)
			})
		}
	})

}

func TestParser(t *testing.T) {
	t.Run("Parser", func(t *testing.T) {
		type TestCase struct {
			name    string
			script  string
			args    []any
			want    *dslResult
			wantErr bool
		}

		c := func(name string, script string, args []any, want *dslResult, wantErr bool) TestCase {
			return TestCase{name, script, args, want, wantErr}
		}
		imgNRGBA64 := image.NewNRGBA64(image.Rect(0, 0, 100, 100))
		imgRGBA64 := image.NewRGBA64(image.Rect(0, 0, 100, 100))
		tests := []TestCase{
			c("empty string as named argument", `test-function-1(str="")`, []any{}, &dslResult{0, nil}, false),
			c("argument out of range", `$3`, []any{1, 2}, &dslResult{nil, fmt.Errorf("argument $3 out of range")}, true),
			c("basic argument usage", `$1`, []any{42}, &dslResult{42, nil}, false),
			c("invalid function name", "unknown-function(1 2)", []any{}, nil, true),
			c("invalid named argument", `test-function-1(x=1 y=2 invalid=3)`, []any{}, nil, true),
			c("just a bool", `true`, []any{}, &dslResult{true, nil}, false),
			c("just a float", `42.1`, []any{}, &dslResult{42.1, nil}, false),
			c("just a string", `"hello"`, []any{}, &dslResult{"hello", nil}, false),
			c("just an int", `42`, []any{}, &dslResult{int64(42), nil}, false),
			c("mixed types", `concat($1 $2)`, []any{"hello", 42}, &dslResult{"hello42", nil}, false),
			c("multiple arguments", `add($1 $2)`, []any{5, 3}, &dslResult{8, nil}, false),
			c("named arguments mixed with positional arguments", `test-function-1(x=1 y=test-function-2(1 2) str="hello")`, []any{}, nil, true),
			c("named arguments out of order", `test-function-1(str="hello" y=2 x=1)`, []any{}, &dslResult{3, nil}, false),
			c("named arguments", `test-function-1(x=1 y=2 str="hello")`, []any{}, &dslResult{3, nil}, false),
			c("named optional arguments 1", `test-function-1(x=1 y=2)`, []any{}, &dslResult{3, nil}, false),
			c("named optional arguments 2", `test-function-1(x=1)`, []any{}, &dslResult{1, nil}, false),
			c("nested argument usage", `add($1 mul($2 $3))`, []any{1, 2, 3}, &dslResult{7, nil}, false),
			c("nested function calls", `test-function-2(test-function-1(1 10 "Hello \" World") 0)`, []any{}, &dslResult{true, nil}, false),
			c("optional arguments with defaults 1", `test-function-1(1 2)`, []any{}, &dslResult{3, nil}, false),
			c("optional arguments with defaults 2", `test-function-1(1)`, []any{}, &dslResult{1, nil}, false),
			c("optional arguments with defaults 3", `test-function-1()`, []any{}, &dslResult{0, nil}, false),
			c("simple function call", `test-function-1(1 2 "hello \" mean\"world!\"")`, []any{}, &dslResult{3, nil}, false),
			c("image processing (RGBA64-RGBA64)", `img-rgba64($1)`, []any{imgRGBA64}, &dslResult{imgRGBA64, nil}, false),
			c("image processing (RGBA64-NRGBA64)", `img-nrgba64($1)`, []any{imgRGBA64}, &dslResult{imgNRGBA64, nil}, false),
			c("image processing (NRGBA64-RGBA64)", `img-rgba64($1)`, []any{imgNRGBA64}, &dslResult{imgRGBA64, nil}, false),
			c("image processing (NRGBA64-NRGBA64)", `img-nrgba64($1)`, []any{imgNRGBA64}, &dslResult{imgNRGBA64, nil}, false),
		}

		createTestLanguage()
		exportPath, _ := filepath.Abs("../../LANGUAGE.vsix")
		if err := dsl.exportVSCodeExtension(exportPath); err != nil {
			t.Errorf("could not generate VSCode extension: %v", err)
		}

		if err := flo.File("../../LANGUAGE.md").StoreString(dsl.docMarkdown()); err != nil {
			t.Errorf("could not generate Makrdown documentation: %v", err)
		}

		if err := flo.File("../../LANGUAGE.html").StoreString(dsl.docHTML()); err != nil {
			t.Errorf("could not generate HTML documentation: %v", err)
		}

		if err := flo.File("../../LANGUAGE.txt").StoreString(dsl.docText()); err != nil {
			t.Errorf("could not generate Text documentation: %v", err)
		}

		for _, tt := range tests {
			dsl.restoreState()
			t.Run(tt.name, func(t *testing.T) {
				got, err := dsl.run(tt.script, false, tt.args...)
				testResult(t, tt.name, tt.want, tt.wantErr, got, err)
			})
		}
	})
}

func TestShell(t *testing.T) {
	t.Run("Shell", func(t *testing.T) {
		createTestLanguage()
		dsl.shell()
	})
}
