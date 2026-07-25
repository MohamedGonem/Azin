package types

// GetPrimitiveTypeByName returns the TypeInfo for a given primitive TypeKind or nil if the kind is not a primitive type.
func GetPrimitiveTypeByName(name string) *TypeInfo {
	if t, ok := builtinTypes[name]; ok {
		return t
	}
	return nil
}

/* Utility functions */

// IsUnknown checks if the TypeInfo represents an unknown type (needs type inference).
func (t *TypeInfo) IsUnknown() bool {
	return t != nil && t.Kind == Unknown
}

// IsValid checks if the TypeInfo is valid (not nil and not Invalid).
func (t *TypeInfo) IsValid() bool {
	return t != nil && t.Kind != Error
}

// IsUnit checks if the TypeInfo represents the unit type.
func (t *TypeInfo) IsUnit() bool {
	return t != nil && t.Kind == Unit
}

// IsPrimitive checks if the type is a primitive type (bool, int, float, string, char).
func (t *TypeInfo) IsPrimitive() bool {
	return t != nil && GetPrimitiveTypeByName(t.Name) != nil
}

// IsNominal checks if the type is a nominal (user-defined) type.
func (t *TypeInfo) IsNominal() bool {
	return t != nil && t.Kind == Nominal
}

// IsInt checks if the type is an integer type.
func (t *TypeInfo) IsInt() bool {
	return t != nil && t.Kind == Int
}

// IsFloat checks if the type is a floating-point type.
func (t *TypeInfo) IsFloat() bool {
	return t != nil && t.Kind == Float
}

// IsNumeric checks if the type is either an integer or a floating-point type.
func (t *TypeInfo) IsNumeric() bool {
	return t != nil && (t.Kind == Int || t.Kind == Float)
}

// IsBool checks if the type is a boolean type.
func (t *TypeInfo) IsBool() bool {
	return t != nil && t.Kind == Bool
}

// IsString checks if the type is a string type.
func (t *TypeInfo) IsString() bool {
	return t != nil && t.Kind == String
}

// IsChar checks if the type is a character type.
func (t *TypeInfo) IsChar() bool {
	return t != nil && t.Kind == Char
}
