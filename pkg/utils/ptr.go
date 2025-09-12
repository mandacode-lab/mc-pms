package utils

// String returns a pointer to the string value
func String(s string) *string {
	return &s
}

// StringNil returns a pointer to string if not empty, otherwise nil
func StringNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Int64 returns a pointer to the int64 value
func Int64(i int64) *int64 {
	return &i
}

// Int64Nil returns a pointer to int64 if not zero, otherwise nil
func Int64Nil(i int64) *int64 {
	if i == 0 {
		return nil
	}
	return &i
}

// Bool returns a pointer to the bool value
func Bool(b bool) *bool {
	return &b
}

// StringValue returns the value of the string pointer or empty string if nil
func StringValue(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// Int64Value returns the value of the int64 pointer or zero if nil
func Int64Value(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}

// BoolValue returns the value of the bool pointer or false if nil
func BoolValue(p *bool) bool {
	if p == nil {
		return false
	}
	return *p
}
