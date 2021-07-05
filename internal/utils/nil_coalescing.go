package utils

func CoalesceStr(source *string, fallback string) string {
	if source != nil {
		return *source
	} else {
		return fallback
	}
}

func CoalesceInt(source *int, fallback int) int {
	if source != nil {
		return *source
	} else {
		return fallback
	}
}
