package postgres

import (
	"context"
	"fmt"
	"stock/internal/service"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

func NewEmittenStore(
	db *Client,
) service.EmittenStore {
	return &emittenStore{
		db: db,
	}
}

type emittenStore struct {
	db *Client
}

func (s *emittenStore) GetEmittens(ctx context.Context) ([]string, error) {
	rows, err := s.db.Leader.QueryContext(ctx, `
		SELECT symbol FROM emittens;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	symbols := make([]string, 0)
	for rows.Next() {
		var symbol string
		err := rows.Scan(&symbol)
		if err != nil {
			return nil, err
		}
		symbols = append(symbols, symbol)
	}

	return symbols, nil
}

func (s *emittenStore) UpsertEmittenProfile(ctx context.Context, symbol string, profile *service.EmittenProfile) error {
	tx, err := s.db.Leader.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i := range profile.ShareholderNumbers {
		parsedDate, err := time.Parse("2 Jan 2006", profile.ShareholderNumbers[i].ShareholderDate)
		if err != nil {
			return fmt.Errorf("failed to parse shareholder_date: %w", err)
		}

		totalShareDecimal, err := decimal.NewFromString(strings.ReplaceAll(profile.ShareholderNumbers[i].TotalShare, ",", ""))
		if err != nil {
			return fmt.Errorf("failed to parse total_share: %w", err)
		}
		totalShare := totalShareDecimal.IntPart()

		_, err = tx.ExecContext(ctx, `
			INSERT INTO emitten_shareholder_numbers (
				symbol,
				shareholder_date,
				total_share,
				change
			) VALUES ($1, $2, $3, $4)
			ON CONFLICT (symbol, shareholder_date) 
			DO NOTHING;
		`,
			symbol,
			parsedDate.Format(time.DateOnly),
			totalShare,
			profile.ShareholderNumbers[i].Change,
		)
		if err != nil {
			return fmt.Errorf("failed to insert shareholder number: %w", err)
		}
	}

	_, err = tx.ExecContext(ctx, `
		DELETE FROM
			emitten_shareholders
		WHERE
			symbol = $1;
	`,
		symbol,
	)
	if err != nil {
		return fmt.Errorf("failed to delete old shareholders: %w", err)
	}

	// Aggregate shareholders by name, summing percentages and merging badges
	type aggregatedShareholder struct {
		percentage decimal.Decimal
		badges     map[string]bool // Using map for unique badges
	}
	shareholderMap := make(map[string]*aggregatedShareholder)

	// Process regular shareholders
	for i := range profile.Shareholder {
		percentageStr := strings.TrimSuffix(profile.Shareholder[i].Percentage, "%")
		percentageStr = strings.TrimPrefix(percentageStr, "<")

		percentage, err := decimal.NewFromString(percentageStr)
		if err != nil {
			return fmt.Errorf("failed to parse shareholder percentage: %w", err)
		}
		percentage = percentage.Div(decimal.NewFromInt(100))

		name := profile.Shareholder[i].Name
		if _, exists := shareholderMap[name]; !exists {
			shareholderMap[name] = &aggregatedShareholder{
				percentage: decimal.Zero,
				badges:     make(map[string]bool),
			}
		}

		shareholderMap[name].percentage = shareholderMap[name].percentage.Add(percentage)
		for _, badge := range profile.Shareholder[i].Badges {
			shareholderMap[name].badges[badge] = true
		}
	}

	// Process director/commissioner shareholders
	for i := range profile.ShareholderDirectorCommissioner {
		percentageStr := strings.TrimSuffix(profile.ShareholderDirectorCommissioner[i].Percentage, "%")
		percentageStr = strings.TrimPrefix(percentageStr, "<")

		percentage, err := decimal.NewFromString(percentageStr)
		if err != nil {
			return fmt.Errorf("failed to parse director/commissioner percentage: %w", err)
		}
		percentage = percentage.Div(decimal.NewFromInt(100))

		name := profile.ShareholderDirectorCommissioner[i].Name
		if _, exists := shareholderMap[name]; !exists {
			shareholderMap[name] = &aggregatedShareholder{
				percentage: decimal.Zero,
				badges:     make(map[string]bool),
			}
		}

		shareholderMap[name].percentage = shareholderMap[name].percentage.Add(percentage)
		for _, badge := range profile.ShareholderDirectorCommissioner[i].Badges {
			shareholderMap[name].badges[badge] = true
		}
	}

	// Insert aggregated shareholders
	for name, shareholder := range shareholderMap {
		// Convert badges map to slice
		badges := make([]string, 0, len(shareholder.badges))
		for badge := range shareholder.badges {
			badges = append(badges, badge)
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO emitten_shareholders (
				symbol,
				shareholder_name,
				shareholder_percentage,
				shareholder_badges
			) VALUES ($1, $2, $3, $4);
		`,
			symbol,
			name,
			shareholder.percentage,
			badges,
		)
		if err != nil {
			return fmt.Errorf("failed to insert shareholder: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *emittenStore) UpsertEmittenProfileInfo(
	ctx context.Context,
	symbol string,
	info *service.EmittenInfo,
	profile *service.EmittenProfile,
) error {
	tx, err := s.db.Leader.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO emitten_profiles (
			symbol,
			name,
			description,
			underwriters,
			free_float,
			subsector_id
		) VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (symbol)
		DO NOTHING;
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		symbol,
		info.Name,
		profile.Description,
		pq.Array(func() []string {
			normalized := make([]string, len(profile.History.Underwriters))
			for i, u := range profile.History.Underwriters {
				s := strings.TrimPrefix(u, "PT.")
				if strings.HasSuffix(s, "Tbk") && !strings.HasSuffix(s, "Tbk.") {
					s += "."
				}
				normalized[i] = strings.TrimSpace(s)
			}
			return normalized
		}()),
		func() *float64 {
			freeFloat, err := decimal.NewFromString(strings.TrimSuffix(profile.History.FreeFloat, "%"))
			if err != nil {
				return nil
			}
			return GetPointer(freeFloat.Div(decimal.NewFromInt(100)).InexactFloat64())
		}(),
		func() int {
			for i := range info.Catalogs {
				if info.Catalogs[i].CompanyType == "sub_sector" {
					subsectorID, err := strconv.Atoi(info.Catalogs[i].ID)
					if err != nil {
						return 0
					}
					return subsectorID
				}
			}
			return 0
		}(),
	)
	if err != nil {
		return fmt.Errorf("failed to insert emitten profile: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func GetPointer[T any](val T) *T {
	return &val
}

func (s *emittenStore) GetEmittensUnderwriters(ctx context.Context) ([]service.EmittenUnderwriters, error) {
	rows, err := s.db.Leader.QueryContext(ctx, `
		SELECT
			symbol,
			underwriters
		FROM emitten_profiles
		ORDER BY
			symbol;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]service.EmittenUnderwriters, 0)
	for rows.Next() {
		var result service.EmittenUnderwriters
		err := rows.Scan(
			&result.Symbol,
			pq.Array(&result.Underwriters),
		)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func (s *emittenStore) UpdateEmittenUnderwriterCode(ctx context.Context, symbol string, brokers []service.Broker) error {
	tx, err := s.db.Leader.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		UPDATE
			emitten_profiles
		SET
			underwriters_code = $1
		WHERE
			symbol = $2;
	`,
		pq.Array(func() []string {
			codes := make([]string, len(brokers))
			for i := range brokers {
				codes[i] = brokers[i].Code
			}
			return codes
		}()),
		symbol,
	)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
