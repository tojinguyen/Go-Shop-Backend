package time

import (
	"testing"
	"time"
)

func TestGetUTCUnixTime(t *testing.T) {
	unixTime := GetUTCUnixTime()
	now := time.Now().UTC().Unix()

	// Should be within 1 second difference
	if unixTime < now-1 || unixTime > now+1 {
		t.Errorf("GetUTCUnixTime() returned %d, expected around %d", unixTime, now)
	}
}

func TestCalculateDurationFromSeconds(t *testing.T) {
	tests := []struct {
		seconds  int64
		expected time.Duration
	}{
		{0, 0},
		{1, time.Second},
		{60, time.Minute},
		{3600, time.Hour},
	}

	for _, test := range tests {
		result := CalculateDurationFromSeconds(test.seconds)
		if result != test.expected {
			t.Errorf("CalculateDurationFromSeconds(%d) = %v, expected %v", test.seconds, result, test.expected)
		}
	}
}

func TestGetSecondsUntilExpiry(t *testing.T) {
	now := GetUTCUnixTime()

	tests := []struct {
		name     string
		expiry   int64
		expected int64
	}{
		{"future expiry", now + 3600, 3600},
		{"past expiry", now - 3600, 0},
		{"exact now", now, 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := GetSecondsUntilExpiry(test.expiry)
			if result != test.expected {
				t.Errorf("GetSecondsUntilExpiry(%d) = %d, expected %d", test.expiry, result, test.expected)
			}
		})
	}
}

func TestParseAndFormatTime(t *testing.T) {
	timeStr := "2025-06-28T10:30:00Z"

	parsedTime, err := ParseTime(timeStr)
	if err != nil {
		t.Errorf("ParseTime failed: %v", err)
	}

	formatted := FormatTime(parsedTime)
	if formatted != timeStr {
		t.Errorf("FormatTime returned %s, expected %s", formatted, timeStr)
	}
}
