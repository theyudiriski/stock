package postgres

import (
	"context"
	"stock/internal/service"
	"time"
)

func NewHolidayStore(
	db *Client,
) service.HolidayStore {
	return &holidayStore{
		db: db,
	}
}

type holidayStore struct {
	db *Client
}

func (s *holidayStore) GetHolidaySet(ctx context.Context, from, to time.Time) (map[string]bool, error) {
	tx, err := s.db.Leader.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, `
	SELECT
		date
	FROM
		holidays
	WHERE
		date BETWEEN $1 AND $2;
	`,
		from.Format(time.DateOnly),
		to.Format(time.DateOnly),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]bool)
	for rows.Next() {
		var d time.Time
		if err := rows.Scan(&d); err != nil {
			return nil, err
		}
		result[d.Format(time.DateOnly)] = true
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}

func (s *holidayStore) AddHoliday(ctx context.Context, date time.Time) error {
	tx, err := s.db.Leader.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO holidays (
			date
		) VALUES ($1) 
		ON CONFLICT (date)
		DO NOTHING;
	`,
		date.Format(time.DateOnly),
	)

	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *holidayStore) DeleteHoliday(ctx context.Context, date time.Time) error {
	tx, err := s.db.Leader.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		DELETE	FROM holidays WHERE date = $1;
	`,
		date.Format(time.DateOnly),
	)

	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
