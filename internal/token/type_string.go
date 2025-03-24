// Code generated by "stringer -type=Type"; DO NOT EDIT.

package token

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TypeLeftParen-0]
	_ = x[TypeRightParen-1]
	_ = x[TypeLeftBrace-2]
	_ = x[TypeRightBrace-3]
	_ = x[TypeComma-4]
	_ = x[TypeDot-5]
	_ = x[TypeMinus-6]
	_ = x[TypePlus-7]
	_ = x[TypeSemicolon-8]
	_ = x[TypeSlash-9]
	_ = x[TypeStar-10]
	_ = x[TypeBang-11]
	_ = x[TypeBangEqual-12]
	_ = x[TypeEqual-13]
	_ = x[TypeEqualEqual-14]
	_ = x[TypeGreater-15]
	_ = x[TypeGreaterEqual-16]
	_ = x[TypeLess-17]
	_ = x[TypeLessEqual-18]
	_ = x[TypeIdentifier-19]
	_ = x[TypeString-20]
	_ = x[TypeNumber-21]
	_ = x[TypeAnd-22]
	_ = x[TypeClass-23]
	_ = x[TypeElse-24]
	_ = x[TypeFalse-25]
	_ = x[TypeFun-26]
	_ = x[TypeFor-27]
	_ = x[TypeIf-28]
	_ = x[TypeNil-29]
	_ = x[TypeOr-30]
	_ = x[TypePrint-31]
	_ = x[TypeReturn-32]
	_ = x[TypeSuper-33]
	_ = x[TypeThis-34]
	_ = x[TypeTrue-35]
	_ = x[TypeVar-36]
	_ = x[TypeWhile-37]
	_ = x[TypeEOF-38]
}

const _Type_name = "TypeLeftParenTypeRightParenTypeLeftBraceTypeRightBraceTypeCommaTypeDotTypeMinusTypePlusTypeSemicolonTypeSlashTypeStarTypeBangTypeBangEqualTypeEqualTypeEqualEqualTypeGreaterTypeGreaterEqualTypeLessTypeLessEqualTypeIdentifierTypeStringTypeNumberTypeAndTypeClassTypeElseTypeFalseTypeFunTypeForTypeIfTypeNilTypeOrTypePrintTypeReturnTypeSuperTypeThisTypeTrueTypeVarTypeWhileTypeEOF"

var _Type_index = [...]uint16{0, 13, 27, 40, 54, 63, 70, 79, 87, 100, 109, 117, 125, 138, 147, 161, 172, 188, 196, 209, 223, 233, 243, 250, 259, 267, 276, 283, 290, 296, 303, 309, 318, 328, 337, 345, 353, 360, 369, 376}

func (i Type) String() string {
	if i < 0 || i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
