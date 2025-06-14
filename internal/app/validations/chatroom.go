package validations

import "errors"

func ValidateChatRoom(roomName, from string) error {
	if roomName == "" {
		return errors.New("empty name")
	}

	

	return nil
}
