package deje

import (
	"github.com/jcelliott/turnpike"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_NewClient(t *testing.T) {
	topic := "com.example.deje.some-doc"
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
	topic := "com.example.deje.some-doc"
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
	topic := "com.example.deje.some-doc"
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
