package rest

import "unicode"

func MaskAccessToken(token string) string {
	if len(token) <= 11 {
		return token
	}
	masked := []rune(token[:11])
	for _, r := range token[11:] {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			masked = append(masked, '*')
		} else {
			masked = append(masked, r)
		}
	}
	return string(masked)
}
