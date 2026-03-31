package rate

import (
	"errors"
	"sort"
)

type Change struct {
	Month      int
	AnnualRate float64
}

type Segment struct {
	StartMonth int
	EndMonth   int
	AnnualRate float64
}

var (
	errInvalidTotalMonths = errors.New("total months must be greater than 0")
	errInvalidChangeMonth = errors.New("change month out of range")
	errInvalidRate        = errors.New("annual rate must be >= 0")
	errDuplicateMonth     = errors.New("duplicate change month")
)

func BuildSegments(totalMonths int, initialRate float64, changes []Change) ([]Segment, error) {
	if totalMonths <= 0 {
		return nil, errInvalidTotalMonths
	}
	if initialRate < 0 {
		return nil, errInvalidRate
	}

	sorted := append([]Change(nil), changes...)
	sort.Slice(sorted, func(i int, j int) bool {
		return sorted[i].Month < sorted[j].Month
	})

	segments := make([]Segment, 0, len(sorted)+1)
	currentStart := 1
	currentRate := initialRate
	seenMonth := map[int]struct{}{}

	for _, change := range sorted {
		if change.Month < 1 || change.Month > totalMonths {
			return nil, errInvalidChangeMonth
		}
		if change.AnnualRate < 0 {
			return nil, errInvalidRate
		}
		if _, ok := seenMonth[change.Month]; ok {
			return nil, errDuplicateMonth
		}
		seenMonth[change.Month] = struct{}{}

		if change.Month == currentStart {
			currentRate = change.AnnualRate
			continue
		}

		segments = append(segments, Segment{
			StartMonth: currentStart,
			EndMonth:   change.Month - 1,
			AnnualRate: currentRate,
		})
		currentStart = change.Month
		currentRate = change.AnnualRate
	}

	segments = append(segments, Segment{
		StartMonth: currentStart,
		EndMonth:   totalMonths,
		AnnualRate: currentRate,
	})

	return segments, nil
}
