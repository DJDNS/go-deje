package deje

import (
	"github.com/jcelliott/turnpike"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestClient_NewClient(t *testing.T) {
	topic := "http://example.com/deje/some-doc"
	client := NewClient(topic)

	if client.Doc.Topic != topic {
		t.Fatalf("Expected topic '%s', got '%s'", topic, client.Doc.Topic)
	}

	if client.tpClient == nil {
		t.Fatal("client.tpClient should not be nil pointer")
	}
}

func setupServer() (string, func()) {
	tp_server := turnpike.NewServer()
	server := httptest.NewServer(tp_server.Handler)
	server_addr := strings.Replace(server.URL, "http", "ws", 1)
	return server_addr, func() {
		// Need extra Oomph to not block on open client connections
		server.CloseClientConnections()
		server.Close()
	}
}

func TestClient_Connect(t *testing.T) {
	topic := "http://example.com/deje/some-doc"
	client := NewClient(topic)
	server_addr, server_closer := setupServer()
	defer server_closer()

	err := client.Connect("foo")
	if err == nil {
		t.Fatal("foo is not a real server - should not 'succeed'")
	}

	err = client.Connect(server_addr)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_Connect_WithCallback(t *testing.T) {
	topic := "http://example.com/deje/some-doc"
	client := NewClient(topic)
	server_addr, server_closer := setupServer()
	defer server_closer()

	var callback_called bool
	var callback_called_with string
	client.SetConnectCallback(func(sessionId string) {
		callback_called = true
		callback_called_with = sessionId
	})

	err := client.Connect(server_addr)
	if err != nil {
		t.Fatal(err)
	}
	if !callback_called {
		t.Fatal("Callback was not called")
	}
	if callback_called_with == "" {
		t.Fatal("Callback was not given a sessionId value")
	}
}

func TestClient_PubSub(t *testing.T) {
	topic := "http://example.com/deje/some-doc"
	client1 := NewClient(topic)
	client2 := NewClient(topic)
	server_addr, server_closer := setupServer()
	defer server_closer()

	// Connect both clients
	if err := client1.Connect(server_addr); err != nil {
		t.Fatal(err)
	}
	if err := client2.Connect(server_addr); err != nil {
		t.Fatal(err)
	}

	// Give client1 a subscription callback
	response_chan := make(chan interface{}, 1)
	client1.SetEventCallback(func(event interface{}) {
		response_chan <- event
	})

	// Publish via client2
	sent := map[string]interface{}{
		"foo":  "bar",
		"fire": []interface{}{"ant", "place", "nation"},
	}
	if err := client2.Publish(sent); err != nil {
		t.Fatal(err)
	}

	select {
	case recvd := <-response_chan:
		if !reflect.DeepEqual(recvd, sent) {
			t.Fatalf("Expected %v, got %v", sent, recvd)
		}
	case <-time.After(15 * time.Millisecond):
		t.Fatal("Recv timed out")
	}

	// Confirm that client2 doesn't break if onEvent == nil
	client1.Publish(sent)

	// And that client1 does not receive its own events
	if len(response_chan) > 0 {
		t.Fatal("client1 recvd response from self")
	}

	// Finally, confirm that callbacks can be set after the fact
	client2.SetEventCallback(func(event interface{}) {
		response_chan <- "I can recv too!"
	})
	if err := client1.Publish(sent); err != nil {
		t.Fatal(err)
	}

	select {
	case recvd := <-response_chan:
		if !reflect.DeepEqual(recvd, "I can recv too!") {
			t.Fatalf("Expected %v, got %v", sent, recvd)
		}
	case <-time.After(15 * time.Millisecond):
		t.Fatal("Recv timed out")
	}

}
