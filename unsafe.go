package pg

import "unsafe"

// s2b converts string to a byte slice without memory allocation.
func s2b(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// b2s converts byte slice to a string without memory allocation.
func b2s(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
