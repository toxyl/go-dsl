package main

type dslMetaVar struct {
	name string
	typ  string
	min  any
	max  any
	def  any
	unit string
	desc string
}

type dslMetaVarType struct {
	meta dslMetaVar
	data any
	get  func() any
	set  func(any) error
}

func (v *dslMetaVarType) validate(value any) error {
	switch v.meta.typ {
	case "int":
		val, ok := value.(int)
		if !ok {
			return errors.REG_VALIDATION_WRONG_TYPE("variable", v.meta.name, "int", value)
		}
		if v.meta.min != nil && val < v.meta.min.(int) {
			return errors.REG_VALIDATION_OUT_OF_BOUNDS("variable", v.meta.name, v.meta.min, v.meta.max, val)
		}
		if v.meta.max != nil && val > v.meta.max.(int) {
			return errors.REG_VALIDATION_OUT_OF_BOUNDS("variable", v.meta.name, v.meta.min, v.meta.max, val)
		}
	case "float":
		val, ok := value.(float64)
		if !ok {
			return errors.REG_VALIDATION_WRONG_TYPE("variable", v.meta.name, "float64", value)
		}
		if v.meta.min != nil && val < v.meta.min.(float64) {
			return errors.REG_VALIDATION_OUT_OF_BOUNDS("variable", v.meta.name, v.meta.min, v.meta.max, val)
		}
		if v.meta.max != nil && val > v.meta.max.(float64) {
			return errors.REG_VALIDATION_OUT_OF_BOUNDS("variable", v.meta.name, v.meta.min, v.meta.max, val)
		}
	case "bool":
		_, ok := value.(bool)
		if !ok {
			return errors.REG_VALIDATION_WRONG_TYPE("variable", v.meta.name, "bool", value)
		}
	case "string":
		val, ok := value.(string)
		if !ok {
			return errors.REG_VALIDATION_WRONG_TYPE("variable", v.meta.name, "string", value)
		}
		if v.meta.min != nil && len(val) < v.meta.min.(int) {
			return errors.REG_VALIDATION_OUT_OF_BOUNDS_LENGTH("variable", v.meta.name, v.meta.min, v.meta.max, val)
		}
		if v.meta.max != nil && len(val) > v.meta.max.(int) {
			return errors.REG_VALIDATION_OUT_OF_BOUNDS_LENGTH("variable", v.meta.name, v.meta.min, v.meta.max, val)
		}
	}
	return nil
}
