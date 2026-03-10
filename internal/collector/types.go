package collector

import "time"

// Snapshot is the data payload sent to the PGVitals server per collection tick.
type Snapshot struct {
	CollectedAt  time.Time        `json:"collected_at"`
	InstanceID   string           `json:"instance_id"`
	Statements   []StatementEntry `json:"statements"`
	Tables       []TableEntry     `json:"tables"`
	Indexes      []IndexEntry     `json:"indexes"`
	Locks        []LockEntry      `json:"locks"`
	Activities   []ActivityEntry  `json:"activities"`
	Deadlocks    int64            `json:"deadlocks"`
	AgentVersion string           `json:"agent_version"`
}

type StatementEntry struct {
	QueryFingerprint string  `json:"query_fingerprint"`
	QueryText        string  `json:"query_text"`
	Calls            int64   `json:"calls"`
	TotalExecTime    float64 `json:"total_exec_time_ms"`
	MeanExecTime     float64 `json:"mean_exec_time_ms"`
	MinExecTime      float64 `json:"min_exec_time_ms"`
	MaxExecTime      float64 `json:"max_exec_time_ms"`
	StddevExecTime   float64 `json:"stddev_exec_time_ms"`
	Rows             int64   `json:"rows"`
	SharedBlksHit    int64   `json:"shared_blks_hit"`
	SharedBlksRead   int64   `json:"shared_blks_read"`
}

type TableEntry struct {
	SchemaName     string     `json:"schema_name"`
	TableName      string     `json:"table_name"`
	SeqScan        int64      `json:"seq_scan"`
	SeqTupRead     int64      `json:"seq_tup_read"`
	IdxScan        int64      `json:"idx_scan"`
	IdxTupFetch    int64      `json:"idx_tup_fetch"`
	NTupIns        int64      `json:"n_tup_ins"`
	NTupUpd        int64      `json:"n_tup_upd"`
	NTupDel        int64      `json:"n_tup_del"`
	NLiveTup       int64      `json:"n_live_tup"`
	NDeadTup       int64      `json:"n_dead_tup"`
	LastVacuum     *time.Time `json:"last_vacuum"`
	LastAutovacuum *time.Time `json:"last_autovacuum"`
	LastAnalyze    *time.Time `json:"last_analyze"`
}

type IndexEntry struct {
	SchemaName     string `json:"schema_name"`
	TableName      string `json:"table_name"`
	IndexName      string `json:"index_name"`
	IdxScan        int64  `json:"idx_scan"`
	IdxTupRead     int64  `json:"idx_tup_read"`
	IdxTupFetch    int64  `json:"idx_tup_fetch"`
	IndexSizeBytes int64  `json:"index_size_bytes"`
}

type LockEntry struct {
	BlockedPID      int32   `json:"blocked_pid"`
	BlockedQuery    string  `json:"blocked_query"`
	BlockingPID     int32   `json:"blocking_pid"`
	BlockingQuery   string  `json:"blocking_query"`
	BlockingState   string  `json:"blocking_state"`
	LockedTable     string  `json:"locked_table"`
	LockMode        string  `json:"lock_mode"`
	WaitDurationSec float64 `json:"wait_duration_seconds"`
}

type ActivityEntry struct {
	WaitEventType string `json:"wait_event_type"`
	WaitEvent     string `json:"wait_event"`
	Count         int64  `json:"count"`
}
