package util

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

const (
	EvUserID = "Ev-UserID"
)

func GetEvUserID(ctx *gin.Context) int {
	userId, err := strconv.Atoi(ctx.GetHeader(EvUserID))
	if err != nil {
		return 0
	}
	return userId
}
