package postgres

import (
	"context"
	"stock/internal/service"
)

func NewBrokerStore(
	db *Client,
) service.BrokerStore {
	return &brokerStore{
		db: db,
	}
}

type brokerStore struct {
	db *Client
}

func (s *brokerStore) GetBrokerByName(
	ctx context.Context,
	name string,
) (*service.Broker, error) {
	tx, err := s.db.Leader.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, `
		SELECT
			code,
			$1 as name
		FROM
			brokers
		WHERE
			name = $1
		OR $1 = ANY(previous_names)
		ORDER BY
			CASE WHEN name = $1 THEN 0 ELSE 1 END
		LIMIT 1;
	`, name)

	var result service.Broker
	if err = row.Scan(
		&result.Code,
		&result.Name,
	); err != nil {
		return nil, err
	}

	return &result, nil
}
