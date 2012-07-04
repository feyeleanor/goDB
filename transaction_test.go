package goDB

import "sqlite3"
import "testing"

func TestSimpleTransaction(t *testing.T) {
	sqlite3.TransientSession(func(db *sqlite3.Database) {
		(*TestDatabase)(db).createTestTables(t, FOO, BAR)
		doNothing := func(d TransactionalDatabase) {}
		raiseOK := func(d TransactionalDatabase) { panic(sqlite3.OK) }
		raiseErrno := func(d TransactionalDatabase) { panic(sqlite3.MISUSE) }

		fatalOnError(t, Transaction{}.Execute(db), "empty transaction")
		fatalOnError(t, Transaction{doNothing}.Execute(db), "inconsequential transaction")
		fatalOnError(t, Transaction{raiseOK}.Execute(db), "transaction raises OK")
		fatalOnSuccess(t, Transaction{raiseErrno}.Execute(db), "transaction raises Errno")
	})
}
