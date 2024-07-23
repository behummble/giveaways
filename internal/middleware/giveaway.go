package middleware

import (
	"net/http"
	"time"
	"strconv"

	"github.com/behummble/giveaway/internal/service/giveaway"
	"github.com/gin-gonic/gin"
)

const (
	apiKeyNotValid = "Api key not valid"
	badGiveaway = "Giveaway does not contain all the necessary fields"
	idNotProvided = "ID not provided"
	idNotValid = "ID not valid"
)

type GiveawayService interface {
	ClientExist(key string) bool
}

type giveaway struct {
	Description string
	EventDate time.Time
	PublishDate time.Time
	WinnersNumber int
	ParticipantNumber int
	Terms []string
}

func ApiKeyIsValid(giveaway *giveawayservice.Giveaway) gin.HandlerFunc {
	return func (c *gin.Context) {
		value := c.GetHeader("Api_Key")
		if value == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": apiKeyNotValid})
		}
		
		if !giveaway.ClientExist(value) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": apiKeyNotValid})
		}
	}
}

func GiveawayIsValid() gin.HandlerFunc {
	return func (c *gin.Context) {
		giveaway := giveaway{}
		if err := c.ShouldBindBodyWithJSON(&giveaway); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": badGiveaway})
		}
	}
}

func GiveawayIDIsValid() gin.HandlerFunc {
	return func (c *gin.Context) {
		id, found := c.Params.Get("id")
		if found {
			if id == "" {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": idNotProvided})
				return
			}
			_, err := strconv.Atoi(id)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": idNotValid})
				return
			}
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": idNotProvided})
		}
	}
}