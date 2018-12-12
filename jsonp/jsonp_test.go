package jsonp

import (
	"encoding/json"
	"testing"
)

func Test_json(t *testing.T) {
	test1 := []byte("a test data")
	if json.Valid(test1) == true {
		t.Error("unexpected value")
	}

	test2, _ := json.Marshal(jsonpWrap{string(test1)})
	if json.Valid(test2) == false {
		t.Error("warp data error")
	}

	temp := &jsonpWrap{}
	err := json.Unmarshal(test2, &temp)
	if err != nil {
		t.Error(err)
	}
}
