package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"stock/internal/service"
	"strings"

	"github.com/shopspring/decimal"
)

func NewSubsidiaryStore(
	db *Client,
	logger *slog.Logger,
) service.SubsidiaryStore {
	return &subsidiaryStore{
		db:     db,
		logger: logger,
	}
}

type subsidiaryStore struct {
	db     *Client
	logger *slog.Logger
}

func (s *subsidiaryStore) GetSubsidiaryCompanyHistory(ctx context.Context, symbol string) (lastUpdatedPeriod *string, err error) {
	tx, err := s.db.Leader.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, `
		SELECT
			last_updated_period
		FROM
			emitten_subsidiary_companies_history
		WHERE
			symbol = $1
	`, symbol)

	if err = row.Scan(&lastUpdatedPeriod); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get subsidiary company history: %w", err)
	}

	return lastUpdatedPeriod, nil
}

func (s *subsidiaryStore) UpsertSubsidiaryCompanies(ctx context.Context, symbol string, data *service.SubsidiaryData) error {
	tx, err := s.db.Leader.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		DELETE FROM
			emitten_subsidiary_companies
		WHERE
			symbol = $1;
	`, symbol)
	if err != nil {
		return fmt.Errorf("failed to delete subsidiary companies: %w", err)
	}

	if len(data.Subsidiaries) > 0 {
		stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO emitten_subsidiary_companies (
			symbol,
			subsidiary_company_name,
			subsidiary_company_percentage,
			subsidiary_company_type
		) VALUES ($1, $2, $3, $4)
		ON CONFLICT (symbol, subsidiary_company_name)
		DO NOTHING;
	`)
		if err != nil {
			return fmt.Errorf("failed to prepare statement: %w", err)
		}
		defer stmt.Close()

		for i := range data.Subsidiaries {
			name := strings.TrimSpace(data.Subsidiaries[i].CompanyName)
			if name != "" && name != "N/A" && name != "-" && name != "0" {
				_, err = stmt.ExecContext(ctx,
					symbol,
					name,
					func() *float64 {
						percentage, err := decimal.NewFromString(data.Subsidiaries[i].Percentage)
						if err != nil {
							s.logger.Error("failed to parse subsidiary percentage", "symbol", symbol, "error", err)
							return nil
						}
						result := percentage.Div(decimal.NewFromInt(100)).InexactFloat64()
						if result > 1 {
							return nil
						}
						return &result
					}(),
					data.Subsidiaries[i].BusinessType,
				)
				if err != nil {
					return fmt.Errorf("failed to insert subsidiary company: %w", err)
				}
			}
		}
	}

	if data.LastUpdatedPeriod != "" {
		_, err = tx.ExecContext(ctx, `
		INSERT INTO emitten_subsidiary_companies_history (
			symbol,
			last_updated_period
		) VALUES ($1, $2)
		ON CONFLICT (symbol)
		DO NOTHING;
	`, symbol, data.LastUpdatedPeriod)
		if err != nil {
			return fmt.Errorf("failed to insert subsidiary company history: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
