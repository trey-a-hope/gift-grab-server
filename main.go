package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
)

func main() {}

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	initStart := time.Now()

	// Leaderboard information.
	id := "weekly_leaderboard"
	authoritative := false
	sortOrder := "desc"
	operator := "best"
	resetSchedule := "0 0 * * 1"
	metadata := make(map[string]interface{})

	// Create leaderboard.
	if err := nk.LeaderboardCreate(ctx, id, authoritative, sortOrder, operator, resetSchedule, metadata); err != nil {
		logger.Error("error creating leaderboard")
	}

	// Register RPC functions.
	if err := initializer.RegisterRpc("account_delete_id", AccountDeleteId); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	logger.Info("Module loaded in %dms", time.Since(initStart).Milliseconds())
	return nil
}

// Delete a user account by uid.
func AccountDeleteId(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	if err := nk.AccountDeleteId(ctx, userID, true); err != nil {
		logger.Error("Unable to delete account: %v", err)
		return "error", err
	}
	return "success", nil
}
