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

	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func Test_dbAvailable(t *testing.T) {
	if !dbAvailable() {
		t.Errorf("dbAvailable() = false")
	}
}
