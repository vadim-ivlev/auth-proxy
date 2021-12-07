package mail

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	ReadConfig("../../configs/mail.yaml", "dev")
	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

// func Test_SendPassword(t *testing.T) {
// 	err := SendMessage("new_password", "vadim", "ivlev@rg.ru", "123456")
// 	assert.NilError(t, err)
// }
