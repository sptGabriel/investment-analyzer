package assets

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func round(val float64, places int) float64 {
	factor := math.Pow(10, float64(places))
	return math.Round(val*factor) / factor
}

func calculateAccumulatedReturn(initialValue, currentValue float64) float64 {
	if initialValue == 0 {
		return 0
	}
	return round((currentValue/initialValue)-1, 5)
}

func generateReport(
	trades []trade,
	pricesA, pricesB []price,
	start, end time.Time,
	interval time.Duration,
) []report {
	portfolio := map[string]int{"A": 0, "B": 0}

	initialCash := 100000.0
	cash := initialCash

	reports := []report{
		{
			Timestamp:              start,
			PatrimonioTotal:        cash,
			RentabilidadeAcumulada: 0,
		},
	}

	currentTime := start
	for currentTime.Before(end) {
		endDate := currentTime.Add(interval)
		if endDate.After(end) {
			endDate = end
		}

		intervalTrades := filterTrades(trades, currentTime, endDate)

		for _, tr := range intervalTrades {
			if tr.Side == "BUY" {
				cash -= tr.Price * float64(tr.Quantity)
				portfolio[tr.Symbol] += tr.Quantity

				continue
			}

			cash += tr.Price * float64(tr.Quantity)
			portfolio[tr.Symbol] -= tr.Quantity
		}

		lastPriceA := getPriceAtTime(pricesA, endDate)
		lastPriceB := getPriceAtTime(pricesB, endDate)

		aQuantity := float64(portfolio["A"])
		bQuantity := float64(portfolio["B"])

		totalAmountInA := lastPriceA * aQuantity
		totalAmountInB := lastPriceB * bQuantity

		patrimonioTotal := cash + totalAmountInA + totalAmountInB

		rentabilidadeAcumulada := calculateAccumulatedReturn(initialCash, patrimonioTotal)

		reports = append(reports, report{
			Timestamp:              endDate,
			PatrimonioTotal:        patrimonioTotal,
			RentabilidadeAcumulada: rentabilidadeAcumulada,
		})

		currentTime = currentTime.Add(interval)
	}

	return reports
}

func filterTrades(trades []trade, start, end time.Time) []trade {
	var filteredTrades []trade
	for _, tr := range trades {
		if !tr.Time.Before(start) && tr.Time.Before(end) {
			filteredTrades = append(filteredTrades, tr)
		}
	}
	return filteredTrades
}

func TestAssets2(t *testing.T) {
	pricesA, err := readPrices("march_2021_pricesA.csv")
	require.NoError(t, err, "on read prices a csv")

	pricesB, err := readPrices("march_2021_pricesB.csv")
	require.NoError(t, err, "on read prices b csv")

	trades, err := readTrades("march_2021_trades.csv")
	require.NoError(t, err, "on read trades csv")

	start, err := time.Parse("2006-01-02 15:04:05", "2021-03-01 10:00:00")
	require.NoError(t, err, "on parse start date to correct layout")

	end, err := time.Parse("2006-01-02 15:04:05", "2021-03-07 17:50:00")
	require.NoError(t, err, "on parse end date to correct layout")

	interval := 10 * time.Minute

	reports := generateReport(trades, pricesA, pricesB, start, end, interval)

	expectedFirst := []report{
		{Timestamp: start, PatrimonioTotal: 100000.0, RentabilidadeAcumulada: 0.00000},
		{Timestamp: start.Add(10 * time.Minute), PatrimonioTotal: 100024.0, RentabilidadeAcumulada: 0.00024},
		{Timestamp: start.Add(20 * time.Minute), PatrimonioTotal: 99919.0, RentabilidadeAcumulada: -0.00081},
	}

	expectedLast := []report{
		{Timestamp: end.Add(-20 * time.Minute), PatrimonioTotal: 99575.0, RentabilidadeAcumulada: -0.00425},
		{Timestamp: end.Add(-10 * time.Minute), PatrimonioTotal: 98972.0, RentabilidadeAcumulada: -0.01028},
		{Timestamp: end, PatrimonioTotal: 99397.0, RentabilidadeAcumulada: -0.00603},
	}

	assert.Equal(t, expectedFirst, reports[:3], "As primeiras 3 linhas do relatório estão incorretas")
	assert.Equal(t, expectedLast, reports[len(reports)-3:], "As últimas 3 linhas do relatório estão incorretas")
}

func getPriceAtTime(prices []price, t time.Time) float64 {
	for _, p := range prices {
		if p.Time.Time.Equal(t) {
			return p.Price
		}
	}

	return 0.0
}
