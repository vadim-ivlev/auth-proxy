package server

import "testing"

func TestFullNameValidate(t *testing.T) {
	// fullName := "Иван Иванов"
	// fullNameValidate(fullName, "Тест пройден")

	// fullNameFail := "Иван Иванов"
	fullNameFail := "Иван https://clck.ru/34Epbu"
	fullNameValidate(fullNameFail, "Тест не пройден")
	// if !DbAvailable() {
	// 	t.Errorf("dbAvailable() = false")
	// }
}
