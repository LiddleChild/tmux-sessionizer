package utils

func Clamp(value, mn, mx int) int {
	return max(mn, min(mx, value))
}
