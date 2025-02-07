package main

import (
	"context"
	"database/sql"
	"encoding/json"
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

	// Tournament information.
	// tournament_id := uuid.Must(uuid.NewV4())
	tournament_id := "tournament"
	tournament_authoritative := false        // true by default
	tournament_sortOrder := "desc"           // one of: "desc", "asc"
	tournament_operator := "best"            // one of: "best", "set", "incr"
	tournament_resetSchedule := "0 12 * * *" // noon UTC each day
	tournament_metadata := map[string]interface{}{}
	tournament_title := "Daily Dash"
	tournament_description := "Dash past your opponents for high scores and big rewards!"
	tournament_category := 1
	tournament_startTime := int(time.Now().UTC().Unix()) // start now
	tournament_endTime := 0                              // never end, repeat the tournament each day forever
	tournament_duration := 3600                          // in seconds
	tournament_maxSize := 10000                          // first 10,000 players who join
	tournament_maxNumScore := 3                          // each player can have 3 attempts to score
	tournament_joinRequired := true

	// Create leaderboard.
	if err := nk.TournamentCreate(ctx, tournament_id, tournament_authoritative, tournament_sortOrder, tournament_operator, tournament_resetSchedule, tournament_metadata, tournament_title, tournament_description, tournament_category, tournament_startTime, tournament_endTime, tournament_duration, tournament_maxSize, tournament_maxNumScore, tournament_joinRequired); err != nil {
		logger.Error("error creating tournament")
	}

	// Register RPC functions.
	if err := initializer.RegisterRpc("account_delete_id", AccountDeleteId); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	if err := initializer.RegisterRpc("notification_send", NotificationSend); err != nil {
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

// Send a notification.
func NotificationSend(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	
	// First unmarshal into a RawMessage
	var raw json.RawMessage
	if err := json.Unmarshal([]byte(payload), &raw); err != nil {
		logger.Error("Failed first unmarshal: %v", err)
		return "", err
	}

	// Then unmarshal the RawMessage into our map
	var payloadMap map[string]interface{}
	if err := json.Unmarshal(raw, &payloadMap); err != nil {
		logger.Error("Failed second unmarshal: %v", err)
		return "", err
	}

	subject := payloadMap["subject"].(string)
	content := map[string]interface{}{
		"message": subject,
	}

	code := 1
	senderID := ""  
	persistent := true

	// Send the notification with all required parameters
 	if err := nk.NotificationSend(ctx, userID, subject, content, code, senderID, persistent); 
	err != nil {
		logger.Error("Failed to send notification: %v", err)
		return "", err
	}

	return "success", nil
}
