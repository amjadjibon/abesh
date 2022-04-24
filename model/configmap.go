package model

import (
	"strconv"
	"strings"
	"time"
)

type ConfigMap map[string]string

func (v ConfigMap) String(key string, defaultValue string) string {
	if o, ok := v[key]; ok {
		return o
	}
	return defaultValue
}

func (v ConfigMap) StringList(key string, sep string, defaultValue []string) []string {
	if o, ok := v[key]; ok {
		return strings.Split(o, sep)
	}
	return defaultValue
}

func (v ConfigMap) StringMap(key string, defaultValue ConfigMap) ConfigMap {
	if o, ok := v[key]; ok {
		data := make(ConfigMap)
		for _, v := range strings.Split(o, ";") {
			sd := strings.Split(v, "=")
			if len(sd) == 2 {
				data[strings.TrimSpace(sd[0])] = strings.TrimSpace(sd[1])
			}
		}
		return data
	}
	return defaultValue
}

func (v ConfigMap) Bytes(key string, defaultValue []byte) []byte {
	if o, ok := v[key]; ok {
		return []byte(o)
	}
	return defaultValue
}

func (v ConfigMap) Float32(key string, defaultValue float32) float32 {
	if o, ok := v[key]; ok {
		i, err := strconv.ParseFloat(o, 32)
		if err != nil {
			return defaultValue
		}
		return float32(i)
	}
	return defaultValue
}

func (v ConfigMap) Float64(key string, defaultValue float64) float64 {
	if o, ok := v[key]; ok {
		i, err := strconv.ParseFloat(o, 64)
		if err != nil {
			return defaultValue
		}
		return i
	}
	return defaultValue
}

func (v ConfigMap) Int(key string, defaultValue int) int {
	if o, ok := v[key]; ok {
		i, err := strconv.Atoi(o)
		if err != nil {
			return defaultValue
		}

		return i
	}
	return defaultValue
}

func (v ConfigMap) Uint(key string, defaultValue uint) uint {
	if o, ok := v[key]; ok {
		i, err := strconv.ParseUint(o, 10, 32)
		if err != nil {
			return defaultValue
		}
		return uint(i)
	}
	return defaultValue
}

func (v ConfigMap) Int8(key string, defaultValue int8) int8 {
	if o, ok := v[key]; ok {
		i, err := strconv.ParseInt(o, 10, 8)
		if err != nil {
			return defaultValue
		}

		return int8(i)
	}
	return defaultValue
}

func (v ConfigMap) Uint8(key string, defaultValue uint8) uint8 {
	if o, ok := v[key]; ok {
		i, err := strconv.ParseUint(o, 10, 8)
		if err != nil {
			return defaultValue
		}

		return uint8(i)
	}
	return defaultValue
}

func (v ConfigMap) Int16(key string, defaultValue int16) int16 {
	if o, ok := v[key]; ok {
		i, err := strconv.ParseInt(o, 10, 16)
		if err != nil {
			return defaultValue
		}

		return int16(i)
	}
	return defaultValue
}

func (v ConfigMap) Uint16(key string, defaultValue uint16) uint16 {
	if o, ok := v[key]; ok {
		i, err := strconv.ParseUint(o, 10, 16)
		if err != nil {
			return defaultValue
		}
		return uint16(i)
	}
	return defaultValue
}

func (v ConfigMap) Int32(key string, defaultValue int32) int32 {
	if o, ok := v[key]; ok {
		i, err := strconv.ParseInt(o, 10, 32)
		if err != nil {
			return defaultValue
		}
		return int32(i)
	}
	return defaultValue
}

func (v ConfigMap) Uint32(key string, defaultValue uint32) uint32 {
	if o, ok := v[key]; ok {
		i, err := strconv.ParseUint(o, 10, 32)
		if err != nil {
			return defaultValue
		}
		return uint32(i)
	}
	return defaultValue
}

func (v ConfigMap) Int64(key string, defaultValue int64) int64 {
	if o, ok := v[key]; ok {
		i, err := strconv.ParseInt(o, 10, 64)
		if err != nil {
			return defaultValue
		}
		return i
	}
	return defaultValue
}

func (v ConfigMap) Uint64(key string, defaultValue uint64) uint64 {
	if o, ok := v[key]; ok {
		i, err := strconv.ParseUint(o, 10, 64)
		if err != nil {
			return defaultValue
		}
		return i
	}
	return defaultValue
}

func (v ConfigMap) Duration(key string, defaultValue time.Duration) time.Duration {
	if o, ok := v[key]; ok {
		val, err := time.ParseDuration(o)
		if err != nil {
			return defaultValue
		}
		return val
	}
	return defaultValue
}

func (v ConfigMap) Time(key string, defaultValue time.Time) time.Time {
	if o, ok := v[key]; ok {
		val, err := time.Parse(time.RFC3339, o)
		if err != nil {
			return defaultValue
		}
		return val
	}
	return defaultValue
}

func (v ConfigMap) Bool(key string, defaultValue bool) bool {
	if o, ok := v[key]; ok {
		val, err := strconv.ParseBool(o)
		if err != nil {
			return defaultValue
		}
		return val
	}
	return defaultValue
}

func (v ConfigMap) IntList(key string, sep string, defaultValue []int) []int {
	if o, ok := v[key]; ok {
		var stringList = strings.Split(o, sep)
		var result []int
		for _, val := range stringList {
			i, err := strconv.Atoi(val)
			if err != nil {
				return defaultValue
			}
			result = append(result, i)
		}
		return result
	}
	return defaultValue
}

func (v ConfigMap) Int8List(key string, sep string, defaultValue []int8) []int8 {
	if o, ok := v[key]; ok {
		var stringList = strings.Split(o, sep)
		var result []int8
		for _, val := range stringList {
			i, err := strconv.Atoi(val)
			if err != nil {
				return defaultValue
			}
			result = append(result, int8(i))
		}
		return result
	}
	return defaultValue
}

func (v ConfigMap) Int16List(key string, sep string, defaultValue []int16) []int16 {
	if o, ok := v[key]; ok {
		var stringList = strings.Split(o, sep)
		var result []int16
		for _, val := range stringList {
			i, err := strconv.Atoi(val)
			if err != nil {
				return defaultValue
			}
			result = append(result, int16(i))
		}
		return result
	}
	return defaultValue
}

func (v ConfigMap) Int32List(key string, sep string, defaultValue []int32) []int32 {
	if o, ok := v[key]; ok {
		var stringList = strings.Split(o, sep)
		var result []int32
		for _, val := range stringList {
			i, err := strconv.Atoi(val)
			if err != nil {
				return defaultValue
			}
			result = append(result, int32(i))
		}
		return result
	}
	return defaultValue
}

func (v ConfigMap) Int64List(key string, sep string, defaultValue []int64) []int64 {
	if o, ok := v[key]; ok {
		var stringList = strings.Split(o, sep)
		var result []int64
		for _, val := range stringList {
			i, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return defaultValue
			}
			result = append(result, i)
		}
		return result
	}
	return defaultValue
}

func (v ConfigMap) UintList(key string, sep string, defaultValue []uint) []uint {
	if o, ok := v[key]; ok {
		var stringList = strings.Split(o, sep)
		var result []uint
		for _, val := range stringList {
			i, err := strconv.Atoi(val)
			if err != nil {
				return defaultValue
			}
			result = append(result, uint(i))
		}
		return result
	}
	return defaultValue
}

func (v ConfigMap) Uint8List(key string, sep string, defaultValue []uint8) []uint8 {
	if o, ok := v[key]; ok {
		var stringList = strings.Split(o, sep)
		var result []uint8
		for _, val := range stringList {
			i, err := strconv.Atoi(val)
			if err != nil {
				return defaultValue
			}
			result = append(result, uint8(i))
		}
		return result
	}
	return defaultValue
}

func (v ConfigMap) Uint16List(key string, sep string, defaultValue []uint16) []uint16 {
	if o, ok := v[key]; ok {
		var stringList = strings.Split(o, sep)
		var result []uint16
		for _, val := range stringList {
			i, err := strconv.Atoi(val)
			if err != nil {
				return defaultValue
			}
			result = append(result, uint16(i))
		}
		return result
	}
	return defaultValue
}

func (v ConfigMap) Uint32List(key string, sep string, defaultValue []uint32) []uint32 {
	if o, ok := v[key]; ok {
		var stringList = strings.Split(o, sep)
		var result []uint32
		for _, val := range stringList {
			i, err := strconv.Atoi(val)
			if err != nil {
				return defaultValue
			}
			result = append(result, uint32(i))
		}
		return result
	}
	return defaultValue
}

func (v ConfigMap) Uint64List(key string, sep string, defaultValue []uint64) []uint64 {
	if o, ok := v[key]; ok {
		var stringList = strings.Split(o, sep)
		var result []uint64
		for _, val := range stringList {
			i, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				return defaultValue
			}
			result = append(result, i)
		}
		return result
	}
	return defaultValue
}

func (v ConfigMap) Float32List(key string, sep string, defaultValue []float32) []float32 {
	if o, ok := v[key]; ok {
		var stringList = strings.Split(o, sep)
		var result []float32
		for _, val := range stringList {
			f, err := strconv.ParseFloat(val, 32)
			if err != nil {
				return defaultValue
			}
			result = append(result, float32(f))
		}
		return result
	}
	return defaultValue
}

func (v ConfigMap) Float64List(key string, sep string, defaultValue []float64) []float64 {
	if o, ok := v[key]; ok {
		var stringList = strings.Split(o, sep)
		var result []float64
		for _, val := range stringList {
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return defaultValue
			}
			result = append(result, f)
		}
		return result
	}
	return defaultValue
}
