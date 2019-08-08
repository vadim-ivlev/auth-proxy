package db

import (
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	ReadConfig("../../configs/db.yaml", "dev")
	ReadSQLiteConfig("../../configs/sqlite.yaml", "dev")
	SQLite = false

	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func Test_dbAvailable(t *testing.T) {
	if !dbAvailable() {
		t.Errorf("dbAvailable() = false")
	}
}

// Локальная БД
func Benchmark_local_DB(b *testing.B) {
	ReadConfig("../../configs/db.yaml", "dev")
	UsePool = false
	SQLite = false
	for i := 0; i < b.N; i++ {
		QueryRowMap("select $1", i)
	}
}

// Локальная БД с пулом
func Benchmark_local_DB_pool(b *testing.B) {
	ReadConfig("../../configs/db.yaml", "dev")
	UsePool = true
	SQLite = false
	for i := 0; i < b.N; i++ {
		QueryRowMap("select $1", i)
	}
}

// Удаленная БД
func Benchmark_remote_DB(b *testing.B) {
	ReadConfig("../../configs/db.yaml", "prod")
	UsePool = false
	SQLite = false
	for i := 0; i < b.N; i++ {
		QueryRowMap("select $1", i)
	}
}

// Удаленная БД с пулом
func Benchmark_remote_DB_pool(b *testing.B) {
	ReadConfig("../../configs/db.yaml", "prod")
	UsePool = true
	SQLite = false
	for i := 0; i < b.N; i++ {
		QueryRowMap("select $1", i)
	}
}

// SQLite
func Benchmark_SQLite(b *testing.B) {
	UsePool = false
	SQLite = true
	for i := 0; i < b.N; i++ {
		QueryRowMap("select $1", i)
	}
}

// SQLite с пулом
func Benchmark_SQLite_pool(b *testing.B) {
	UsePool = true
	SQLite = true
	for i := 0; i < b.N; i++ {
		QueryRowMap("select $1", i)
	}
}

var args = map[string]interface{}{
	"a": 1,
	"b": 2,
	"c": 3,
	"d": 4,
	"e": 5,
	"f": 6,
	"g": 7,
	"h": 8,
}

var fields []string
var values []interface{}
var placeholders []string

func Benchmark_getKeysAndValues(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fields, values, placeholders = getKeysAndValues(args)
	}

}
