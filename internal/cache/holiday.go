package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"stock/internal/service"
	"time"

	redis "github.com/redis/go-redis/v9"
)

const (
	yearCacheTTL = 365 * 24 * time.Hour
)

func NewHolidayCache(
	redisClient RedisClient,
	holidayStore service.HolidayStore,
) service.HolidayStore {
	return &holidayCache{
		redisClient:  redisClient,
		holidayStore: holidayStore,
	}
}

type holidayCache struct {
	redisClient  RedisClient
	holidayStore service.HolidayStore
}

func (c *holidayCache) GetHolidaySet(ctx context.Context, from, to time.Time) (map[string]bool, error) {
	result := make(map[string]bool)

	for year := from.Year(); year <= to.Year(); year++ {
		yearSet, err := c.getYearHolidays(ctx, year)
		if err != nil {
			return nil, err
		}

		for d := range yearSet {
			// only include dates within the requested [from, to] window
			t, err := time.Parse(time.DateOnly, d)
			if err != nil {
				return nil, err
			}

			if !t.Before(from) && !t.After(to) {
				result[d] = true
			}
		}
	}

	return result, nil
}

// AddHoliday writes to postgres then invalidates the affected year's cache.
func (c *holidayCache) AddHoliday(ctx context.Context, date time.Time) error {
	if err := c.holidayStore.AddHoliday(ctx, date); err != nil {
		return err
	}

	if err := c.redisClient.Del(ctx, yearCacheKey(date.Year())); err != nil {
		return err
	}

	return nil
}

// DeleteHoliday deletes from postgres then invalidates the affected year's cache.
func (c *holidayCache) DeleteHoliday(ctx context.Context, date time.Time) error {
	if err := c.holidayStore.DeleteHoliday(ctx, date); err != nil {
		return err
	}

	if err := c.redisClient.Del(ctx, yearCacheKey(date.Year())); err != nil {
		return err
	}

	return nil
}

// getYearHolidays returns all holidays for a given year, using cache when available.
func (c *holidayCache) getYearHolidays(ctx context.Context, year int) (map[string]bool, error) {
	key := yearCacheKey(year)

	// try to get from cache
	val, err := c.redisClient.Get(ctx, key)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return nil, err
		}
	} else {
		var dates []string
		if err := json.Unmarshal([]byte(val), &dates); err != nil {
			return nil, err
		}
		result := make(map[string]bool, len(dates))
		for _, d := range dates {
			result[d] = true
		}
		return result, nil
	}

	yearFrom := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	yearTo := time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC)

	// fallback to postgres
	set, err := c.holidayStore.GetHolidaySet(ctx, yearFrom, yearTo)
	if err != nil {
		return nil, err
	}

	dates := make([]string, 0, len(set))
	for d := range set {
		dates = append(dates, d)
	}

	data, err := json.Marshal(dates)
	if err != nil {
		return nil, err
	}

	// store to cache
	if err := c.redisClient.Set(ctx, key, string(data), yearCacheTTL); err != nil {
		return nil, err
	}

	return set, nil
}

func yearCacheKey(year int) string {
	return fmt.Sprintf("holidays:%d", year)
}
