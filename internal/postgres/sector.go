package postgres

import (
	"context"
	"fmt"
	"stock/internal/service"
)

func NewSectorStore(
	db *Client,
) service.SectorStore {
	return &sectorStore{
		db: db,
	}
}

type sectorStore struct {
	db *Client
}

func (s *sectorStore) UpsertSubsectors(ctx context.Context, sectors []service.Sector) error {
	tx, err := s.db.Leader.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO emitten_sectors (
			id,
			name,
			parent_id
		)
		VALUES ($1, $2, $3)
		ON CONFLICT (id)
		DO NOTHING;
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for i := range sectors {
		_, err = stmt.ExecContext(ctx,
			sectors[i].ID,
			sectors[i].Name,
			sectors[i].ParentID,
		)
		if err != nil {
			return fmt.Errorf("failed to insert subsector: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
