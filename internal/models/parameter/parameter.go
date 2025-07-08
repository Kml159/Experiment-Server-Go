package parameter

import (
	"experiment-server/internal/config"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	base    = 10000
)

// Arguments needed to run the program:

type Parameter struct {
	ID                   string  `json:"id"`
	GenerationLimit      int     `json:"generation_limit"`
	PopulationSize       int     `json:"population_size"`
	KClosestCsPercantage float32 `json:"k_closest_cs_percantage"` // 1.0 always // 3. M_K_CLOSEST_CHARGE_STATIONS_TO_CS_AMOUNT_RATIO: Ratio for closest charge stations.
	A                    int     `json:"A"`
	B                    int     `json:"B"`
	EtaAr                int     `json:"eta_ar"`
	GammaAr              int     `json:"gamma_ar"`
	GammaMgr             float64 `json:"gamma_mgr"`
	DistanceMethod       int     `json:"distance_method"`
	DeviatonEraParameter int     `json:"deviation_era_parameter"`
	DatasetIndex         int     `json:"dataset_index"`
}

func (p Parameter) Print() {
	log.Printf(
		"Parameter{ID: %s, GenerationLimit: %d, PopulationSize: %d, KClosestCsPercantage: %f, A: %d, B: %d, EtaAr: %d, GammaAr: %d, GammaMgr: %.2f, DistanceMethod: %d, DeviatonEraParameter: %d, DatasetIndex: %d}\n",
		p.ID, p.GenerationLimit, p.PopulationSize, p.KClosestCsPercantage, p.A, p.B, p.EtaAr, p.GammaAr, p.GammaMgr, p.DistanceMethod, p.DeviatonEraParameter, p.DatasetIndex,
	)
}

func (p Parameter) String() string {
	return fmt.Sprintf("Parameter{ID: %s, GenerationLimit: %d, PopulationSize: %d, KClosestCsPercantage: %f, A: %d, B: %d, EtaAr: %d, GammaAr: %d, GammaMgr: %.2f, DistanceMethod: %d, DeviatonEraParameter: %d, DatasetIndex: %d}\n",
		p.ID, p.GenerationLimit, p.PopulationSize, p.KClosestCsPercantage, p.A, p.B, p.EtaAr, p.GammaAr, p.GammaMgr, p.DistanceMethod, p.DeviatonEraParameter, p.DatasetIndex)
}
func GenerateParamCombinations(duplicate int, cfg *config.Config) map[string]Parameter {

	populationSizes := []int{1000, 500}
	kClosestCs := []float32{1.0}
	etaArs := []int{400}
	gammaArs := []int{10}
	gammaMgrs := []float64{0.9}
	datasetIndexes := []int{0, 1, 4}
	AnB := [][2]int{{1, 1}, {2, 1}, {3, 2}, {3, 3}, {4, 2}, {5, 3}}
	DistanceMethod := [2]int{0, 1}
	DeviationEraParameter := []int{30}

	rand.Seed(time.Now().UnixNano())
	ExperimentSessionID := cfg.ExperimentBaseId * 10000

	type combo struct {
		pop    int
		k      float32
		eta    int
		gammaA int
		gammaM float64
		dset   int
		a      int
		b      int
		dm     int
		devEra int
	}

	var combos []combo
	for _, pop := range populationSizes {
		for _, k := range kClosestCs {
			for _, eta := range etaArs {
				for _, gammaA := range gammaArs {
					for _, gammaM := range gammaMgrs {
						for _, dset := range datasetIndexes {
							for _, ab := range AnB {
								for _, dm := range DistanceMethod {
									for _, devEra := range DeviationEraParameter {
										combos = append(combos, combo{
											pop, k, eta, gammaA, gammaM, dset, ab[0], ab[1], dm, devEra,
										})
									}
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
			ID:                   fmt.Sprint(ExperimentSessionID*base + expSetID),
			GenerationLimit:      int(float64(cfg.ProductGenerationPopulation) / float64(c.pop)),
			PopulationSize:       c.pop,
			KClosestCsPercantage: c.k,
			A:                    c.a,
			B:                    c.b,
			EtaAr:                c.eta,
			GammaAr:              c.gammaA,
			GammaMgr:             c.gammaM,
			DatasetIndex:         c.dset,
			DistanceMethod:       c.dm,
			DeviatonEraParameter: c.devEra,
		}
		finalCombinations = append(finalCombinations, param)
	}

	for i := range finalCombinations {
		switch finalCombinations[i].DatasetIndex {
		case 0: // Default dataset
			finalCombinations[i].GammaAr = 10
			finalCombinations[i].GammaMgr = 0.9
			finalCombinations[i].EtaAr = 400
		case 1: // Baylands
			finalCombinations[i].GammaAr = 10
			finalCombinations[i].GammaMgr = 1.2
			finalCombinations[i].EtaAr = 600
		case 4: // Marmaracık
			finalCombinations[i].GammaAr = 15
			finalCombinations[i].GammaMgr = 1.2
			finalCombinations[i].EtaAr = 600
		}
	}

	duplicated := make(map[string]Parameter)
	counter := 0
	for _, combination := range finalCombinations {
		for i := 0; i < duplicate; i++ {
			dup := combination
			dup.ID = combination.ID + "x" + fmt.Sprintf("%02d", i)
			duplicated[dup.ID] = dup
			counter++
		}
	}

	return duplicated
}

func SubtractCompleted(experimentsMap *map[string]Parameter, receivedOutputPath string, experimentBaseId int) error {

	files, err := os.ReadDir(receivedOutputPath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %v", receivedOutputPath, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		delete(*experimentsMap, strings.Split(strings.Split(file.Name(), "_")[1], ".")[0])
	}

	return nil
}
