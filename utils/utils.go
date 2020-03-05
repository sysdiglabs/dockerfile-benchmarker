package utils

func ArrayToMap(arr []string) map[string]bool {
	m := map[string]bool{}

	for _, i := range arr {
		m[i] = true
	}

	return m
}

func MapToArray(m map[string]bool) []string {
	arr := []string{}

	for k := range m {
		arr = append(arr, k)
	}

	return arr
}
