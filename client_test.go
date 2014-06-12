package deje

import "testing"

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
