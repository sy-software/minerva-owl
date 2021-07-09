package utils

// CoalesceStr check a string pointer if nil returns the fallback value
func CoalesceStr(source *string, fallback string) string {
	if source != nil {
		return *source
	} else {
		return fallback
	}
}

// CoalesceInt check a string pointer if nil returns the fallback value
func CoalesceInt(source *int, fallback int) int {
	if source != nil {
		return *source
	} else {
		return fallback
	}
}
