package strutil

// StringInSlice returns whether or not a value is in the given slice.
func StringInSlice(val string, list []string) bool {
	for _, s := range list {
		if s == val {
			return true
		}
	}
	return false
}
