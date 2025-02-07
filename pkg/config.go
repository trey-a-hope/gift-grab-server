package config

const (
	// Leaderboard configs
	LeaderboardID            = "monthly_leaderboard"
	LeaderboardAuthoritative = false
	LeaderboardSortOrder     = "desc"
	LeaderboardOperator      = "best"
	LeaderboardResetSchedule = "0 0 1 * *" // 1st of each month at midnight.

	// Tournament configs
	TournamentID            = "daily_tournament"
	TournamentAuthoritative = false
	TournamentSortOrder     = "desc"
	TournamentOperator      = "best"
	TournamentResetSchedule = "0 12 * * *" // Noon UTC each day
	TournamentTitle         = "Daily Dash"
	TournamentDescription   = "Dash past your opponents for high scores and big rewards!"
	TournamentCategory      = 1
	TournamentDuration      = 3600 // 1 hour in seconds
	TournamentMaxSize       = 10000
	TournamentMaxAttempts   = 3
	TournamentJoinRequired  = false
)
