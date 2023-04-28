package server

import "testing"

func TestFullNameValidate(t *testing.T) {
	// fullName := "Иван Иванов"
	// fullNameValidate(fullName, "Тест пройден")

	fullNameFail := "Иван http Иванов"
	fullNameValidate(fullNameFail, "Тест не пройден")
	// if !DbAvailable() {
	// 	t.Errorf("dbAvailable() = false")
	// }
}
