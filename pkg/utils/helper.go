package utils

func StrOrNil(s *string, fallback string) string {
	if s == nil {
		return fallback
	}

	return *s
}
