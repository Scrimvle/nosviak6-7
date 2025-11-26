package tui

import "bytes"

type KeyPress struct {
	code []byte
	onEvent func() bool
}

// AddKeyPressListener will add a brand new listener for when that byte is found
func (t *Terminal) AddKeyPressListener(code string, event func() bool) {
	t.entities = append(t.entities, &KeyPress{
		code: []byte(code),
		onEvent: event,
	})
}

// checkListener will check if the data contains the object.code
func (t *Terminal) checkListener(data []byte) bool {
	for _, item := range t.entities {
		object, ok := item.(*KeyPress)
		if !ok || object == nil || !bytes.Equal(data, object.code) {
			continue
		}

		return object.onEvent()
	}

	return false
}