package data

import (
	"encoding/json"
	"fmt"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/require"
	"github.com/thetkpark/lineman-wongnai-intern/covid"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
)

func NewTestFileOnDisk(content []byte, fileName string) error {
	return ioutil.WriteFile(fileName, content, os.ModePerm)
}

func TestReadJsonDataFile(t *testing.T) {
	testFilePath := fmt.Sprintf("test-%d.json", rand.Intn(1000))
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
	var casesData covid.CasesJSONData
	require.NoError(t, json.Unmarshal(jsonString, &casesData))
	require.NoError(t, NewTestFileOnDisk(jsonString, testFilePath))
	defer func() {
		require.NoError(t, os.Remove(testFilePath))
	}()
	store := &DiskCovidCasesDataStore{dataFilePath: testFilePath}
	cases, err := store.Read()
	require.NoError(t, err)
	require.Nil(t, deep.Equal(cases, casesData.Data))

}

func TestReadNotExistFile(t *testing.T) {
	store := NewDiskCovidCasesDataStore("not-existed.json")
	cases, err := store.Read()
	require.Error(t, err)
	require.Nil(t, cases)
}

func TestReadInvalidJsonDataFile(t *testing.T) {
	testFilePath := fmt.Sprintf("test-%d.json", rand.Intn(1000))
	jsonString := []byte(`{"Data": error`)

	err := ioutil.WriteFile(testFilePath, jsonString, os.ModePerm)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.Remove(testFilePath))
	}()
	store := NewDiskCovidCasesDataStore(testFilePath)

	got, err := store.Read()
	require.Error(t, err)
	require.Nil(t, got)
}
