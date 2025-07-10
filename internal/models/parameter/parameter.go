package parameter

import (
	"experiment-server/internal/config"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	base = 10000
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
	IsMultiModal         int     `json:"is_multi_modal"`
	Done                 bool    `json:"-"`
}

func (p Parameter) Print() {
	log.Printf(
		"Parameter{ID: %s, GenerationLimit: %d, PopulationSize: %d, KClosestCsPercantage: %f, A: %d, B: %d, EtaAr: %d, GammaAr: %d, GammaMgr: %.2f, DistanceMethod: %d, DeviatonEraParameter: %d, DatasetIndex: %d, IsMultiModal: %d}\n",
		p.ID, p.GenerationLimit, p.PopulationSize, p.KClosestCsPercantage, p.A, p.B, p.EtaAr, p.GammaAr, p.GammaMgr, p.DistanceMethod, p.DeviatonEraParameter, p.DatasetIndex, p.IsMultiModal,
	)
}

func (p Parameter) String() string {
	return fmt.Sprintf("Parameter{ID: %s, GenerationLimit: %d, PopulationSize: %d, KClosestCsPercantage: %f, A: %d, B: %d, EtaAr: %d, GammaAr: %d, GammaMgr: %.2f, DistanceMethod: %d, DeviatonEraParameter: %d, DatasetIndex: %d, IsMultiModal: %d}\n",
		p.ID, p.GenerationLimit, p.PopulationSize, p.KClosestCsPercantage, p.A, p.B, p.EtaAr, p.GammaAr, p.GammaMgr, p.DistanceMethod, p.DeviatonEraParameter, p.DatasetIndex, p.IsMultiModal,
	)
}

func GenerateParamCombinations(duplicate int, cfg *config.Config) map[string]Parameter {

	// Cartesian elements
	populationSizes := []int{500, 250}
	kClosestCs := []float32{1.0}									// Ignore
	etaArs := []int{1200}
	gammaArs := []int{10}
	gammaMgrs := []float64{1.0}
	datasetIndexes := []int{0, 2, 5}                                // 0: Kilyos, 1: Baylands, 2: Baylands2408, 3: - , 4: Marmaracık, 5: Marmaracık4096
	AnB := [][2]int{{1, 1}, {2, 1}, {3, 2}, {3, 3}, {4, 2}, {5, 3}} // Ar to Mgr
	DistanceMethod := []int{1}                                      // 0: Normal, 1: Alternative
	DeviationEraParameter := []int{-1, 30}                          // -1: No Dev Era, 30: Normal
	IsMultiModal := []int{0, 1}										// 0: N-Best, 1: MultiModal

	ExperimentSessionID := cfg.ExperimentBaseId * base

	type combo struct {
		pop          int
		k            float32
		eta          int
		gammaA       int
		gammaM       float64
		dset         int
		a            int
		b            int
		dm           int
		devEra       int
		IsMultiModal int
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
										for _, isMulti := range IsMultiModal {
											combos = append(combos, combo{
												pop, k, eta, gammaA, gammaM, dset, ab[0], ab[1], dm, devEra, isMulti,
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
			IsMultiModal:         c.IsMultiModal,
			Done: 				  false,
		}
		finalCombinations = append(finalCombinations, param)
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
