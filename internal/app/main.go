package app

import "container/heap"

func (rm *RouteManager) DetermineSavingList(lr *LocationRegistry) {
	var value float64
	for i := range lr.DistanceMatrix {
		for j := range lr.DistanceMatrix[i] {
			if i == 0 || j == 0 || i == j {
				continue
			}
			value = lr.retreiveValueFromPair(i, j)
			rm.addToSavingsList(i, j, value)
		}
	}
}

func (lr *LocationRegistry) retreiveValueFromPair(i, j int) float64 {

	// Clarke-Wright Algorithm Formula: d(D, i) + d(D, j) - d(i, j)
	depotToI := (lr.DistanceMatrix)[0][i]
	depotToJ := (lr.DistanceMatrix)[0][j]
	iToJ := (lr.DistanceMatrix)[i][j]

	result := depotToI + depotToJ - iToJ

	return result
}

func (rm *RouteManager) addToSavingsList(first, second int, value float64) {
	newSaving := saving{
		i:     first,
		j:     second,
		value: value,
	}
	heap.Push(&rm.SavingList, newSaving)
}

func (rm *RouteManager) TestRemoveAll() {
	rm.SavingList.popAll()
}
