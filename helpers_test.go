package goDB

import "fmt"
import "encoding/gob"

import "sqlite3"
import "testing"

type TestDatabase sqlite3.Database

var FOO *sqlite3.Table
var BAR *sqlite3.Table

func init() {
	FOO = &sqlite3.Table{"foo", "number INTEGER, text VARCHAR(20)"}
	BAR = &sqlite3.Table{"bar", "number INTEGER, value BLOB"}
}

type TwoItems struct {
	Number string
	Text   string
}

func (t *TwoItems) String() string {
	return "[" + t.Number + " : " + t.Text + "]"
}

func fatalOnError(t *testing.T, e error, message string, parameters ...interface{}) {
	if e != nil {
		t.Fatalf("%v : %v", e, fmt.Sprintf(message, parameters...))
	}
}

func fatalOnSuccess(t *testing.T, e error, message string, parameters ...interface{}) {
	if e == nil {
		t.Fatalf("%v : %v", e, fmt.Sprintf(message, parameters...))
	}
}

func (db *TestDatabase) stepThroughRows(t *testing.T, table *sqlite3.Table) (c int) {
	var e error
	sql := fmt.Sprintf("SELECT * from %v;", table.Name)
	c, e = (*sqlite3.Database)(db).Execute(sql, func(st *sqlite3.Statement, values ...interface{}) {
		data := values[1]
		switch data := data.(type) {
		case *gob.Decoder:
			blob := &TwoItems{}
			data.Decode(blob)
			t.Logf("BLOB =>   %v: %v, %v: %v\n", sqlite3.ResultColumn(0).Name(st), sqlite3.ResultColumn(0).Value(st), st.ColumnName(1), blob)
		default:
			t.Logf("TEXT => %v: %v, %v: %v\n", sqlite3.ResultColumn(0).Name(st), sqlite3.ResultColumn(0).Value(st), st.ColumnName(1), st.Column(1))
		}
	})
	fatalOnError(t, e, "%v failed on step %v", sql, c)
	if rows, _ := table.Rows((*sqlite3.Database)(db)); rows != c {
		t.Fatalf("%v: %v rows expected, %v rows found", table.Name, rows, c)
	}
	return
}

func (db *TestDatabase) runQuery(t *testing.T, sql string, params ...interface{}) {
	st, e := (*sqlite3.Database)(db).Prepare(sql, params...)
	fatalOnError(t, e, st.SQLSource())
	st.Step()
	st.Finalize()
}

func (db *TestDatabase) createTestTables(t *testing.T, tables ...*sqlite3.Table) {
	for _, table := range tables {
		testdb := (*sqlite3.Database)(db)
		table.Drop(testdb)
		table.Create(testdb)
		if c, _ := table.Rows(testdb); c != 0 {
			t.Fatalf("%v already contains data", table.Name)
		}
	}
}

func (db *TestDatabase) populate(t *testing.T, table *sqlite3.Table) {
	switch table.Name {
	case "foo":
		db.runQuery(t, "INSERT INTO foo values (1, 'this is a test')")
		db.runQuery(t, "INSERT INTO foo values (?, ?)", 2, "holy moly")
		if c, _ := table.Rows((*sqlite3.Database)(db)); c != 2 {
			t.Fatal("Failed to populate %v", table.Name)
		}
	case "bar":
		db.runQuery(t, "INSERT INTO bar values (1, 'this is a test')")
		db.runQuery(t, "INSERT INTO bar values (?, ?)", 2, "holy moly")
		db.runQuery(t, "INSERT INTO bar values (?, ?)", 3, TwoItems{"holy moly", "guacomole"})
		if c, _ := table.Rows((*sqlite3.Database)(db)); c != 3 {
			t.Fatal("Failed to populate %v", table.Name)
		}
	}
}
