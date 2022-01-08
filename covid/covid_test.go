package covid

import (
	"encoding/json"
	"errors"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type MockDataStore struct {
	mock.Mock
}

func (m *MockDataStore) Read() ([]Case, error) {
	args := m.Called()
	return args.Get(0).([]Case), args.Error(1)
}

type MockContext struct {
	mock.Mock
	response map[string]interface{}
	status   int
}

func (m *MockContext) JSON(status int, v interface{}) {
	_ = m.Called(status, v)
	m.response = v.(map[string]interface{})
	m.status = status
}

func (m *MockContext) Error(status int, err error) {
	_ = m.Called(status, err)
	m.response = map[string]interface{}{
		"error": err.Error(),
	}
	m.status = status
}

func TestCovidSummarizeHandler(t *testing.T) {
	suite.Run(t, new(CovidHandlerTestSuite))
}

type CovidHandlerTestSuite struct {
	suite.Suite
	cases     []Case
	handler   *Handler
	ctx       *MockContext
	store     *MockDataStore
	summerize map[string]interface{}
}

func (s *CovidHandlerTestSuite) TestSummarize() {
	s.store.On("Read").Return(s.cases, nil)
	s.ctx.On("JSON", http.StatusOK, s.summerize)
	s.handler.Summarize(s.ctx)
}

func (s *CovidHandlerTestSuite) TestSummarizeReadCasesError() {
	err := errors.New("error reading data")
	s.store.On("Read").Return([]Case{}, err)
	s.ctx.On("Error", http.StatusInternalServerError, err)
	s.handler.Summarize(s.ctx)
	require.Nil(s.T(), deep.Equal(s.ctx.response, map[string]interface{}{"error": err.Error()}))
}

func (s *CovidHandlerTestSuite) SetupTest() {
	s.ctx = new(MockContext)
	s.store = new(MockDataStore)
	s.handler = NewHandler(s.store)
	casesJson := []byte(`[{
    "ConfirmDate": "2021-05-01",
    "No": null,
    "Age": 94,
    "Gender": "ชาย",
    "GenderEn": "Male",
    "Nation": null,
    "NationEn": null,
    "Province": "Roi Et",
    "ProvinceId": 53,
    "District": null,
    "ProvinceEn": "Roi Et",
    "StatQuarantine": 12
  }, {
    "ConfirmDate": null,
    "No": null,
    "Age": 64,
    "Gender": "ชาย",
    "GenderEn": "Male",
    "Nation": null,
    "NationEn": null,
    "Province": "Songkhla",
    "ProvinceId": 63,
    "District": null,
    "ProvinceEn": "Songkhla",
    "StatQuarantine": 2
  }, {
    "ConfirmDate": "2021-05-01",
    "No": null,
    "Age": null,
    "Gender": "หญิง",
    "GenderEn": "Female",
    "Nation": null,
    "NationEn": "India",
    "Province": "Chai Nat",
    "ProvinceId": 6,
    "District": null,
    "ProvinceEn": "Chai Nat",
    "StatQuarantine": 1
  }, {
    "ConfirmDate": "2021-05-01",
    "No": null,
    "Age": 35,
    "Gender": null,
    "GenderEn": null,
    "Nation": null,
    "NationEn": "USA",
    "Province": null,
    "ProvinceId": null,
    "District": null,
    "ProvinceEn": null,
    "StatQuarantine": 0
  }, {
    "ConfirmDate": "2021-05-02",
    "No": null,
    "Age": 26,
    "Gender": null,
    "GenderEn": null,
    "Nation": null,
    "NationEn": "Thailand",
    "Province": "Kamphaeng Phet",
    "ProvinceId": 14,
    "District": null,
    "ProvinceEn": "Kamphaeng Phet",
    "StatQuarantine": 17
  }]`)
	_ = json.Unmarshal(casesJson, &s.cases)

	s.summerize = map[string]interface{}{
		"Province": map[string]uint{
			"Chai Nat":       1,
			"Songkhla":       1,
			"Roi Et":         1,
			"Kamphaeng Phet": 1,
			"N/A":            1,
		},
		"AgeGroup": map[string]uint{
			"0-30":  1,
			"31-60": 1,
			"61+":   2,
			"N/A":   1,
		},
	}
}
