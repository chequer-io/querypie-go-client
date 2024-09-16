package utils

func Optional(value string) string {
	if value == "" {
		return "-"
	}
	return value
}

func OptionalPtr(ptr *string) string {
	if ptr == nil {
		return "-"
	} else if *ptr == "" {
		return "-"
	}
	return *ptr
}
