package covid

import (
	"net/http"
)

type Context interface {
	JSON(status int, v interface{})
	Error(status int, err error)
}

type DataStore interface {
	Read() ([]Case, error)
}

type Handler struct {
	store DataStore
}

func NewHandler(store DataStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) Summarize(c Context) {
	covidCases, err := h.store.Read()
	if err != nil {
		c.Error(http.StatusInternalServerError, err)
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

	c.JSON(http.StatusOK, map[string]interface{}{
		"Province": provinces,
		"AgeGroup": ageGroup,
	})
}
