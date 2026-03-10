package collector

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func collectActivity(ctx context.Context, pool *pgxpool.Pool) ([]ActivityEntry, error) {
	query := `
		SELECT
			wait_event_type,
			wait_event,
			count(*)
		FROM pg_stat_activity
		WHERE wait_event IS NOT NULL
		GROUP BY wait_event_type, wait_event
		ORDER BY count(*) DESC
		LIMIT 100
	`
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []ActivityEntry
	for rows.Next() {
		var e ActivityEntry
		err := rows.Scan(&e.WaitEventType, &e.WaitEvent, &e.Count)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func collectDeadlocks(ctx context.Context, pool *pgxpool.Pool) (int64, error) {
	var deadlocks int64
	err := pool.QueryRow(ctx,
		`SELECT deadlocks FROM pg_stat_database WHERE datname = current_database()`,
	).Scan(&deadlocks)
	if err != nil {
		return 0, err
	}
	return deadlocks, nil
}
