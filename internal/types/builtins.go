package types

const (
	UnitTypeName    = "unit"
	BoolTypeName    = "bool"
	IntTypeName     = "int"
	FloatTypeName   = "float"
	StringTypeName  = "string"
	CharTypeName    = "char"
	UnknownTypeName = "<unknown>"
	ErrorTypeName   = "<error>"
)

var builtinTypes = map[string]*TypeInfo{
	UnknownTypeName: {Kind: Unknown, Name: UnknownTypeName},
	ErrorTypeName:   {Kind: Error, Name: ErrorTypeName},
	UnitTypeName:    {Kind: Unit, Name: UnitTypeName},
	BoolTypeName:    {Kind: Bool, Name: BoolTypeName},
	IntTypeName:     {Kind: Int, Name: IntTypeName},
	FloatTypeName:   {Kind: Float, Name: FloatTypeName},
	StringTypeName:  {Kind: String, Name: StringTypeName},
	CharTypeName:    {Kind: Char, Name: CharTypeName},
}

// UnknownType returns a TypeInfo representing an unknown type, used for type inference before the actual type is determined.
func UnknownType() *TypeInfo {
	return builtinTypes[UnknownTypeName]
}

// ErrorType returns a TypeInfo representing an invalid type.
func ErrorType() *TypeInfo {
	return builtinTypes[ErrorTypeName]
}

// UnitType returns the TypeInfo for the unit type.
func UnitType() *TypeInfo {
	return builtinTypes[UnitTypeName]
}

// IntType returns the TypeInfo for the integer type.
func IntType() *TypeInfo {
	return builtinTypes[IntTypeName]
}

// FloatType returns the TypeInfo for the floating-point type.
func FloatType() *TypeInfo {
	return builtinTypes[FloatTypeName]
}

// BoolType returns the TypeInfo for the boolean type.
func BoolType() *TypeInfo {
	return builtinTypes[BoolTypeName]
}

// CharType returns the TypeInfo for the character type.
func CharType() *TypeInfo {
	return builtinTypes[CharTypeName]
}

// StringType returns the TypeInfo for the string type.
func StringType() *TypeInfo {
	return builtinTypes[StringTypeName]
}
