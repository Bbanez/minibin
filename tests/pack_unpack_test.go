package tests

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	minibin "github.com/bbanez/minibin/dist"
	"github.com/bbanez/minibin/src/utils"
)

func TestPack(t *testing.T) {
	entry := minibin.Entry{
		Id:        "123",
		CreatedAt: uint64(time.Now().UnixMilli()),
		UpdatedAt: 2,
		Name:      utils.StringRef("My Entry"),
		Props: []*minibin.EntryProp{
			{
				Id:  "1",
				Typ: minibin.ENTRY_PROP_TYP_STRING,
			},
			{
				Id:  "2",
				Typ: minibin.ENTRY_PROP_TYP_NUMBER,
			},
		},
		Tags: []*string{utils.StringRef("Blog"), utils.StringRef("L")},
	}
	bytes := entry.Pack()
	jsonB, err := json.Marshal(entry)
	fmt.Println("Packed bytes:", bytes)
	fmt.Println("Size:", len(bytes), len(jsonB))
	e, err := minibin.UnpackEntry(bytes)
	if err != nil {
		t.Fatal("Failed to unpack entry:", err)
	}
	fmt.Println(utils.SerializeJsonPretty(e))
}
