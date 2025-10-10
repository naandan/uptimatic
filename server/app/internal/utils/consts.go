package utils

var AllowedIntervals = []int{60, 120, 300, 600, 900, 1800, 3600}

func ContainsInt(target int) bool {
	for _, allowed := range AllowedIntervals {
		if allowed == target {
			return true
		}
	}
	return false
}
