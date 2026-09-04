package innodblock_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

// TestForUpdateOnChildSelectTakesGapLocks is B-1: the implementations that
// must not be used. Each case is expected to take a gap lock, which is what
// makes it a rule rather than a preference.
func TestForUpdateOnChildSelectTakesGapLocks(t *testing.T) {
	resetSchema(t)
	requireRepeatableRead(t)

	tests := map[string]struct {
		table string
		stmt  string
		want  []lock
	}{
		// B-1-a. Locking the child SELECT turns the read into the same
		// range scan the first report measured for a range DELETE.
		"child select for update": {
			table: "job_attachment",
			stmt:  "SELECT ord FROM job_attachment WHERE job_id = 1 FOR UPDATE",
			want:  []lock{"X 1, 1", "X 1, 2", "X 1, 3", "X,GAP 2, 1"},
		},
		// B-1-b. Picking the next ord in SQL. The optimizer reads only the
		// last matching row, so fewer locks -- but still a next-key lock and
		// still the gap in front of the next job's row.
		"select max ord for update": {
			table: "job_attachment",
			stmt:  "SELECT MAX(ord) FROM job_attachment WHERE job_id = 1 FOR UPDATE",
			want:  []lock{"X 1, 3", "X,GAP 2, 1"},
		},
		// B-1-c. Locking a parent row that does not exist. There is no
		// record to lock, so the scan ends on the supremum pseudo-record and
		// takes a next-key lock there: the whole range past the last job.
		"parent select for update on missing row": {
			table: "job",
			stmt:  "SELECT id FROM job WHERE id = 9999 FOR UPDATE",
			want:  []lock{"X supremum pseudo-record"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := locksHeldOn(t, tt.table, tt.stmt)
			if !equalLocks(got, tt.want) {
				t.Errorf("locks:\n got: %s\nwant: %s", render(got), render(tt.want))
			}
		})
	}
}

// TestForUpdateOnMissingParentBlocksInserts is the second half of B-1-c: how
// far the gap lock on a missing parent actually reaches.
func TestForUpdateOnMissingParentBlocksInserts(t *testing.T) {
	resetSchema(t)
	requireRepeatableRead(t)
	ctx := t.Context()

	a, b := newSession(t), newSession(t)

	if err := a.exec(ctx, "SELECT id FROM job WHERE id = 9999 FOR UPDATE"); err != nil {
		t.Fatalf("session A select for update: %v", err)
	}
	dumpLocks(t, "session A holding the missing-row gap", "job")

	// A new job takes the next auto-increment id, which falls inside the
	// locked gap.
	err := b.exec(ctx, "INSERT INTO job (name) VALUES ('z')")
	if got := mysqlErrNo(err); got != errLockWaitTime {
		t.Errorf("session B creating a job: error %d (%v), want %d (lock wait timeout)", got, err, errLockWaitTime)
	}
}

// seedAttachments replaces job 1's attachments with n rows, ord 1..n.
func seedAttachments(t *testing.T, n int) {
	t.Helper()
	ctx := t.Context()
	if _, err := db.ExecContext(ctx, "DELETE FROM job_attachment WHERE job_id = 1"); err != nil {
		t.Fatalf("clear attachments: %v", err)
	}
	values := make([]string, n)
	for i := range values {
		values[i] = fmt.Sprintf("(1,%d,'k%d')", i+1, i+1)
	}
	stmt := "INSERT INTO job_attachment (job_id, ord, object_key) VALUES " + strings.Join(values, ",")
	if _, err := db.ExecContext(ctx, stmt); err != nil {
		t.Fatalf("seed attachments: %v", err)
	}
}

func ordList(n int) string {
	parts := make([]string, n)
	for i := range parts {
		parts[i] = strconv.Itoa(i + 1)
	}
	return strings.Join(parts, ",")
}

// summarize reduces a lock set to the distinct modes it contains, so a large
// result can be reported in one line.
func summarize(locks []lock) (modes []string, count int) {
	seen := map[string]bool{}
	for _, l := range locks {
		mode := string(l)
		if i := strings.Index(mode, " "); i > 0 {
			mode = mode[:i]
		}
		if !seen[mode] {
			seen[mode] = true
			modes = append(modes, mode)
		}
	}
	return modes, len(locks)
}

// seedOtherJobs adds jobs jobs, each with perJob attachments, so that job 1
// is a small share of the table rather than nearly all of it.
func seedOtherJobs(t *testing.T, jobs, perJob int) {
	t.Helper()
	ctx := t.Context()

	var jobRows, attRows []string
	for j := 0; j < jobs; j++ {
		id := 1000 + j
		jobRows = append(jobRows, fmt.Sprintf("(%d,'j%d')", id, id))
		for o := 1; o <= perJob; o++ {
			attRows = append(attRows, fmt.Sprintf("(%d,%d,'k')", id, o))
		}
	}
	insertInBatches(t, ctx, "INSERT INTO job (id, name) VALUES ", jobRows)
	insertInBatches(t, ctx, "INSERT INTO job_attachment (job_id, ord, object_key) VALUES ", attRows)
	if _, err := db.ExecContext(ctx, "ANALYZE TABLE job_attachment"); err != nil {
		t.Fatalf("analyze table: %v", err)
	}
}

func insertInBatches(t *testing.T, ctx context.Context, prefix string, values []string) {
	t.Helper()
	const batch = 500
	for start := 0; start < len(values); start += batch {
		end := min(start+batch, len(values))
		if _, err := db.ExecContext(ctx, prefix+strings.Join(values[start:end], ",")); err != nil {
			t.Fatalf("seed rows: %v", err)
		}
	}
}

// TestLargeInListLockModes is B-2: does a long IN list change the lock modes?
//
// Two table shapes are measured, because the answer depends on the shape
// rather than on the list length. eq_range_index_dive_limit (default 200)
// only changes how a range is costed; what decides the locks is whether the
// optimizer still picks the primary key at all.
func TestLargeInListLockModes(t *testing.T) {
	requireRepeatableRead(t)

	var diveLimit int
	if err := db.QueryRowContext(t.Context(), "SELECT @@eq_range_index_dive_limit").Scan(&diveLimit); err != nil {
		t.Fatalf("read eq_range_index_dive_limit: %v", err)
	}
	t.Logf("eq_range_index_dive_limit = %d", diveLimit)

	sizes := []int{10, 100, 199, 200, 201, 230, 250}

	// Job 1 holds 300 of the table's 304 rows. A full scan is cheaper than
	// walking the index, so the optimizer abandons the primary key and the
	// delete takes next-key locks over the entire table, other jobs
	// included. Nothing about the IN list changed -- only the selectivity.
	t.Run("job is most of the table", func(t *testing.T) {
		resetSchema(t)
		seedAttachments(t, 300)
		for _, size := range sizes {
			t.Run(strconv.Itoa(size), func(t *testing.T) {
				accessType, _, count := reportInListLocks(t, size)
				if size < 100 {
					return
				}
				if accessType != "ALL" {
					t.Errorf("access type = %q, want ALL", accessType)
				}
				if count <= size {
					t.Errorf("locks = %d, want more than the %d deleted rows", count, size)
				}
			})
		}
	})

	// Job 1 holds 300 rows out of roughly 5300: the realistic shape, where
	// one job's attachments are a small share of the table.
	t.Run("job is a small share of the table", func(t *testing.T) {
		resetSchema(t)
		seedAttachments(t, 300)
		seedOtherJobs(t, 1000, 5)
		for _, size := range sizes {
			t.Run(strconv.Itoa(size), func(t *testing.T) {
				accessType, modes, count := reportInListLocks(t, size)
				if accessType != "range" {
					t.Errorf("access type = %q, want range", accessType)
				}
				// One extra next-key lock on a page supremum shows up at
				// some sizes; see the test below for how far it reaches.
				for _, m := range modes {
					if m != "X,REC_NOT_GAP" && m != "X" {
						t.Errorf("unexpected lock mode %q", m)
					}
				}
				if count != size && count != size+1 {
					t.Errorf("locks = %d, want %d (or %d with the supremum lock)", count, size, size+1)
				}
			})
		}
	})
}

// reportInListLocks measures one IN list size and logs the plan alongside the
// lock summary.
func reportInListLocks(t *testing.T, size int) (accessType string, modes []string, count int) {
	t.Helper()
	stmt := fmt.Sprintf("DELETE FROM job_attachment WHERE job_id = 1 AND ord IN (%s)", ordList(size))
	accessType, key, keyLen := explain(t, "EXPLAIN "+stmt)

	got := locksQuietly(t, "job_attachment", stmt)
	modes, count = summarize(got)
	t.Logf("size=%d type=%s key=%q key_len=%q locks=%d modes=%v",
		size, accessType, key, keyLen, count, modes)
	return accessType, modes, count
}

// TestSupremumLockFromLongInListStaysWithinTheJob pins down the one lock that
// is not a plain record lock in the realistic shape above. It is a next-key
// lock on a page supremum inside the job's own rows, not on the end of the
// table, so it blocks no other job.
func TestSupremumLockFromLongInListStaysWithinTheJob(t *testing.T) {
	resetSchema(t)
	requireRepeatableRead(t)
	seedAttachments(t, 300)
	seedOtherJobs(t, 1000, 5)
	ctx := t.Context()

	a := newSession(t)
	stmt := fmt.Sprintf("DELETE FROM job_attachment WHERE job_id = 1 AND ord IN (%s)", ordList(230))
	if err := a.exec(ctx, stmt); err != nil {
		t.Fatalf("session A delete: %v", err)
	}
	modes, count := summarize(readLocks(t, "job_attachment", connID(t, a.tx)))
	t.Logf("session A holds %d locks, modes=%v", count, modes)

	// Every one of these is outside the deleted range; none may block.
	for _, ins := range []string{
		"INSERT INTO job_attachment (job_id, ord, object_key) VALUES (1999,99,'b')",
		"INSERT INTO job_attachment (job_id, ord, object_key) VALUES (1500,99,'b')",
		"INSERT INTO job_attachment (job_id, ord, object_key) VALUES (2,99,'b')",
		"INSERT INTO job_attachment (job_id, ord, object_key) VALUES (1,400,'b')",
	} {
		b := newSession(t)
		if err := b.exec(ctx, ins); err != nil {
			t.Errorf("%s: %v", ins, err)
		}
		if err := b.tx.Rollback(); err != nil {
			t.Fatalf("rollback: %v", err)
		}
	}
}

// TestInListBelowDiveLimitKeepsRecordLocks forces the statistics-based path
// at a small IN list size, isolating the dive limit from the list length.
func TestInListBelowDiveLimitKeepsRecordLocks(t *testing.T) {
	resetSchema(t)
	requireRepeatableRead(t)
	seedAttachments(t, 300)
	seedOtherJobs(t, 1000, 5)

	stmt := fmt.Sprintf("DELETE FROM job_attachment WHERE job_id = 1 AND ord IN (%s)", ordList(10))
	got := locksQuietly(t, "job_attachment",
		"SET SESSION eq_range_index_dive_limit = 2",
		stmt,
	)
	modes, count := summarize(got)
	t.Logf("eq_range_index_dive_limit=2 locks=%d modes=%v", count, modes)

	if count != 10 {
		t.Errorf("locks = %d, want 10", count)
	}
	if len(modes) != 1 || modes[0] != "X,REC_NOT_GAP" {
		t.Errorf("lock modes = %v, want only [X,REC_NOT_GAP]", modes)
	}
}
