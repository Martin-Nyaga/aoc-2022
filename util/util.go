package util

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

func SumInts(arr []int) int {
	total := 0
	for _, i := range arr {
		total += i
	}
	return total
}
