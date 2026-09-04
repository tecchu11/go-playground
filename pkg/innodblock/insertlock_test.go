package innodblock_test

import "testing"

// TestInsertLockModes covers the INSERT half of the delete-then-insert
// pattern: after the DELETE has been shown to take only record locks, does
// the INSERT that follows it in the same transaction add a gap lock?
//
// Every case runs DELETE and INSERT in one transaction and then reads the
// locks that transaction still holds.
//
// Note on reading the results: a row a transaction inserted itself carries no
// explicit lock. InnoDB marks the record with the inserting transaction id
// and lets a later reader promote that to a real lock only if it actually
// conflicts, so a freshly inserted key contributes no row to data_locks at
// all. Absence here therefore means "no explicit lock", not "unprotected" --
// TestInsertedRowIsStillProtected below pins down what it does block.
func TestInsertLockModes(t *testing.T) {
	resetSchema(t)
	requireRepeatableRead(t)

	tests := map[string]struct {
		stmts []string
		want  []lock
	}{
		// A-1-a. Same number of rows back in, same ord values. Only the
		// DELETE's record locks remain; re-inserting the same keys adds
		// nothing.
		"reinsert same ord values": {
			stmts: []string{
				"DELETE FROM job_attachment WHERE job_id = 1 AND ord IN (1,2,3)",
				"INSERT INTO job_attachment (job_id, ord, object_key) VALUES (1,1,'new1'), (1,2,'new2'), (1,3,'new3')",
			},
			want: []lock{"X,REC_NOT_GAP 1, 1", "X,REC_NOT_GAP 1, 2", "X,REC_NOT_GAP 1, 3"},
		},
		// A-1-b. More rows back in than were deleted. (1,4) and (1,5) go
		// into the gap between (1,3) and (2,1) -- the very gap the range
		// delete of job_id=1 locked in the first report. No X,GAP appears:
		// the two new keys add no lock row of any kind.
		"append beyond deleted rows": {
			stmts: []string{
				"DELETE FROM job_attachment WHERE job_id = 1 AND ord IN (1,2,3)",
				"INSERT INTO job_attachment (job_id, ord, object_key) VALUES (1,1,'n1'), (1,2,'n2'), (1,3,'n3'), (1,4,'n4'), (1,5,'n5')",
			},
			want: []lock{"X,REC_NOT_GAP 1, 1", "X,REC_NOT_GAP 1, 2", "X,REC_NOT_GAP 1, 3"},
		},
		// A-1-c. The last job in the table: (3,2) lands in the gap before
		// the supremum pseudo-record, which is the lock that caused the
		// original deadlock. Nothing is locked on supremum.
		"append at table tail": {
			stmts: []string{
				"DELETE FROM job_attachment WHERE job_id = 3 AND ord IN (1)",
				"INSERT INTO job_attachment (job_id, ord, object_key) VALUES (3,1,'n1'), (3,2,'n2')",
			},
			want: []lock{"X,REC_NOT_GAP 3, 1"},
		},
		// Control: an INSERT with no preceding DELETE holds no explicit
		// lock at all.
		"insert only": {
			stmts: []string{
				"INSERT INTO job_attachment (job_id, ord, object_key) VALUES (3,2,'n2')",
			},
			want: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := locksHeldOn(t, "job_attachment", tt.stmts...)
			if !equalLocks(got, tt.want) {
				t.Errorf("locks:\n got: %s\nwant: %s", render(got), render(tt.want))
			}
		})
	}
}

// TestInsertedRowIsStillProtected shows what the missing lock rows above do
// and do not mean. A second transaction inserting a different key into the
// same gap goes straight through, which is the property the delete-then-
// insert pattern needs. Inserting the same key blocks, which is duplicate key
// protection, not a gap lock.
func TestInsertedRowIsStillProtected(t *testing.T) {
	requireRepeatableRead(t)

	t.Run("other session may insert into the same gap", func(t *testing.T) {
		resetSchema(t)
		ctx := t.Context()
		a, b := newSession(t), newSession(t)

		// Session A replays A-1-b: delete job 1's rows and write five back,
		// two of them into the gap in front of (2,1).
		if err := a.exec(ctx, "DELETE FROM job_attachment WHERE job_id = 1 AND ord IN (1,2,3)"); err != nil {
			t.Fatalf("session A delete: %v", err)
		}
		if err := a.exec(ctx, "INSERT INTO job_attachment (job_id, ord, object_key) VALUES (1,1,'n1'), (1,2,'n2'), (1,3,'n3'), (1,4,'n4'), (1,5,'n5')"); err != nil {
			t.Fatalf("session A insert: %v", err)
		}
		dumpLocks(t, "after session A delete+insert", "job_attachment")

		// Session B writes into the same gap, one ord further along.
		if err := b.exec(ctx, "INSERT INTO job_attachment (job_id, ord, object_key) VALUES (1,6,'b6')"); err != nil {
			t.Errorf("session B insert into the same gap: %v", err)
		}
	})

	t.Run("other session blocks on the same key", func(t *testing.T) {
		resetSchema(t)
		ctx := t.Context()
		a, b := newSession(t), newSession(t)

		if err := a.exec(ctx, "INSERT INTO job_attachment (job_id, ord, object_key) VALUES (1,4,'a4')"); err != nil {
			t.Fatalf("session A insert: %v", err)
		}
		err := b.exec(ctx, "INSERT INTO job_attachment (job_id, ord, object_key) VALUES (1,4,'b4')")
		if got := mysqlErrNo(err); got != errLockWaitTime {
			t.Errorf("session B inserting the same key: error %d (%v), want %d (lock wait timeout)", got, err, errLockWaitTime)
		}
	})

	t.Run("other session may insert at the table tail", func(t *testing.T) {
		resetSchema(t)
		ctx := t.Context()
		a, b := newSession(t), newSession(t)

		// A-1-c in two sessions: A appends past the last row of the table,
		// B appends past that.
		if err := a.exec(ctx, "DELETE FROM job_attachment WHERE job_id = 3 AND ord IN (1)"); err != nil {
			t.Fatalf("session A delete: %v", err)
		}
		if err := a.exec(ctx, "INSERT INTO job_attachment (job_id, ord, object_key) VALUES (3,1,'n1'), (3,2,'n2')"); err != nil {
			t.Fatalf("session A insert: %v", err)
		}
		dumpLocks(t, "after session A tail append", "job_attachment")

		if err := b.exec(ctx, "INSERT INTO job_attachment (job_id, ord, object_key) VALUES (3,3,'b3')"); err != nil {
			t.Errorf("session B insert at the tail: %v", err)
		}
	})
}
