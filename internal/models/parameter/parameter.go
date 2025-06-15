package parameter

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	base    = 10000
	product = 1e3
)

type Parameter struct {
	ID                   string // 4234234320000x0
	GenerationLimit      int
	PopulationSize       int
	ProcreatorSize       int
	DeviationEraInterval int
	KClosestCs           int
	A                    int
	B                    int
	EtaAr                int
	GammaAr              int
	GammaMgr             float64
	DatasetIndex         int
}

func (p Parameter) Print() {
	fmt.Printf(
		"Parameter{ID: %s, GenerationLimit: %d, PopulationSize: %d, ProcreatorSize: %d, DeviationEraInterval: %d, KClosestCs: %d, A: %d, B: %d, EtaAr: %d, GammaAr: %d, GammaMgr: %.2f, DatasetIndex: %d}\n",
		p.ID, p.GenerationLimit, p.PopulationSize, p.ProcreatorSize, p.DeviationEraInterval, p.KClosestCs, p.A, p.B, p.EtaAr, p.GammaAr, p.GammaMgr, p.DatasetIndex,
	)
}

func GenerateParamCombinations(duplicate int) map[string]Parameter {
	populationSizes := []int{500}
	procreatorRatios := []float64{0.2}
	kClosestCs := []int{3}
	etaArs := []int{400}
	gammaArs := []int{10}
	gammaMgrs := []float64{0.9}
	datasetIndexes := []int{0}
	AnB := [][2]int{{1, 1}}

	rand.Seed(time.Now().UnixNano())
	ExperimentSessionID := rand.Intn(100000) * 10000

	type combo struct {
		pop    int
		ratio  float64
		k      int
		eta    int
		gammaA int
		gammaM float64
		dset   int
		a      int
		b      int
	}

	var combos []combo
	for _, pop := range populationSizes {
		for _, ratio := range procreatorRatios {
			for _, k := range kClosestCs {
				for _, eta := range etaArs {
					for _, gammaA := range gammaArs {
						for _, gammaM := range gammaMgrs {
							for _, dset := range datasetIndexes {
								for _, ab := range AnB {
									combos = append(combos, combo{
										pop, ratio, k, eta, gammaA, gammaM, dset, ab[0], ab[1],
									})
								}
							}
						}
					}
				}
			}
		}
	}

	finalCombinations := make([]Parameter, 0, len(combos))
	for expSetID, c := range combos {
		param := Parameter{
			ID:                   fmt.Sprint(ExperimentSessionID*base) + "x" + fmt.Sprint(expSetID),
			GenerationLimit:      int(product / float64(c.pop)),
			PopulationSize:       c.pop,
			ProcreatorSize:       int(float64(c.pop) * c.ratio),
			DeviationEraInterval: int((product / float64(c.pop)) * 0.01),
			KClosestCs:           c.k,
			A:                    c.a,
			B:                    c.b,
			EtaAr:                c.eta,
			GammaAr:              c.gammaA,
			GammaMgr:             c.gammaM,
			DatasetIndex:         c.dset,
		}
		finalCombinations = append(finalCombinations, param)
	}

	duplicated := make(map[string]Parameter)
	for _, combination := range finalCombinations {
		for i := 0; i < duplicate; i++ {
			dup := combination
			dup.ID = combination.ID + fmt.Sprint(i)
			duplicated[dup.ID] = dup
		}
	}

	return duplicated
}
