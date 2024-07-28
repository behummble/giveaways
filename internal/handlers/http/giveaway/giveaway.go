package giveaway

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/behummble/giveaway/internal/middleware"
	"github.com/behummble/giveaway/internal/service/giveaway"
	"github.com/gin-gonic/gin"
)

const (
	newGiveaway = "/giveaways"
	//giveaway = "/giveaways/:id"
	//winners = "/giveaways/:id/winners"
	giveaway = "/:id"
	winners = "/:id/winners"
)

/*type GiveawayService interface {
	AddGiveaway()
	Giveaway()
	UpdateGiveaway()
	DeleteGiveaway()
	Winners()
	Client()
} */

type Router struct {
	log *slog.Logger
	mux *gin.Engine
	giveaway *giveawayservice.Giveaway
}

func New(log *slog.Logger, giveaway *giveawayservice.Giveaway) *Router {
	m := gin.Default()
	return &Router{
		log: log,
		mux: m,
		giveaway: giveaway,
	}
}

func (r *Router) Register() {
	giveawaysUpdate := r.mux.Group(newGiveaway)
	giveawaysUpdate.Use(middleware.GiveawayIsValid()) 
	{
		giveawaysUpdate.POST("", r.newGiveaway)
		existingGiveaway := giveawaysUpdate.Group(giveaway)
		existingGiveaway.Use(middleware.GiveawayIDIsValid())
		existingGiveaway.PUT("", r.updateGiveaway)
	}

	giveaways := r.mux.Group(fmt.Sprintf("%s/%s", newGiveaway, giveaway))
	giveaways.Use(middleware.GiveawayIDIsValid())
	{
		giveaways.GET("", r.giveawayInfo)
		giveaways.DELETE("", r.deleteGiveaway)
		giveaways.GET(winners, r.winnersInfo)
	} 
	
	r.mux.POST("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "%s", "pong")
	})

	r.registerMiddleware()
}

func (r *Router) Serve(host string, port int) {
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), r.mux)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func (r *Router) registerMiddleware() {
	r.mux.Use(middleware.ApiKeyIsValid(r.giveaway))
}

func (r *Router) newGiveaway(context *gin.Context) {
	body, err := context.GetRawData()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp, err := r.giveaway.AddGiveaway(body)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error":err})
	} else {
		context.String(http.StatusCreated, string(resp))
	}

}

func (r *Router) giveawayInfo(context *gin.Context) {
	idStr, _ := context.Params.Get("id")
	id, _ := strconv.Atoi(idStr)
	resp, err := r.giveaway.Giveaway(id)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error":err})
	} else {
		context.String(http.StatusOK, string(resp))
	}
}

func (r *Router) updateGiveaway(context *gin.Context) {
	body, err := context.GetRawData()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	idStr, _ := context.Params.Get("id")
	id, _ := strconv.Atoi(idStr)

	resp, err := r.giveaway.UpdateGiveaway(id, body)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error":err})
	} else {
		context.String(http.StatusOK, string(resp))
	}
}

func (r *Router) deleteGiveaway(context *gin.Context) {
	idStr, _ := context.Params.Get("id")
	id, _ := strconv.Atoi(idStr)
	err := r.giveaway.DeleteGiveaway(id)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error":err})
	} else {
		context.Status(http.StatusOK)
	}
}

func (r *Router) winnersInfo(context *gin.Context) {
	idStr, _ := context.Params.Get("id")
	id, _ := strconv.Atoi(idStr)
	
	resp, err := r.giveaway.Winners(id)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error":err})
	} else {
		context.String(http.StatusOK, string(resp))
	}
}