package utils

func Optional(value string) string {
	if value == "" {
		return "-"
	}
	return value
}
