package db

import (
	"testing"
	"time"
)

func date(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

func TestPawukonEpochIsSunday(t *testing.T) {
	if pawukonEpoch.Weekday() != time.Sunday {
		t.Fatalf("epoch should be Redite (Sunday), got %s", pawukonEpoch.Weekday())
	}
}

func TestPawukonDatesMatchKnownHolidays(t *testing.T) {
	cases := []struct {
		name   string
		offset int
		want   time.Time
	}{
		{"Galungan", 73, date(2024, time.February, 28)},
		{"Galungan next cycle", 73, date(2024, time.September, 25)},
		{"Kuningan", 83, date(2024, time.March, 9)},
	}
	for _, c := range cases {
		got := pawukonDates(c.offset, c.want, c.want.AddDate(0, 0, 1))
		if len(got) != 1 || !got[0].Equal(c.want) {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}

func TestPawukonHolidaysFallOnTheRightWeekday(t *testing.T) {
	want := map[string]time.Weekday{
		"Hari Raya Galungan":  time.Wednesday,
		"Hari Raya Kuningan":  time.Saturday,
		"Hari Raya Saraswati": time.Saturday,
		"Hari Raya Pagerwesi": time.Wednesday,
		"Tumpek Landep":       time.Saturday,
		"Tumpek Wariga":       time.Saturday,
	}
	for _, h := range pawukonHolidays {
		for _, d := range pawukonDates(h.offset, date(2026, 1, 1), date(2028, 1, 1)) {
			if d.Weekday() != want[h.name] {
				t.Errorf("%s on %s is a %s", h.name, d.Format("2006-01-02"), d.Weekday())
			}
		}
	}
}

func TestPawukonDatesRepeatEvery210Days(t *testing.T) {
	got := pawukonDates(73, date(2026, 1, 1), date(2028, 1, 1))
	for i := 1; i < len(got); i++ {
		if diff := got[i].Sub(got[i-1]).Hours() / 24; int(diff+0.5) != 210 {
			t.Errorf("gap %d: %v days", i, diff)
		}
	}
	if len(got) < 3 {
		t.Fatalf("expected at least 3 Galungan dates in 2026-2027, got %d", len(got))
	}
}
