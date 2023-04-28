package server

import "testing"

func TestFullNameValidate(t *testing.T) {
	// fullName := "Иван Иванов"
	// fullNameValidate(fullName, "Тест пройден")

	fullNameFail := "Иван Иванов и тут какой-то длинный текст"
	fullNameValidate(fullNameFail, "Тест не пройден")
	// if !DbAvailable() {
	// 	t.Errorf("dbAvailable() = false")
	// }
}
