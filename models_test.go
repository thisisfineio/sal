package sal

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestInlineApplicationMapper_PrintJSON(t *testing.T) {
	f := &InlineApplicationMapper{}
	m, err := f.LoadApplicationMappings()
	if err != nil {
		t.Fail()
	}

	data, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		t.Fail()
	}
	fmt.Println(string(data))
}
