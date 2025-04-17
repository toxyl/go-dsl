package main

import (
	"reflect"

	"github.com/toxyl/math"
)

// deepEqual is a utility function that extends reflect.DeepEqual with special handling
// for edge cases like NaN values in floating-point comparisons and numeric type conversions.
func (dsl *dslCollection) deepEqual(x, y interface{}) bool {
	// Handle nil cases
	if x == nil || y == nil {
		return x == y
	}

	// Get the types of both values
	v1 := reflect.ValueOf(x)
	v2 := reflect.ValueOf(y)

	// Special handling for numeric types
	if dsl.isNumeric(v1.Kind()) && dsl.isNumeric(v2.Kind()) {
		// Convert both values to float64 for comparison
		f1 := dsl.toFloat64(v1)
		f2 := dsl.toFloat64(v2)

		// Handle NaN cases
		if math.IsNaN(f1) && math.IsNaN(f2) {
			return true
		}

		// Handle infinity cases
		if math.IsInf(f1, 1) && math.IsInf(f2, 1) {
			return true
		}
		if math.IsInf(f1, -1) && math.IsInf(f2, -1) {
			return true
		}

		// For non-special cases, compare the float64 values
		return f1 == f2
	}

	// For all other types, use reflect.DeepEqual
	return reflect.DeepEqual(x, y)
}

// isNumeric returns true if the kind is a numeric type
func (dsl *dslCollection) isNumeric(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	}
	return false
}

// toFloat64 converts any numeric value to float64
func (dsl *dslCollection) toFloat64(v reflect.Value) float64 {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint())
	case reflect.Float32, reflect.Float64:
		return v.Float()
	}
	return 0
}
