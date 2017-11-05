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

	return &Pubsub{
		client:       client,
		topic:        t,
		subscription: s,
	}, nil
}

func (p *Pubsub) Read() (chan *Message, error) {
	c := make(chan *Message)

	// We need to wrap this so we don't lose context here.
	func(ic chan *Message) {
		p.subscription.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
			// Decode our message
			var message *Message
			b := bytes.Buffer{}
			b.Write(msg.Data)
			d := gob.NewDecoder(&b)
			d.Decode(&message)

			// Create our response channel
			message.Ack = make(chan bool)

			// Send our decoded message + response channel
			ic <- message

			// Wait for a response and handle it.
			ack := <-message.Ack
			if ack {
				msg.Ack()
			} else {
				msg.Nack()
			}
		})
	}(c)

	return c, nil
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

	_, err = res.Get(context.Background())
	if err != nil {
		return err
	}
	return nil
}
