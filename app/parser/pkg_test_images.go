package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"path/filepath"
)

const (
	imageTestOutputDir = "test_output/test_images" // Combined path for all test image outputs
)

// GenerateTestImages creates a set of test images for blend mode and filter testing
func GenerateTestImages() error {
	dsl := &dslCollection{}

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
		generateEdgeTestImage,      // Edge detection test
		generateAlphaGradientImage, // Pure alpha gradient
	}

	for _, generator := range generators {
		if err := generator(); err != nil {
			return err
		}
	}

	// Generate and save test images
	testImages := []struct {
		name string
		img  *image.NRGBA
	}{
		{"gradient", dsl.GenerateGradient()},
		{"checkerboard", dsl.GenerateCheckerboard()},
		{"color_wheel", dsl.GenerateColorWheel()},
		{"noise", dsl.GenerateNoise()},
		{"alpha_gradient", dsl.GenerateAlphaGradient()},
		{"color_bands", dsl.GenerateColorBands()},
		{"edge_cases", dsl.GenerateEdgeCases()},
	}

	for _, ti := range testImages {
		file, err := os.Create(filepath.Join(imageTestOutputDir, fmt.Sprintf("%s.png", ti.name)))
		if err != nil {
			return err
		}
		defer file.Close()

		if err := png.Encode(file, ti.img); err != nil {
			return err
		}
	}

	return nil
}

// generateGradientImage creates a gradient image with transparency
func generateGradientImage() error {
	const width, height = 512, 512
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

	return saveImage(img, "test_images/gradient.png")
}

// generateCheckerboardImage creates a checkerboard pattern with transparency
func generateCheckerboardImage() error {
	const width, height = 512, 512
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

	return saveImage(img, "test_images/checkerboard.png")
}

// generateColorWheelImage creates a color wheel with transparency
func generateColorWheelImage() error {
	const width, height = 512, 512
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

	return saveImage(img, "test_images/color_wheel.png")
}

// generateNoiseImage creates a noise pattern with transparency
func generateNoiseImage() error {
	const width, height = 512, 512
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

	return saveImage(img, "test_images/noise.png")
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

// saveImage saves an image to a file
func saveImage(img image.Image, filename string) error {
	f, err := os.Create(filepath.Join(imageTestOutputDir, filepath.Base(filename)))
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

// generateTextPatternImage creates a pattern that simulates text-like structures
func generateTextPatternImage() error {
	const width, height = 512, 512
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

	return saveImage(img, "test_images/text_pattern.png")
}

// generateConcentricsImage creates concentric circles with varying colors
func generateConcentricsImage() error {
	const width, height = 512, 512
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

	return saveImage(img, "test_images/concentrics.png")
}

// generateHighContrastImage creates a high contrast pattern
func generateHighContrastImage() error {
	const width, height = 512, 512
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

	return saveImage(img, "test_images/high_contrast.png")
}

// generateRainbowStripsImage creates horizontal rainbow strips
func generateRainbowStripsImage() error {
	const width, height = 512, 512
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

	return saveImage(img, "test_images/rainbow_strips.png")
}

// generateEdgeTestImage creates patterns good for testing edge detection
func generateEdgeTestImage() error {
	const width, height = 512, 512
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

	return saveImage(img, "test_images/edge_test.png")
}

// generateAlphaGradientImage creates a pattern with pure alpha gradients
func generateAlphaGradientImage() error {
	const width, height = 512, 512
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

	return saveImage(img, "test_images/alpha_gradient.png")
}

// GenerateAlphaGradient creates an image with varying alpha values
func (dsl *dslCollection) GenerateAlphaGradient() *image.NRGBA {
	width, height := 256, 256
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		alpha := uint8((y * 255) / height)
		for x := 0; x < width; x++ {
			img.Set(x, y, color.NRGBA{
				R: 255,
				G: 255,
				B: 255,
				A: alpha,
			})
		}
	}
	return img
}

// GenerateColorBands creates an image with distinct color bands
func (dsl *dslCollection) GenerateColorBands() *image.NRGBA {
	width, height := 256, 256
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	bandHeight := height / 8
	colors := []color.NRGBA{
		{255, 0, 0, 255},     // Red
		{0, 255, 0, 255},     // Green
		{0, 0, 255, 255},     // Blue
		{255, 255, 0, 255},   // Yellow
		{255, 0, 255, 255},   // Magenta
		{0, 255, 255, 255},   // Cyan
		{255, 255, 255, 255}, // White
		{0, 0, 0, 255},       // Black
	}

	for i, c := range colors {
		for y := i * bandHeight; y < (i+1)*bandHeight; y++ {
			for x := 0; x < width; x++ {
				img.Set(x, y, c)
			}
		}
	}
	return img
}

// GenerateEdgeCases creates an image with various edge cases
func (dsl *dslCollection) GenerateEdgeCases() *image.NRGBA {
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

	return img
}

// GenerateGradient creates a gradient image with transparency
func (dsl *dslCollection) GenerateGradient() *image.NRGBA {
	width, height := 256, 256
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Calculate normalized positions
			xNorm := float64(x) / float64(width)
			yNorm := float64(y) / float64(height)

			// Create smooth gradients for each channel
			r := uint8(255 * xNorm)
			g := uint8(255 * yNorm)
			b := uint8(128)
			// Alpha gradient from top-left to bottom-right
			a := uint8(255 * (1 - (xNorm+yNorm)/2))

			// Pre-multiply RGB values by alpha
			r = uint8(float64(r) * float64(a) / 255)
			g = uint8(float64(g) * float64(a) / 255)
			b = uint8(float64(b) * float64(a) / 255)

			img.Set(x, y, color.NRGBA{r, g, b, a})
		}
	}
	return img
}

// GenerateCheckerboard creates a checkerboard pattern with transparency
func (dsl *dslCollection) GenerateCheckerboard() *image.NRGBA {
	width, height := 256, 256
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	squareSize := 32

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			squareX := x / squareSize
			squareY := y / squareSize
			// Alternate between opaque and semi-transparent squares
			if (squareX+squareY)%2 == 0 {
				img.Set(x, y, color.NRGBA{0, 0, 0, 255}) // Opaque black
			} else {
				img.Set(x, y, color.NRGBA{255, 255, 255, 128}) // Semi-transparent white
			}
		}
	}
	return img
}

// GenerateColorWheel creates a color wheel image with transparency
func (dsl *dslCollection) GenerateColorWheel() *image.NRGBA {
	width, height := 256, 256
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	centerX, centerY := width/2, height/2
	maxRadius := float64(width) / 2

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dx := float64(x - centerX)
			dy := float64(y - centerY)
			radius := math.Sqrt(dx*dx + dy*dy)

			if radius <= maxRadius {
				angle := math.Atan2(dy, dx)
				if angle < 0 {
					angle += 2 * math.Pi
				}

				// Calculate alpha based on distance from center
				alpha := uint8(255 * (1 - radius/maxRadius))

				// Convert angle to hue (0-1)
				hue := angle / (2 * math.Pi)
				r, g, b := hslToRGB(hue, 1.0, 0.5)

				// Pre-multiply RGB values by alpha
				r = uint8(float64(r) * float64(alpha) / 255)
				g = uint8(float64(g) * float64(alpha) / 255)
				b = uint8(float64(b) * float64(alpha) / 255)

				img.Set(x, y, color.NRGBA{r, g, b, alpha})
			} else {
				img.Set(x, y, color.NRGBA{0, 0, 0, 0}) // Fully transparent outside the wheel
			}
		}
	}
	return img
}

// GenerateNoise creates a noise pattern with varying transparency
func (dsl *dslCollection) GenerateNoise() *image.NRGBA {
	width, height := 256, 256
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	rand.Seed(42) // Fixed seed for reproducibility

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Generate random values including alpha
			r := uint8(rand.Intn(256))
			g := uint8(rand.Intn(256))
			b := uint8(rand.Intn(256))
			a := uint8(rand.Intn(256)) // Random transparency

			// Pre-multiply RGB values by alpha
			r = uint8(float64(r) * float64(a) / 255)
			g = uint8(float64(g) * float64(a) / 255)
			b = uint8(float64(b) * float64(a) / 255)

			img.Set(x, y, color.NRGBA{r, g, b, a})
		}
	}
	return img
}
