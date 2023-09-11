package rates

import (
	"encoding/json"
	"github.com/mprokocki/interview-o/getclient"
	"math"
	"time"
)

const epsilon = 1e-9 // or another small value, depending on the required precision

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < epsilon
}

type exchangeRates struct {
	Rates []rate `json:"rates"`
}

type rate struct {
	No            string  `json:"no"`
	EffectiveDate string  `json:"effectiveDate"`
	Mid           float64 `json:"mid"`
}

func (r rate) inRange(min float64, max float64) bool {
	return (r.Mid > min || almostEqual(r.Mid, min)) && (r.Mid < max || almostEqual(r.Mid, max))
}

func NotInRange(min, max float64) (map[time.Time]float64, error) {
	client := getclient.NewClient(nil)

	meta, err := client.Get("http://api.nbp.pl/api/exchangerates/rates/a/eur/last/100/?format=json")
	if err != nil {
		return nil, err
	}

	rates := &exchangeRates{}
	err = json.Unmarshal([]byte(meta.Content), rates)
	if err != nil {
		return nil, err
	}

	notInRange := make(map[time.Time]float64, 0)
	for _, r := range rates.Rates {
		if !r.inRange(min, max) {
			date, err := time.Parse("2006-01-02", r.EffectiveDate)
			if err != nil {
				return nil, err
			}

			notInRange[date] = r.Mid
		}
	}

	return notInRange, nil
}
