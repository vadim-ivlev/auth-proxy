package server

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestValidateXReqID(t *testing.T) {

	h := sha256.New()
	h.Write([]byte("Для тестирования хеша"))
	fmt.Printf("%x\n", h.Sum(nil))

	checkReqID := hex.EncodeToString(h.Sum(nil))

	fmt.Println(checkReqID)

}
