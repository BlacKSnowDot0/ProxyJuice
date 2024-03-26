package Utility

import "unsafe"

// FastStrAtoi converts a string to an integer quickly.
func FastStrAtoi(s string) int {
	n := 0
	for _, c := range s {
		n = n*10 + int(c-'0')
	}
	return n
}

func itoa(i int) []byte {
	// Special case for zero since the rest of the algorithm doesn't handle it
	if i == 0 {
		return []byte{'0'}
	}

	// Account for negative numbers
	negative := i < 0
	if negative {
		i = -i
	}

	// Extract digits in reverse order
	var b [20]byte // int64 has at most 19 digits in base 10, plus sign
	cur := len(b)
	for i > 0 {
		cur--
		b[cur] = byte(i%10 + '0') // Convert int digit to ASCII
		i /= 10
	}

	// Add negative sign if needed
	if negative {
		cur--
		b[cur] = '-'
	}

	return b[cur:]
}

// FastIntToStr converts an integer to its string representation.
func FastIntToStr(i int) string {
	bytes := itoa(i)
	return *(*string)(unsafe.Pointer(&bytes)) // Convert bytes to string without copying
}
