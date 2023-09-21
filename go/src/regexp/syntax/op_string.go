// Code generated by "stringer -type Op -trimprefix Op"; DO NOT EDIT.

package syntax

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OpNoMatch-1]
	_ = x[OpEmptyMatch-2]
	_ = x[OpLiteral-3]
	_ = x[OpCharClass-4]
	_ = x[OpAnyCharNotNL-5]
	_ = x[OpAnyChar-6]
	_ = x[OpBeginLine-7]
	_ = x[OpEndLine-8]
	_ = x[OpBeginText-9]
	_ = x[OpEndText-10]
	_ = x[OpWordBoundary-11]
	_ = x[OpNoWordBoundary-12]
	_ = x[OpCapture-13]
	_ = x[OpStar-14]
	_ = x[OpPlus-15]
	_ = x[OpQuest-16]
	_ = x[OpRepeat-17]
	_ = x[OpConcat-18]
	_ = x[OpAlternate-19]
	_ = x[opPseudo-128]
}

const (
	_Op_name_0 = "NoMatchEmptyMatchLiteralCharClassAnyCharNotNLAnyCharBeginLineEndLineBeginTextEndTextWordBoundaryNoWordBoundaryCaptureStarPlusQuestRepeatConcatAlternate"
	_Op_name_1 = "opPseudo"
)

var (
	_Op_index_0 = [...]uint8{0, 7, 17, 24, 33, 45, 52, 61, 68, 77, 84, 96, 110, 117, 121, 125, 130, 136, 142, 151}
)

func (i Op) String() string {
	switch {
	case 1 <= i && i <= 19:
		i -= 1
		return _Op_name_0[_Op_index_0[i]:_Op_index_0[i+1]]
	case i == 128:
		return _Op_name_1
	default:
		return "Op(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}