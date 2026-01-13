package helper

func DerefString(s *string, fallback string) string {
	if s == nil {
		return fallback
	}
	return *s
}
