package innodblock_test

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"
)

// schemaDDL rebuilds the fixture from scratch. job_attachment rows are
// contiguous across job_id boundaries on purpose: a gap lock taken while
// deleting job_id=1 can only be observed if some row of another job sits
// immediately after (1,3) in primary key order.
const schemaDDL = `
DROP TABLE IF EXISTS job_attachment;
DROP TABLE IF EXISTS job;
DROP TABLE IF EXISTS t_single;

CREATE TABLE job (
  id   BIGINT NOT NULL AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB;

CREATE TABLE job_attachment (
  job_id     BIGINT NOT NULL,
  ord        INT    NOT NULL,
  object_key VARCHAR(255) NOT NULL,
  PRIMARY KEY (job_id, ord)
) ENGINE=InnoDB;

CREATE TABLE t_single (
  id BIGINT NOT NULL PRIMARY KEY,
  v  VARCHAR(32)
) ENGINE=InnoDB;

INSERT INTO job (id, name) VALUES (1,'a'), (2,'b'), (3,'c'), (10,'x'), (11,'y');

INSERT INTO job_attachment (job_id, ord, object_key) VALUES
  (1,1,'k11'), (1,2,'k12'), (1,3,'k13'),
  (2,1,'k21'), (2,2,'k22'),
  (3,1,'k31');

INSERT INTO t_single VALUES (1,'a'),(2,'b'),(3,'c'),(4,'d'),(5,'e');
`

func resetSchema(t *testing.T) {
	t.Helper()
	if _, err := db.ExecContext(t.Context(), schemaDDL); err != nil {
		t.Fatalf("reset schema: %v", err)
	}
}

// requireRepeatableRead guards the whole file: under READ COMMITTED InnoDB
// takes no gap locks at all, so every case below would pass vacuously.
func requireRepeatableRead(t *testing.T) {
	t.Helper()
	var isolation string
	if err := db.QueryRowContext(t.Context(), "SELECT @@transaction_isolation").Scan(&isolation); err != nil {
		t.Fatalf("read isolation level: %v", err)
	}
	if isolation != "REPEATABLE-READ" {
		t.Fatalf("isolation level is %q, want REPEATABLE-READ", isolation)
	}
}

// lock is one row of performance_schema.data_locks, rendered as
// "<LOCK_MODE> <LOCK_DATA>" so expectations read like the mysql client output.
type lock string

// locksHeldOn opens a transaction on its own connection, runs stmt, and
// reports the record locks it left behind on table. Table-level intention
// locks are dropped: they are always present and never the answer.
func locksHeldOn(t *testing.T, table, stmt string) []lock {
	t.Helper()
	ctx := t.Context()

	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatalf("acquire connection: %v", err)
	}
	defer func() { _ = conn.Close() }()

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, stmt); err != nil {
		t.Fatalf("exec %q: %v", stmt, err)
	}
	return readLocks(t, table)
}

func readLocks(t *testing.T, table string) []lock {
	t.Helper()
	rows, err := db.QueryContext(t.Context(), `
		SELECT LOCK_MODE, IFNULL(LOCK_DATA, '')
		FROM performance_schema.data_locks
		WHERE OBJECT_NAME = ? AND LOCK_TYPE = 'RECORD'
		ORDER BY LOCK_DATA, LOCK_MODE`, table)
	if err != nil {
		t.Fatalf("read data_locks: %v", err)
	}
	defer func() { _ = rows.Close() }()

	var got []lock
	for rows.Next() {
		var mode, data string
		if err := rows.Scan(&mode, &data); err != nil {
			t.Fatalf("scan data_locks: %v", err)
		}
		got = append(got, lock(mode+" "+data))
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate data_locks: %v", err)
	}
	return got
}

func TestDeleteLockModes(t *testing.T) {
	resetSchema(t)
	requireRepeatableRead(t)

	tests := map[string]struct {
		table string
		stmt  string
		want  []lock
	}{
		// Baseline. Deleting by job_id alone is a range scan, so InnoDB takes
		// next-key locks (bare "X") on the matched rows plus the gap in front
		// of the next row, which belongs to a different job.
		"range delete by job_id": {
			table: "job_attachment",
			stmt:  "DELETE FROM job_attachment WHERE job_id = 1",
			want:  []lock{"X 1, 1", "X 1, 2", "X 1, 3", "X,GAP 2, 1"},
		},
		// Reference point: equality on a single-column primary key never
		// needs a gap.
		"single column primary key IN": {
			table: "t_single",
			stmt:  "DELETE FROM t_single WHERE id IN (2,3)",
			want:  []lock{"X,REC_NOT_GAP 2", "X,REC_NOT_GAP 3"},
		},
		// The question this package exists to answer.
		"composite primary key row constructor IN": {
			table: "job_attachment",
			stmt:  "DELETE FROM job_attachment WHERE (job_id, ord) IN ((1,1),(1,2))",
			want:  []lock{"X,REC_NOT_GAP 1, 1", "X,REC_NOT_GAP 1, 2"},
		},
		// Same set of rows, written as a leading-column equality plus an IN
		// on the second column. Locks identically.
		"composite primary key AND IN": {
			table: "job_attachment",
			stmt:  "DELETE FROM job_attachment WHERE job_id = 1 AND ord IN (1,2)",
			want:  []lock{"X,REC_NOT_GAP 1, 1", "X,REC_NOT_GAP 1, 2"},
		},
		// Rows either side of a job_id boundary still lock as plain records:
		// what matters is that every tuple is a full primary key, not that
		// the tuples are adjacent.
		"composite primary key across job boundary": {
			table: "job_attachment",
			stmt:  "DELETE FROM job_attachment WHERE (job_id, ord) IN ((1,3),(2,1))",
			want:  []lock{"X,REC_NOT_GAP 1, 3", "X,REC_NOT_GAP 2, 1"},
		},
		// The caveat. A tuple that matches no row cannot be locked as a
		// record, so InnoDB locks the gap where it would have been. Naming a
		// primary key is only gap-free while the row actually exists.
		"composite primary key with missing row": {
			table: "job_attachment",
			stmt:  "DELETE FROM job_attachment WHERE (job_id, ord) IN ((1,1),(1,9))",
			want:  []lock{"X,REC_NOT_GAP 1, 1", "X,GAP 2, 1"},
		},
		// A job with no attachments at all: there is no record to lock, so
		// the scan ends on the supremum pseudo-record and takes a next-key
		// lock there, which is the gap after the last real row. Two
		// transactions doing this for two different empty jobs end up holding
		// the very same lock.
		"range delete matching no rows": {
			table: "job_attachment",
			stmt:  "DELETE FROM job_attachment WHERE job_id = 10",
			want:  []lock{"X supremum pseudo-record"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := locksHeldOn(t, tt.table, tt.stmt)
			if !equalLocks(got, tt.want) {
				t.Errorf("locks after %q:\n got: %s\nwant: %s", tt.stmt, render(got), render(tt.want))
			}
		})
	}
}

func TestExplainDeletePlans(t *testing.T) {
	resetSchema(t)

	stmts := []string{
		"EXPLAIN DELETE FROM job_attachment WHERE (job_id, ord) IN ((1,1),(1,2))",
		"EXPLAIN DELETE FROM job_attachment WHERE job_id = 1 AND ord IN (1,2)",
		"EXPLAIN DELETE FROM job_attachment WHERE job_id = 1",
	}
	for _, stmt := range stmts {
		accessType, key := explain(t, stmt)
		t.Logf("%s\n  type=%s key=%s", stmt, accessType, key)
		if key != "PRIMARY" {
			t.Errorf("%s: key = %q, want PRIMARY", stmt, key)
		}
	}
}

// explain returns the access type and chosen index of a single-row EXPLAIN.
func explain(t *testing.T, stmt string) (accessType, key string) {
	t.Helper()
	rows, err := db.QueryContext(t.Context(), stmt)
	if err != nil {
		t.Fatalf("explain %q: %v", stmt, err)
	}
	defer func() { _ = rows.Close() }()

	cols, err := rows.Columns()
	if err != nil {
		t.Fatalf("explain columns: %v", err)
	}
	if !rows.Next() {
		t.Fatalf("explain %q returned no rows", stmt)
	}
	cells := make([]sql.NullString, len(cols))
	dest := make([]any, len(cols))
	for i := range cells {
		dest[i] = &cells[i]
	}
	if err := rows.Scan(dest...); err != nil {
		t.Fatalf("scan explain: %v", err)
	}
	for i, c := range cols {
		switch c {
		case "type":
			accessType = cells[i].String
		case "key":
			key = cells[i].String
		}
	}
	return accessType, key
}

func equalLocks(got, want []lock) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

func render(locks []lock) string {
	parts := make([]string, len(locks))
	for i, l := range locks {
		parts[i] = fmt.Sprintf("[%s]", l)
	}
	return strings.Join(parts, " ")
}

// waitForBlockedLock blocks until some transaction is waiting on a lock on
// table, so a test can order "session A is stuck" before session B proceeds.
func waitForBlockedLock(ctx context.Context, table string) error {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		var n int
		err := db.QueryRowContext(ctx, `
			SELECT COUNT(*)
			FROM performance_schema.data_locks
			WHERE OBJECT_NAME = ? AND LOCK_STATUS = 'WAITING'`, table).Scan(&n)
		if err != nil {
			return err
		}
		if n > 0 {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
