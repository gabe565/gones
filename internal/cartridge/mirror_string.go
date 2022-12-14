// Code generated by "stringer -type Mirror"; DO NOT EDIT.

package cartridge

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Horizontal-0]
	_ = x[Vertical-1]
	_ = x[SingleLower-2]
	_ = x[SingleUpper-3]
	_ = x[FourScreen-4]
}

const _Mirror_name = "HorizontalVerticalSingleLowerSingleUpperFourScreen"

var _Mirror_index = [...]uint8{0, 10, 18, 29, 40, 50}

func (i Mirror) String() string {
	if i >= Mirror(len(_Mirror_index)-1) {
		return "Mirror(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Mirror_name[_Mirror_index[i]:_Mirror_index[i+1]]
}
