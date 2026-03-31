package rate_test

import (
	"testing"

	"calcLoanAmount/internal/rate"
)

func TestBuildSegments(t *testing.T) {
	segments, err := rate.BuildSegments(12, 0.12, []rate.Change{
		{Month: 7, AnnualRate: 0.06},
	})
	if err != nil {
		t.Fatalf("BuildSegments returned error: %v", err)
	}

	if len(segments) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(segments))
	}

	if segments[0].StartMonth != 1 || segments[0].EndMonth != 6 || segments[0].AnnualRate != 0.12 {
		t.Fatalf("unexpected first segment: %+v", segments[0])
	}

	if segments[1].StartMonth != 7 || segments[1].EndMonth != 12 || segments[1].AnnualRate != 0.06 {
		t.Fatalf("unexpected second segment: %+v", segments[1])
	}
}

func TestBuildSegments_InvalidChangeMonth(t *testing.T) {
	_, err := rate.BuildSegments(12, 0.12, []rate.Change{
		{Month: 13, AnnualRate: 0.06},
	})
	if err == nil {
		t.Fatalf("expected error for out-of-range change month")
	}
}
