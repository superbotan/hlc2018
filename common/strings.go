package common

func SubString(s string, start int, len int) string {
	r := []rune(s)
	res := string(r[start : start+len])
	return res
}

func Abs(a int32) int32 {
	if a < 0 {
		return -a
	}
	return a
}

func Min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
