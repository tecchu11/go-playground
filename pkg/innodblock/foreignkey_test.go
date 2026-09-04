package innodblock_test

import (
	"context"
	"errors"
	"testing"
	"time"
)

// addForeignKey turns job_attachment into a real child table. With the
// constraint in place, writing a child row makes InnoDB check the parent and
// lock it, which is what these tests measure.
func addForeignKey(t *testing.T) {
	t.Helper()
	_, err := db.ExecContext(t.Context(), `
		ALTER TABLE job_attachment
		  ADD CONSTRAINT fk_job_attachment_job
		  FOREIGN KEY (job_id) REFERENCES job(id)`)
	if err != nil {
		t.Fatalf("add foreign key: %v", err)
	}
}

// TestForeignKeyChildInsertLocksParent establishes the mechanism the rest of
// C-1 depends on: an insert into the child takes a shared lock on the parent
// row it references.
func TestForeignKeyChildInsertLocksParent(t *testing.T) {
	resetSchema(t)
	addForeignKey(t)
	requireRepeatableRead(t)

	got := locksHeldOn(t, "job",
		"INSERT INTO job_attachment (job_id, ord, object_key) VALUES (1,9,'k')")
	want := []lock{"S,REC_NOT_GAP 1"}
	if !equalLocks(got, want) {
		t.Errorf("locks on job:\n got: %s\nwant: %s", render(got), render(want))
	}
}

// TestForeignKeySequenceOnAdjacentJobs is C-1-a: the production sequence,
// which locks the parent first, run under the constraint.
func TestForeignKeySequenceOnAdjacentJobs(t *testing.T) {
	resetSchema(t)
	addForeignKey(t)
	requireRepeatableRead(t)
	ctx := t.Context()

	a, b := newSession(t), newSession(t)

	lockParent(t, a, 1)
	ordsA := readOrds(t, a, 1)
	if err := updateJobName(a, ctx, 1, "a2"); err != nil {
		t.Fatalf("session A update: %v", err)
	}
	deleteOrds(t, a, 1, ordsA)

	lockParent(t, b, 2)
	ordsB := readOrds(t, b, 2)
	if err := updateJobName(b, ctx, 2, "b2"); err != nil {
		t.Fatalf("session B update: %v", err)
	}
	deleteOrds(t, b, 2, ordsB)

	if err := insertAttachments(a, ctx, 1, "a-k1", "a-k2"); err != nil {
		t.Errorf("session A insert: %v", err)
	}
	if err := insertAttachments(b, ctx, 2, "b-k1", "b-k2"); err != nil {
		t.Errorf("session B insert: %v", err)
	}

	// Each session should hold its own exclusive lock on its own parent row
	// and nothing else: the foreign key check finds the lock already held
	// and adds no shared lock on top.
	dumpLocks(t, "both sessions before commit", "job")
	if got := readLocks(t, "job", connID(t, a.tx)); !equalLocks(got, []lock{"X,REC_NOT_GAP 1"}) {
		t.Errorf("session A locks on job: %s, want [X,REC_NOT_GAP 1]", render(got))
	}

	if err := a.commit(); err != nil {
		t.Errorf("session A commit: %v", err)
	}
	if err := b.commit(); err != nil {
		t.Errorf("session B commit: %v", err)
	}
}

// TestForeignKeyWithoutParentLock is C-1-b: the same work with the parent
// SELECT ... FOR UPDATE removed.
func TestForeignKeyWithoutParentLock(t *testing.T) {
	requireRepeatableRead(t)

	// Different jobs never meet, constraint or not.
	t.Run("different jobs", func(t *testing.T) {
		resetSchema(t)
		addForeignKey(t)
		ctx := t.Context()
		a, b := newSession(t), newSession(t)

		ordsA := readOrds(t, a, 1)
		if err := updateJobName(a, ctx, 1, "a2"); err != nil {
			t.Fatalf("session A update: %v", err)
		}
		deleteOrds(t, a, 1, ordsA)

		ordsB := readOrds(t, b, 2)
		if err := updateJobName(b, ctx, 2, "b2"); err != nil {
			t.Fatalf("session B update: %v", err)
		}
		deleteOrds(t, b, 2, ordsB)

		if err := insertAttachments(a, ctx, 1, "a-k1"); err != nil {
			t.Errorf("session A insert: %v", err)
		}
		if err := insertAttachments(b, ctx, 2, "b-k1"); err != nil {
			t.Errorf("session B insert: %v", err)
		}
	})

	// Same job, statements in the order the sequence specifies: the parent
	// UPDATE comes first, so it serialises the two transactions exactly as
	// the explicit parent lock would have.
	t.Run("same job, parent updated first", func(t *testing.T) {
		resetSchema(t)
		addForeignKey(t)
		ctx := t.Context()
		a, b := newSession(t), newSession(t)

		if err := updateJobName(a, ctx, 1, "a2"); err != nil {
			t.Fatalf("session A update: %v", err)
		}
		err := updateJobName(b, ctx, 1, "b2")
		if got := mysqlErrNo(err); got != errLockWaitTime {
			t.Errorf("session B update: error %d (%v), want %d (lock wait timeout)", got, err, errLockWaitTime)
		}
	})

	// Same job, child written before the parent. Each transaction takes a
	// shared lock on the parent row through the foreign key check, then asks
	// to upgrade it to exclusive. Both are waiting on the other's shared
	// lock, which is a deadlock.
	t.Run("same job, child written first", func(t *testing.T) {
		resetSchema(t)
		addForeignKey(t)
		ctx := t.Context()
		a, b := newSession(t), newSession(t)

		if err := a.exec(ctx, "INSERT INTO job_attachment (job_id, ord, object_key) VALUES (1,101,'a')"); err != nil {
			t.Fatalf("session A child insert: %v", err)
		}
		if err := b.exec(ctx, "INSERT INTO job_attachment (job_id, ord, object_key) VALUES (1,102,'b')"); err != nil {
			t.Fatalf("session B child insert: %v", err)
		}
		dumpLocks(t, "both sessions after the child inserts", "job")

		doneA := make(chan error, 1)
		go func() {
			doneA <- updateJobName(a, ctx, 1, "a2")
		}()

		waitCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		if err := waitForBlockedLock(waitCtx, "job"); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("wait for blocked lock: %v", err)
		}

		errB := updateJobName(b, ctx, 1, "b2")
		var errA error
		select {
		case errA = <-doneA:
		case <-time.After(30 * time.Second):
			t.Fatal("session A update never returned")
		}
		t.Logf("session A update: %v", errA)
		t.Logf("session B update: %v", errB)

		if mysqlErrNo(errA) != errDeadlock && mysqlErrNo(errB) != errDeadlock {
			t.Errorf("expected one update to fail with error %d (deadlock), got A=%v B=%v", errDeadlock, errA, errB)
		} else {
			t.Logf("SHOW ENGINE INNODB STATUS:\n%s", latestDeadlock(t))
		}
	})
}
