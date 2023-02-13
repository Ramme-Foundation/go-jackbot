package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"go-jackbot/prisma/db"
	"net/http"
	"strings"
)

type Controller struct {
	db *db.PrismaClient
}

func NewController(db *db.PrismaClient) *Controller {
	return &Controller{
		db: db,
	}
}

func (c *Controller) CommandRoute() func(c *gin.Context) {
	return func(ctx *gin.Context) {
		c.Command(ctx)
	}
}

func (c *Controller) Command(ctx *gin.Context) {
	s, err := slack.SlashCommandParse(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status:": "error", "message": err.Error()})
		return
	}
	commands := strings.Split(s.Text, " ")
	if len(commands) == 0 {
		c.help(ctx, s.Text)
		return
	}
	switch commands[0] {
	case SlackCommandHelp:
		c.help(ctx, strings.Join(commands[1:], " "))
		return
	case SlackCommandJackpot:
		c.jackpot(ctx)
		return
	case SlackCommandChoose:
		c.choose(ctx, &s)
		return
	default:
		c.help(ctx, strings.Join(commands[1:], " "))
		return
	}
}
