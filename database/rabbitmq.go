package database

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

type Connection struct {
	host       string
	errors     chan *amqp.Error
	connection *amqp.Connection
	channels   []*amqp.Channel
}

type RabbitMQ struct {
	connections []*Connection
}

func SetupRabbitMQ() {
	hosts := viper.GetStringSlice("rabbitmq.hosts")
	NewRabbitMQ(hosts)
}

func NewRabbitMQ(hosts []string) *RabbitMQ {
	rmq := &RabbitMQ{}
	for _, host := range hosts {
		rmq.connections = append(rmq.connections, newConnection(host))
	}
	return rmq
}

func newConnection(host string) *Connection {
	c := &Connection{
		host:   host,
		errors: make(chan *amqp.Error),
	}
	go c.watch()
	c.dial()
	return c
}

// dial a connection until we connect, and redial on error (with backoff)
func (c *Connection) dial() {
	connection, err := amqp.Dial(c.host)
	if err != nil {
		log.Error().
			Err(err).
			Str("host", c.host).
			Msg("could not connect to rabbitmq")
		c.errors <- &amqp.Error{
			Reason: err.Error(),
		}
		return
	}
	c.connection = connection
	c.connection.NotifyClose(c.errors)
}

// watch a connection handler
func (c *Connection) watch() {
	<-c.errors
	// Everything is invalid, reboot.
	c.reset()
	<-time.After(time.Second * 3)
	c.dial()
	go c.watch()
}

func (c *Connection) reset() {
	if c.connection != nil {
		c.connection.Close()
		c.connection = nil
	}
	c.channels = nil
	close(c.errors)
	c.errors = make(chan *amqp.Error)
}
