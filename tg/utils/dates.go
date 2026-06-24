package utils

import "time"

func YearEnd(year int) time.Time {
	return time.Date(year+1, 1, 1, 0, 0, 0, 0, time.Local).Add(-time.Second)
}

func YearBegin(year int) time.Time {
	return time.Date(year+1, 1, 1, 0, 0, 0, 0, time.Local)
}
