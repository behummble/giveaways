package middleware

import (
	//"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/behummble/giveaway/internal/service/giveaway"

)

type GiveawayService interface {
	ClientExist(key string) bool
}

func ApiKeyIsValid(c *gin.Context, giveaway *giveawayservice.Giveaway) (bool, int) {
	value := c.GetHeader("Api_Key")
	if value == "" {
		return false, http.StatusUnauthorized
	}
	
	if giveaway.ClientExist(value) {
		return true, http.StatusOK
	}

	return false, http.StatusUnauthorized
}