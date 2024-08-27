package util

import (
	"net/http"
	"strconv"
)

const (
	EvRoleID = "Ev-RoleID"
	EvUserID = "Ev-UserID"
)

func GetEvUserID(req *http.Request) int {
	userId, err := strconv.Atoi(req.Header.Get(EvUserID))
	if err != nil {
		return 0
	}
	return userId
}

func GetEvRoleID(req *http.Request) int {
	roleId, err := strconv.Atoi(req.Header.Get(EvRoleID))
	if err != nil {
		return 0
	}
	return roleId
}
