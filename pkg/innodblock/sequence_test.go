package innodblock_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

// The sequence under test, as it would be issued by the application:
//
//	BEGIN;
//	SELECT id, name FROM job WHERE id = ? FOR UPDATE;
//	SELECT ord, object_key FROM job_attachment WHERE job_id = ? ORDER BY ord;  -- no FOR UPDATE
//	UPDATE job SET name = ? WHERE id = ?;
//	DELETE FROM job_attachment WHERE job_id = ? AND ord IN (...);              -- only the ords just read
//	INSERT INTO job_attachment (job_id, ord, object_key) VALUES (...);         -- bulk
//	COMMIT;
//
// It is split into steps so two sessions can be interleaved statement by
// statement.

// lockParent takes the parent row lock that orders the whole sequence.
func lockParent(t *testing.T, s *session, jobID int) {
	t.Helper()
	if err := s.exec(t.Context(), fmt.Sprintf("SELECT id, name FROM job WHERE id = %d FOR UPDATE", jobID)); err != nil {
		t.Fatalf("select job %d for update: %v", jobID, err)
	}
}

// readOrds runs the plain child SELECT and returns the ord values found.
// Deliberately no FOR UPDATE: see TestForUpdateOnChildSelectTakesGapLocks.
func readOrds(t *testing.T, s *session, jobID int) []int {
	t.Helper()
	rows, err := s.tx.QueryContext(t.Context(),
		"SELECT ord, object_key FROM job_attachment WHERE job_id = ? ORDER BY ord", jobID)
	if err != nil {
		t.Fatalf("select attachments of job %d: %v", jobID, err)
	}
	defer func() { _ = rows.Close() }()

	var ords []int
	for rows.Next() {
		var ord int
		var key string
		if err := rows.Scan(&ord, &key); err != nil {
			t.Fatalf("scan attachment: %v", err)
		}
		ords = append(ords, ord)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate attachments: %v", err)
	}
	return ords
}

// deleteOrds issues the DELETE for exactly the ords that were read, and
// issues nothing at all when there were none. Skipping the empty DELETE is
// the fix established in the first report.
func deleteOrds(t *testing.T, s *session, jobID int, ords []int) {
	t.Helper()
	if len(ords) == 0 {
		return
	}
	list := make([]string, len(ords))
	for i, o := range ords {
		list[i] = strconv.Itoa(o)
	}
	stmt := fmt.Sprintf("DELETE FROM job_attachment WHERE job_id = %d AND ord IN (%s)", jobID, strings.Join(list, ","))
	if err := s.exec(t.Context(), stmt); err != nil {
		t.Fatalf("delete attachments of job %d: %v", jobID, err)
	}
}

// insertAttachments writes the new attachment set in one bulk statement.
func insertAttachments(s *session, ctx context.Context, jobID int, keys ...string) error {
	values := make([]string, len(keys))
	for i, k := range keys {
		values[i] = fmt.Sprintf("(%d,%d,'%s')", jobID, i+1, k)
	}
	return s.exec(ctx, "INSERT INTO job_attachment (job_id, ord, object_key) VALUES "+strings.Join(values, ","))
}

func updateJobName(s *session, ctx context.Context, jobID int, name string) error {
	return s.exec(ctx, fmt.Sprintf("UPDATE job SET name = '%s' WHERE id = %d", name, jobID))
}

// TestSequenceOnAdjacentJobs is A-2-a: two sessions run the full sequence on
// two different jobs whose attachment rows are adjacent in primary key order.
func TestSequenceOnAdjacentJobs(t *testing.T) {
	resetSchema(t)
	requireRepeatableRead(t)
	ctx := t.Context()

	a, b := newSession(t), newSession(t)

	// Step 1: session A up to and including its DELETE.
	lockParent(t, a, 1)
	ordsA := readOrds(t, a, 1)
	if err := updateJobName(a, ctx, 1, "a2"); err != nil {
		t.Fatalf("session A update: %v", err)
	}
	deleteOrds(t, a, 1, ordsA)

	// Step 2: session B up to and including its DELETE.
	lockParent(t, b, 2)
	ordsB := readOrds(t, b, 2)
	if err := updateJobName(b, ctx, 2, "b2"); err != nil {
		t.Fatalf("session B update: %v", err)
	}
	deleteOrds(t, b, 2, ordsB)

	dumpLocks(t, "both sessions after DELETE", "job_attachment")

	// Step 3 and 4: the bulk inserts, A first.
	if err := insertAttachments(a, ctx, 1, "a-k1", "a-k2"); err != nil {
		t.Errorf("session A insert: %v", err)
	}
	if err := insertAttachments(b, ctx, 2, "b-k1", "b-k2", "b-k3"); err != nil {
		t.Errorf("session B insert: %v", err)
	}

	dumpLocks(t, "both sessions after INSERT", "job_attachment")

	// Step 5.
	if err := a.commit(); err != nil {
		t.Errorf("session A commit: %v", err)
	}
	if err := b.commit(); err != nil {
		t.Errorf("session B commit: %v", err)
	}
}

// TestSequenceOnEmptyJobs is A-2-b: the exact shape that deadlocked before,
// now with the empty DELETE skipped. Jobs 10 and 11 have no attachments.
func TestSequenceOnEmptyJobs(t *testing.T) {
	resetSchema(t)
	requireRepeatableRead(t)
	ctx := t.Context()

	a, b := newSession(t), newSession(t)

	lockParent(t, a, 10)
	ordsA := readOrds(t, a, 10)
	if len(ordsA) != 0 {
		t.Fatalf("job 10 should have no attachments, got %v", ordsA)
	}

	lockParent(t, b, 11)
	ordsB := readOrds(t, b, 11)
	if len(ordsB) != 0 {
		t.Fatalf("job 11 should have no attachments, got %v", ordsB)
	}

	dumpLocks(t, "both sessions after the empty child SELECT", "job_attachment")

	if err := updateJobName(a, ctx, 10, "x2"); err != nil {
		t.Fatalf("session A update: %v", err)
	}
	deleteOrds(t, a, 10, ordsA) // issues nothing
	if err := insertAttachments(a, ctx, 10, "k"); err != nil {
		t.Errorf("session A insert: %v", err)
	}

	if err := updateJobName(b, ctx, 11, "y2"); err != nil {
		t.Fatalf("session B update: %v", err)
	}
	deleteOrds(t, b, 11, ordsB) // issues nothing
	if err := insertAttachments(b, ctx, 11, "k"); err != nil {
		t.Errorf("session B insert: %v", err)
	}

	dumpLocks(t, "both sessions after INSERT", "job_attachment")

	if err := a.commit(); err != nil {
		t.Errorf("session A commit: %v", err)
	}
	if err := b.commit(); err != nil {
		t.Errorf("session B commit: %v", err)
	}
}

// TestSequenceUpdateVersusCreate is A-2-c: session A updates the highest job
// so its insert lands at the table tail, while session B creates a brand new
// job whose attachment lands past it.
func TestSequenceUpdateVersusCreate(t *testing.T) {
	resetSchema(t)
	requireRepeatableRead(t)
	ctx := t.Context()

	a, b := newSession(t), newSession(t)

	lockParent(t, a, 3)
	ordsA := readOrds(t, a, 3)
	if err := updateJobName(a, ctx, 3, "c2"); err != nil {
		t.Fatalf("session A update: %v", err)
	}
	deleteOrds(t, a, 3, ordsA)
	if err := insertAttachments(a, ctx, 3, "a-k1", "a-k2"); err != nil {
		t.Fatalf("session A insert: %v", err)
	}
	dumpLocks(t, "session A holding the tail", "job_attachment")

	// Session B creates a job and its first attachment while A is open.
	if err := b.exec(ctx, "INSERT INTO job (name) VALUES ('new')"); err != nil {
		t.Errorf("session B create job: %v", err)
	}
	var newID int
	if err := b.tx.QueryRowContext(ctx, "SELECT LAST_INSERT_ID()").Scan(&newID); err != nil {
		t.Fatalf("read LAST_INSERT_ID: %v", err)
	}
	t.Logf("session B created job id %d", newID)
	if err := insertAttachments(b, ctx, newID, "b-k1"); err != nil {
		t.Errorf("session B insert attachment: %v", err)
	}

	if err := b.commit(); err != nil {
		t.Errorf("session B commit: %v", err)
	}
	if err := a.commit(); err != nil {
		t.Errorf("session A commit: %v", err)
	}
}

// TestSequenceOnSameJobSerializes is A-2-d: two sessions on the same job must
// queue behind the parent row lock rather than deadlock.
func TestSequenceOnSameJobSerializes(t *testing.T) {
	resetSchema(t)
	requireRepeatableRead(t)
	ctx := t.Context()

	a, b := newSession(t), newSession(t)

	lockParent(t, a, 1)
	ordsA := readOrds(t, a, 1)
	if err := updateJobName(a, ctx, 1, "a2"); err != nil {
		t.Fatalf("session A update: %v", err)
	}
	deleteOrds(t, a, 1, ordsA)
	if err := insertAttachments(a, ctx, 1, "a-k1", "a-k2"); err != nil {
		t.Fatalf("session A insert: %v", err)
	}

	// Session B starts the same sequence and must stop at the parent lock.
	blocked := make(chan error, 1)
	go func() {
		blocked <- b.exec(ctx, "SELECT id, name FROM job WHERE id = 1 FOR UPDATE")
	}()

	waitCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := waitForBlockedLock(waitCtx, "job"); err != nil {
		t.Fatalf("session B never blocked on the parent row: %v", err)
	}
	if n := countLockWaits(t); n != 1 {
		t.Errorf("outstanding lock waits = %d, want exactly 1", n)
	}
	dumpLocks(t, "session B waiting on the parent row", "job")

	if err := a.commit(); err != nil {
		t.Fatalf("session A commit: %v", err)
	}

	select {
	case err := <-blocked:
		if err != nil {
			t.Errorf("session B parent lock after A committed: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("session B never acquired the parent row lock")
	}

	// B now sees A's committed rows and completes its own sequence.
	ordsB := readOrds(t, b, 1)
	t.Logf("session B read ords %v after A committed", ordsB)
	if err := updateJobName(b, ctx, 1, "b2"); err != nil {
		t.Errorf("session B update: %v", err)
	}
	deleteOrds(t, b, 1, ordsB)
	if err := insertAttachments(b, ctx, 1, "b-k1"); err != nil {
		t.Errorf("session B insert: %v", err)
	}
	if err := b.commit(); err != nil {
		t.Errorf("session B commit: %v", err)
	}
}
