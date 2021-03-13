package helpers

func SliceSummation(theSlice []int) int {
	summation := 0
	for _, value := range theSlice {
		summation = summation + value
	}

	return summation
}
