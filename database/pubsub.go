package database

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	"cloud.google.com/go/pubsub"
)

type Pubsub struct {
	client       *pubsub.Client
	topic        *pubsub.Topic
	subscription *pubsub.Subscription
}

// SetupPubsub global client
func NewPubsub(project, topic, subscription string) (*Pubsub, error) {
	client, err := pubsub.NewClient(context.Background(), project)

	if err != nil {
		return nil, err
	}

	// Create our topic
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
	// Create our subscription
	s := client.Subscription(subscription)
	exists, err = s.Exists(context.Background())

	if err != nil {
		return nil, err
	}

	if !exists {
		_, err = client.CreateSubscription(context.Background(), subscription, pubsub.SubscriptionConfig{
			Topic: t,
		})
		if err != nil {
			fmt.Printf("%s", err.Error())
			return nil, err
		}
	}

	// TODO: configure this
	s.ReceiveSettings.NumGoroutines = 1
	s.ReceiveSettings.MaxOutstandingMessages = 10

	return &Pubsub{
		client:       client,
		topic:        t,
		subscription: s,
	}, nil
}

func (p *Pubsub) Read(rctx context.Context, cb MessageFn) error {

	// Use a callback to fire a message up the ladder to the caller
	// and the caller returns an ack/nack response.
	return p.subscription.Receive(rctx, func(ctx context.Context, msg *pubsub.Message) {
		// Decode our message
		var message Message
		var b bytes.Buffer
		b.Write(msg.Data)
		d := gob.NewDecoder(&b)
		d.Decode(&message)
		if cb(&message) {
			msg.Ack()
		} else {
			msg.Nack()
		}

	})
}

func (p *Pubsub) Save(message *Message) error {
	var b bytes.Buffer
	e := gob.NewEncoder(&b)
	err := e.Encode(message)

	if err != nil {
		return err
	}

	res := p.topic.Publish(context.Background(), &pubsub.Message{
		Data: b.Bytes(),
	})

	_, err = res.Get(context.Background())
	if err != nil {
		return err
	}
	return nil
}
