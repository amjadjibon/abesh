package utility

// IsIn check if the value exists in the dataList
// don't use it for large data set
// time complexity O(n) because of linear scan
func IsIn(dataList []string, value string) bool {
	for _, v := range dataList {
		if v == value {
			return true
		}
	}
	return false
}

func IsInInt(dataList []int, value int) bool {
	for _, v := range dataList {
		if v == value {
			return true
		}
	}
	return false
}

func IsInInt8(dataList []int8, value int8) bool {
	for _, v := range dataList {
		if v == value {
			return true
		}
	}
	return false
}

func IsInInt16(dataList []int16, value int16) bool {
	for _, v := range dataList {
		if v == value {
			return true
		}
	}
	return false
}

func IsInInt32(dataList []int32, value int32) bool {
	for _, v := range dataList {
		if v == value {
			return true
		}
	}
	return false
}

func IsInInt64(dataList []int64, value int64) bool {
	for _, v := range dataList {
		if v == value {
			return true
		}
	}
	return false
}

func IsInUint(dataList []uint, value uint) bool {
	for _, v := range dataList {
		if v == value {
			return true
		}
	}
	return false
}

func IsInUint8(dataList []uint8, value uint8) bool {
	for _, v := range dataList {
		if v == value {
			return true
		}
	}
	return false
}

func IsInUint16(dataList []uint16, value uint16) bool {
	for _, v := range dataList {
		if v == value {
			return true
		}
	}
	return false
}

func IsInUint32(dataList []uint32, value uint32) bool {
	for _, v := range dataList {
		if v == value {
			return true
		}
	}
	return false
}

func IsInUint64(dataList []uint64, value uint64) bool {
	for _, v := range dataList {
		if v == value {
			return true
		}
	}
	return false
}

func IsInFloat32(dataList []float32, value float32) bool {
	for _, v := range dataList {
		if v == value {
			return true
		}
	}
	return false
}

func IsInFloat64(dataList []float64, value float64) bool {
	for _, v := range dataList {
		if v == value {
			return true
		}
	}
	return false
}

// IsInGeneric check if the value exists in the dataList
// don't use it for large data set
// time complexity O(n) because of linear scan
func IsInGeneric[T comparable](dataList []T, value T) bool {
	for _, v := range dataList {
		if v == value {
			return true
		}
	}
	return false
}
