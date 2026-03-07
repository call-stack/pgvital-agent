package collector

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func collectIndexes(ctx context.Context, pool *pgxpool.Pool) ([]IndexEntry, error) {
	query := `
		SELECT
			s.schemaname,
			s.relname,
			s.indexrelname,
			COALESCE(s.idx_scan, 0),
			COALESCE(s.idx_tup_read, 0),
			COALESCE(s.idx_tup_fetch, 0),
			pg_relation_size(s.indexrelid)
		FROM pg_stat_user_indexes s
		ORDER BY pg_relation_size(s.indexrelid) DESC
		LIMIT 200
	`
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []IndexEntry
	for rows.Next() {
		var e IndexEntry
		err := rows.Scan(
			&e.SchemaName, &e.TableName, &e.IndexName,
			&e.IdxScan, &e.IdxTupRead, &e.IdxTupFetch,
			&e.IndexSizeBytes,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}
