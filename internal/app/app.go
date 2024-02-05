package app

import (
	"github.com/OksidGen/grpc_thumbnail/server/internal/delivery"
	"github.com/OksidGen/grpc_thumbnail/server/internal/repository"
	"github.com/OksidGen/grpc_thumbnail/server/internal/usecase"
	"github.com/OksidGen/grpc_thumbnail/server/proto"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	db, err := sqlx.Open("sqlite3", "./thumbnail.db")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to close database")
		}
	}(db)

	thumbnailRepo := repository.NewThumbnailRepository(db)
	thumbnailUsecase := usecase.NewThumbnailUsecase(thumbnailRepo)
	thumbnailHandler := delivery.NewThumbnailHandler(thumbnailUsecase)

	server := delivery.NewGrpcServer()
	proto.RegisterThumbnailServiceServer(server.Server, thumbnailHandler)

	errs := make(chan error, 2)
	go func() {
		errs <- server.Start()
	}()

	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		errs <- server.WaitForShutdown(signalChan)
	}()

	err = <-errs
	log.Fatal().Err(err).Msg("Failed to start gRPC server")
}
