package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2006, 1, 02, 15, 04, 0, 0, time.UTC),
			want: "02 Jan 2006 at 15:04",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2006, 1, 02, 15, 04, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "02 Jan 2006 at 14:04", // Adjust to CET timezone
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			if hd != tt.want {
				t.Errorf("got %q, want %q", hd, tt.want)
			}
		})
	}
}
