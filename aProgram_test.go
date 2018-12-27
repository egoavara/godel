package godel

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	var temp interface{} = [2]float32{1, 1}
	if v, ok := temp.([2]float32); ok {
		fmt.Println(v)
	} else {
		panic("E")
	}
}
