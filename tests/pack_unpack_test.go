package tests

import (
	"fmt"
	"testing"

	"github.com/bbanez/minibin/src/utils"
	minibin "github.com/bbanez/minibin/tests/dist"
)

func TestPack(t *testing.T) {
	obj1 := minibin.Obj1{
		Str:    "str",
		StrArr: []string{"1", "2", "3"},
		I32:    1234,
		I32Arr: []int32{1, 2, 3, 4},
	}
	s1 := utils.SerializeJson(obj1)
	obj1Bytes := obj1.Pack()
	fmt.Println("Packed size:", len(obj1Bytes), len(s1), "-> Packed bytes:", obj1Bytes)
	e, err := minibin.UnpackObj1(obj1Bytes)
	if err != nil {
		t.Fatal("Failed to unpack entry:", err)
	}
	s2 := utils.SerializeJson(e)
	fmt.Println(s1)
	fmt.Println(s2)
	if s1 != s2 {
		t.Errorf("Object do not match after pack/unpack")
	}
}
