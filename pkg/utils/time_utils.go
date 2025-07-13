package utils

import "time"

const (
	TimeFormat = "2006-01-02 15:04:05"
)

func GetCurrentTime() string {
	return time.Now().Format(TimeFormat)
}

func AddToCurrentTime(duration time.Duration) string {
	return time.Now().Add(duration).Format(TimeFormat)
}

func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse(TimeFormat, timeStr)
}

func IsTimeExpired(expiryTime string) (bool, error) {
	parsedTime, err := ParseTime(expiryTime)
	if err != nil {
		return false, err
	}
	return time.Now().After(parsedTime), nil
}
