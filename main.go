package main

import (
	"context"
	"database/sql"
	"gift-grab-server/internal/module"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
)

func main() {}

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	initStart := time.Now()

	if err := module.Initialize(ctx, logger, db, nk, initializer); err != nil {
		logger.Error("Failed to initialize module: %v", err)
		return err
	}

	logger.Info("Module loaded in %dms", time.Since(initStart).Milliseconds())
	return nil
}
