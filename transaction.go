package goDB

import "os"

type TransactionalDatabase interface {
	Begin() os.Error
	Rollback() os.Error
	Commit() os.Error
}

type MarkableDatabase interface {
	Mark(id interface{}) os.Error
	MergeSteps(id interface{}) os.Error
	Release(id interface{}) os.Error
}

type Transaction []func(db TransactionalDatabase)

func (t Transaction) Execute(db TransactionalDatabase) (e os.Error) {
	defer func() {
		switch r := recover().(type) {
		case nil:
			e = db.Commit()
		case os.Error:
			if db.Rollback() != nil {
				panic(e)
			} else {
				e = r
			}
		default:
			panic(r)
		}
	}()

	if e = db.Begin(); e == nil {
		for _, f := range t {
			f(db)
		}
	}
	return
}

func (t Transaction) Step(db MarkableDatabase, id interface{}) (e os.Error) {
	defer func() {
		switch r := recover().(type) {
		case nil:
		case os.Error:
			if db.Release(id) != nil {
				panic(e)
			} else {
				e = r
			}
		default:
			panic(r)
		}
	}()

	if e = db.Mark(id); e == nil {
		switch tdb := db.(type) {
		case TransactionalDatabase:
			for _, f := range t {
				f(tdb)
			}
		default:
			panic(db)
		}
	}
	return
}