package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"

	pb "github.com/Cidan/sheep/api/v1"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/Cidan/sheep/config"
	"github.com/Cidan/sheep/database"
	pubsub "github.com/Cidan/sheep/database/pubsub"
	spanner "github.com/Cidan/sheep/database/spanner"
	"github.com/Cidan/sheep/util"
	"github.com/labstack/echo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"google.golang.org/grpc"
)

var e *echo.Echo

func main() {
	setupLogging()
	config.Setup("")

	stream, err := setupQueue()
	if err != nil {
		log.Panic().Err(err).Msg("Could not setup queue")
	}

	database, err := setupDatabase()
	if err != nil {
		log.Panic().Err(err).Msg("Could not setup database")
	}

	if viper.GetBool("worker") {
		go setupWorker(stream, database)
	}

	if viper.GetBool("master") {
		go startGrpc(stream, database)
		go startWeb(stream, database)
	}

	util.WaitForSigInt()

	log.Info().Msg("Shutting down...")
}

// TODO: This should be a function of the Stream
func setupWorker(stream database.Stream, db database.Database) {
	// TODO: Add context
	// TODO: This can drop into an infinite loop if an error is permanently fatal
	// but max retry breaks our promise of 100% correct -- message must be validated
	// on PUT into queue before HTTP return.
	// TODO: Check stream.Read return error
	// For loop here is a guard in case read fails, read blocks/loops.
	for {
		stream.Read(context.Background(), func(msg *database.Message) bool {
			err := db.Save(msg)
			if err != nil {
				log.Error().Err(err).Msg("Unable to save stream message")
				return false
			}
			return true
		})
	}
}

func setupDatabase() (database.Database, error) {
	if viper.GetBool("cockroachdb.enabled") {
		return nil, errors.New("CockroachDB support is not ready")
		/*
			return database.NewCockroachDB(
				viper.GetString("cockroachdb.host"),
				viper.GetString("cockroachdb.username"),
				viper.GetString("cockroachdb.password"),
				viper.GetString("cockroachdb.dbname"),
				viper.GetString("cockroachdb.sslmode"),
				viper.GetInt("cockroachdb.port"),
			)
		*/
	}

	if viper.GetBool("spanner.enabled") {
		log.Info().Msg("Setting up Spanner connection and schema")
		return spanner.New(
			viper.GetString("spanner.project"),
			viper.GetString("spanner.instance"),
			viper.GetString("spanner.database"),
		)
	}

	return nil, fmt.Errorf("no database enabled")
}

func setupQueue() (database.Stream, error) {
	if viper.GetBool("rabbitmq.enabled") {
		return nil, errors.New("RabbitMQ support is not ready")
		/*
			return database.NewRabbitMQ(
				viper.GetStringSlice("rabbitmq.hosts"),
			)
		*/
	}

	if viper.GetBool("pubsub.enabled") {
		log.Info().Msg("Setting up Pub/Sub connection, topic, and subscription")
		return pubsub.New(
			viper.GetString("pubsub.project"),
			viper.GetString("pubsub.topic"),
			viper.GetString("pubsub.subscription"),
		)
	}

	return nil, fmt.Errorf("no queue setup")
}

func startGrpc(stream database.Stream, database database.Database) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", viper.GetInt("service.port")))
	if err != nil {
		log.Panic().
			Err(err).
			Msg("error listening on gRPC port")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterV1Server(grpcServer, &pb.API{
		Stream:   stream,
		Database: database,
	})
	log.Info().
		Int("port", viper.GetInt("service.port")).
		Msg("Started gRPC server, you're good to go!")

	grpcServer.Serve(lis)
}

func startWeb(stream database.Stream, database database.Database) {
	mux := runtime.NewServeMux()
	err := pb.RegisterV1HandlerFromEndpoint(
		context.Background(),
		mux,
		fmt.Sprintf("localhost:%d", viper.GetInt("service.port")),
		[]grpc.DialOption{grpc.WithInsecure()})

	if err != nil {
		log.Panic().
			Err(err).
			Msg("error setting up gRPC gateway")
	}
	http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("service.rest")), mux)
}

func setupLogging() {
	// If we're in a terminal, pretty print
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		var level zerolog.Level
		switch viper.GetString("level") {
		case "info":
			level = zerolog.InfoLevel
		case "warn":
			level = zerolog.WarnLevel
		case "error":
			level = zerolog.ErrorLevel
		case "debug":
			level = zerolog.DebugLevel
		default:
			level = zerolog.InfoLevel
		}
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(level)
		log.Info().Msg("Detected terminal, pretty logging enabled.")
	}
}
