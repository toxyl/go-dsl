package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/toxyl/math"

	"github.com/toxyl/flo"
)

const (
	testOutputDir = "test_output"
)

func init() {
	// Create test output directory if it doesn't exist
	if err := os.MkdirAll(testOutputDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create test output directory: %v", err))
	}
}

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
	dsl.funcs.register(
		"img-nrgba64", "This is a function to process an NRGBA64 image",
		[]dslParamMeta{{name: "img", typ: "*image.NRGBA64", desc: "The image to process"}},
		[]dslParamMeta{{name: "res", typ: "*image.NRGBA64", def: false, desc: "The converted image"}},
		func(a ...any) (any, error) {
			img := a[0].(*image.NRGBA64)
			file, err := os.Create(filepath.Join(testOutputDir, "LANGUAGE-NRGBA64.png"))
			if err != nil {
				return nil, err
			}
			defer file.Close()
			png.Encode(file, img)
			return img, nil
		},
	)
	dsl.funcs.register(
		"img-rgba64", "This is another function to process an RGBA64 image",
		[]dslParamMeta{{name: "img", typ: "*image.RGBA64", desc: "The image to process"}},
		[]dslParamMeta{{name: "res", typ: "*image.RGBA64", def: false, desc: "The converted image"}},
		func(a ...any) (any, error) {
			img := a[0].(*image.RGBA64)
			file, err := os.Create(filepath.Join(testOutputDir, "LANGUAGE-RGBA64.png"))
			if err != nil {
				return nil, err
			}
			defer file.Close()
			png.Encode(file, img)
			return img, nil
		},
	)
	dsl.funcs.register(
		"img-nrgba", "This is a function to process an NRGBA image",
		[]dslParamMeta{{name: "img", typ: "*image.NRGBA", desc: "The image to process"}},
		[]dslParamMeta{{name: "res", typ: "*image.NRGBA", def: false, desc: "The converted image"}},
		func(a ...any) (any, error) {
			img := a[0].(*image.NRGBA)
			file, err := os.Create(filepath.Join(testOutputDir, "LANGUAGE-NRGBA64-NRGBA.png"))
			if err != nil {
				return nil, err
			}
			defer file.Close()
			png.Encode(file, img)
			return img, nil
		},
	)
	dsl.funcs.register(
		"img-rgba", "This is another function to process an RGBA image",
		[]dslParamMeta{{name: "img", typ: "*image.RGBA", desc: "The image to process"}},
		[]dslParamMeta{{name: "res", typ: "*image.RGBA", def: false, desc: "The converted image"}},
		func(a ...any) (any, error) {
			img := a[0].(*image.RGBA)
			file, err := os.Create(filepath.Join(testOutputDir, "LANGUAGE-RGBA64-RGBA.png"))
			if err != nil {
				return nil, err
			}
			defer file.Close()
			png.Encode(file, img)
			return img, nil
		},
	)
	dsl.funcs.register(
		"load", "This is a function to load RGBA images",
		[]dslParamMeta{{name: "src", typ: "string", desc: "The image file to process"}},
		[]dslParamMeta{{name: "res", typ: "*image.NRGBA", def: false, desc: "The image"}},
		func(a ...any) (any, error) {
			path := a[0].(string)
			file, err := os.Open(path)
			if err != nil {
				return nil, err
			}
			defer file.Close()

			img, err := png.Decode(file)
			if err != nil {
				return nil, err
			}

			// Try to get as NRGBA first
			if nrgba, ok := img.(*image.NRGBA); ok {
				return nrgba, nil
			}

			// If it's RGBA, convert it to NRGBA
			if rgba, ok := img.(*image.RGBA); ok {
				return dsl.convertRGBAToNRGBA(rgba), nil
			}

			// If we get here, the image is neither RGBA nor NRGBA
			return nil, fmt.Errorf("unsupported image format: %T (must be NRGBA or RGBA)", img)
		},
	)
	dsl.funcs.register(
		"save", "This is a function to save RGBA images",
		[]dslParamMeta{{name: "img", typ: "*image.NRGBA", desc: "The image to save"}, {name: "path", typ: "string", desc: "The file to write to"}},
		[]dslParamMeta{{name: "res", typ: "bool", def: false, desc: "Whether saving was successful"}},
		func(a ...any) (any, error) {
			img := a[0].(*image.NRGBA)
			path := a[1].(string)
			file, err := os.Create(path)
			if err != nil {
				return nil, err
			}
			defer file.Close()
			png.Encode(file, img)
			return img, nil
		},
	)
	dsl.funcs.register(
		"invert", "Inverts the image",
		[]dslParamMeta{{name: "img", typ: "*image.NRGBA", desc: "The image to process"}},
		[]dslParamMeta{{name: "res", typ: "*image.NRGBA", def: false, desc: "The inverted image"}},
		func(a ...any) (any, error) {
			img := a[0].(*image.NRGBA)
			bounds := img.Bounds()
			inverted := image.NewNRGBA(bounds)

			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					c := img.NRGBAAt(x, y)
					inverted.Set(x, y, color.NRGBA{
						R: 255 - c.R,
						G: 255 - c.G,
						B: 255 - c.B,
						A: c.A, // Keep original alpha
					})
				}
			}
			return inverted, nil
		},
	)
	dsl.funcs.register(
		"grayscale", "Converts the image to grayscale",
		[]dslParamMeta{{name: "img", typ: "*image.NRGBA", desc: "The image to process"}},
		[]dslParamMeta{{name: "res", typ: "*image.NRGBA", def: false, desc: "The grayscale image"}},
		func(a ...any) (any, error) {
			img := a[0].(*image.NRGBA)
			bounds := img.Bounds()
			grayscaled := image.NewNRGBA(bounds)

			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					c := img.NRGBAAt(x, y)
					// Using luminosity method: 0.21 R + 0.72 G + 0.07 B
					gray := uint8(float64(c.R)*0.21 + float64(c.G)*0.72 + float64(c.B)*0.07)
					grayscaled.Set(x, y, color.NRGBA{
						R: gray,
						G: gray,
						B: gray,
						A: c.A,
					})
				}
			}
			return grayscaled, nil
		},
	)
	dsl.funcs.register(
		"brightness", "Adjusts the brightness of the image",
		[]dslParamMeta{
			{name: "img", typ: "*image.NRGBA", desc: "The image to process"},
			{name: "factor", typ: "float64", desc: "Brightness adjustment factor (0.0 to 2.0, 1.0 is original)"},
		},
		[]dslParamMeta{{name: "res", typ: "*image.NRGBA", def: false, desc: "The brightness-adjusted image"}},
		func(a ...any) (any, error) {
			img := a[0].(*image.NRGBA)
			factor := a[1].(float64)

			if factor < 0.0 || factor > 2.0 {
				return nil, fmt.Errorf("brightness factor must be between 0.0 and 2.0")
			}

			bounds := img.Bounds()
			adjusted := image.NewNRGBA(bounds)

			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					c := img.NRGBAAt(x, y)
					adjusted.Set(x, y, color.NRGBA{
						R: uint8(math.Min(float64(c.R)*factor, 255)),
						G: uint8(math.Min(float64(c.G)*factor, 255)),
						B: uint8(math.Min(float64(c.B)*factor, 255)),
						A: c.A,
					})
				}
			}
			return adjusted, nil
		},
	)
	dsl.funcs.register(
		"sepia", "Applies a sepia tone effect to the image",
		[]dslParamMeta{{name: "img", typ: "*image.NRGBA", desc: "The image to process"}},
		[]dslParamMeta{{name: "res", typ: "*image.NRGBA", def: false, desc: "The sepia-toned image"}},
		func(a ...any) (any, error) {
			img := a[0].(*image.NRGBA)
			bounds := img.Bounds()
			sepia := image.NewNRGBA(bounds)

			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					c := img.NRGBAAt(x, y)
					r := float64(c.R)
					g := float64(c.G)
					b := float64(c.B)

					newR := uint8(math.Min((r*0.393)+(g*0.769)+(b*0.189), 255))
					newG := uint8(math.Min((r*0.349)+(g*0.686)+(b*0.168), 255))
					newB := uint8(math.Min((r*0.272)+(g*0.534)+(b*0.131), 255))

					sepia.Set(x, y, color.NRGBA{
						R: newR,
						G: newG,
						B: newB,
						A: c.A,
					})
				}
			}
			return sepia, nil
		},
	)
	dsl.funcs.register(
		"blur", "Applies a simple box blur to the image",
		[]dslParamMeta{
			{name: "img", typ: "*image.NRGBA", desc: "The image to process"},
			{name: "radius", typ: "int", desc: "Blur radius (1-10)"},
		},
		[]dslParamMeta{{name: "res", typ: "*image.NRGBA", def: false, desc: "The blurred image"}},
		func(a ...any) (any, error) {
			img := a[0].(*image.NRGBA)
			radius := a[1].(int)

			if radius < 1 || radius > 10 {
				return nil, fmt.Errorf("blur radius must be between 1 and 10")
			}

			bounds := img.Bounds()
			blurred := image.NewNRGBA(bounds)

			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					var rSum, gSum, bSum, aSum, count float64

					for dy := -radius; dy <= radius; dy++ {
						for dx := -radius; dx <= radius; dx++ {
							nx, ny := x+dx, y+dy
							if nx >= bounds.Min.X && nx < bounds.Max.X && ny >= bounds.Min.Y && ny < bounds.Max.Y {
								c := img.NRGBAAt(nx, ny)
								rSum += float64(c.R)
								gSum += float64(c.G)
								bSum += float64(c.B)
								aSum += float64(c.A)
								count++
							}
						}
					}

					blurred.Set(x, y, color.NRGBA{
						R: uint8(rSum / count),
						G: uint8(gSum / count),
						B: uint8(bSum / count),
						A: uint8(aSum / count),
					})
				}
			}
			return blurred, nil
		},
	)
	dsl.funcs.register(
		"blend-multiply", "Overlays two images using the multiply blendmode",
		[]dslParamMeta{
			{name: "imgA", typ: "*image.RGBA", desc: "The lower image"},
			{name: "imgB", typ: "*image.RGBA", desc: "The upper image"},
		},
		[]dslParamMeta{{name: "res", typ: "*image.RGBA", def: false, desc: "The blended images"}},
		func(a ...any) (any, error) {
			imgA := a[0].(*image.RGBA)
			imgB := a[1].(*image.RGBA)

			// Create a new image with the same bounds
			bounds := imgA.Bounds()
			result := image.NewRGBA(bounds)

			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					// Get premultiplied colors
					c1 := imgA.RGBAAt(x, y)
					c2 := imgB.RGBAAt(x, y)

					// Work with 16-bit precision for intermediate calculations
					a1 := uint32(c1.A)
					r1 := uint32(c1.R)
					g1 := uint32(c1.G)
					b1 := uint32(c1.B)

					a2 := uint32(c2.A)
					r2 := uint32(c2.R)
					g2 := uint32(c2.G)
					b2 := uint32(c2.B)

					// Porter-Duff alpha compositing
					aOut := a1 + a2 - ((a1 * a2) / 255)

					var rOut, gOut, bOut uint32
					if aOut > 0 {
						// Multiply blend mode with minimal conversions
						// For each channel: (c1 * c2) + (c1 * (255 - a2)) + (c2 * (255 - a1)) all in fixed point
						rOut = ((r1 * r2) + (r1 * (255 - a2)) + (r2 * (255 - a1))) / 255
						gOut = ((g1 * g2) + (g1 * (255 - a2)) + (g2 * (255 - a1))) / 255
						bOut = ((b1 * b2) + (b1 * (255 - a2)) + (b2 * (255 - a1))) / 255
					}

					// Clamp to valid range
					rOut = min(rOut, 255)
					gOut = min(gOut, 255)
					bOut = min(bOut, 255)
					aOut = min(aOut, 255)

					result.SetRGBA(x, y, color.RGBA{
						R: uint8(rOut),
						G: uint8(gOut),
						B: uint8(bOut),
						A: uint8(aOut),
					})
				}
			}

			return result, nil
		},
	)
	dsl.funcs.register(
		"blend-screen", "Overlays two images using the screen blendmode",
		[]dslParamMeta{
			{name: "imgA", typ: "*image.RGBA", desc: "The lower image"},
			{name: "imgB", typ: "*image.RGBA", desc: "The upper image"},
		},
		[]dslParamMeta{{name: "res", typ: "*image.RGBA", def: false, desc: "The blended images"}},
		func(a ...any) (any, error) {
			imgA := a[0].(*image.RGBA)
			imgB := a[1].(*image.RGBA)

			// Create a new image with the same bounds
			bounds := imgA.Bounds()
			resultPre := image.NewRGBA(bounds)

			// Iterate through each pixel
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					// Get colors from both images (already premultiplied)
					c1 := imgA.RGBAAt(x, y)
					c2 := imgB.RGBAAt(x, y)

					// Screen blend mode in premultiplied space: 1 - (1 - a) * (1 - b)
					r := uint8(255 - ((255 - uint32(c1.R)) * (255 - uint32(c2.R)) / 255))
					g := uint8(255 - ((255 - uint32(c1.G)) * (255 - uint32(c2.G)) / 255))
					b := uint8(255 - ((255 - uint32(c1.B)) * (255 - uint32(c2.B)) / 255))
					// Alpha compositing
					a := uint8(255 - ((255 - uint32(c1.A)) * (255 - uint32(c2.A)) / 255))

					// Set the resulting color
					resultPre.Set(x, y, color.RGBA{r, g, b, a})
				}
			}

			// Convert back to non-premultiplied alpha
			return resultPre, nil
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

func TestImageProcessing(t *testing.T) {
	// Helper function to process an image through all filters
	processedImages := map[string]struct{}{}
	processImageWithAllFilters := func(t *testing.T, inputPathBottom, inputPathTop string) {
		// Create a fresh DSL instance for each test
		createTestLanguage()

		baseNameBottom := filepath.Base(inputPathBottom)
		nameBottomWithoutExt := strings.TrimSuffix(baseNameBottom, filepath.Ext(baseNameBottom))
		baseNameTop := filepath.Base(inputPathTop)
		nameTopWithoutExt := strings.TrimSuffix(baseNameTop, filepath.Ext(baseNameTop))

		// Basic filter tests
		filters := []struct {
			name   string
			script string
		}{
			{"invert", `save(invert(load(%q)) %q)`},
			{"grayscale", `save(grayscale(load(%q)) %q)`},
			{"sepia", `save(sepia(load(%q)) %q)`},
			{"blur-light", `save(blur(load(%q) 2) %q)`},
			{"blur-medium", `save(blur(load(%q) 5) %q)`},
			{"blur-heavy", `save(blur(load(%q) 8) %q)`},
			{"brightness-dark", `save(brightness(load(%q) 0.5) %q)`},
			{"brightness-normal", `save(brightness(load(%q) 1.0) %q)`},
			{"brightness-bright", `save(brightness(load(%q) 1.5) %q)`},
		}

		// Process basic filters
		for _, filter := range filters {
			dsl.restoreState() // Restore clean state for each test
			outputDir := filepath.Join(testOutputDir, filter.name)
			_ = flo.Dir(outputDir).Mkdir(0755)
			outputPath := filepath.Join(outputDir, fmt.Sprintf("%s.png", nameBottomWithoutExt))
			outputPathSprite := fmt.Sprintf("./test_output/%s-%s.png", filter.name, nameBottomWithoutExt)
			if _, ok := processedImages[outputPath]; ok {
				continue
			}
			processedImages[outputPath] = struct{}{}
			script := fmt.Sprintf(filter.script, inputPathBottom, outputPath)
			_, err := dsl.run(script, false)
			if err != nil {
				t.Errorf("Failed to process %s with %s: %v", baseNameBottom, filter.name, err)
			}
			fmt.Println(inputPathBottom, "-", outputPath, "-", outputPathSprite)
			saveImage(generateImageSprite(inputPathBottom, outputPath), outputPathSprite)
		}

		// Basic blendmode tests
		blendmodes := []struct {
			name   string
			script string
		}{
			{"blend-multiply", `save(blend-multiply(load(%q) load(%q)) %q)`},
			{"blend-screen", `save(blend-screen(load(%q) load(%q)) %q)`},
		}

		// Process basic blendmodes
		for _, filter := range blendmodes {
			dsl.restoreState() // Restore clean state for each test
			outputDir := filepath.Join(testOutputDir, filter.name)
			_ = flo.Dir(outputDir).Mkdir(0755)
			outputPath := filepath.Join(outputDir, fmt.Sprintf("%s-%s.png", nameBottomWithoutExt, nameTopWithoutExt))
			outputPathSprite := fmt.Sprintf("./test_output/%s-%s-%s.png", filter.name, nameBottomWithoutExt, nameTopWithoutExt)
			if _, ok := processedImages[outputPath]; ok {
				continue
			}
			processedImages[outputPath] = struct{}{}
			script := fmt.Sprintf(filter.script, inputPathBottom, inputPathTop, outputPath)
			_, err := dsl.run(script, false)
			if err != nil {
				t.Errorf("Failed to process %s with %s: %v", baseNameBottom, filter.name, err)
			}
			saveImage(generateImageSprite(inputPathBottom, inputPathTop, outputPath), outputPathSprite)
		}

		// Combined filter tests
		combinations := []struct {
			name   string
			script string
		}{
			{"invert-blur", `save(blur(invert(load(%q)) 3) %q)`},
			{"grayscale-bright", `save(brightness(grayscale(load(%q)) 1.4) %q)`},
			{"sepia-blur", `save(blur(sepia(load(%q)) 3) %q)`},
			{"invert-sepia", `save(sepia(invert(load(%q))) %q)`},
			{"blur-bright", `save(brightness(blur(load(%q) 3) 1.3) %q)`},
			{"grayscale-sepia", `save(sepia(grayscale(load(%q))) %q)`},
			{"triple-effect", `save(blur(sepia(brightness(load(%q) 1.2)) 2) %q)`},
		}

		// Process combinations
		for _, combo := range combinations {
			dsl.restoreState() // Restore clean state for each test
			outputDir := filepath.Join(testOutputDir, combo.name)
			_ = flo.Dir(outputDir).Mkdir(0755)
			outputPath := filepath.Join(outputDir, fmt.Sprintf("%s.png", nameBottomWithoutExt))
			outputPathSprite := fmt.Sprintf("./test_output/%s-%s.png", combo.name, nameBottomWithoutExt)
			if _, ok := processedImages[outputPath]; ok {
				continue
			}
			processedImages[outputPath] = struct{}{}
			script := fmt.Sprintf(combo.script, inputPathBottom, outputPath)
			_, err := dsl.run(script, false)
			if err != nil {
				t.Errorf("Failed to process %s with %s: %v", baseNameBottom, combo.name, err)
			}
			saveImage(generateImageSprite(inputPathBottom, outputPath), outputPathSprite)
		}

		// Error case tests
		errorCases := []struct {
			name   string
			script string
		}{
			{"brightness-error", `save(brightness(load(%q) 2.5) %q)`},
			{"blur-error", `save(blur(load(%q) 11) %q)`},
		}

		// Process error cases
		for _, errCase := range errorCases {
			dsl.restoreState() // Restore clean state for each test
			outputDir := filepath.Join(testOutputDir, "errors", errCase.name)
			_ = flo.Dir(outputDir).Mkdir(0755)
			outputPath := filepath.Join(outputDir, fmt.Sprintf("%s.png", nameBottomWithoutExt))
			if _, ok := processedImages[outputPath]; ok {
				continue
			}
			processedImages[outputPath] = struct{}{}
			script := fmt.Sprintf(errCase.script, inputPathBottom, outputPath)
			_, err := dsl.run(script, false)
			if err == nil {
				t.Errorf("Expected error for %s with %s but got none", baseNameBottom, errCase.name)
			}
		}
	}

	t.Run("Comprehensive Filter Tests", func(t *testing.T) {
		// First ensure test images are generated
		if err := GenerateTestImages(); err != nil {
			t.Fatalf("Failed to generate test images: %v", err)
		}

		// Read all test images from the directory
		entries, err := os.ReadDir(imageTestOutputDir)
		if err != nil {
			t.Fatalf("Failed to read test images directory: %v", err)
		}

		testImages := []string{}

		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
				inputPath := filepath.Join(imageTestOutputDir, entry.Name())
				t.Run(entry.Name(), func(t *testing.T) {
					testImages = append(testImages, inputPath)
				})
			}
		}

		// Process each test image
		for _, inA := range testImages {
			for _, inB := range testImages {
				processImageWithAllFilters(t, inA, inB)
			}
		}
	})
}

func TestShell(t *testing.T) {
	t.Run("Shell", func(t *testing.T) {
		createTestLanguage()
		dsl.shell()
	})
}
