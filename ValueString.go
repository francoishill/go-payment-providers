package gopp

import (
	"strconv"
)

type ValueString string

func (this ValueString) ToString() string {
	return string(this)
}

func (this ValueString) ToInt64() int64 {
	intVal, err := strconv.ParseInt(string(this), 10, 64)
	if err != nil {
		panic(err)
	}
	return intVal
}
func (this ValueString) ToFloat32() float32 {
	floatVal, err := strconv.ParseFloat(string(this), 32)
	if err != nil {
		panic(err)
	}
	return float32(floatVal)
}
