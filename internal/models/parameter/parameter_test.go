package parameter

import (
    "testing"
)

func TestGenerateParamCombinations(t *testing.T) {
    duplicate := 2
    combinations := GenerateParamCombinations(duplicate)

    if len(combinations) == 0 {
        t.Fatal("expected at least one parameter combination")
    }

    ids := make(map[string]struct{})
    for id, param := range combinations {
        if _, exists := ids[id]; exists {
            t.Errorf("duplicate ID found: %s", id)
        }
        ids[id] = struct{}{}


        if param.PopulationSize <= 0 {
            t.Errorf("invalid PopulationSize for ID %s", id)
        }
        if param.ProcreatorSize < 0 {
            t.Errorf("invalid ProcreatorSize for ID %s", id)
        }
        if param.ID == "" {
            t.Errorf("empty ID for parameter: %+v", param)
        }
    }
}

func TestParameterPrint(t *testing.T) {
    param := Parameter{
        ID:                   "test123",
        GenerationLimit:      10,
        PopulationSize:       100,
        ProcreatorSize:       20,
        DeviationEraInterval: 1,
        KClosestCs:           3,
        A:                    1,
        B:                    1,
        EtaAr:                400,
        GammaAr:              10,
        GammaMgr:             0.9,
        DatasetIndex:         0,
    }
    param.Print()
}