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

func TestSignalQueue(t *testing.T) {
	q := newSignalQueue()
	pushSignal := func(i int) {
		q.Push(&Signal{
			Body: []interface{}{i},
		})
	}
	checkSignal := func(sig *Signal, i int) {
		if len(sig.Body) != 1 {
			t.Fatalf("Unexpected signal %v", sig)
		}
		n, ok := sig.Body[0].(int)
		if !ok {
			t.Fatalf("Unexpected signal body content %v", sig.Body[0])
		}
		if i != n {
			t.Fatalf("Expect signal #%d, got #%d", i, n)
		}
	}
	popAndCheck := func(i int) {
		if sig, ok := q.Next(); !ok {
			t.Fatal("Unexpected empty queue in Next call")
		} else {
			checkSignal(sig, i)
		}
		if sig, ok := q.Pop(); !ok {
			t.Fatal("Unexpected empty queue in Pop call")
		} else {
			checkSignal(sig, i)
		}
	}
	checkEmpty := func() {
		if _, ok := q.Next(); ok {
			t.Fatal("Expect Next to fail but succeeded")
		}
		if _, ok := q.Pop(); ok {
			t.Fatal("Expect Pop to fail but succeeded")
		}
	}
	pushSignal(1)
	popAndCheck(1)
	checkEmpty()
	pushSignal(2)
	pushSignal(3)
	popAndCheck(2)
	pushSignal(4)
	popAndCheck(3)
	popAndCheck(4)
	checkEmpty()
}
