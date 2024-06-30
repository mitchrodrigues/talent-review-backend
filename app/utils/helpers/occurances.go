package helpers

import (
	"fmt"
	"sort"
	"strings"

	"github.com/golly-go/golly"
)

type OccurrenceCount struct {
	Value      int
	Percentage float64
}

type Occurrences struct {
	Type  string
	Total int

	Insensitive bool
	AllowEmpty  bool

	Counts map[string]OccurrenceCount
}

type OccurrencesOptions struct {
	Insensitive bool
	AllowEmpty  bool
}

// Helper type for sorting
type OccurrencePair struct {
	Key   string
	Value OccurrenceCount
}

type OccurrencePairList []OccurrencePair

func (p OccurrencePairList) Keys() []string {
	return golly.Map(p, func(pair OccurrencePair) string { return pair.Key })
}

func (p OccurrencePairList) Len() int           { return len(p) }
func (p OccurrencePairList) Less(i, j int) bool { return p[i].Value.Value > p[j].Value.Value } // Sorting in descending order
func (p OccurrencePairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func InsensitiveOccurrences() OccurrencesOptions {
	return OccurrencesOptions{Insensitive: true}
}

func EmptyOccurrences() OccurrencesOptions {
	return OccurrencesOptions{AllowEmpty: true}
}

func NewOccurrences(insensitive ...OccurrencesOptions) Occurrences {
	opts := OccurrencesOptions{}

	if len(insensitive) > 0 {
		opts = insensitive[0]
	}

	return Occurrences{Counts: make(map[string]OccurrenceCount), Insensitive: opts.Insensitive}
}

func (o Occurrences) ToOccurrencePairList() (pairs OccurrencePairList) {
	for key, value := range o.Counts {
		pairs = append(pairs, OccurrencePair{Key: key, Value: value})
	}
	return
}

func (o Occurrences) Keys() (ret []string) {
	for key := range o.Counts {
		ret = append(ret, key)
	}
	return
}

// TopK returns the top K occurrences
func (o *Occurrences) TopK(k int) OccurrencePairList {
	pairs := o.ToOccurrencePairList()

	sort.Sort(pairs)

	if k > len(pairs) {
		k = len(pairs)
	}

	return pairs[:k]
}

func (o *Occurrences) Add(value string) {

	// Do not track empty
	if value == "" {
		return
	}

	if o.Insensitive {
		value = strings.ToLower(value)
	}

	occurrence, ok := o.Counts[value]
	if !ok {
		occurrence = OccurrenceCount{}
	}

	o.Total++
	occurrence.Value++
	o.Counts[value] = occurrence

	o.CalculatePercentages()
}

func (o *Occurrences) Sub(value string) {

	if o.Insensitive {
		value = strings.ToLower(value)
	}

	occurrence, ok := o.Counts[value]
	if !ok {
		return
	}

	occurrence.Value--

	if occurrence.Value == 0 {
		delete(o.Counts, value)
		return
	}

	o.Total--
	o.Counts[value] = occurrence

	o.CalculatePercentages()
}

func (o *Occurrences) Find(value string) OccurrenceCount {
	if o.Insensitive {
		value = strings.ToLower(value)
	}

	counts, ok := o.Counts[value]
	if !ok {
		return OccurrenceCount{}
	}

	return counts
}

func (o *Occurrences) CalculatePercentages() {
	for key := range o.Counts {
		occurrence := o.Counts[key]
		occurrence.Percentage = (float64(occurrence.Value) / float64(o.Total)) * 100
		o.Counts[key] = occurrence
	}
}

func (o *Occurrences) String() string {
	var sb strings.Builder

	for val, occ := range o.Counts {
		fmt.Fprintf(&sb, "%s: %d (%.2f%%)\n", val, occ.Value, occ.Percentage)
	}

	return sb.String()
}
