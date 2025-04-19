package main

import (
	"fmt"
	"image"
	"image/color"
	"reflect"
	"strconv"
	"strings"

	"github.com/toxyl/math"
)

// cast attempts to convert a value to the target type
func (dsl *dslCollection) cast(value any, targetType string) (any, error) {
	if value == nil {
		return nil, dsl.errors.NIL_CAST()
	}

	// Validate input type
	switch value.(type) {
	case *image.NRGBA, *image.RGBA, *image.RGBA64, *image.NRGBA64:
		return dsl.castImage(value, targetType)
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, string:
		// These types are supported
	default:
		return nil, dsl.errors.UNSUPPORTED_SOURCE_TYPE(value)
	}

	// Handle string inputs
	if str, ok := value.(string); ok {
		// If target type is string, return the string as is
		if targetType == "string" {
			return str, nil
		}

		str = strings.TrimSpace(strings.ToLower(str))

		// Try parsing as bool
		if b, err := strconv.ParseBool(str); err == nil {
			return dsl.castToType(b, targetType)
		}

		// Try parsing as number
		if f, err := strconv.ParseFloat(str, 64); err == nil {
			return dsl.castToType(f, targetType)
		}

		return nil, dsl.errors.STRING_CAST(str, targetType)
	}

	return dsl.castToType(value, targetType)
}

func (dsl *dslCollection) castToType(value any, targetType string) (any, error) {
	switch targetType {
	case "bool", "float32", "float64", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "string":
		// Check if the value is already of the target type
		typeStr := reflect.TypeOf(value).String()
		if typeStr == targetType || strings.HasSuffix(typeStr, "."+targetType) {
			return value, nil
		}

		// Handle string target type specially
		if targetType == "string" {
			switch v := value.(type) {
			case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
				return fmt.Sprintf("%v", v), nil
			}
		}

		// Convert the value to the target type
		switch targetType {
		case "bool":
			switch v := value.(type) {
			case bool:
				return v, nil
			case int, int8, int16, int32, int64:
				return reflect.ValueOf(value).Int() != 0, nil
			case uint, uint8, uint16, uint32, uint64:
				return reflect.ValueOf(value).Uint() != 0, nil
			case float32, float64:
				return reflect.ValueOf(value).Float() != 0, nil
			}
		case "int", "int8", "int16", "int32", "int64":
			var result int64
			switch v := value.(type) {
			case bool:
				if v {
					result = 1
				}
			case int:
				result = int64(v)
			case int8:
				result = int64(v)
			case int16:
				result = int64(v)
			case int32:
				result = int64(v)
			case int64:
				result = v
			case uint:
				result = int64(v)
			case uint8:
				result = int64(v)
			case uint16:
				result = int64(v)
			case uint32:
				result = int64(v)
			case uint64:
				if v > math.MaxInt64 {
					return math.MaxInt64, nil
				}
				result = int64(v)
			case float32, float64:
				f := reflect.ValueOf(value).Float()
				if math.IsNaN(f) {
					return 0, nil
				}
				if math.IsInf(f, 1) {
					switch targetType {
					case "int8":
						return int8(math.MaxInt8), nil
					case "int16":
						return int16(math.MaxInt16), nil
					case "int32":
						return int32(math.MaxInt32), nil
					case "int64":
						return math.MaxInt64, nil
					default: // "int"
						return int(math.MaxInt), nil
					}
				}
				if math.IsInf(f, -1) {
					switch targetType {
					case "int8":
						return int8(math.MinInt8), nil
					case "int16":
						return int16(math.MinInt16), nil
					case "int32":
						return int32(math.MinInt32), nil
					case "int64":
						return math.MinInt64, nil
					default: // "int"
						return int(math.MinInt), nil
					}
				}
				// For very large float values, clamp to the target type's range
				switch targetType {
				case "int8":
					if f > float64(math.MaxInt8) {
						return int8(math.MaxInt8), nil
					}
					if f < float64(math.MinInt8) {
						return int8(math.MinInt8), nil
					}
					return int8(f), nil
				case "int16":
					if f > float64(math.MaxInt16) {
						return int16(math.MaxInt16), nil
					}
					if f < float64(math.MinInt16) {
						return int16(math.MinInt16), nil
					}
					return int16(f), nil
				case "int32":
					if f > float64(math.MaxInt32) {
						return int32(math.MaxInt32), nil
					}
					if f < float64(math.MinInt32) {
						return int32(math.MinInt32), nil
					}
					return int32(f), nil
				case "int64":
					if f > float64(math.MaxInt64) {
						return math.MaxInt64, nil
					}
					if f < float64(math.MinInt64) {
						return math.MinInt64, nil
					}
					return int64(f), nil
				default: // "int"
					if f > float64(math.MaxInt) {
						return int(math.MaxInt), nil
					}
					if f < float64(math.MinInt) {
						return int(math.MinInt), nil
					}
					return int(f), nil
				}
			}

			// Check for overflow based on target type
			switch targetType {
			case "int8":
				if result > math.MaxInt8 {
					return int8(math.MaxInt8), nil
				}
				if result < math.MinInt8 {
					return int8(math.MinInt8), nil
				}
				return int8(result), nil
			case "int16":
				if result > math.MaxInt16 {
					return int16(math.MaxInt16), nil
				}
				if result < math.MinInt16 {
					return int16(math.MinInt16), nil
				}
				return int16(result), nil
			case "int32":
				if result > math.MaxInt32 {
					return int32(math.MaxInt32), nil
				}
				if result < math.MinInt32 {
					return int32(math.MinInt32), nil
				}
				return int32(result), nil
			case "int64":
				return result, nil
			default: // "int"
				if result > int64(math.MaxInt) {
					return int(math.MaxInt), nil
				}
				if result < int64(math.MinInt) {
					return int(math.MinInt), nil
				}
				return int(result), nil
			}
		case "uint", "uint8", "uint16", "uint32", "uint64":
			var result uint64
			switch v := value.(type) {
			case bool:
				if v {
					result = 1
				}
			case int, int8, int16, int32, int64:
				i := reflect.ValueOf(value).Int()
				if i < 0 {
					result = uint64(0)
				} else {
					result = uint64(i)
				}
			case uint:
				result = uint64(v)
			case uint8:
				result = uint64(v)
			case uint16:
				result = uint64(v)
			case uint32:
				result = uint64(v)
			case uint64:
				result = v
			case float32, float64:
				f := reflect.ValueOf(value).Float()
				if math.IsNaN(f) || f < 0 {
					return uint64(0), nil
				}
				if math.IsInf(f, 1) {
					switch targetType {
					case "uint8":
						return uint8(math.MaxUint8), nil
					case "uint16":
						return uint16(math.MaxUint16), nil
					case "uint32":
						return uint32(math.MaxUint32), nil
					case "uint64":
						return uint64(math.MaxUint64), nil
					default: // "uint"
						return uint(math.MaxUint), nil
					}
				}
				result = uint64(f)
			}

			// Check for overflow based on target type
			switch targetType {
			case "uint8":
				if result > math.MaxUint8 {
					return uint8(math.MaxUint8), nil
				}
				return uint8(result), nil
			case "uint16":
				if result > math.MaxUint16 {
					return uint16(math.MaxUint16), nil
				}
				return uint16(result), nil
			case "uint32":
				if result > math.MaxUint32 {
					return uint32(math.MaxUint32), nil
				}
				return uint32(result), nil
			case "uint64":
				return result, nil
			default: // "uint"
				if result > math.MaxUint {
					return uint(math.MaxUint), nil
				}
				return uint(result), nil
			}
		case "float32", "float64":
			var result float64
			switch v := value.(type) {
			case bool:
				if v {
					result = 1
				}
			case int, int8, int16, int32, int64:
				result = float64(reflect.ValueOf(value).Int())
			case uint, uint8, uint16, uint32, uint64:
				result = float64(reflect.ValueOf(value).Uint())
			case float32:
				result = float64(v)
			case float64:
				result = v
			}

			if targetType == "float32" {
				if math.IsInf(result, 0) || math.IsNaN(result) {
					return float32(result), nil
				}
				if result > math.MaxFloat32 {
					return float32(math.MaxFloat32), nil
				}
				if result < -math.MaxFloat32 {
					return float32(-math.MaxFloat32), nil
				}
				return float32(result), nil
			}
			return result, nil
		}
	}
	return nil, dsl.errors.UNSUPPORTED_TARGET_TYPE(targetType)
}

// castImage handles conversions between different image types
func (dsl *dslCollection) castImage(value any, targetType string) (any, error) {
	switch v := value.(type) {
	case *image.NRGBA:
		switch targetType {
		case "*image.NRGBA":
			return v, nil
		case "*image.RGBA":
			return dsl.convertNRGBAToRGBA(v), nil
		case "*image.NRGBA64":
			return dsl.convertNRGBAToNRGBA64(v), nil
		case "*image.RGBA64":
			return dsl.convertNRGBAToRGBA64(v), nil
		}
	case *image.RGBA:
		switch targetType {
		case "*image.NRGBA":
			return dsl.convertRGBAToNRGBA(v), nil
		case "*image.RGBA":
			return v, nil
		case "*image.NRGBA64":
			return dsl.convertRGBAToNRGBA64(v), nil
		case "*image.RGBA64":
			return dsl.convertRGBAToRGBA64(v), nil
		}
	case *image.RGBA64:
		switch targetType {
		case "*image.RGBA64":
			return v, nil
		case "*image.NRGBA64":
			return dsl.convertRGBA64ToNRGBA64(v), nil
		case "*image.RGBA":
			return dsl.convertRGBA64ToRGBA(v), nil
		case "*image.NRGBA":
			return dsl.convertRGBA64ToNRGBA(v), nil
		}
	case *image.NRGBA64:
		switch targetType {
		case "*image.RGBA64":
			return dsl.convertNRGBA64ToRGBA64(v), nil
		case "*image.NRGBA64":
			return v, nil
		case "*image.RGBA":
			return dsl.convertNRGBA64ToRGBA(v), nil
		case "*image.NRGBA":
			return dsl.convertNRGBA64ToNRGBA(v), nil
		}
	}
	return nil, dsl.errors.CAST_NOT_POSSIBLE(reflect.TypeOf(value).String(), targetType)
}

// Helper functions for image conversions
func (dsl *dslCollection) convertNRGBAToRGBA(src *image.NRGBA) *image.RGBA {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := src.NRGBAAt(x, y)
			// Convert to premultiplied alpha
			dst.Set(x, y, color.RGBA{
				R: uint8(uint32(c.R) * uint32(c.A) / 255),
				G: uint8(uint32(c.G) * uint32(c.A) / 255),
				B: uint8(uint32(c.B) * uint32(c.A) / 255),
				A: c.A,
			})
		}
	}
	return dst
}

func (dsl *dslCollection) convertRGBAToNRGBA(src *image.RGBA) *image.NRGBA {
	bounds := src.Bounds()
	dst := image.NewNRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := src.RGBAAt(x, y)
			if c.A == 0 {
				dst.Set(x, y, color.NRGBA{0, 0, 0, 0})
				continue
			}
			// Convert from premultiplied to non-premultiplied alpha
			dst.Set(x, y, color.NRGBA{
				R: uint8(uint32(c.R) * 255 / uint32(c.A)),
				G: uint8(uint32(c.G) * 255 / uint32(c.A)),
				B: uint8(uint32(c.B) * 255 / uint32(c.A)),
				A: c.A,
			})
		}
	}
	return dst
}

func (dsl *dslCollection) convertRGBA64ToNRGBA64(src *image.RGBA64) *image.NRGBA64 {
	bounds := src.Bounds()
	dst := image.NewNRGBA64(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := src.RGBA64At(x, y)
			if c.A == 0 {
				dst.Set(x, y, color.NRGBA64{0, 0, 0, 0})
				continue
			}
			// Convert from premultiplied to non-premultiplied alpha
			dst.Set(x, y, color.NRGBA64{
				R: uint16((uint32(c.R) * 0xffff) / uint32(c.A)),
				G: uint16((uint32(c.G) * 0xffff) / uint32(c.A)),
				B: uint16((uint32(c.B) * 0xffff) / uint32(c.A)),
				A: c.A,
			})
		}
	}
	return dst
}

func (dsl *dslCollection) convertNRGBA64ToRGBA64(src *image.NRGBA64) *image.RGBA64 {
	bounds := src.Bounds()
	dst := image.NewRGBA64(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := src.NRGBA64At(x, y)
			// Convert to premultiplied alpha
			dst.Set(x, y, color.RGBA64{
				R: uint16((uint32(c.R) * uint32(c.A)) / 0xffff),
				G: uint16((uint32(c.G) * uint32(c.A)) / 0xffff),
				B: uint16((uint32(c.B) * uint32(c.A)) / 0xffff),
				A: c.A,
			})
		}
	}
	return dst
}

// convertNRGBAToNRGBA64 converts an 8-bit non-premultiplied RGBA image to 16-bit
func (dsl *dslCollection) convertNRGBAToNRGBA64(src *image.NRGBA) *image.NRGBA64 {
	bounds := src.Bounds()
	dst := image.NewNRGBA64(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := src.NRGBAAt(x, y)
			dst.Set(x, y, color.NRGBA64{
				R: uint16(c.R) * 257, // Scale from 8-bit to 16-bit (257 = 65535/255)
				G: uint16(c.G) * 257,
				B: uint16(c.B) * 257,
				A: uint16(c.A) * 257,
			})
		}
	}
	return dst
}

// convertRGBAToRGBA64 converts an 8-bit premultiplied RGBA image to 16-bit
func (dsl *dslCollection) convertRGBAToRGBA64(src *image.RGBA) *image.RGBA64 {
	bounds := src.Bounds()
	dst := image.NewRGBA64(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := src.RGBAAt(x, y)
			dst.Set(x, y, color.RGBA64{
				R: uint16(c.R) * 257,
				G: uint16(c.G) * 257,
				B: uint16(c.B) * 257,
				A: uint16(c.A) * 257,
			})
		}
	}
	return dst
}

// convertRGBAToNRGBA64 converts an 8-bit premultiplied RGBA image to 16-bit non-premultiplied
func (dsl *dslCollection) convertRGBAToNRGBA64(src *image.RGBA) *image.NRGBA64 {
	bounds := src.Bounds()
	dst := image.NewNRGBA64(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := src.RGBAAt(x, y)
			if c.A == 0 {
				dst.Set(x, y, color.NRGBA64{0, 0, 0, 0})
				continue
			}
			// First unpremultiply 8-bit values, then convert to 16-bit
			dst.Set(x, y, color.NRGBA64{
				R: uint16((uint32(c.R) * 0xff * 257) / uint32(c.A)),
				G: uint16((uint32(c.G) * 0xff * 257) / uint32(c.A)),
				B: uint16((uint32(c.B) * 0xff * 257) / uint32(c.A)),
				A: uint16(c.A) * 257,
			})
		}
	}
	return dst
}

// convertNRGBA64ToNRGBA converts a 16-bit non-premultiplied RGBA image to 8-bit
func (dsl *dslCollection) convertNRGBA64ToNRGBA(src *image.NRGBA64) *image.NRGBA {
	bounds := src.Bounds()
	dst := image.NewNRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := src.NRGBA64At(x, y)
			dst.Set(x, y, color.NRGBA{
				R: uint8(c.R >> 8), // Scale from 16-bit to 8-bit
				G: uint8(c.G >> 8),
				B: uint8(c.B >> 8),
				A: uint8(c.A >> 8),
			})
		}
	}
	return dst
}

// convertRGBA64ToRGBA converts a 16-bit premultiplied RGBA image to 8-bit
func (dsl *dslCollection) convertRGBA64ToRGBA(src *image.RGBA64) *image.RGBA {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := src.RGBA64At(x, y)
			dst.Set(x, y, color.RGBA{
				R: uint8(c.R >> 8),
				G: uint8(c.G >> 8),
				B: uint8(c.B >> 8),
				A: uint8(c.A >> 8),
			})
		}
	}
	return dst
}

// convertRGBA64ToNRGBA converts a 16-bit premultiplied RGBA image to 8-bit non-premultiplied
func (dsl *dslCollection) convertRGBA64ToNRGBA(src *image.RGBA64) *image.NRGBA {
	bounds := src.Bounds()
	dst := image.NewNRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := src.RGBA64At(x, y)
			if c.A == 0 {
				dst.Set(x, y, color.NRGBA{0, 0, 0, 0})
				continue
			}
			// Unpremultiply and convert to 8-bit
			dst.Set(x, y, color.NRGBA{
				R: uint8((uint32(c.R) * 0xffff) / uint32(c.A) >> 8),
				G: uint8((uint32(c.G) * 0xffff) / uint32(c.A) >> 8),
				B: uint8((uint32(c.B) * 0xffff) / uint32(c.A) >> 8),
				A: uint8(c.A >> 8),
			})
		}
	}
	return dst
}

// convertNRGBA64ToRGBA converts a 16-bit non-premultiplied RGBA image to 8-bit premultiplied
func (dsl *dslCollection) convertNRGBA64ToRGBA(src *image.NRGBA64) *image.RGBA {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := src.NRGBA64At(x, y)
			// Convert to 8-bit and premultiply
			a8 := uint8(c.A >> 8)
			dst.Set(x, y, color.RGBA{
				R: uint8((uint32(c.R) * uint32(a8)) / 0xff >> 8),
				G: uint8((uint32(c.G) * uint32(a8)) / 0xff >> 8),
				B: uint8((uint32(c.B) * uint32(a8)) / 0xff >> 8),
				A: a8,
			})
		}
	}
	return dst
}

// convertNRGBAToRGBA64 converts an 8-bit non-premultiplied RGBA image to 16-bit premultiplied
func (dsl *dslCollection) convertNRGBAToRGBA64(src *image.NRGBA) *image.RGBA64 {
	bounds := src.Bounds()
	dst := image.NewRGBA64(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := src.NRGBAAt(x, y)
			// First convert to 16-bit, then premultiply
			r16 := uint16(c.R) * 257
			g16 := uint16(c.G) * 257
			b16 := uint16(c.B) * 257
			a16 := uint16(c.A) * 257

			dst.Set(x, y, color.RGBA64{
				R: uint16((uint32(r16) * uint32(a16)) / 0xffff),
				G: uint16((uint32(g16) * uint32(a16)) / 0xffff),
				B: uint16((uint32(b16) * uint32(a16)) / 0xffff),
				A: a16,
			})
		}
	}
	return dst
}
