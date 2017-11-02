package database

import (
	"bytes"
	"context"
	"encoding/gob"

	"cloud.google.com/go/pubsub"
)

type Pubsub struct {
	client *pubsub.Client
	topic  *pubsub.Topic
}

// SetupPubsub global client
func NewPubsub(project, topic string) (*Pubsub, error) {
	client, err := pubsub.NewClient(context.Background(), project)

	if err != nil {
		return nil, err
	}
	t := client.Topic(topic)
	exists, err := t.Exists(context.Background())

	if err != nil {
		return nil, err
	}

	if !exists {
		_, err = client.CreateTopic(context.Background(), topic)
		if err != nil {
			return nil, err
		}
	}

	return &Pubsub{
		client: client,
		topic:  t,
	}, nil
}

func (p *Pubsub) Read() (chan *Message, error) {
	return nil, nil
}

func (p *Pubsub) Save(message *Message) error {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(message)

	if err != nil {
		return err
	}

	res := p.topic.Publish(context.Background(), &pubsub.Message{
		Data: b.Bytes(),
	})

	<-res.Ready()
	return nil
}
