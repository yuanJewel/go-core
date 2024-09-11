package utils

func InSlice(key interface{}, slice interface{}) bool {
	switch key.(type) {
	case int:
		data, ok := slice.([]int)
		if !ok {
			return false
		}
		for _, e := range data {
			if e == key {
				return true
			}
		}
	case string:
		data, ok := slice.([]string)
		if !ok {
			return false
		}
		for _, e := range data {
			if e == key {
				return true
			}
		}
	default:
		return false
	}
	return false
}
