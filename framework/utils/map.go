package utils

func Merge(dst map[string]interface{}, src map[string]interface{}) {
	for k, v := range src {
		_, ok := (dst)[k]
		if !ok {
			(dst)[k] = v
		}
	}
}
