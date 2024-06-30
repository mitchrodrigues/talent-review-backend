package helpers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOccurances(t *testing.T) {
	tests := []struct {
		name        string
		occurrences []string
		wantCounts  map[string]OccurrenceCount
	}{
		{
			name:        "Mixed Occurrences",
			occurrences: []string{"Alice", "Bob", "Alice"},
			wantCounts: map[string]OccurrenceCount{
				"Alice": {Value: 2, Percentage: 66.67},
				"Bob":   {Value: 1, Percentage: 33.33},
			},
		},
		{
			name:        "Single Occurrence",
			occurrences: []string{"Alice"},
			wantCounts: map[string]OccurrenceCount{
				"Alice": {Value: 1, Percentage: 100.00},
			},
		},
		{
			name:        "No Occurrences",
			occurrences: []string{},
			wantCounts:  map[string]OccurrenceCount{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			occ := NewOccurrences()

			for _, occurrence := range tt.occurrences {
				occ.Add(occurrence)
			}

			assert.Equal(t, len(tt.wantCounts), len(occ.Counts))
			for key, wantCount := range tt.wantCounts {
				gotCount, exists := occ.Counts[key]
				assert.True(t, exists)
				assert.InEpsilon(t, wantCount.Value, gotCount.Value, 0.01, "Value mismatch for key: %s", key)
				assert.InEpsilon(t, wantCount.Percentage, gotCount.Percentage, 0.01, "Percentage mismatch for key: %s", key)
			}
		})
	}
}

func TestOccurancesSub(t *testing.T) {
	tests := []struct {
		name       string
		addValues  []string
		subValues  []string
		wantCounts map[string]OccurrenceCount
	}{
		{
			name:      "Subtract Occurrences",
			addValues: []string{"Alice", "Bob", "Alice", "Charlie", "Bob"},
			subValues: []string{"Alice", "Bob"},
			wantCounts: map[string]OccurrenceCount{
				"Alice":   {Value: 1, Percentage: 33.33},
				"Bob":     {Value: 1, Percentage: 33.33},
				"Charlie": {Value: 1, Percentage: 33.33},
			},
		},
		{
			name:      "Subtract Non-Existing Occurrence",
			addValues: []string{"Alice", "Charlie"},
			subValues: []string{"Bob"},
			wantCounts: map[string]OccurrenceCount{
				"Alice":   {Value: 1, Percentage: 50.00},
				"Charlie": {Value: 1, Percentage: 50.00},
			},
		},
		{
			name:       "Subtract All Occurrences",
			addValues:  []string{"Alice", "Alice"},
			subValues:  []string{"Alice", "Alice"},
			wantCounts: map[string]OccurrenceCount{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			occ := NewOccurrences()

			for _, value := range tt.addValues {
				occ.Add(value)
			}

			fmt.Printf("Counts: %+v\n", occ.Counts)

			for _, value := range tt.subValues {
				occ.Sub(value)
			}

			assert.Equal(t, len(tt.wantCounts), len(occ.Counts))

			for key, wantCount := range tt.wantCounts {
				gotCount, exists := occ.Counts[key]

				assert.True(t, exists)
				assert.InEpsilon(t, wantCount.Value, gotCount.Value, 0.01)
				assert.InEpsilon(t, wantCount.Percentage, gotCount.Percentage, 0.01)
			}
		})
	}
}

func TestTopK(t *testing.T) {
	testCases := []struct {
		name         string
		occurrences  Occurrences
		k            int
		expectedTopK OccurrencePairList
	}{
		{
			name: "Top 3 of 4",
			occurrences: Occurrences{
				Counts: map[string]OccurrenceCount{
					"A": {Value: 10, Percentage: 50.0},
					"B": {Value: 15, Percentage: 75.0},
					"C": {Value: 5, Percentage: 25.0},
					"D": {Value: 20, Percentage: 100.0},
				},
			},
			k: 3,
			expectedTopK: OccurrencePairList{
				{"D", OccurrenceCount{Value: 20, Percentage: 100.0}},
				{"B", OccurrenceCount{Value: 15, Percentage: 75.0}},
				{"A", OccurrenceCount{Value: 10, Percentage: 50.0}},
			},
		},
		{
			name: "Top 1 of 2",
			occurrences: Occurrences{
				Counts: map[string]OccurrenceCount{
					"X": {Value: 7, Percentage: 35.0},
					"Y": {Value: 13, Percentage: 65.0},
				},
			},
			k: 1,
			expectedTopK: OccurrencePairList{
				{"Y", OccurrenceCount{Value: 13, Percentage: 65.0}},
			},
		},
		// Add more test cases as necessary...
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualTopK := tc.occurrences.TopK(tc.k)

			assert.Equal(t, tc.expectedTopK, actualTopK, "TopK results do not match expected")
		})
	}
}
