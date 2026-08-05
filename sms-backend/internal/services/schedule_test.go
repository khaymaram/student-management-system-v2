package services

import (
	"reflect"
	"testing"

	"sms-backend/internal/models"
)

func TestNormalizeSchedule(t *testing.T) {
	days, start, err := normalizeSchedule([]string{"F", "M", "M"}, "9:30")
	if err != nil {
		t.Fatalf("normalizeSchedule returned error: %v", err)
	}
	if !reflect.DeepEqual(days, []string{"M", "F"}) {
		t.Fatalf("days = %v, want [M F]", days)
	}
	if start != "09:30" {
		t.Fatalf("start = %q, want 09:30", start)
	}
}

func TestSchedulesOverlap(t *testing.T) {
	tests := []struct {
		name string
		a    models.Course
		b    models.Course
		want bool
	}{
		{
			name: "same day overlapping half hour",
			a:    models.Course{MeetingDays: []string{"M"}, StartTime: "09:00"},
			b:    models.Course{MeetingDays: []string{"M"}, StartTime: "09:30"},
			want: true,
		},
		{
			name: "adjacent classes do not overlap",
			a:    models.Course{MeetingDays: []string{"T"}, StartTime: "09:00"},
			b:    models.Course{MeetingDays: []string{"T"}, StartTime: "10:00"},
			want: false,
		},
		{
			name: "same time on different days",
			a:    models.Course{MeetingDays: []string{"W"}, StartTime: "11:00"},
			b:    models.Course{MeetingDays: []string{"Th"}, StartTime: "11:00"},
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := schedulesOverlap(&test.a, &test.b); got != test.want {
				t.Fatalf("schedulesOverlap() = %v, want %v", got, test.want)
			}
		})
	}
}
