package utils

import "math"

func CheckSliceDiff(inputList, currentList []int) []int {
	exists := make(map[int]bool)
	var output []int
	for _, value := range currentList {
		exists[value] = true
	}

	for _, value := range inputList {
		if !exists[value] {
			output = append(output, value)
		}
	}

	return output
}

func CalcTotalPage(total int64, limit int) int {
	t := int(total) / limit
	if (int(total) % limit) != 0 {
		t += 1
	}
	return t
}

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius in kilometers

	// Convert degrees to radians
	lat1 = lat1 * math.Pi / 180
	lon1 = lon1 * math.Pi / 180
	lat2 = lat2 * math.Pi / 180
	lon2 = lon2 * math.Pi / 180

	// Haversine formula
	dLat := lat2 - lat1
	dLon := lon2 - lon1
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Distance in kilometers
	return R * c
}
