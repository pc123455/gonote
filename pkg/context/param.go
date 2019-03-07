package context

type Param map[string]interface{}

func MergeParam(dst, src Param) {
	for k, v := range src {
		_, ok := (dst)[k]
		if !ok {
			(dst)[k] = v
		}
	}
}
