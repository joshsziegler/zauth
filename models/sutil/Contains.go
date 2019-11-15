package sutil

// Contains returns true if and only if the string s is in the slice of strings.
func Contains(sSlice []string, s string) bool {
	for _, temp := range sSlice {
		if temp == s {
			return true
		}
	}
	return false
}
