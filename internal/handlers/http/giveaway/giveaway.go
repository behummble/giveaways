package giveaway

import (
	"net/http"
	"log/slog"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/behummble/giveaway/internal/service/giveaway"
)

const (
	newGiveaway = "/giveaways"
	giveaway = "/giveaways/:id"
	winners = "/giveaways/:id/winners"
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
	r.mux.POST(newGiveaway, r.newGiveaway)
	r.mux.GET(giveaway, r.giveawayInfo)
	r.mux.PUT(giveaway, r.updateGiveaway)
	r.mux.DELETE(giveaway, r.deleteGiveaway)
	r.mux.GET(winners, r.winnersInfo)

	r.mux.POST("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "%s", "pong")
	})
}

func (r *Router) Serve(host string, port int) {
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), r.mux)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func (r *Router) newGiveaway(context *gin.Context) {
	
}

func (r *Router) giveawayInfo(context *gin.Context) {
	
}

func (r *Router) updateGiveaway(context *gin.Context) {
	
}

func (r *Router) deleteGiveaway(context *gin.Context) {
	
}

func (r *Router) winnersInfo(context *gin.Context) {
	
}