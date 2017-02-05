package utils

import (
	"strconv"
	"time"
)

func UnixMicroStringToTime(s string) (time.Time, error) {
	validate := NewValidate()
	validate.RegisterValidation("unixmicro", IsTimeUnixMicro)
	if err := validate.Var(s, "unixmicro"); err != nil {
		return time.Time{}, err
	}

	sec, _ := strconv.ParseInt(s[:10], 10, 64)
	nsec, _ := strconv.ParseInt(s[10:], 10, 64)

	return time.Unix(sec, nsec), nil
}
