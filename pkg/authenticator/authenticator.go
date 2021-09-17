package authenticator

import "errors"

// IsPinGood верен ли PIN введенный пользователем
func IsPinGood(pin, username string) error {
	if len(pin) != 6 {
		return errors.New("PIN must be 6 symbols")
	}
	return nil
}

func GetBarcode(username string) {

}
