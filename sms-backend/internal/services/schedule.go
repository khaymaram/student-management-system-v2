package services

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"sms-backend/internal/models"
)

var weekdayOrder = map[string]int{"M": 0, "T": 1, "W": 2, "Th": 3, "F": 4}

func normalizeSchedule(days []string, startTime string) ([]string, string, error) {
	if len(days) == 0 {
		return nil, "", errors.New("select at least one meeting day")
	}

	seen := make(map[string]bool, len(days))
	normalizedDays := make([]string, 0, len(days))
	for _, day := range days {
		if _, valid := weekdayOrder[day]; !valid {
			return nil, "", fmt.Errorf("invalid meeting day %q", day)
		}
		if !seen[day] {
			seen[day] = true
			normalizedDays = append(normalizedDays, day)
		}
	}
	sort.Slice(normalizedDays, func(i, j int) bool {
		return weekdayOrder[normalizedDays[i]] < weekdayOrder[normalizedDays[j]]
	})

	minutes, err := scheduleMinutes(startTime)
	if err != nil || minutes < 8*60 || minutes > 16*60 || minutes%30 != 0 {
		return nil, "", errors.New("start time must be on the hour or half hour between 08:00 and 16:00")
	}

	return normalizedDays, fmt.Sprintf("%02d:%02d", minutes/60, minutes%60), nil
}

func scheduleMinutes(value string) (int, error) {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return 0, errors.New("invalid time")
	}
	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}
	minute, err := strconv.Atoi(parts[1])
	if err != nil || minute < 0 || minute > 59 {
		return 0, errors.New("invalid time")
	}
	return hour*60 + minute, nil
}

func schedulesOverlap(a, b *models.Course) bool {
	days := make(map[string]bool, len(a.MeetingDays))
	for _, day := range a.MeetingDays {
		days[day] = true
	}
	sharesDay := false
	for _, day := range b.MeetingDays {
		if days[day] {
			sharesDay = true
			break
		}
	}
	if !sharesDay {
		return false
	}

	aStart, aErr := scheduleMinutes(a.StartTime)
	bStart, bErr := scheduleMinutes(b.StartTime)
	if aErr != nil || bErr != nil {
		return false
	}
	return aStart < bStart+60 && bStart < aStart+60
}
