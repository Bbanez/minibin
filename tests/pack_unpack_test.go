package tests

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/bbanez/minibin/src/utils"
	minibin "github.com/bbanez/minibin/tests/dist"
)

func TestPack(t *testing.T) {
	items := minibin.Obj1Arr{
		Items: []*minibin.Obj1{},
	}
	for i := 0; i < 5; i++ {
		item := &minibin.Obj1{
			Str:    fmt.Sprintf("item_%d", i),
			I32:    int32(i * 10),
			I32Arr: []int32{int32(i * 10), int32(i * 20), int32(i * 30)},
			I64:    int64(i * 100),
			I64Arr: []int64{int64(i * 100), int64(i * 200), int64(i * 300)},
			U32:    uint32(i * 1000),
			U32Arr: []uint32{uint32(i * 1000), uint32(i * 2000), uint32(i * 3000)},
			U64:    uint64(i * 10000),
			U64Arr: []uint64{uint64(i * 10000), uint64(i * 20000), uint64(i * 30000)},
			F32:    float32(i) + 0.123,
			F32Arr: []float32{float32(i) + 0.234, float32(i) + 0.345},
			F64:    float64(i) + 0.456789,
			F64Arr: []float64{float64(i) + 0.567890, float64(i) + 0.678901},
			Obj2: minibin.Obj2{
				Key:   fmt.Sprintf("key_%d", i),
				Value: float32(i) * 1.5,
			},
			Obj2Arr: []*minibin.Obj2{
				{
					Key:   fmt.Sprintf("keyarr_%d_1", i),
					Value: float32(i) * 2.5,
				},
				{
					Key:   fmt.Sprintf("keyarr_%d_2", i),
					Value: float32(i) * 3.5,
				},
			},
			Enum1:    minibin.ENUM1_E1,
			Enum1Arr: []minibin.Enum1{minibin.ENUM1_E2, minibin.ENUM1_E3},
		}
		items.Items = append(items.Items, item)
	}
	timeOffset := time.Now().UnixNano()
	s1 := utils.SerializeJson(items)
	fmt.Println("JS time", time.Now().UnixNano()-timeOffset)
	timeOffset = time.Now().UnixNano()
	bytes := items.Pack()
	fmt.Println("MB time", time.Now().UnixNano()-timeOffset)
	fmt.Println("Packed size:", len(bytes), len(s1))
	timeOffset = time.Now().UnixNano()
	e, err := minibin.UnpackObj1Arr(bytes)
	fmt.Println("MB Unpack time", time.Now().UnixNano()-timeOffset)
	if err != nil {
		t.Fatal("Failed to unpack entry:", err)
	}
	timeOffset = time.Now().UnixNano()
	tmp := minibin.Obj1Arr{}
	err = json.Unmarshal([]byte(s1), &tmp)
	fmt.Println("JS Unpack time", time.Now().UnixNano()-timeOffset)
	if err != nil {
		t.Fatal("Failed to unmarshal json:", err)
	}
	s2 := utils.SerializeJson(e)
	if s1 != s2 {
		t.Errorf("Object do not match after pack/unpack")
	}
}
