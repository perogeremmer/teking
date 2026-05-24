package common

import (
	"net/http"
)

type contextKey string

const (
	ContextKeyOperatorID contextKey = "operator_id"
	ContextKeyRole       contextKey = "role"
)

const CookieSessionName = "session_id"

const (
	RoleSuperadmin = "superadmin"
	RoleAdmin      = "admin"
	RoleUser       = "user"
)

func GetOperatorID(r *http.Request) int64 {
	id, ok := r.Context().Value(ContextKeyOperatorID).(int64)
	if !ok {
		return 0
	}
	return id
}

func GetRole(r *http.Request) string {
	role, ok := r.Context().Value(ContextKeyRole).(string)
	if !ok {
		return ""
	}
	return role
}
