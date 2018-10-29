package acr122u

import "testing"

func TestHandlerFuncServeCard(t *testing.T) {
	var handled bool

	h := HandlerFunc(func(Card) {
		handled = true
	})

	h.ServeCard(nil)

	if !handled {
		t.Fatalf("card was not handled")
	}
}
