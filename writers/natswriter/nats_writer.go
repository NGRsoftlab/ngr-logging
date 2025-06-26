package natswriter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsWriter struct {
	JetStream jetstream.JetStream
	hosts     []string
	user      string
	passwd    string
	stream    string
}

var _ io.Writer = &NatsWriter{}

func NewNatsWriter(hosts []string, stream string, maxMessages, maxBytes int64, maxAge time.Duration) (*NatsWriter, error) {
	js, err := connect(hosts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a stream
	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     stream,
		Subjects: []string{"logger.>"},
		MaxMsgs:  maxMessages,
		MaxBytes: maxBytes,
		MaxAge:   maxAge,
		Replicas: 3,
		// NoAck:    true,
	})
	if err != nil {
		return nil, err
	}

	writer := NatsWriter{
		JetStream: js,
		hosts:     hosts,
		stream:    stream,
	}

	return &writer, nil
}

func connect(hosts []string) (jetstream.JetStream, error) {

	nc, err := nats.Connect(strings.Join(hosts, ","))
	if err != nil {
		return nil, err
	}

	// Create a JetStream management interface
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}

	return js, err
}

func (w *NatsWriter) Write(p []byte) (n int, err error) {
	if w.JetStream == nil {
		w.JetStream, err = connect(w.hosts)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// fmt.Println(string(p))

	row := make(map[string]interface{})
	err = json.Unmarshal(p, &row)
	if err != nil {
		fmt.Println(err)
	}

	if _, err := w.JetStream.Publish(ctx, fmt.Sprintf("logger.%s.%s", row["component"], row["level"]), p); err != nil {
		fmt.Println("pub error: ", err)
		return 0, err
	}

	return len(p), nil
}
