package time

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	tests := []struct {
		in      string
		wantErr bool
		want    time.Duration
	}{
		{"3 hour", false, 3 * time.Hour},
		{"3hour", false, 3 * time.Hour},
		{"3h", false, 3 * time.Hour},
		{"3Hour", false, 3 * time.Hour},
		{"3H", false, 3 * time.Hour},
		{"+3 hour", false, 3 * time.Hour},
		{"- 3 hour", false, -3 * time.Hour},
		{"12 days", false, 12 * 24 * time.Hour},
		{"4min", false, 4 * time.Minute},
		{"5 hour 6 min", false, 5*time.Hour + 6*time.Minute},
		{"7 hour - 9sec", false, 7*time.Hour - 9*time.Second},
		{"7day+4days-3days", false, 8 * 24 * time.Hour},
		{"0.5 hour", false, 30 * time.Minute},
		{"1w1d1h1m1s", false, 8*24*time.Hour + time.Hour + time.Minute + time.Second},
		// want error
		{"2 years", true, time.Duration(0)},
		{"7", true, time.Duration(0)},
	}
	for _, tt := range tests {
		got, err := ParseDuration(tt.in)
		if !tt.wantErr && err != nil {
			t.Fatalf("ParseDuration(\"%s\") error: %v", tt.in, err)
		}
		if got != tt.want {
			t.Errorf("ParseDuration(\"%s\"): got %d want %d", tt.in, got, tt.want)
		}
	}
}
