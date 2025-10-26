package main

import (
	"fmt"
	"image"
	"sync"

	"github.com/toxyl/math"
)

type dslPixelProcessor func(r1, g1, b1, a1 uint32) (r, g, b, a uint32)

func dslParallelProcessImage[T image.Image](img image.Image, processor dslPixelProcessor, numWorkers int) (result image.Image) {
	switch t := any((*T)(nil)).(type) {
	case **image.NRGBA:
		result = image.NewNRGBA(img.Bounds())
	case **image.NRGBA64:
		result = image.NewNRGBA64(img.Bounds())
	case **image.RGBA:
		result = image.NewRGBA(img.Bounds())
	case **image.RGBA64:
		result = image.NewRGBA64(img.Bounds())
	default:
		panic(fmt.Sprintf("image type is unsupported: %T", t))
	}
	var (
		bounds = img.Bounds()
		width  = bounds.Dx()
		height = bounds.Dy()
		minX   = bounds.Min.X
		maxX   = bounds.Max.X
		minY   = bounds.Min.Y
		maxY   = bounds.Max.Y
	)

	if height == 0 || width == 0 {
		return
	}

	var wg sync.WaitGroup
	for i := range numWorkers {
		rowsPerWorker := (height + numWorkers - 1) / numWorkers
		startY := minY + i*rowsPerWorker
		endY := math.Min(startY+rowsPerWorker, maxY)
		if startY >= endY {
			continue
		}

		wg.Add(1)
		go func(startY, endY int) {
			defer wg.Done()
			r, g, b, a := uint32(0), uint32(0), uint32(0), uint32(0)
			for y := startY; y < endY; y++ {
				for x := minX; x < maxX; x++ {
					r, g, b, a = processor(dsl.getColor(img, x, y))
					dsl.setColor(result, x, y, r, g, b, a)
				}
			}
		}(startY, endY)
	}

	wg.Wait()
	return
}

func (dsl *dslCollection) getColor(img image.Image, x, y int) (r, g, b, a uint32) {
	if !(image.Point{x, y}.In(img.Bounds())) {
		return
	}
	switch t := img.(type) {
	case *image.NRGBA:
		i := t.PixOffset(x, y)
		s := t.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
		return uint32(s[0]), uint32(s[1]), uint32(s[2]), uint32(s[3])
	case *image.NRGBA64:
		i := t.PixOffset(x, y)
		s := t.Pix[i : i+8 : i+8] // Small cap improves performance, see https://golang.org/issue/27857
		return uint32(s[0])<<8 | uint32(s[1]),
			uint32(s[2])<<8 | uint32(s[3]),
			uint32(s[4])<<8 | uint32(s[5]),
			uint32(s[6])<<8 | uint32(s[7])
	case *image.RGBA:
		i := t.PixOffset(x, y)
		s := t.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
		return uint32(s[0]), uint32(s[1]), uint32(s[2]), uint32(s[3])
	case *image.RGBA64:
		i := t.PixOffset(x, y)
		s := t.Pix[i : i+8 : i+8] // Small cap improves performance, see https://golang.org/issue/27857
		return uint32(s[0])<<8 | uint32(s[1]),
			uint32(s[2])<<8 | uint32(s[3]),
			uint32(s[4])<<8 | uint32(s[5]),
			uint32(s[6])<<8 | uint32(s[7])
	default:
		panic(fmt.Sprintf("image type is unsupported: %T", t))
	}
}

func (dsl *dslCollection) getColorRedChannel(img image.Image, x, y int) (r uint32) {
	if !(image.Point{x, y}.In(img.Bounds())) {
		return
	}
	switch t := img.(type) {
	case *image.NRGBA:
		return uint32(t.Pix[t.PixOffset(x, y)])
	case *image.NRGBA64:
		i := t.PixOffset(x, y)
		return uint32(t.Pix[i])<<8 | uint32(t.Pix[i+1])
	case *image.RGBA:
		return uint32(t.Pix[t.PixOffset(x, y)])
	case *image.RGBA64:
		i := t.PixOffset(x, y)
		return uint32(t.Pix[i])<<8 | uint32(t.Pix[i+1])
	default:
		panic(fmt.Sprintf("image type is unsupported: %T", t))
	}
}

func (dsl *dslCollection) getColorGreenChannel(img image.Image, x, y int) (g uint32) {
	if !(image.Point{x, y}.In(img.Bounds())) {
		return
	}
	switch t := img.(type) {
	case *image.NRGBA:
		return uint32(t.Pix[t.PixOffset(x, y)+1])
	case *image.NRGBA64:
		i := t.PixOffset(x, y)
		return uint32(t.Pix[i+2])<<8 | uint32(t.Pix[i+3])
	case *image.RGBA:
		return uint32(t.Pix[t.PixOffset(x, y)+1])
	case *image.RGBA64:
		i := t.PixOffset(x, y)
		return uint32(t.Pix[i+2])<<8 | uint32(t.Pix[i+3])
	default:
		panic(fmt.Sprintf("image type is unsupported: %T", t))
	}
}

func (dsl *dslCollection) getColorBlueChannel(img image.Image, x, y int) (b uint32) {
	if !(image.Point{x, y}.In(img.Bounds())) {
		return
	}
	switch t := img.(type) {
	case *image.NRGBA:
		return uint32(t.Pix[t.PixOffset(x, y)+2])
	case *image.NRGBA64:
		i := t.PixOffset(x, y)
		return uint32(t.Pix[i+4])<<8 | uint32(t.Pix[i+5])
	case *image.RGBA:
		return uint32(t.Pix[t.PixOffset(x, y)+2])
	case *image.RGBA64:
		i := t.PixOffset(x, y)
		return uint32(t.Pix[i+4])<<8 | uint32(t.Pix[i+5])
	default:
		panic(fmt.Sprintf("image type is unsupported: %T", t))
	}
}

func (dsl *dslCollection) getColorAlphaChannel(img image.Image, x, y int) (a uint32) {
	if !(image.Point{x, y}.In(img.Bounds())) {
		return
	}
	switch t := img.(type) {
	case *image.NRGBA:
		return uint32(t.Pix[t.PixOffset(x, y)+3])
	case *image.NRGBA64:
		i := t.PixOffset(x, y)
		return uint32(t.Pix[i+6])<<8 | uint32(t.Pix[i+7])
	case *image.RGBA:
		return uint32(t.Pix[t.PixOffset(x, y)+3])
	case *image.RGBA64:
		i := t.PixOffset(x, y)
		return uint32(t.Pix[i+6])<<8 | uint32(t.Pix[i+7])
	default:
		panic(fmt.Sprintf("image type is unsupported: %T", t))
	}
}

func (dsl *dslCollection) setColor(img image.Image, x, y int, r, g, b, a uint32) {
	switch t := img.(type) {
	case *image.NRGBA:
		if !(image.Point{x, y}.In(t.Rect)) {
			return
		}
		i := t.PixOffset(x, y)
		s := t.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857

		// convert NRGBA -> RGBA
		r |= r << 8
		r *= a
		r /= 0xff

		g |= g << 8
		g *= a
		g /= 0xff

		b |= b << 8
		b *= a
		b /= 0xff

		a |= a << 8

		// check edge cases
		if a == 0xffff {
			s[0] = uint8(r >> 8)
			s[1] = uint8(g >> 8)
			s[2] = uint8(b >> 8)
			s[3] = 0xff
			return
		}
		if a == 0 {
			s[0] = 0
			s[1] = 0
			s[2] = 0
			s[3] = 0
			return
		}

		// Since the color is an alpha-premultiplied color, we should have r <= a && g <= a && b <= a.
		r = (r * 0xffff) / a
		g = (g * 0xffff) / a
		b = (b * 0xffff) / a

		s[0] = uint8(r >> 8)
		s[1] = uint8(g >> 8)
		s[2] = uint8(b >> 8)
		s[3] = uint8(a >> 8)
	case *image.NRGBA64:
		if !(image.Point{x, y}.In(t.Rect)) {
			return
		}
		i := t.PixOffset(x, y)
		s := t.Pix[i : i+8 : i+8] // Small cap improves performance, see https://golang.org/issue/27857
		s[0] = uint8(r >> 8)
		s[1] = uint8(r)
		s[2] = uint8(g >> 8)
		s[3] = uint8(g)
		s[4] = uint8(b >> 8)
		s[5] = uint8(b)
		s[6] = uint8(a >> 8)
		s[7] = uint8(a)
	case *image.RGBA:
		if !(image.Point{x, y}.In(t.Rect)) {
			return
		}
		i := t.PixOffset(x, y)
		s := t.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
		s[0] = uint8(r)
		s[1] = uint8(g)
		s[2] = uint8(b)
		s[3] = uint8(a)
	case *image.RGBA64:
		if !(image.Point{x, y}.In(t.Rect)) {
			return
		}
		i := t.PixOffset(x, y)
		s := t.Pix[i : i+8 : i+8] // Small cap improves performance, see https://golang.org/issue/27857
		s[0] = uint8(r >> 8)
		s[1] = uint8(r)
		s[2] = uint8(g >> 8)
		s[3] = uint8(g)
		s[4] = uint8(b >> 8)
		s[5] = uint8(b)
		s[6] = uint8(a >> 8)
		s[7] = uint8(a)
	default:
		panic(fmt.Sprintf("image type is unsupported: %T", t))
	}
}

func (dsl *dslCollection) setColorRedChannel(img image.Image, x, y int, r uint32) {
	if !(image.Point{x, y}.In(img.Bounds())) {
		return
	}
	switch t := img.(type) {
	case *image.NRGBA:
		t.Pix[t.PixOffset(x, y)] = uint8(r)
	case *image.NRGBA64:
		i := t.PixOffset(x, y)
		t.Pix[i] = uint8(r >> 8)
		t.Pix[i+1] = uint8(r)
	case *image.RGBA:
		t.Pix[t.PixOffset(x, y)] = uint8(r)
	case *image.RGBA64:
		i := t.PixOffset(x, y)
		t.Pix[i] = uint8(r >> 8)
		t.Pix[i+1] = uint8(r)
	default:
		panic(fmt.Sprintf("image type is unsupported: %T", t))
	}
}

func (dsl *dslCollection) setColorGreenChannel(img image.Image, x, y int, g uint32) {
	if !(image.Point{x, y}.In(img.Bounds())) {
		return
	}
	switch t := img.(type) {
	case *image.NRGBA:
		t.Pix[t.PixOffset(x, y)+1] = uint8(g)
	case *image.NRGBA64:
		i := t.PixOffset(x, y)
		t.Pix[i+2] = uint8(g >> 8)
		t.Pix[i+3] = uint8(g)
	case *image.RGBA:
		t.Pix[t.PixOffset(x, y)+1] = uint8(g)
	case *image.RGBA64:
		i := t.PixOffset(x, y)
		t.Pix[i+2] = uint8(g >> 8)
		t.Pix[i+3] = uint8(g)
	default:
		panic(fmt.Sprintf("image type is unsupported: %T", t))
	}
}

func (dsl *dslCollection) setColorBlueChannel(img image.Image, x, y int, b uint32) {
	if !(image.Point{x, y}.In(img.Bounds())) {
		return
	}
	switch t := img.(type) {
	case *image.NRGBA:
		t.Pix[t.PixOffset(x, y)+2] = uint8(b)
	case *image.NRGBA64:
		i := t.PixOffset(x, y)
		t.Pix[i+4] = uint8(b >> 8)
		t.Pix[i+5] = uint8(b)
	case *image.RGBA:
		t.Pix[t.PixOffset(x, y)+2] = uint8(b)
	case *image.RGBA64:
		i := t.PixOffset(x, y)
		t.Pix[i+4] = uint8(b >> 8)
		t.Pix[i+5] = uint8(b)
	default:
		panic(fmt.Sprintf("image type is unsupported: %T", t))
	}
}

func (dsl *dslCollection) setColorAlphaChannel(img image.Image, x, y int, a uint32) {
	if !(image.Point{x, y}.In(img.Bounds())) {
		return
	}
	switch t := img.(type) {
	case *image.NRGBA:
		t.Pix[t.PixOffset(x, y)+3] = uint8(a)
	case *image.NRGBA64:
		i := t.PixOffset(x, y)
		t.Pix[i+6] = uint8(a >> 8)
		t.Pix[i+7] = uint8(a)
	case *image.RGBA:
		t.Pix[t.PixOffset(x, y)+3] = uint8(a)
	case *image.RGBA64:
		i := t.PixOffset(x, y)
		t.Pix[i+6] = uint8(a >> 8)
		t.Pix[i+7] = uint8(a)
	default:
		panic(fmt.Sprintf("image type is unsupported: %T", t))
	}
}

func (dsl *dslCollection) parallelProcessNRGBA64(img image.Image, processor dslPixelProcessor, numWorkers int) (result *image.NRGBA64) {
	return dslParallelProcessImage[*image.NRGBA64](img, processor, numWorkers).(*image.NRGBA64)
}

func (dsl *dslCollection) parallelProcessNRGBA(img image.Image, processor dslPixelProcessor, numWorkers int) (result *image.NRGBA) {
	return dslParallelProcessImage[*image.NRGBA](img, processor, numWorkers).(*image.NRGBA)
}

func (dsl *dslCollection) parallelProcessRGBA64(img image.Image, processor dslPixelProcessor, numWorkers int) (result *image.RGBA64) {
	return dslParallelProcessImage[*image.RGBA64](img, processor, numWorkers).(*image.RGBA64)
}

func (dsl *dslCollection) parallelProcessRGBA(img image.Image, processor dslPixelProcessor, numWorkers int) (result *image.RGBA) {
	return dslParallelProcessImage[*image.RGBA](img, processor, numWorkers).(*image.RGBA)
}
