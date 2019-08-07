package db

import (
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	fmt.Println("Тесты DB ******************************************************")
	ReadConfig("../../configs/db.yaml", "dev")
	ReadSQLiteConfig("../../configs/sqlite.yaml", "dev")

	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func Test_dbAvailable(t *testing.T) {
	if !dbAvailable() {
		t.Errorf("dbAvailable() = false")
	}
}

func Benchmark_DB_connection(b *testing.B) {
	// var res map[string]interface{}
	// var err error
	for index := 0; index < 1000; index++ {
		// res, err = QueryRowMap("select * from app where appname = 'auth'")
		QueryRowMap("select * from app where appname = 'auth'")
	}
	// println(res["appname"].(string), err)
}
