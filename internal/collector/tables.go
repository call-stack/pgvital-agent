package collector

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func collectTables(ctx context.Context, pool *pgxpool.Pool) ([]TableEntry, error) {
	query := `
		SELECT
			schemaname,
			relname,
			COALESCE(seq_scan, 0),
			COALESCE(seq_tup_read, 0),
			COALESCE(idx_scan, 0),
			COALESCE(idx_tup_fetch, 0),
			COALESCE(n_tup_ins, 0),
			COALESCE(n_tup_upd, 0),
			COALESCE(n_tup_del, 0),
			COALESCE(n_live_tup, 0),
			COALESCE(n_dead_tup, 0),
			last_vacuum,
			last_autovacuum,
			last_analyze
		FROM pg_stat_user_tables
		ORDER BY n_live_tup DESC
		LIMIT 200
	`
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []TableEntry
	for rows.Next() {
		var e TableEntry
		err := rows.Scan(
			&e.SchemaName, &e.TableName,
			&e.SeqScan, &e.SeqTupRead, &e.IdxScan, &e.IdxTupFetch,
			&e.NTupIns, &e.NTupUpd, &e.NTupDel,
			&e.NLiveTup, &e.NDeadTup,
			&e.LastVacuum, &e.LastAutovacuum, &e.LastAnalyze,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}
