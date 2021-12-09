package mail

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// ReadConfig("../../configs/mail.yaml", "dev")
	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}
