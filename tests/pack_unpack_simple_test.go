package tests

import (
	"encoding/json"
	"fmt"
	"testing"

	m "github.com/bbanez/minibin/tests/dist/go"
)

func TestPackSimple(t *testing.T) {
	obj := m.NewObjS(
		-10.1234,
	)
	s1, err := json.Marshal(obj)
	if err != nil {
		t.Fatal(err)
	}
	packed := obj.Pack()
	fmt.Println("Packed", packed)
	unpacked, err := m.UnpackObjS(packed, nil)
	if err != nil {
		t.Fatal(err)
	}
	s2, err := json.Marshal(unpacked)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("S1", string(s1))
	fmt.Println("S2", string(s2))
}
