package lexorank

import "strings"

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
const base = len(alphabet)

// Between returns a rank string lexicographically between lo and hi.
// If lo is empty, ranks before hi. If hi is empty, ranks after lo.
// If both empty, returns "m" (middle of alphabet).
func Between(lo, hi string) string {
	if lo == "" && hi == "" {
		return "m"
	}
	if lo == "" {
		return before(hi)
	}
	if hi == "" {
		return after(lo)
	}

	// Find midpoint between lo and hi
	lo = normalize(lo, len(hi))
	hi = normalize(hi, len(lo))

	mid := midpoint(lo, hi)
	if mid == lo || mid == hi {
		// No room — append a middle char to lo
		return lo + "m"
	}
	return mid
}

// Initial returns the starting rank for a fresh list.
func Initial() string { return "m" }

// After returns a rank guaranteed to sort after s.
func After(s string) string { return after(s) }

// Before returns a rank guaranteed to sort before s.
func Before(s string) string { return before(s) }

func after(s string) string {
	if s == "" {
		return "m"
	}
	last := s[len(s)-1]
	idx := strings.IndexByte(alphabet, last)
	if idx < base-1 {
		return s[:len(s)-1] + string(alphabet[idx+1])
	}
	return s + "m"
}

func before(s string) string {
	if s == "" {
		return "m"
	}
	last := s[len(s)-1]
	idx := strings.IndexByte(alphabet, last)
	if idx > 0 {
		return s[:len(s)-1] + string(alphabet[idx-1])
	}
	return s + "0"
}

func normalize(s string, length int) string {
	for len(s) < length {
		s += "0"
	}
	return s
}

func midpoint(lo, hi string) string {
	result := make([]byte, len(lo))
	carry := 0
	for i := len(lo) - 1; i >= 0; i-- {
		li := strings.IndexByte(alphabet, lo[i])
		hi2 := strings.IndexByte(alphabet, hi[i])
		sum := li + hi2 + carry
		carry = 0
		mid := sum / 2
		if sum%2 != 0 {
			carry = base
		}
		result[i] = alphabet[mid]
	}
	return string(result)
}
