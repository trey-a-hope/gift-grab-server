package leaderboard

import (
	"context"
	config "gift-grab-server/pkg"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
)

func CreateMonthlyLeaderboard(ctx context.Context, nk runtime.NakamaModule, logger runtime.Logger) error {
	return nk.LeaderboardCreate(ctx,
		config.LeaderboardID,
		config.LeaderboardAuthoritative,
		config.LeaderboardSortOrder,
		config.LeaderboardOperator,
		config.LeaderboardResetSchedule,
		//TODO: Is this needed?
		make(map[string]interface{}))
}

func CreateDailyTournament(ctx context.Context, nk runtime.NakamaModule, logger runtime.Logger) error {
	startTime := int(time.Now().UTC().Unix())
	endTime := 0

	return nk.TournamentCreate(ctx,
		config.TournamentID,
		config.TournamentAuthoritative,
		config.TournamentSortOrder,
		config.TournamentOperator,
		config.TournamentResetSchedule,
		//TODO: Is this needed?
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
