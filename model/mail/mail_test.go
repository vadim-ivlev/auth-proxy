package mail

import (
	"os"
	"testing"

	"gotest.tools/assert"
)

func TestMain(m *testing.M) {
	ReadConfig("../../configs/mail.yaml", "dev")
	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func Test_SendPassword(t *testing.T) {
	err := SendPassword("ivlev@rg.ru", "123456")
	assert.NilError(t, err)
}
