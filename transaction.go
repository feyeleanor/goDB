package goDB

type TransactionalDatabase interface {
	Begin() error
	Rollback() error
	Commit() error
}

type MarkableDatabase interface {
	Mark(id interface{}) error
	MergeSteps(id interface{}) error
	Release(id interface{}) error
}

type Transaction []func(db TransactionalDatabase)

func (t Transaction) Execute(db TransactionalDatabase) (e error) {
	defer func() {
		switch r := recover().(type) {
		case nil:
			e = db.Commit()
		case error:
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

func (t Transaction) Step(db MarkableDatabase, id interface{}) (e error) {
	defer func() {
		switch r := recover().(type) {
		case nil:
		case error:
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
