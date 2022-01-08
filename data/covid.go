package data

import (
	"encoding/json"
	"github.com/thetkpark/lineman-wongnai-intern/covid"
	"os"
)

type DiskCovidCasesDataStore struct {
	dataFilePath string
}

func NewDiskCovidCasesDataStore(filePath string) *DiskCovidCasesDataStore {
	return &DiskCovidCasesDataStore{dataFilePath: filePath}
}

func (s *DiskCovidCasesDataStore) Read() ([]covid.Case, error) {
	jsonFile, err := os.Open(s.dataFilePath)
	if err != nil {
		return nil, err
	}
	var covidCaseData covid.CasesJSONData
	if err := json.NewDecoder(jsonFile).Decode(&covidCaseData); err != nil {
		return nil, err
	}
	return covidCaseData.Data, nil
}
