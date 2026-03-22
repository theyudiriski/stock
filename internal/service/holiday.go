package service

import (
	"context"
	"time"
)

type HolidayStore interface {
	GetHolidaySet(ctx context.Context, from, to time.Time) (map[string]bool, error)
	AddHoliday(ctx context.Context, date time.Time) error
	DeleteHoliday(ctx context.Context, date time.Time) error
}
