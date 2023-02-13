package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"net/http"
	"strings"
)

var allCommands = []SlackCommand{
	SlackCommandHelp,
	SlackCommandJackpot,
	SlackCommandChoose,
	SlackCommandRandom,
}

func (c *Controller) help(ctx *gin.Context, command SlackCommand) {
	msg := &slack.Msg{
		Text: fmt.Sprintf("Unknown command '%s', Available commands are %s", command, strings.Join(allCommands, " | ")),
	}
	if command == SlackCommandHelp {
		msg.Text = fmt.Sprintf("Available commands are %s", strings.Join(allCommands, " | "))
	}
	ctx.JSON(http.StatusOK, msg)
	return
}
