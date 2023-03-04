package resolver

func strValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
