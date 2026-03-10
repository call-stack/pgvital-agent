package collector

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvitals/pgvitals/agent/internal/scrubber"
)

func collectLocks(ctx context.Context, pool *pgxpool.Pool) ([]LockEntry, error) {
	query := `
		SELECT
			blocked_activity.pid,
			blocked_activity.query,
			blocking_activity.pid,
			blocking_activity.query,
			COALESCE(blocking_activity.state, 'unknown'),
			COALESCE(blocked_locks.relation::regclass::text, ''),
			blocked_locks.mode,
			EXTRACT(EPOCH FROM (now() - blocked_activity.query_start))
		FROM pg_locks blocked_locks
		JOIN pg_stat_activity blocked_activity
			ON blocked_activity.pid = blocked_locks.pid
		JOIN pg_locks blocking_locks
			ON blocking_locks.locktype = blocked_locks.locktype
			AND blocking_locks.relation IS NOT DISTINCT FROM blocked_locks.relation
			AND blocking_locks.page IS NOT DISTINCT FROM blocked_locks.page
			AND blocking_locks.tuple IS NOT DISTINCT FROM blocked_locks.tuple
			AND blocking_locks.virtualxid IS NOT DISTINCT FROM blocked_locks.virtualxid
			AND blocking_locks.transactionid IS NOT DISTINCT FROM blocked_locks.transactionid
			AND blocking_locks.classid IS NOT DISTINCT FROM blocked_locks.classid
			AND blocking_locks.objid IS NOT DISTINCT FROM blocked_locks.objid
			AND blocking_locks.objsubid IS NOT DISTINCT FROM blocked_locks.objsubid
			AND blocking_locks.pid != blocked_locks.pid
			AND blocking_locks.granted
		JOIN pg_stat_activity blocking_activity
			ON blocking_activity.pid = blocking_locks.pid
		WHERE NOT blocked_locks.granted
		ORDER BY EXTRACT(EPOCH FROM (now() - blocked_activity.query_start)) DESC
		LIMIT 50
	`
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []LockEntry
	for rows.Next() {
		var e LockEntry
		var blockedQuery, blockingQuery string
		err := rows.Scan(
			&e.BlockedPID, &blockedQuery,
			&e.BlockingPID, &blockingQuery,
			&e.BlockingState, &e.LockedTable,
			&e.LockMode, &e.WaitDurationSec,
		)
		if err != nil {
			return nil, err
		}
		e.BlockedQuery = scrubber.Scrub(blockedQuery)
		e.BlockingQuery = scrubber.Scrub(blockingQuery)
		entries = append(entries, e)
	}
	return entries, rows.Err()
}
