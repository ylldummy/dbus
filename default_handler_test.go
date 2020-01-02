package dbus

import (
	"testing"
)

func TestSignalOrder(t *testing.T) {
	const sigCount = 20
	sigHandler := NewDefaultSignalHandler()
	defer sigHandler.Terminate()
	receiver := make(chan *Signal)
	sigHandler.AddSignal(receiver)
	for i := 0; i < sigCount; i++ {
		sig := &Signal{
			Body: []interface{}{i},
		}
		sigHandler.DeliverSignal("", "", sig)
	}
	for i := 0; i < sigCount; i++ {
		sig := <-receiver
		if len(sig.Body) != 1 {
			t.Fatalf("Unexpect signal %v", sig)
		}
		n, ok := sig.Body[0].(int)
		if !ok {
			t.Fatalf("Unexpect signal body content %v", sig.Body[0])
		}
		if i != n {
			t.Fatalf("Expect signal #%d, got #%d", i, n)
		}
	}
}
