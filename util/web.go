package util

import (
	"net/http"
	"strconv"
)

const (
	EvUserID = "Ev-UserID"
)

func GetEvUserID(req *http.Request) int {
	userId, err := strconv.Atoi(req.Header.Get(EvUserID))
	if err != nil {
		return 0
	}
	return userId
}
