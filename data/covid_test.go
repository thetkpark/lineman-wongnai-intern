package data

import (
	"encoding/json"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/require"
	"github.com/thetkpark/lineman-wongnai-intern/covid"
	"io/ioutil"
	"os"
	"testing"
)

const testFilePath = "test.json"

func TestReadJsonDataFile(t *testing.T) {
	jsonString := []byte(`{ "Data": [{
    "ConfirmDate": "2021-05-01",
    "No": null,
    "Age": 91,
    "Gender": "ชาย",
    "GenderEn": "Male",
    "Nation": null,
    "NationEn": "India",
    "Province": "Kamphaeng Phet",
    "ProvinceId": 14,
    "District": null,
    "ProvinceEn": "Kamphaeng Phet",
    "StatQuarantine": 12
  }, {
    "ConfirmDate": null,
    "No": null,
    "Age": 92,
    "Gender": "ชาย",
    "GenderEn": "Male",
    "Nation": null,
    "NationEn": "China",
    "Province": "Nonthaburi",
    "ProvinceId": 35,
    "District": null,
    "ProvinceEn": "Nonthaburi",
    "StatQuarantine": 1
  }]}`)
	var expect covid.CasesJSONData
	err := json.Unmarshal(jsonString, &expect)
	require.NoError(t, err)

	err = ioutil.WriteFile(testFilePath, jsonString, os.ModePerm)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.Remove(testFilePath))
	}()
	store := NewDiskCovidCasesDataStore(testFilePath)

	got, err := store.Read()
	require.NoError(t, err)
	require.Nil(t, deep.Equal(got, expect.Data))
}
