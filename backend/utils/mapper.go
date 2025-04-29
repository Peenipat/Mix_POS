package utils

// MapSlice สร้าง slice ใหม่จาก input โดยเรียก mapper กับแต่ละ element
// - T คือ ชนิดของข้อมูลใน input slice
// - R คือ ชนิดของข้อมูลใน output slice
func MapSlice [T any , R any](input []T, mapper func(T) R) []R{
	output := make([]R , len(input))
	for i,item := range input {
		output[i] = mapper(item)
	} 
	return output
}