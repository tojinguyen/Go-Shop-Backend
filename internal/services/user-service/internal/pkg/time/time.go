package time

import "time"

// GetUtcTime returns the current UTC time
func GetUtcTime() time.Time {
	return time.Now().UTC()
}

// GetUtcTimeString returns the current UTC time in RFC3339 format
func GetUtcTimeString() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// GetUTCUnixTime returns the current UTC time as Unix timestamp
func GetUTCUnixTime() int64 {
	return time.Now().UTC().Unix()
}

// ParseTime parses a time string in RFC3339 format
func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}

// FormatTime formats a time to RFC3339 format
func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// CalculateDurationFromSeconds calculates duration from seconds
func CalculateDurationFromSeconds(seconds int64) time.Duration {
	return time.Duration(seconds) * time.Second
}

// GetSecondsUntilExpiry calculates seconds until expiry from current time
func GetSecondsUntilExpiry(expiryUnix int64) int64 {
	now := GetUTCUnixTime()
	if expiryUnix > now {
		return expiryUnix - now
	}
	return 0
}
