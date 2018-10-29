package acr122u

import (
	"bytes"
	"testing"

	"github.com/ebfe/scard"
)

func TestNewStatus(t *testing.T) {
	for _, tc := range []struct {
		scs  *scard.CardStatus
		err  error
		want Status
	}{
		{
			nil,
			scard.ErrUnknownCard,
			Status{},
		},
		{
			&scard.CardStatus{
				Reader:         "Test",
				State:          scard.Present,
				ActiveProtocol: scard.ProtocolAny,
				Atr:            []byte{0x56, 0x78},
			},
			nil,
			Status{
				Reader:         "Test",
				State:          0x4,
				ActiveProtocol: 0x3,
				Atr:            []byte{0x56, 0x78},
			},
		},
	} {
		got, err := newStatus(tc.scs)
		if err != tc.err {
			t.Fatalf("unexpected error: %v", err)
		}

		if !got.equal(tc.want) {
			t.Fatalf("%#v != %#v", got, tc.want)
		}
	}
}

func (s Status) equal(o Status) bool {
	return s.Reader == o.Reader &&
		s.State == o.State &&
		s.ActiveProtocol == o.ActiveProtocol &&
		bytes.Equal(s.Atr, o.Atr)
}
