package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success":   true,
			"timestamp": time.Now(),
		})
	})

	router.GET("/covid", func(c *gin.Context) {
		covidCases, err := ReadCovidCaseDataFile()
		if err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}

		provinces := make(map[string]uint)
		ageGroup := map[string]uint{
			"0-30":  0,
			"31-60": 0,
			"61+":   0,
			"N/A":   0,
		}

		for _, covidCase := range covidCases {
			provinceName := covidCase.Province
			if len(provinceName) == 0 {
				provinceName = "N/A"
			}
			provinces[provinceName]++

			age := "N/A"
			if covidCase.Age != nil {
				if *covidCase.Age >= 0 && *covidCase.Age <= 30 {
					age = "0-30"
				} else if *covidCase.Age >= 31 && *covidCase.Age <= 60 {
					age = "31-60"
				} else if *covidCase.Age >= 61 {
					age = "61+"
				}
			}

			ageGroup[age]++
		}

		c.JSON(200, gin.H{
			"Province": provinces,
			"AgeGroup": ageGroup,
		})
	})

	router.Run()
}

func ReadCovidCaseDataFile() ([]CovidCase, error) {
	jsonFile, err := os.Open("covid-case.json")
	if err != nil {
		return nil, err
	}
	var covidCaseData CovidCaseJSONData
	if err := json.NewDecoder(jsonFile).Decode(&covidCaseData); err != nil {
		return nil, err
	}
	return covidCaseData.Data, nil
}

type CovidCase struct {
	ConfirmDate    string `json:"ConfirmDate"`
	No             int    `json:"No"`
	Age            *int   `json:"Age"`
	Gender         string `json:"Gender"`
	GenderEn       string `json:"GenderEn"`
	Nation         string `json:"Nation"`
	NationEn       string `json:"NationEn"`
	Province       string `json:"Province"`
	ProvinceID     int    `json:"ProvinceId"`
	District       string `json:"District"`
	ProvinceEn     string `json:"ProvinceEn"`
	StatQuarantine int    `json:"StatQuarantine"`
}

type CovidCaseJSONData struct {
	Data []CovidCase `json:"Data"`
}
