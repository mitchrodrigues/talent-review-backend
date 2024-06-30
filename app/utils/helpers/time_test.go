package helpers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHumanizeDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{"Less than an hour", time.Minute * 30, "30 minutes ago"},
		{"Exactly one hour", time.Hour, "1 hour ago"},
		{"Multiple hours", time.Hour * 5, "5 hours ago"},
		{"One day", time.Hour * 24, "1 day ago"},
		{"Multiple days", time.Hour * 48, "2 days ago"},
		{"One month", time.Hour * 731, "1 month ago"},
		{"Multiple months", time.Hour * 731 * 3, "3 months ago"},
		{"One year", time.Hour * 24 * 366, "1 year ago"},
		{"One year and a few months", time.Hour * 24 * 400, "1 year and 1 month ago"},
		{"Multiple years", time.Hour * 24 * (366 * 2), "2 years ago"},
		// Add more test cases as necessary
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := HumanizeDuration(tc.duration)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGenerateRandomDate(t *testing.T) {
	testCases := []struct {
		name   string
		years  int
		before time.Time
		after  time.Time
	}{
		{
			name:   "20 years ago",
			years:  -20,
			before: time.Now(),
			after:  time.Now().AddDate(-20, -1, 0),
		},
		{
			name:   "30 years ago",
			years:  -30,
			before: time.Now(),
			after:  time.Now().AddDate(-30, -1, 0),
		},
		{
			name:   "40 years ago",
			years:  -40,
			before: time.Now(),
			after:  time.Now().AddDate(-40, -1, 0),
		},
		{
			name:   "20 years in the future",
			years:  20,
			before: time.Now().AddDate(20, 1, 0),
			after:  time.Now(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			generatedDate := GenerateRandomDate(tc.years)

			if tc.years < 0 {
				// Check if the date is in the past
				assert.True(t, generatedDate.Before(tc.before))
				assert.True(t, generatedDate.After(tc.after))
			} else {
				// Check if the date is in the future
				assert.True(t, generatedDate.After(tc.after))
				assert.True(t, generatedDate.Before(tc.before))
			}
		})
	}
}
