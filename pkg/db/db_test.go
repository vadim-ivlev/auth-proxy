package db

import (
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	ReadEnvConfig("../../configs/db.env")

	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func Test_dbAvailable(t *testing.T) {
	if !DbAvailable() {
		t.Errorf("dbAvailable() = false")
	}
}

// Локальная БД
func Benchmark_local_DB(b *testing.B) {
	ReadEnvConfig("../../configs/db.env")
	for i := 0; i < b.N; i++ {
		_, _ = QueryRowMap("select $1", i)
	}
}

// Удаленная БД
func Benchmark_remote_DB(b *testing.B) {
	ReadEnvConfig("../../configs/db.env.prod")
	for i := 0; i < b.N; i++ {
		_, _ = QueryRowMap("select $1", i)
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
