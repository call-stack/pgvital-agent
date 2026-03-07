package collector

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvitals/pgvitals/agent/internal/scrubber"
)

func collectStatements(ctx context.Context, pool *pgxpool.Pool) ([]StatementEntry, error) {
	query := `
		SELECT
			query,
			calls,
			total_exec_time,
			mean_exec_time,
			min_exec_time,
			max_exec_time,
			stddev_exec_time,
			rows,
			shared_blks_hit,
			shared_blks_read
		FROM pg_stat_statements
		ORDER BY total_exec_time DESC
		LIMIT 200
	`
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []StatementEntry
	for rows.Next() {
		var e StatementEntry
		var rawQuery string
		err := rows.Scan(
			&rawQuery, &e.Calls, &e.TotalExecTime, &e.MeanExecTime,
			&e.MinExecTime, &e.MaxExecTime, &e.StddevExecTime,
			&e.Rows, &e.SharedBlksHit, &e.SharedBlksRead,
		)
		if err != nil {
			return nil, err
		}
		e.QueryText = scrubber.Scrub(rawQuery)
		e.QueryFingerprint = scrubber.Fingerprint(e.QueryText)
		entries = append(entries, e)
	}
	return entries, rows.Err()
}
