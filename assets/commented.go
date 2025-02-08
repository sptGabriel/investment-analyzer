package assets

// func generateReport(
// 	trades []trade,
// 	pricesA, pricesB []price,
// 	start, end time.Time,
// 	interval time.Duration,
// ) []report {
// 	valorInicial := 100000.0
// 	cash := valorInicial
// 	portfolio := map[string]int{"A": 0, "B": 0}

// 	var reports []report

// 	reports = append(reports, report{
// 		Timestamp:              start,
// 		PatrimonioTotal:        valorInicial,
// 		RentabilidadeAcumulada: 0.0,
// 	})

// 	windowStart := start
// 	for windowStart.Before(end) {
// 		windowEnd := windowStart.Add(interval)
// 		if windowEnd.After(end) {
// 			windowEnd = end
// 		}

// 		for _, tr := range trades {
// 			if tr.Time.After(windowEnd) {
// 				break
// 			}
// 			if tr.Time.Before(windowStart) {
// 				continue
// 			}

// 			if tr.Side == "BUY" {
// 				cash -= tr.Price * float64(tr.Quantity)
// 				portfolio[tr.Symbol] += tr.Quantity
// 			} else {
// 				cash += tr.Price * float64(tr.Quantity)
// 				portfolio[tr.Symbol] -= tr.Quantity
// 			}
// 		}

// 		lastPriceA := getPriceAtTime(pricesA, windowEnd)
// 		lastPriceB := getPriceAtTime(pricesB, windowEnd)

// 		aQuantity := float64(portfolio["A"])
// 		bQuantity := float64(portfolio["B"])

// 		totalAmountInA := lastPriceA * aQuantity
// 		totalAmountInB := lastPriceB * bQuantity

// 		patrimonioTotal := cash + totalAmountInA + totalAmountInB

// 		rentabilidadeAcumulada := (patrimonioTotal - valorInicial) / valorInicial

// 		r := report{
// 			Timestamp:              windowEnd,
// 			PatrimonioTotal:        patrimonioTotal,
// 			RentabilidadeAcumulada: rentabilidadeAcumulada,
// 		}

// 		reports = append(reports, r)

// 		windowStart = windowEnd
// 	}

// 	return reports
// }
