package module

import (
	"context"
	"database/sql"
	"gift-grab-server/internal/account"
	"gift-grab-server/internal/leaderboard"
	notification "gift-grab-server/internal/notifications"

	"github.com/heroiclabs/nakama-common/runtime"
)

func Initialize(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	if err := leaderboard.CreateMonthlyLeaderboard(ctx, nk, logger); err != nil {
		return err
	}

	if err := leaderboard.CreateDailyTournament(ctx, nk, logger); err != nil {
		return err
	}

	if err := registerRPCs(initializer, logger, db, nk); err != nil {
		return err
	}

	return nil
}

func registerRPCs(initializer runtime.Initializer, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) error {
	if err := initializer.RegisterRpc("account_delete_id", account.AccountDeleteHandler); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	if err := initializer.RegisterRpc("notification_send", notification.SendNotificationHandler); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	if err := initializer.RegisterTournamentEnd(leaderboard.DistributeRewards); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	return nil
}
