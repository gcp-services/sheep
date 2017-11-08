# sheep
[![Build Status](https://travis-ci.org/Cidan/sheep.svg?branch=master)](https://travis-ci.org/Cidan/sheep)

- [Sheep?](#sheep)
- [How It Works](#how-it-works)
- [Docs](#docs)

## Sheep?

Sheep is a distributed, idempotent, and eventually consistent counter service backed by different backends. When properly configured, Sheep will guarantee accurate counts, so long as the underlying databases, i.e. Pub/Sub, Spanner, are durable.

Sheep is built to scale. By leveraging scalable systems and designs, Sheep allows you to keep fully idempotent and accurate counts at a rate of hudreds of thousands of requests a second.

Backends are pluggable, with current support for Google Spanner and CockroachDB(WIP) on the storage side, and with Google Pub/Sub and RabbitMQ(WIP) support on the transport side.

## How It Works

Sheep runs in two modes: master, and worker.

When in master mode, Sheep accepts a very simple REST API in order to submit transactions and to get counter data. Each transaction (incr, decr) must contain a `keyspace`, `key`, `name`, and a `uuid`. When a transaction is submitted, the format is validated, then committed to the configured streaming backend. It's important the caller to Sheep retries until a `200` is returned by Sheep. Multiple submissions of the same transaction using the same UUID are okay -- they will be deduped by the worker.

When in worker mode, Sheep pulls transactions off of the configured stream and applies the changes to permanent storage. Only once a transaction has been committed to disk will the service acknowledge the transaction from the stream. Retries are okay -- the storage system will keep a transaction log and ensure that any transaction only happens once.

## Docs

Sheep reads in a configuration YAML called `sheep.yml` from one of the following locations:
```
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/")
	viper.AddConfigPath("/etc/sheep/")
```

A sample config is located in the `config` directory.

The API looks like this:

```
# Get a counter value
curl "localhost:5309/v1/get?keyspace=test&key=test&name=some%20counter"

# Increment a counter
curl -X PUT -H "Content-Type: application/json" "localhost:5309/v1/incr" -d \ 
  '{"UUID": "aabbcc", "Keyspace": "users", "Key": "some-user-id", "Name": "counter-name"}'

# Decrement a counter
curl -X PUT -H "Content-Type: application/json" "localhost:5309/v1/decr" -d \ 
  '{"UUID": "aabbcc", "Keyspace": "users", "Key": "some-user-id", "Name": "counter-name"}'
```