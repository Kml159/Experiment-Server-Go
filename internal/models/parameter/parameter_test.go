package parameter

import (
	"fmt"
	"strings"

	"testing"
)

func TestGenerateParamCombinations(t *testing.T) {
    params := GenerateParamCombinations(2)
    if len(params) == 0 {
        t.Errorf("expected non-empty parameter combinations")
    }
    for _, p := range params {
        fmt.Println(p)
    }

    if strings.Split(params[0].ID, "x")[1] != "00" || strings.Split(params[1].ID, "x")[1] != "01"{
        t.Errorf("wrong set id")
    }
}
