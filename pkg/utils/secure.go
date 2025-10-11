package utils

// ZeroBytes zeroes out the provided byte slice to prevent sensitive data from remaining in memory
func ZeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
