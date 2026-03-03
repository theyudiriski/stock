package postgres

import (
	"context"
	"fmt"
	"stock/internal/service"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

var (
	loc, _ = time.LoadLocation("Asia/Jakarta")
)

func NewCorpactionStore(
	db *Client,
) service.CorpactionStore {
	return &corpactionStore{
		db: db,
	}
}

type corpactionStore struct {
	db *Client
}

func (s *corpactionStore) UpsertCorpactions(
	ctx context.Context,
	symbol string,
	corpaction *service.Corpaction,
) error {
	tx, err := s.db.Leader.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if corpaction.Dividends != nil {
		for _, data := range *corpaction.Dividends {
			_, err = tx.ExecContext(ctx, `
				INSERT INTO corpaction_dividend (
					symbol,
					amount,
					currency,
					cum_date,
					ex_date,
					recording_date,
					payment_date
				) VALUES ($1, $2, $3, $4, $5, $6, $7)
				ON CONFLICT (symbol, cum_date, ex_date, recording_date, payment_date)
				DO NOTHING;
			`,
				symbol,
				func() float64 {
					amount, err := decimal.NewFromString(data.Amount)
					if err != nil {
						return 0
					}
					return amount.InexactFloat64()
				}(),
				func() string {
					parts := strings.Split(data.Currency, "_")
					return parts[len(parts)-1]
				}(),
				data.CumDate,
				data.ExDate,
				data.RecordingDate,
				data.PaymentDate,
			)
			if err != nil {
				return fmt.Errorf("failed to insert dividend: %w", err)
			}
		}
	}

	if corpaction.RUPS != nil {
		for _, data := range *corpaction.RUPS {
			_, err = tx.ExecContext(ctx, `
				INSERT INTO corpaction_rups (
					symbol,
					datetime,
					venue,
					eligible_date
				) VALUES ($1, $2, $3, $4)
				ON CONFLICT (symbol, eligible_date)
				DO NOTHING;
			`,
				symbol,
				func() *time.Time {
					if data.Date != "0001-01-01" {
						datetime, err := time.ParseInLocation("2006-01-02 15:04:05",
							data.Date+" "+data.Time+":00", loc)
						if err != nil {
							return nil
						}
						return &datetime
					}
					return nil
				}(),
				data.Venue,
				data.EligibleDate,
			)
			if err != nil {
				return fmt.Errorf("failed to insert rups: %w", err)
			}
		}
	}

	if corpaction.PublicExpose != nil {
		for _, data := range *corpaction.PublicExpose {
			_, err = tx.ExecContext(ctx, `
				INSERT INTO corpaction_public_expose (
					symbol,
					datetime,
					venue
				) VALUES ($1, $2, $3)
				ON CONFLICT (symbol, datetime)
				DO NOTHING;
		`,
				symbol,
				func() time.Time {
					datetime, err := time.ParseInLocation("2006-01-02 15:04:05",
						data.Date+" "+data.Time, loc)
					if err != nil {
						return time.Time{}
					}
					return datetime
				}(),
				data.Venue,
			)
			if err != nil {
				return fmt.Errorf("failed to insert public expose: %w", err)
			}
		}
	}

	if corpaction.RightIssue != nil {
		for _, data := range *corpaction.RightIssue {
			_, err = tx.ExecContext(ctx, `
				INSERT INTO corpaction_right_issue (
					symbol,
					price,
					old,
					new,
					cum_date,
					ex_date,
					recording_date,
					start_date,
					end_date
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
				ON CONFLICT (symbol, cum_date)
				DO NOTHING;
		`,
				symbol,
				data.Price,
				func() int {
					old, err := strconv.Atoi(data.Old)
					if err != nil {
						return 0
					}
					return old
				}(),
				func() int {
					new, err := strconv.Atoi(data.New)
					if err != nil {
						return 0
					}
					return new
				}(),
				data.CumDate,
				data.ExDate,
				data.RecordingDate,
				data.StartDate,
				data.EndDate,
			)
			if err != nil {
				return fmt.Errorf("failed to insert right issue: %w", err)
			}
		}
	}

	if corpaction.StockSplit != nil {
		for _, data := range *corpaction.StockSplit {
			_, err = tx.ExecContext(ctx, `
			INSERT INTO corpaction_stock_split (
				symbol,
				old,
				new,
				cum_date,
				ex_date,
				recording_date
			) VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (symbol, cum_date)
			DO NOTHING;
	`,
				symbol,
				func() int {
					old, err := strconv.Atoi(data.Old)
					if err != nil {
						return 0
					}
					return old
				}(),
				func() int {
					new, err := strconv.Atoi(data.New)
					if err != nil {
						return 0
					}
					return new
				}(),
				data.CumDate,
				data.ExDate,
				data.RecordingDate,
			)
			if err != nil {
				return fmt.Errorf("failed to insert stock split: %w", err)
			}
		}
	}

	if corpaction.ReverseSplit != nil {
		for _, data := range *corpaction.ReverseSplit {
			_, err = tx.ExecContext(ctx, `
			INSERT INTO corpaction_reverse_split (
				symbol,
				old,
				new,
				cum_date,
				ex_date,
				recording_date
			) VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (symbol, cum_date)
			DO NOTHING;
	`,
				symbol,
				func() int {
					old, err := strconv.Atoi(data.Old)
					if err != nil {
						return 0
					}
					return old
				}(),
				func() int {
					new, err := strconv.Atoi(data.New)
					if err != nil {
						return 0
					}
					return new
				}(),
				data.CumDate,
				data.ExDate,
				data.RecordingDate,
			)
			if err != nil {
				return fmt.Errorf("failed to insert reverse split: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
