package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"stock/internal/service"
	"strconv"
	"time"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

func NewShareholderStore(
	logger *slog.Logger,
	db *Client,
) service.ShareholderStore {
	return &shareholderStore{
		logger: logger,
		db:     db,
	}
}

type shareholderStore struct {
	logger *slog.Logger
	db     *Client
}

var (
	investorCategoryMap = map[string][]service.InvestorCategory{
		"Individual":              {service.InvestorCategoryIndividual},
		"Perusahaan":              {service.InvestorCategoryCorporate},
		"Reksadana":               {service.InvestorCategoryMutualFund},
		"Bank":                    {service.InvestorCategoryFinancialInstitution},
		"Sekuritas":               {service.InvestorCategorySecuritiesCompany},
		"Dana Pensiun & Asuransi": {service.InvestorCategoryPensionFund, service.InvestorCategoryInsurance},
		"Dana Pensiun":            {service.InvestorCategoryPensionFund},
		"Asuransi":                {service.InvestorCategoryInsurance},
		"Lainnya":                 {service.InvestorCategoryOthers},
	}
)

func (s *shareholderStore) GetShareholderChartHistory(ctx context.Context, symbol string) (lastUpdatedDate *string, err error) {
	tx, err := s.db.Leader.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, `
		SELECT
			last_updated_date
		FROM
			emitten_shareholder_chart_history
		WHERE
			symbol = $1
	`, symbol)

	if err = row.Scan(&lastUpdatedDate); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get shareholder chart history: %w", err)
	}

	return lastUpdatedDate, nil
}

func (s *shareholderStore) UpsertShareholderChart(
	ctx context.Context,
	symbol string,
	data *service.ShareholderChartData,
) error {
	tx, err := s.db.Leader.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO emitten_shareholder_chart (
			symbol,
			investor_category_codes,
			investor_type,
			date,
			percentage
		) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (symbol, investor_category_codes, investor_type, date)
		DO NOTHING;
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for i := range data.Shareholder {
		for j := range data.Shareholder[i].ChartData {
			item := data.Shareholder[i].ChartData[j]

			_, err = stmt.ExecContext(ctx,
				symbol,
				pq.Array(investorCategoryMap[data.Shareholder[i].ItemName]),
				data.Shareholder[i].InvestorType,
				func() string {
					sec, err := strconv.ParseInt(item.UnixDate, 10, 64)
					if err != nil {
						s.logger.Error("failed to parse unix date", "unixDate", item.UnixDate, "error", err)
						return item.UnixDate
					}
					return time.Unix(sec, 0).In(loc).Format("2006-01-02")
				}(),
				decimal.NewFromFloat(item.Value).Div(decimal.NewFromInt(100)).InexactFloat64(),
			)
			if err != nil {
				return fmt.Errorf("failed to insert shareholder chart: %w", err)
			}
		}
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO emitten_shareholder_chart_history (
			symbol,
			last_updated_date
		) VALUES ($1, $2)
		ON CONFLICT (symbol)
		DO NOTHING;
	`, symbol, data.LastUpdate)
	if err != nil {
		return fmt.Errorf("failed to insert shareholder chart history: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
