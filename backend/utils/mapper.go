package utils

func MapSlice [T any , R any](input []T, mapper func(T) R) []R{
	output := make([]R , len(input))
	for i,item := range input {
		output[i] = mapper(item)
	} 
	return output
}