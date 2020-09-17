package utils

// Bool returns a pointer to the bool value passed in
func Bool(v bool) *bool {
	return &v
}

// BoolValue returns the value of the bool pointer passed in or
func BoolValue(v *bool) bool {
	if v != nil {
		return *v
	}

	return false
}

// Float64 returns a pointer to the float64 value passed in
func Float64(v float64) *float64 {
	return &v
}

// Float64Value returns the value of the float64 pointer passed in or
func Float64Value(v *float64) float64 {
	if v != nil {
		return *v
	}

	return 0
}

// Float32 returns a pointer to the float32 value passed in
func Float32(v float64) *float64 {
	return &v
}

// Float32Value returns the value of the float32 pointer passed in or
func Float32Value(v *float64) float64 {
	if v != nil {
		return *v
	}

	return 0
}

// Int64 ptr
func Int64(v int64) *int64 {
	return &v
}

// Int64Value value
func Int64Value(v *int64) int64 {
	if v != nil {
		return *v
	}

	return 0
}

// Int ptr
func Int(v int) *int {
	return &v
}

// IntValue value
func IntValue(v *int) int {
	if v != nil {
		return *v
	}

	return 0
}

// String returns a pointer to the string value passed in.
func String(v string) *string {
	return &v
}

// StringValue returns the value of the string pointer passed in or
// "" if the pointer is nil.
func StringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

// StringSlice returns a slice of string pointers given a slice of strings.
func StringSlice(v []string) []*string {
	out := make([]*string, len(v))
	for i := range v {
		out[i] = &v[i]
	}
	return out
}

// StringSliceValue returns a slice of string given a slice of strings pointers.
func StringSliceValue(v []*string) []string {
	out := make([]string, len(v))
	for i := range v {
		if v[i] != nil {
			out[i] = *v[i]
		} else {
			out[i] = ""
		}

	}
	return out
}
