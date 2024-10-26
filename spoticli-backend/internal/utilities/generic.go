package utilities

func Flatten(arr [][]byte) []byte {
	result := []byte{}
	for _, subarr := range arr {
		result = append(result, subarr...)
	}
	return result
}
