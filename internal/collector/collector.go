package collector

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Collector struct {
	pool     *pgxpool.Pool
	interval time.Duration
}

func New(pool *pgxpool.Pool, interval time.Duration) *Collector {
	return &Collector{pool: pool, interval: interval}
}

// Collect runs one round of stat collection and returns a Snapshot.
func (c *Collector) Collect(ctx context.Context) (*Snapshot, error) {
	snap := &Snapshot{CollectedAt: time.Now().UTC()}

	stmts, err := collectStatements(ctx, c.pool)
	if err != nil {
		return nil, fmt.Errorf("statements: %w", err)
	}
	snap.Statements = stmts

	tables, err := collectTables(ctx, c.pool)
	if err != nil {
		return nil, fmt.Errorf("tables: %w", err)
	}
	snap.Tables = tables

	indexes, err := collectIndexes(ctx, c.pool)
	if err != nil {
		return nil, fmt.Errorf("indexes: %w", err)
	}
	snap.Indexes = indexes

	return snap, nil
}
