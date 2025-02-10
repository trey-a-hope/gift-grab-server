package leaderboard

import (
	"context"
	"database/sql"
	config "gift-grab-server/pkg"
	"time"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

func DistributeRewards(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, tournament *api.Tournament, end int64, reset int64) error {
	wallets := []*runtime.WalletUpdate{}
	notifications := []*runtime.NotificationSend{}
	content := map[string]interface{}{}
	changeset := map[string]int64{"coins": 100}
	records, _, _, _, err := nk.LeaderboardRecordsList(ctx, tournament.Id, []string{}, 10, "", reset)
	for _, record := range records {
		wallets = append(wallets, &runtime.WalletUpdate{record.OwnerId, changeset, content})
		notifications = append(notifications, &runtime.NotificationSend{record.OwnerId, "Leaderboard winner", content, 1, "", true})
	}
	_, err = nk.WalletsUpdate(ctx, wallets, false)
	if err != nil {
		logger.Error("failed to update winner wallets: %v", err)
		return err
	}
	err = nk.NotificationsSend(ctx, notifications)
	if err != nil {
		logger.Error("failed to send winner notifications: %v", err)
		return err
	}
	return nil
}

func CreateMonthlyLeaderboard(ctx context.Context, nk runtime.NakamaModule, logger runtime.Logger) error {
	return nk.LeaderboardCreate(ctx,
		config.LeaderboardID,
		config.LeaderboardAuthoritative,
		config.LeaderboardSortOrder,
		config.LeaderboardOperator,
		config.LeaderboardResetSchedule,
		make(map[string]interface{}))
}

// A tournament created with:

// 1. Only duration starts immediately and ends after the set duration
// 2. A duration and resetSchedule starts immediately and closes after the set duration, then resets and starts again on the defined schedule
// 3. A duration, resetSchedule, and endTime starts immediately and closes after the set duration, then resets and starts again on the defined schedule until the end time is reached

// If an endTime is set, that timestamp marks the definitive end of the tournament, regardless of any resetSchedule or duration values.

func CreateDailyTournament(ctx context.Context, nk runtime.NakamaModule, logger runtime.Logger) error {
	startTime := int(time.Now().UTC().Unix())
	endTime := 0 // never end, repeat the tournament each day forever

	return nk.TournamentCreate(ctx,
		config.TournamentID,
		config.TournamentAuthoritative,
		config.TournamentSortOrder,
		config.TournamentOperator,
		config.TournamentResetSchedule,
		map[string]interface{}{},
		config.TournamentTitle,
		config.TournamentDescription,
		config.TournamentCategory,
		startTime,
		endTime,
		config.TournamentDuration,
		config.TournamentMaxSize,
		config.TournamentMaxAttempts,
		config.TournamentJoinRequired)
}
