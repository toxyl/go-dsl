package main

import (
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/toxyl/flo"
	"github.com/toxyl/math"
)

const (
	imageTestOutputDir = "test_output/test_images" // Combined path for all test image outputs
)

// GenerateTestImages creates a set of test images for blend mode and filter testing
func GenerateTestImages() error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(imageTestOutputDir, 0755); err != nil {
		return err
	}

	// Generate each test image
	generators := []func() error{
		generateGradientImage,
		generateCheckerboardImage,
		generateColorWheelImage,
		generateNoiseImage,
		generateTextPatternImage,   // Text-like patterns
		generateConcentricsImage,   // Concentric circles
		generateHighContrastImage,  // High contrast patterns
		generateRainbowStripsImage, // Rainbow strips
		generateEdgeDetectionImage, // Edge detection test
		generateEdgeCasesImage,     // Edge cases test
		generateAlphaGradientImage, // Pure alpha gradient
	}

	for _, generator := range generators {
		if err := generator(); err != nil {
			return err
		}
	}

	return nil
}

// hslToRGB converts HSL to RGB
func hslToRGB(h, s, l float64) (uint8, uint8, uint8) {
	var r, g, b float64

	if s == 0 {
		r, g, b = l, l, l
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p := 2*l - q

		r = hueToRGB(p, q, h+1.0/3.0)
		g = hueToRGB(p, q, h)
		b = hueToRGB(p, q, h-1.0/3.0)
	}

	return uint8(r * 255), uint8(g * 255), uint8(b * 255)
}

// hueToRGB helper function for HSL to RGB conversion
func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}

// loadImage loads an image from a file and converts it to NRGBA format
func loadImage(filename string) (*image.NRGBA, error) {
	fabs, _ := filepath.Abs(filename)
	f, err := os.Open(fabs)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	// Convert to NRGBA if it isn't already
	if nrgba, ok := img.(*image.NRGBA); ok {
		return nrgba, nil
	}

	// Convert the image to NRGBA format
	bounds := img.Bounds()
	nrgba := image.NewNRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			nrgba.Set(x, y, img.At(x, y))
		}
	}
	return nrgba, nil
}

// saveImage saves an image to a file
func saveImage(img image.Image, filename string) error {
	fabs, _ := filepath.Abs(filename)
	flo.Dir(filepath.Dir(fabs)).Mkdir(0755)
	f, err := os.Create(fabs)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func generateImageSprite(filenames ...string) *image.NRGBA {
	totalWidth := 0
	maxHeight := 0
	images := []*image.NRGBA{}
	for _, file := range filenames {
		img, err := loadImage(file)
		if err != nil {
			panic("could not load sprite image: " + file + "; err: " + err.Error())
		}
		images = append(images, img)
		totalWidth += img.Rect.Dx()
		maxHeight = math.Max(maxHeight, img.Rect.Dy())
	}
	sp := image.NewNRGBA(image.Rect(0, 0, totalWidth, maxHeight))
	x := 0
	for _, img := range images {
		w := img.Rect.Dx()
		h := img.Rect.Dy()

		// Copy each pixel from the source image to the sprite sheet
		for y := range h {
			for sx := range w {
				sp.Set(x+sx, y, img.At(sx, y))
			}
		}

		x += w
	}
	return sp
}

func drawGrid(img *image.NRGBA, gridSize int) *image.NRGBA {
	bounds := img.Bounds()
	debug := image.NewNRGBA(bounds)

	// Copy original image
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			debug.Set(x, y, img.NRGBAAt(x, y))
		}
	}

	// Draw grid
	for y := bounds.Min.Y; y < bounds.Max.Y; y += gridSize {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			debug.Set(x, y, color.NRGBA{0, 0, 0, 128})
		}
	}
	for x := bounds.Min.X; x < bounds.Max.X; x += gridSize {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			debug.Set(x, y, color.NRGBA{0, 0, 0, 128})
		}
	}
	return debug
}

func drawBounds(img *image.NRGBA) *image.NRGBA {
	bounds := img.Bounds()
	debug := image.NewNRGBA(bounds)

	// Copy original image
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			debug.Set(x, y, img.NRGBAAt(x, y))
		}
	}

	// Draw bounds
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		debug.Set(x, bounds.Min.Y, color.NRGBA{255, 0, 0, 255})   // Top
		debug.Set(x, bounds.Max.Y-1, color.NRGBA{255, 0, 0, 255}) // Bottom
	}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		debug.Set(bounds.Min.X, y, color.NRGBA{255, 0, 0, 255})   // Left
		debug.Set(bounds.Max.X-1, y, color.NRGBA{255, 0, 0, 255}) // Right
	}

	// Draw center cross
	centerX := (bounds.Min.X + bounds.Max.X) / 2
	centerY := (bounds.Min.Y + bounds.Max.Y) / 2
	for i := -5; i <= 5; i++ {
		if centerX+i >= bounds.Min.X && centerX+i < bounds.Max.X {
			debug.Set(centerX+i, centerY, color.NRGBA{0, 255, 0, 255})
		}
		if centerY+i >= bounds.Min.Y && centerY+i < bounds.Max.Y {
			debug.Set(centerX, centerY+i, color.NRGBA{0, 255, 0, 255})
		}
	}

	return debug
}

// generateTextPatternImage creates a pattern that simulates text-like structures
func generateTextPatternImage() error {
	const width, height = 256, 256
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	// Create "text-like" horizontal lines with varying thickness and spacing
	for y := 0; y < height; y++ {
		// Determine if we're in a "text line" region
		lineRegion := (y/40)%3 == 0
		if lineRegion {
			wordStart := 0
			for x := 0; x < width; x++ {
				// Create "word-like" blocks
				if x >= wordStart && x < wordStart+rand.Intn(50)+20 {
					img.Set(x, y, color.NRGBA{0, 0, 0, 255})
				} else {
					img.Set(x, y, color.NRGBA{255, 255, 255, 255})
					if x >= wordStart+70 {
						wordStart = x + 20
					}
				}
			}
		} else {
			for x := 0; x < width; x++ {
				img.Set(x, y, color.NRGBA{255, 255, 255, 255})
			}
		}
	}
	return saveImage(drawBounds(drawGrid(img, 20)), "./test_output/test_images/text_pattern.png")
}

// LoadTextPatternImage loads the text pattern image
func LoadTextPatternImage() *image.NRGBA {
	img, err := loadImage("./test_output/test_images/text_pattern.png")
	if err != nil {
		return nil
	}
	return img
}

// generateConcentricsImage creates concentric circles with varying colors
func generateConcentricsImage() error {
	const width, height = 256, 256
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	centerX, centerY := width/2, height/2

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dx := float64(x - centerX)
			dy := float64(y - centerY)
			distance := math.Sqrt(dx*dx + dy*dy)

			// Create rings with different colors
			ring := int(distance) / 20
			switch ring % 4 {
			case 0:
				img.Set(x, y, color.NRGBA{255, 0, 0, 255})
			case 1:
				img.Set(x, y, color.NRGBA{0, 255, 0, 255})
			case 2:
				img.Set(x, y, color.NRGBA{0, 0, 255, 255})
			case 3:
				img.Set(x, y, color.NRGBA{255, 255, 0, 255})
			}
		}
	}

	return saveImage(drawBounds(drawGrid(img, 20)), "./test_output/test_images/concentrics.png")
}

// LoadConcentricsImage loads the concentrics image
func LoadConcentricsImage() *image.NRGBA {
	img, err := loadImage("./test_output/test_images/concentrics.png")
	if err != nil {
		return nil
	}
	return img
}

// generateHighContrastImage creates a high contrast pattern
func generateHighContrastImage() error {
	const width, height = 256, 256
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Create a complex pattern with high contrast
			val := math.Sin(float64(x)/10) * math.Cos(float64(y)/10)
			if val > 0 {
				img.Set(x, y, color.NRGBA{255, 255, 255, 255})
			} else {
				img.Set(x, y, color.NRGBA{0, 0, 0, 255})
			}
		}
	}

	return saveImage(drawBounds(drawGrid(img, 20)), "./test_output/test_images/high_contrast.png")
}

// LoadHighContrastImage loads the high contrast image
func LoadHighContrastImage() *image.NRGBA {
	img, err := loadImage("./test_output/test_images/high_contrast.png")
	if err != nil {
		return nil
	}
	return img
}

// generateRainbowStripsImage creates horizontal rainbow strips
func generateRainbowStripsImage() error {
	const width, height = 256, 256
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	colors := []color.NRGBA{
		{255, 0, 0, 255},   // Red
		{255, 127, 0, 255}, // Orange
		{255, 255, 0, 255}, // Yellow
		{0, 255, 0, 255},   // Green
		{0, 0, 255, 255},   // Blue
		{75, 0, 130, 255},  // Indigo
		{148, 0, 211, 255}, // Violet
	}

	stripHeight := height / len(colors)
	for y := 0; y < height; y++ {
		colorIndex := y / stripHeight
		if colorIndex >= len(colors) {
			colorIndex = len(colors) - 1
		}
		for x := 0; x < width; x++ {
			img.Set(x, y, colors[colorIndex])
		}
	}

	return saveImage(drawBounds(drawGrid(img, 20)), "./test_output/test_images/rainbow_strips.png")
}

// LoadRainbowStripsImage loads the rainbow strips image
func LoadRainbowStripsImage() *image.NRGBA {
	img, err := loadImage("./test_output/test_images/rainbow_strips.png")
	if err != nil {
		return nil
	}
	return img
}

// generateEdgeDetectionImage creates patterns good for testing edge detection
func generateEdgeDetectionImage() error {
	const width, height = 256, 256
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	// Draw various geometric shapes
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Default to white
			c := color.NRGBA{255, 255, 255, 255}

			// Draw a triangle
			if y > x && y < height/2 {
				c = color.NRGBA{0, 0, 0, 255}
			}

			// Draw a circle
			dx := float64(x - width*3/4)
			dy := float64(y - height/2)
			if math.Sqrt(dx*dx+dy*dy) < 100 {
				c = color.NRGBA{0, 0, 0, 255}
			}

			// Draw some lines
			if (x+y)%50 == 0 {
				c = color.NRGBA{0, 0, 0, 255}
			}

			img.Set(x, y, c)
		}
	}

	return saveImage(drawBounds(drawGrid(img, 20)), "./test_output/test_images/edge_test.png")
}

// LoadEdgeDetectionImage loads the edge test image
func LoadEdgeDetectionImage() *image.NRGBA {
	img, err := loadImage("./test_output/test_images/edge_test.png")
	if err != nil {
		return nil
	}
	return img
}

// generateEdgeCasesImage creates an image with various edge cases
func generateEdgeCasesImage() error {
	width, height := 256, 256
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	// Fill with white
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.NRGBA{255, 255, 255, 255})
		}
	}

	// Add various edge cases
	// 1. Single pixel in each corner
	img.Set(0, 0, color.NRGBA{255, 0, 0, 255})                // Top-left
	img.Set(width-1, 0, color.NRGBA{0, 255, 0, 255})          // Top-right
	img.Set(0, height-1, color.NRGBA{0, 0, 255, 255})         // Bottom-left
	img.Set(width-1, height-1, color.NRGBA{255, 255, 0, 255}) // Bottom-right

	// 2. Vertical and horizontal lines
	for y := 0; y < height; y++ {
		img.Set(width/2, y, color.NRGBA{0, 0, 0, 255}) // Vertical line
	}
	for x := 0; x < width; x++ {
		img.Set(x, height/2, color.NRGBA{0, 0, 0, 255}) // Horizontal line
	}
	// 3. Checkerboard pattern in center
	for y := height / 4; y < 3*height/4; y++ {
		for x := width / 4; x < 3*width/4; x++ {
			if (x+y)%2 == 0 {
				img.Set(x, y, color.NRGBA{0, 0, 0, 255})
			}
		}
	}

	return saveImage(drawBounds(drawGrid(img, 20)), "./test_output/test_images/edge_cases.png")
}

// LoadEdgeCasesImage loads the edge test image
func LoadEdgeCasesImage() *image.NRGBA {
	img, err := loadImage("./test_output/test_images/edge_cases.png")
	if err != nil {
		return nil
	}
	return img
}

// generateAlphaGradientImage creates a pattern with pure alpha gradients
func generateAlphaGradientImage() error {
	const width, height = 256, 256
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Create a radial alpha gradient
			dx := float64(x - width/2)
			dy := float64(y - height/2)
			distance := math.Sqrt(dx*dx + dy*dy)
			maxDistance := math.Sqrt(float64(width*width+height*height)) / 2

			alpha := uint8((1 - distance/maxDistance) * 255)
			// Use pure white with varying alpha
			img.Set(x, y, color.NRGBA{255, 255, 255, alpha})
		}
	}

	return saveImage(drawBounds(drawGrid(img, 20)), "./test_output/test_images/alpha_gradient.png")
}

// LoadAlphaGradientImage loads the alpha gradient image
func LoadAlphaGradientImage() *image.NRGBA {
	img, err := loadImage("./test_output/test_images/alpha_gradient.png")
	if err != nil {
		return nil
	}
	return img
}

// generateGradientImage creates a gradient image with transparency
func generateGradientImage() error {
	const width, height = 256, 256
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Calculate gradient values
			r := uint8((float64(x) / float64(width)) * 255)
			g := uint8((float64(y) / float64(height)) * 255)
			b := uint8(128)
			a := uint8((float64(x+y) / float64(width+height)) * 255)

			img.Set(x, y, color.NRGBA{r, g, b, a})
		}
	}

	return saveImage(drawBounds(drawGrid(img, 20)), "./test_output/test_images/gradient.png")
}

// LoadGradientImage loads the gradient image
func LoadGradientImage() *image.NRGBA {
	img, err := loadImage("./test_output/test_images/gradient.png")
	if err != nil {
		return nil
	}
	return img
}

// generateCheckerboardImage creates a checkerboard pattern with transparency
func generateCheckerboardImage() error {
	const width, height = 256, 256
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	const squareSize = 64

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Calculate checkerboard pattern
			squareX := x / squareSize
			squareY := y / squareSize
			isBlack := (squareX+squareY)%2 == 0

			// Create semi-transparent black and white squares
			if isBlack {
				img.Set(x, y, color.NRGBA{0, 0, 0, 192})
			} else {
				img.Set(x, y, color.NRGBA{255, 255, 255, 128})
			}
		}
	}

	return saveImage(drawBounds(drawGrid(img, 20)), "./test_output/test_images/checkerboard.png")
}

// LoadGradientImage loads the checkerboard image
func LoadCheckerboardImage() *image.NRGBA {
	img, err := loadImage("./test_output/test_images/checkerboard.png")
	if err != nil {
		return nil
	}
	return img
}

// generateColorWheelImage creates a color wheel with transparency
func generateColorWheelImage() error {
	const width, height = 256, 256
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	centerX, centerY := width/2, height/2
	maxRadius := math.Min(float64(width), float64(height)) / 2

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Calculate polar coordinates
			dx := float64(x - centerX)
			dy := float64(y - centerY)
			radius := math.Sqrt(dx*dx + dy*dy)
			angle := math.Atan2(dy, dx)

			if radius <= maxRadius {
				// Convert angle to hue (0-360)
				hue := (angle + math.Pi) / (2 * math.Pi)
				// Convert to RGB
				r, g, b := hslToRGB(hue, 1.0, 0.5)
				// Calculate alpha based on radius
				alpha := uint8((1 - radius/maxRadius) * 255)
				img.Set(x, y, color.NRGBA{r, g, b, alpha})
			} else {
				img.Set(x, y, color.NRGBA{0, 0, 0, 0})
			}
		}
	}

	return saveImage(drawBounds(drawGrid(img, 20)), "./test_output/test_images/color_wheel.png")
}

// LoadColorWheelImage loads the color wheel image
func LoadColorWheelImage() *image.NRGBA {
	img, err := loadImage("./test_output/test_images/color_wheel.png")
	if err != nil {
		return nil
	}
	return img
}

// generateNoiseImage creates a noise pattern with transparency
func generateNoiseImage() error {
	const width, height = 256, 256
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	rand.Seed(42) // Fixed seed for reproducibility

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Generate random values for RGB
			r := uint8(rand.Intn(256))
			g := uint8(rand.Intn(256))
			b := uint8(rand.Intn(256))
			// Generate random alpha with bias towards transparency
			a := uint8(rand.Intn(192)) // Max 75% opacity

			img.Set(x, y, color.NRGBA{r, g, b, a})
		}
	}

	return saveImage(drawBounds(drawGrid(img, 20)), "./test_output/test_images/noise.png")
}

// LoadNoiseImage loads the noise image
func LoadNoiseImage() *image.NRGBA {
	img, err := loadImage("./test_output/test_images/noise.png")
	if err != nil {
		return nil
	}
	return img
}
