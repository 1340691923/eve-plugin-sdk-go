package util

import (
	"github.com/1340691923/eve-plugin-sdk-go/enum"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetEvUserID(ctx *gin.Context) int {
	userId, err := strconv.Atoi(ctx.GetHeader(enum.EvUserID))
	if err != nil {
		return 0
	}
	return userId
}
