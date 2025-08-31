package friend

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
)

//	{
//	    "user_b": "b316d7c3-40d1-4a00-84b6-86a46acfcf8d"
//	}

func GetFriendState(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	// Parse the JSON payload
	var request struct {
		UserB string `json:"user_b"`
	}

	if err := json.Unmarshal([]byte(payload), &request); err != nil {
		logger.Error("Failed to unmarshal payload: %v", err)
		return "", errors.New("invalid JSON payload")
	}

	userIDToBeChecked := request.UserB

	// Basic validation
	if userIDToBeChecked == "" {
		return "", errors.New("user_b cannot be empty")
	}

	// Use QueryRow to get the result
	query := `
		SELECT state
		FROM user_edge
		WHERE (destination_id = $1 AND source_id = $2)
		LIMIT 1`

	var state int

	err := db.QueryRow(query, userID, userIDToBeChecked).Scan(&state)

	if err != nil {
		if err == sql.ErrNoRows {
			return "-1", nil // Not friends - return -1
		}
		logger.Error("Database query error: %v", err)
		return "", err
	}

	return fmt.Sprintf("%d", state), nil
}
