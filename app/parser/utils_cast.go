package main

import (
	"fmt"
	"image"
	"image/color"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/toxyl/math"
)

var (
	NumColorConversionWorkers = runtime.NumCPU()
)

// castSelfOnly handles casting for types that can only be cast to themselves
func castSelfOnly[T any](value T, targetType, typeName string) (any, error) {
	if targetType == typeName || targetType == "language."+typeName {
		return value, nil
	}
	if targetType == "*"+typeName || targetType == "*language."+typeName {
		return value, nil
	}
	return nil, errors.CAST_NOT_POSSIBLE(reflect.TypeOf(value).String(), targetType)
}

// cast attempts to convert a value to the target type
func (dsl *dslCollection) cast(value any, targetType string) (any, error) {
	if value == nil {
		return nil, errors.NIL_CAST()
	}

	// Handle slice types before the main switch
	if strings.HasPrefix(targetType, "[]") {
		valType := reflect.TypeOf(value)
		if valType != nil && valType.Kind() == reflect.Slice {
			// Target is a slice, source is also a slice
			if targetType == "[]any" {
				// Convert any slice to []any
				val := reflect.ValueOf(value)
				result := make([]any, val.Len())
				for i := 0; i < val.Len(); i++ {
					result[i] = val.Index(i).Interface()
				}
				return result, nil
			}
			// Slice-to-slice conversions (except []any) are not allowed
			return nil, errors.CAST_NOT_POSSIBLE(reflect.TypeOf(value).String(), targetType)
		}
	}

	// Validate input type
	switch value := value.(type) {
	// These types are supported
	case *image.NRGBA, *image.RGBA, *image.RGBA64, *image.NRGBA64:
		return dsl.castImage(value, targetType)
	case color.RGBA, color.RGBA64:
		return dsl.castColor(value, targetType)
	case Point:
		return castSelfOnly(value, targetType, "Point")
	case Rect:
		return castSelfOnly(value, targetType, "Rect")
	case NGon:
		return castSelfOnly(value, targetType, "NGon")
	case Triangle:
		return castSelfOnly(value, targetType, "Triangle")
	case Quad:
		return castSelfOnly(value, targetType, "Quad")
	case Ellipse:
		return castSelfOnly(value, targetType, "Ellipse")
	case Vector:
		return castSelfOnly(value, targetType, "Vector")
	case Text:
		return castSelfOnly(value, targetType, "Text")
	case LineStyle:
		return castSelfOnly(value, targetType, "LineStyle")
	case FillStyle:
		return castSelfOnly(value, targetType, "FillStyle")
	case TextStyle:
		return castSelfOnly(value, targetType, "TextStyle")
	// TODO: NEW TYPES: add additional types
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, string:
	default:
		// If we get here the type is not supported
		return nil, errors.UNSUPPORTED_SOURCE_TYPE(value)
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

		return nil, errors.STRING_CAST(str, targetType)
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
	return nil, errors.UNSUPPORTED_TARGET_TYPE(targetType)
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
	return nil, errors.CAST_NOT_POSSIBLE(reflect.TypeOf(value).String(), targetType)
}

// Helper functions for image conversions
func (dsl *dslCollection) convertNRGBAToRGBA(src *image.NRGBA) *image.RGBA {
	return dsl.parallelProcessRGBA(src, func(r1, g1, b1, a1 uint32) (r uint32, g uint32, b uint32, a uint32) {
		r = r1 * a1 / 255
		g = g1 * a1 / 255
		b = b1 * a1 / 255
		a = a1
		return
	}, NumColorConversionWorkers)
}

func (dsl *dslCollection) convertRGBAToNRGBA(src *image.RGBA) *image.NRGBA {
	return dsl.parallelProcessNRGBA(src, func(r1, g1, b1, a1 uint32) (r uint32, g uint32, b uint32, a uint32) {
		if a1 == 0 {
			return
		}
		r = r1 * 255 / a1
		g = g1 * 255 / a1
		b = b1 * 255 / a1
		a = a1
		return
	}, NumColorConversionWorkers)
}

func (dsl *dslCollection) convertRGBA64ToNRGBA64(src *image.RGBA64) *image.NRGBA64 {
	return dsl.parallelProcessNRGBA64(src, func(r1, g1, b1, a1 uint32) (r uint32, g uint32, b uint32, a uint32) {
		if a1 == 0 {
			return
		}
		r = (r1 * 0xFFFF) / a1
		g = (g1 * 0xFFFF) / a1
		b = (b1 * 0xFFFF) / a1
		a = a1
		return
	}, NumColorConversionWorkers)
}

func (dsl *dslCollection) convertNRGBA64ToRGBA64(src *image.NRGBA64) *image.RGBA64 {
	return dsl.parallelProcessRGBA64(src, func(r1, g1, b1, a1 uint32) (r uint32, g uint32, b uint32, a uint32) {
		if a1 == 0 {
			return
		}
		r = (r1 * a1) / 0xFFFF
		g = (g1 * a1) / 0xFFFF
		b = (b1 * a1) / 0xFFFF
		a = a1
		return
	}, NumColorConversionWorkers)
}

// convertNRGBAToNRGBA64 converts an 8-bit non-premultiplied RGBA image to 16-bit
func (dsl *dslCollection) convertNRGBAToNRGBA64(src *image.NRGBA) *image.NRGBA64 {
	return dsl.parallelProcessNRGBA64(src, func(r1, g1, b1, a1 uint32) (r uint32, g uint32, b uint32, a uint32) {
		if a1 == 0 {
			return
		}
		r = r1 * 257 // Scale from 8-bit to 16-bit (257 = 65535/255)
		g = g1 * 257
		b = b1 * 257
		a = a1
		return
	}, NumColorConversionWorkers)
}

// convertRGBAToRGBA64 converts an 8-bit premultiplied RGBA image to 16-bit
func (dsl *dslCollection) convertRGBAToRGBA64(src *image.RGBA) *image.RGBA64 {
	return dsl.parallelProcessRGBA64(src, func(r1, g1, b1, a1 uint32) (r uint32, g uint32, b uint32, a uint32) {
		if a1 == 0 {
			return
		}
		r = r1 * 257 // Scale from 8-bit to 16-bit (257 = 65535/255)
		g = g1 * 257
		b = b1 * 257
		a = a1
		return
	}, NumColorConversionWorkers)
}

// convertRGBAToNRGBA64 converts an 8-bit premultiplied RGBA image to 16-bit non-premultiplied
func (dsl *dslCollection) convertRGBAToNRGBA64(src *image.RGBA) *image.NRGBA64 {
	return dsl.parallelProcessNRGBA64(src, func(r1, g1, b1, a1 uint32) (r uint32, g uint32, b uint32, a uint32) {
		if a1 == 0 {
			return
		}
		r = (r1 * 0xFFFF) / a1
		g = (g1 * 0xFFFF) / a1
		b = (b1 * 0xFFFF) / a1
		a = a1 * 257
		return
	}, NumColorConversionWorkers)
}

// convertNRGBA64ToNRGBA converts a 16-bit non-premultiplied RGBA image to 8-bit
func (dsl *dslCollection) convertNRGBA64ToNRGBA(src *image.NRGBA64) *image.NRGBA {
	return dsl.parallelProcessNRGBA(src, func(r1, g1, b1, a1 uint32) (r uint32, g uint32, b uint32, a uint32) {
		if a1 == 0 {
			return
		}
		r = (r1 >> 8)
		g = (g1 >> 8)
		b = (b1 >> 8)
		a = (a1 >> 8)
		return
	}, NumColorConversionWorkers)
}

// convertRGBA64ToRGBA converts a 16-bit premultiplied RGBA image to 8-bit
func (dsl *dslCollection) convertRGBA64ToRGBA(src *image.RGBA64) *image.RGBA {
	return dsl.parallelProcessRGBA(src, func(r1, g1, b1, a1 uint32) (r uint32, g uint32, b uint32, a uint32) {
		if a1 == 0 {
			return
		}
		r = (r1 >> 8)
		g = (g1 >> 8)
		b = (b1 >> 8)
		a = (a1 >> 8)
		return
	}, NumColorConversionWorkers)
}

// convertRGBA64ToNRGBA converts a 16-bit premultiplied RGBA image to 8-bit non-premultiplied
func (dsl *dslCollection) convertRGBA64ToNRGBA(src *image.RGBA64) *image.NRGBA {
	return dsl.parallelProcessNRGBA(src, func(r1, g1, b1, a1 uint32) (r uint32, g uint32, b uint32, a uint32) {
		if a1 == 0 {
			return
		}
		r = (r1 * 0xFFFF) / (a1 >> 8)
		g = (g1 * 0xFFFF) / (a1 >> 8)
		b = (b1 * 0xFFFF) / (a1 >> 8)
		a = a1 >> 8
		return
	}, NumColorConversionWorkers)
}

// convertNRGBA64ToRGBA converts a 16-bit non-premultiplied RGBA image to 8-bit premultiplied
func (dsl *dslCollection) convertNRGBA64ToRGBA(src *image.NRGBA64) *image.RGBA {
	return dsl.parallelProcessRGBA(src, func(r1, g1, b1, a1 uint32) (r uint32, g uint32, b uint32, a uint32) {
		if a1 == 0 {
			return
		}
		a8 := a1 >> 8
		r = (r1 * a8) / 0xFF >> 8
		g = (g1 * a8) / 0xFF >> 8
		b = (b1 * a8) / 0xFF >> 8
		a = a8
		return
	}, NumColorConversionWorkers)
}

// convertNRGBAToRGBA64 converts an 8-bit non-premultiplied RGBA image to 16-bit premultiplied
func (dsl *dslCollection) convertNRGBAToRGBA64(src *image.NRGBA) *image.RGBA64 {
	return dsl.parallelProcessRGBA64(src, func(r1, g1, b1, a1 uint32) (r uint32, g uint32, b uint32, a uint32) {
		if a1 == 0 {
			return
		}
		a16 := a1 * 257
		r = ((r1 * 257) * a16) / 0xFFFF
		g = ((g1 * 257) * a16) / 0xFFFF
		b = ((b1 * 257) * a16) / 0xFFFF
		a = a16
		return
	}, NumColorConversionWorkers)
}

// castColor handles conversions between different color types
func (dsl *dslCollection) castColor(value any, targetType string) (any, error) {
	switch v := value.(type) {
	case color.RGBA:
		switch targetType {
		case "color.RGBA":
			return v, nil
		case "color.RGBA64":
			return color.RGBA64{
				R: uint16(v.R) * 257,
				G: uint16(v.G) * 257,
				B: uint16(v.B) * 257,
				A: uint16(v.A) * 257,
			}, nil
		case "color.NRGBA":
			// Convert from premultiplied to non-premultiplied alpha
			if v.A == 0 {
				return color.NRGBA{0, 0, 0, 0}, nil
			}
			return color.NRGBA{
				R: uint8((uint32(v.R) * 255) / uint32(v.A)),
				G: uint8((uint32(v.G) * 255) / uint32(v.A)),
				B: uint8((uint32(v.B) * 255) / uint32(v.A)),
				A: v.A,
			}, nil
		case "color.NRGBA64":
			// Convert from premultiplied to non-premultiplied alpha
			if v.A == 0 {
				return color.NRGBA64{0, 0, 0, 0}, nil
			}
			r16 := uint16(v.R) * 257
			g16 := uint16(v.G) * 257
			b16 := uint16(v.B) * 257
			a16 := uint16(v.A) * 257
			return color.NRGBA64{
				R: uint16((uint32(r16) * 0xffff) / uint32(a16)),
				G: uint16((uint32(g16) * 0xffff) / uint32(a16)),
				B: uint16((uint32(b16) * 0xffff) / uint32(a16)),
				A: a16,
			}, nil
		}
	case color.RGBA64:
		switch targetType {
		case "color.RGBA":
			return color.RGBA{
				R: uint8(v.R >> 8),
				G: uint8(v.G >> 8),
				B: uint8(v.B >> 8),
				A: uint8(v.A >> 8),
			}, nil
		case "color.RGBA64":
			return v, nil
		case "color.NRGBA":
			// Convert from premultiplied to non-premultiplied alpha
			if v.A == 0 {
				return color.NRGBA{0, 0, 0, 0}, nil
			}
			return color.NRGBA{
				R: uint8((uint32(v.R) * 255) / uint32(v.A>>8)),
				G: uint8((uint32(v.G) * 255) / uint32(v.A>>8)),
				B: uint8((uint32(v.B) * 255) / uint32(v.A>>8)),
				A: uint8(v.A >> 8),
			}, nil
		case "color.NRGBA64":
			// Convert from premultiplied to non-premultiplied alpha
			if v.A == 0 {
				return color.NRGBA64{0, 0, 0, 0}, nil
			}
			return color.NRGBA64{
				R: uint16((uint32(v.R) * 0xffff) / uint32(v.A)),
				G: uint16((uint32(v.G) * 0xffff) / uint32(v.A)),
				B: uint16((uint32(v.B) * 0xffff) / uint32(v.A)),
				A: v.A,
			}, nil
		}
	case color.NRGBA:
		switch targetType {
		case "color.NRGBA":
			return v, nil
		case "color.NRGBA64":
			return color.NRGBA64{
				R: uint16(v.R) * 257,
				G: uint16(v.G) * 257,
				B: uint16(v.B) * 257,
				A: uint16(v.A) * 257,
			}, nil
		case "color.RGBA":
			// Convert from non-premultiplied to premultiplied alpha
			return color.RGBA{
				R: uint8((uint32(v.R) * uint32(v.A)) / 255),
				G: uint8((uint32(v.G) * uint32(v.A)) / 255),
				B: uint8((uint32(v.B) * uint32(v.A)) / 255),
				A: v.A,
			}, nil
		case "color.RGBA64":
			// Convert from non-premultiplied to premultiplied alpha
			r16 := uint16(v.R) * 257
			g16 := uint16(v.G) * 257
			b16 := uint16(v.B) * 257
			a16 := uint16(v.A) * 257
			return color.RGBA64{
				R: uint16((uint32(r16) * uint32(a16)) / 0xffff),
				G: uint16((uint32(g16) * uint32(a16)) / 0xffff),
				B: uint16((uint32(b16) * uint32(a16)) / 0xffff),
				A: a16,
			}, nil
		}
	case color.NRGBA64:
		switch targetType {
		case "color.NRGBA64":
			return v, nil
		case "color.NRGBA":
			return color.NRGBA{
				R: uint8(v.R >> 8),
				G: uint8(v.G >> 8),
				B: uint8(v.B >> 8),
				A: uint8(v.A >> 8),
			}, nil
		case "color.RGBA":
			// Convert from non-premultiplied to premultiplied alpha
			a8 := uint8(v.A >> 8)
			return color.RGBA{
				R: uint8((uint32(v.R) * uint32(a8)) / 0xff >> 8),
				G: uint8((uint32(v.G) * uint32(a8)) / 0xff >> 8),
				B: uint8((uint32(v.B) * uint32(a8)) / 0xff >> 8),
				A: a8,
			}, nil
		case "color.RGBA64":
			// Convert from non-premultiplied to premultiplied alpha
			return color.RGBA64{
				R: uint16((uint32(v.R) * uint32(v.A)) / 0xffff),
				G: uint16((uint32(v.G) * uint32(v.A)) / 0xffff),
				B: uint16((uint32(v.B) * uint32(v.A)) / 0xffff),
				A: v.A,
			}, nil
		}
	}
	return nil, errors.CAST_NOT_POSSIBLE(reflect.TypeOf(value).String(), targetType)
}

// toFloat64 converts any numeric value to float64.
// Supports int, int64, float32, float64, and string numeric types.
// Returns an error if the value cannot be converted to float64.
func (dsl *dslCollection) toFloat64(value any) (float64, error) {
	if value == nil {
		return 0, nil
	}

	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, errors.CAST_NOT_POSSIBLE("string", "float64")
		}
		return f, nil
	default:
		return 0, errors.UNSUPPORTED_SOURCE_TYPE(value)
	}
}

// TODO: NEW TYPES: add additional cast* functions if needed
