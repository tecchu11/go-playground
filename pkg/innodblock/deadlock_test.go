package innodblock_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
)

const (
	errDeadlock     = 1213
	errLockWaitTime = 1205
)

// session is one pinned connection running one transaction, so that two
// concurrent transactions can be interleaved statement by statement.
type session struct {
	conn *sql.Conn
	tx   *sql.Tx
}

func newSession(t *testing.T) *session {
	t.Helper()
	ctx := t.Context()

	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatalf("acquire connection: %v", err)
	}
	// Keep a stuck statement from holding the test for the 50s default.
	if _, err := conn.ExecContext(ctx, "SET SESSION innodb_lock_wait_timeout = 5"); err != nil {
		t.Fatalf("set lock wait timeout: %v", err)
	}
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	s := &session{conn: conn, tx: tx}
	t.Cleanup(func() {
		_ = s.tx.Rollback()
		_ = s.conn.Close()
	})
	return s
}

func (s *session) exec(ctx context.Context, stmt string) error {
	_, err := s.tx.ExecContext(ctx, stmt)
	return err
}

func (s *session) commit() error {
	return s.tx.Commit()
}

func mysqlErrNo(err error) uint16 {
	var myErr *mysql.MySQLError
	if errors.As(err, &myErr) {
		return myErr.Number
	}
	return 0
}

// runInterleaved plays out the delete-then-insert pattern in two transactions
// that touch different jobs, and reports the error each insert returned.
//
//	A: BEGIN; deleteA
//	B: BEGIN; deleteB
//	A: INSERT (10,1)   -- may block
//	B: INSERT (11,1)   -- closes the cycle
//
// An empty deleteX skips that statement, which models a caller that read the
// current rows first and found nothing to delete.
func runInterleaved(t *testing.T, deleteA, deleteB string) (errInsertA, errInsertB error) {
	t.Helper()
	ctx := t.Context()

	a, b := newSession(t), newSession(t)

	if deleteA != "" {
		if err := a.exec(ctx, deleteA); err != nil {
			t.Fatalf("session A %q: %v", deleteA, err)
		}
	}
	if deleteB != "" {
		if err := b.exec(ctx, deleteB); err != nil {
			t.Fatalf("session B %q: %v", deleteB, err)
		}
	}

	doneA := make(chan error, 1)
	go func() {
		doneA <- a.exec(ctx, "INSERT INTO job_attachment VALUES (10,1,'k')")
	}()

	// Let session A's insert either finish or come to rest as a waiter before
	// session B inserts; otherwise the interleaving is a coin flip.
	waitCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	waitErr := waitForBlockedLock(waitCtx, "job_attachment")
	if waitErr != nil && !errors.Is(waitErr, context.DeadlineExceeded) {
		t.Fatalf("wait for blocked lock: %v", waitErr)
	}

	errInsertB = b.exec(ctx, "INSERT INTO job_attachment VALUES (11,1,'k')")

	select {
	case errInsertA = <-doneA:
	case <-time.After(30 * time.Second):
		t.Fatal("session A insert never returned")
	}
	return errInsertA, errInsertB
}

// TestDeadlockOnRangeDelete reproduces the failure the whole investigation is
// about: two transactions delete attachments of two different jobs that have
// no attachments, each takes the same gap lock, and each then tries to insert
// into that gap.
func TestDeadlockOnRangeDelete(t *testing.T) {
	resetSchema(t)
	requireRepeatableRead(t)

	errA, errB := runInterleaved(t,
		"DELETE FROM job_attachment WHERE job_id = 10",
		"DELETE FROM job_attachment WHERE job_id = 11",
	)
	t.Logf("session A insert: %v", errA)
	t.Logf("session B insert: %v", errB)

	if mysqlErrNo(errA) != errDeadlock && mysqlErrNo(errB) != errDeadlock {
		t.Errorf("expected one insert to fail with error %d (deadlock), got A=%v B=%v", errDeadlock, errA, errB)
	}
}

// TestNoDeadlockWhenNothingIsDeleted is the fix: the caller reads the current
// attachments first and issues no DELETE when there is nothing to delete, so
// neither transaction ever holds a gap lock.
func TestNoDeadlockWhenNothingIsDeleted(t *testing.T) {
	resetSchema(t)
	requireRepeatableRead(t)

	errA, errB := runInterleaved(t, "", "")
	if errA != nil || errB != nil {
		t.Errorf("inserts should both succeed, got A=%v B=%v", errA, errB)
	}
}

// TestDeadlockOnPrimaryKeyDeleteOfMissingRow is the trap. Switching to a
// primary key equality DELETE is not enough on its own: naming a row that
// does not exist still takes the gap lock, and the deadlock comes back.
func TestDeadlockOnPrimaryKeyDeleteOfMissingRow(t *testing.T) {
	resetSchema(t)
	requireRepeatableRead(t)

	errA, errB := runInterleaved(t,
		"DELETE FROM job_attachment WHERE (job_id, ord) IN ((10,1))",
		"DELETE FROM job_attachment WHERE (job_id, ord) IN ((11,1))",
	)
	t.Logf("session A insert: %v", errA)
	t.Logf("session B insert: %v", errB)

	if mysqlErrNo(errA) != errDeadlock && mysqlErrNo(errB) != errDeadlock {
		t.Errorf("expected one insert to fail with error %d (deadlock), got A=%v B=%v", errDeadlock, errA, errB)
	}
}

// TestNoDeadlockOnPrimaryKeyDeleteOfExistingRows deletes rows that do exist,
// in two transactions working on different jobs. Only record locks are taken,
// so the transactions never meet.
func TestNoDeadlockOnPrimaryKeyDeleteOfExistingRows(t *testing.T) {
	resetSchema(t)
	requireRepeatableRead(t)

	ctx := t.Context()
	a, b := newSession(t), newSession(t)

	if err := a.exec(ctx, "DELETE FROM job_attachment WHERE (job_id, ord) IN ((1,1),(1,2))"); err != nil {
		t.Fatalf("session A delete: %v", err)
	}
	if err := b.exec(ctx, "DELETE FROM job_attachment WHERE (job_id, ord) IN ((2,1),(2,2))"); err != nil {
		t.Fatalf("session B delete: %v", err)
	}
	if err := a.exec(ctx, "INSERT INTO job_attachment VALUES (1,1,'k11-new')"); err != nil {
		t.Errorf("session A insert: %v", err)
	}
	if err := b.exec(ctx, "INSERT INTO job_attachment VALUES (2,1,'k21-new')"); err != nil {
		t.Errorf("session B insert: %v", err)
	}
}
