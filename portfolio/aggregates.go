package portfolio

import "strings"

func GetHoldingsByUser(userId uint, category string) []HoldingSecurity {
	var holdings []HoldingSecurity
	category = strings.ToLower(category)
	switch category {
	case "stock":
		stockHoldings := GetStockHoldingsByUser(userId)
		for _, sh := range stockHoldings {
			holdings = append(holdings, sh.getHoldings())
		}
	case "mf":
		mfHoldings := GetMFHoldingsByUser(userId)
		for _, mh := range mfHoldings {
			holdings = append(holdings, mh.getHoldings())
		}
	case "fd":
		fdHoldings := GetFDsByUser(userId)
		for _, mh := range fdHoldings {
			holdings = append(holdings, mh.getHoldings())
		}
	case "listed-ncd":
		listedNcdHoldings := GetListedNCDByUser(userId)
		for _, mh := range listedNcdHoldings {
			holdings = append(holdings, mh.getHoldings())
		}
	default:
		stockHoldings := GetStockHoldingsByUser(userId)
		for _, sh := range stockHoldings {
			holdings = append(holdings, sh.getHoldings())
		}
		mfHoldings := GetMFHoldingsByUser(userId)
		for _, mh := range mfHoldings {
			holdings = append(holdings, mh.getHoldings())
		}
		fdHoldings := GetFDsByUser(userId)
		for _, mh := range fdHoldings {
			holdings = append(holdings, mh.getHoldings())
		}
		listedNcdHoldings := GetListedNCDByUser(userId)
		for _, mh := range listedNcdHoldings {
			holdings = append(holdings, mh.getHoldings())
		}
	}

	return holdings
}
func GetHoldingsAggregatesByUser(userId uint) map[string]float64 {
	var holdingsAggregates = make(map[string]float64)
	var stTotal, mfTotal, fdTotal, listedNcdTotal float64 = 0, 0, 0, 0
	stockHoldings := GetStockHoldingsByUser(userId)
	for _, sh := range stockHoldings {
		stTotal += (sh.Price * float64(sh.Quantity))
	}
	mfHoldings := GetMFHoldingsByUser(userId)
	for _, mh := range mfHoldings {
		mfTotal += (mh.Price * float64(mh.Quantity))
	}
	fdHoldings := GetFDsByUser(userId)
	for _, mh := range fdHoldings {
		fdTotal += mh.FixedDeposit.Amount
	}
	listedNcdHoldings := GetListedNCDByUser(userId)
	for _, mh := range listedNcdHoldings {
		listedNcdTotal += (mh.DirtyPrice * mh.Quantity)
	}
	holdingsAggregates["stock"] = stTotal
	holdingsAggregates["fd"] = fdTotal
	holdingsAggregates["listed-ncd"] = listedNcdTotal
	holdingsAggregates["mf"] = mfTotal
	holdingsAggregates["total"] = mfTotal + stTotal + listedNcdTotal + fdTotal
	return holdingsAggregates
}
