package collector

import (
	"context"
	"fmt"
	"log/slog"
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

	// Lock and activity collection is best-effort — the monitoring user
	// may not have access to pg_stat_activity or pg_locks.
	locks, err := collectLocks(ctx, c.pool)
	if err != nil {
		slog.Warn("lock collection failed (non-fatal)", "err", err)
	} else {
		snap.Locks = locks
	}

	activities, err := collectActivity(ctx, c.pool)
	if err != nil {
		slog.Warn("activity collection failed (non-fatal)", "err", err)
	} else {
		snap.Activities = activities
	}

	deadlocks, err := collectDeadlocks(ctx, c.pool)
	if err != nil {
		slog.Warn("deadlock collection failed (non-fatal)", "err", err)
	} else {
		snap.Deadlocks = deadlocks
	}

	return snap, nil
}
