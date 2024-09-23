package models

import "time"

type UserSession struct {
	ID                  int        `json:"id"`
	UserID              int        `json:"user_id"`
	LoginAt             time.Time  `json:"login_at"`
	LogoutAt            *time.Time `json:"logout_at"`
	InvalidLogoutReason string     `json:"invalid_logout_reason"`
	CrashReasonType     *int       `json:"crash_reason_type"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}

const KInvalidLogoutReasonUndefined string = "undefined"

type ECrashReason int

const (
	KCrashReasonUndefined ECrashReason = iota // 0
	KCrashReasonSoftware                      // 1
	KCrashReasonSystem                        // 2
)

type UserSessionParams struct {
	PaginateParams
	UserID                      int  `json:"user_id"`
	OnlyUnmarkedInvalidSessions bool `json:"only_unmarked_invalid_sessions"`
}
