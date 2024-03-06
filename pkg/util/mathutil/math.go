package mathutil

func Max32(x, y int32) int32 {
	if x > y {
		return x
	}
	return y
}

func Min32(x, y int32) int32 {
	if x > y {
		return y
	}
	return x
}
