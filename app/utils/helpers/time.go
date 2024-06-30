package helpers

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/gertd/go-pluralize"
)

var (
	pclient = pluralize.NewClient()

	HumanTimeFormat = "01/02/2006 15:04:05"
)

// HumanizeDuration takes a duration and returns a humanized string ie: 1 minute ago, 2 hours ago, 3 days ago, 1 year and 2 months ago
func HumanizeDuration(duration time.Duration) string {
	if duration.Minutes() < 60.0 {
		return fmt.Sprintf("%s ago", pclient.Pluralize("minute", int(duration.Minutes()), true))
	}

	if duration.Hours() < 24.0 {
		return fmt.Sprintf("%s ago",
			pclient.Pluralize("hours", int(duration.Hours()), true))
	}

	days := duration.Hours() / 24
	if days < 30.4167 {
		return fmt.Sprintf("%s ago",
			pclient.Pluralize("days", int(days), true))
	}

	months := days / 30.4167
	if months < 12 {
		return fmt.Sprintf("%s ago",
			pclient.Pluralize("months", int(months), true))
	}

	years := math.Floor(months / 12)
	remainingMonths := int(math.Mod(months, 12))

	if remainingMonths == 0 {
		return fmt.Sprintf("%s ago",
			pclient.Pluralize("years", int(years), true))
	}

	return fmt.Sprintf("%s and %s ago",
		pclient.Pluralize("years", int(years), true),
		pclient.Pluralize("months", int(remainingMonths), true),
	)
}

func GenerateRandomDate(years int) time.Time {
	if years < 0 {
		years = years + 1
	} else {
		years = years - 1
	}

	// Adjust the date by the specified number of years
	result := time.Now().AddDate(years, 0, 0)

	// Generate a random number of months (from 0 to 11)
	randomMonths := rand.Intn(12) // 0 to 11

	// Generate a random number of days (from 0 to daysInMonth-1)
	year, month, _ := result.Date()
	daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, result.Location()).Day()
	randomDays := rand.Intn(daysInMonth) // 0 to daysInMonth-1

	// Adjust the month and day
	result = result.AddDate(0, -randomMonths, -randomDays)

	return result
}
