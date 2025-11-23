package main

import "reflect"

type dslFnMeta struct {
	name    string
	desc    string
	params  []dslParamMeta
	returns []dslParamMeta
}

type dslParamMeta struct {
	name string
	typ  string
	min  any
	max  any
	def  any
	unit string
	desc string
}

type dslFnType struct {
	meta dslFnMeta
	data func(...any) (any, error)
}

func (fn *dslFnType) validate(args ...any) error {
	if len(args) != len(fn.meta.params) {
		if len(args) < len(fn.meta.params) {
			for i := len(args); i < len(fn.meta.params); i++ {
				args = append(args, fn.meta.params[i].def)
			}
		}
	}

	for i, param := range fn.meta.params {
		arg := args[i]
		switch param.typ {
		case "int":
			val, ok := arg.(int)
			if !ok {
				return errors.REG_VALIDATION_WRONG_TYPE("parameter", param.name, "int", arg)
			}
			if min, ok := param.min.(int); ok && val < min {
				return errors.REG_VALIDATION_OUT_OF_BOUNDS("parameter", param.name, param.min, param.max, val)
			}
			if max, ok := param.max.(int); ok && val > max {
				return errors.REG_VALIDATION_OUT_OF_BOUNDS("parameter", param.name, param.min, param.max, val)
			}
		case "float":
			val, ok := arg.(float64)
			if !ok {
				return errors.REG_VALIDATION_WRONG_TYPE("parameter", param.name, "float64", arg)
			}
			if min, ok := param.min.(float64); ok && val < min {
				return errors.REG_VALIDATION_OUT_OF_BOUNDS("parameter", param.name, param.min, param.max, val)
			}
			if max, ok := param.max.(float64); ok && val > max {
				return errors.REG_VALIDATION_OUT_OF_BOUNDS("parameter", param.name, param.min, param.max, val)
			}
		case "bool":
			_, ok := arg.(bool)
			if !ok {
				return errors.REG_VALIDATION_WRONG_TYPE("parameter", param.name, "bool", arg)
			}
		case "string":
			val, ok := arg.(string)
			if !ok {
				return errors.REG_VALIDATION_WRONG_TYPE("parameter", param.name, "string", arg)
			}
			if min, ok := param.min.(int); ok && len(val) < min {
				return errors.REG_VALIDATION_OUT_OF_BOUNDS_LENGTH("parameter", param.name, param.min, param.max, val)
			}
			if max, ok := param.max.(int); ok && len(val) > max {
				return errors.REG_VALIDATION_OUT_OF_BOUNDS_LENGTH("parameter", param.name, param.min, param.max, val)
			}

		}
	}
	return nil
}

func (f *dslFnType) call(vars *dslVarRegistry, args ...any) (any, error) {
	// Make a copy of args to avoid modifying the original
	callArgs := make([]any, len(args))
	copy(callArgs, args)

	// Handle variable references and type conversions
	for i, arg := range callArgs {
		if str, ok := arg.(string); ok {
			// Check if it's a variable reference
			if vars.has(str) {
				// Get the variable value in a thread-safe way
				varVal := vars.get(str)
				if varVal != nil {
					callArgs[i] = varVal.get()
				}
			}
		}

		// Handle type conversions
		if f.meta.params[i].typ != "" && f.meta.params[i].typ != "any" {
			// Check if types already match exactly before casting
			argType := reflect.TypeOf(callArgs[i])
			if argType != nil && argType.String() == f.meta.params[i].typ {
				// Types match exactly, no conversion needed
				continue
			}
			converted, err := dsl.cast(callArgs[i], f.meta.params[i].typ)
			if err != nil {
				return nil, err
			}
			callArgs[i] = converted
		} else {
			callArgs[i] = arg
		}
	}

	// Validate arguments
	if err := f.validate(callArgs...); err != nil {
		return nil, err
	}

	// Call the function
	return f.data(callArgs...)
}
