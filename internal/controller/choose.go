package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"go-jackbot/internal/utils"
	"go-jackbot/prisma/db"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const MaxJackpotNumber = 50
const MaxPowerballNumber = 12
const MaxNumbers = 5
const MaxPowerballs = 2

var TotalMaxNumbers = MaxNumbers + MaxPowerballs

var PowerballTriggers = []string{"powerball", "pb", "power ball", "power-ball", "power ball", "bonus", "extra"}

func (c *Controller) choose(ctx *gin.Context, s *slack.SlashCommand) {
	if len(s.Text) == 0 {
		c.help(ctx, SlackCommandChoose)
		return
	}

	commands := strings.Split(s.Text, " ")
	number, numberType, err := parseNumber(strings.Join(commands[1:], " "))
	if err != nil {
		msg := &slack.Msg{
			Text: err.Error(),
		}
		ctx.JSON(http.StatusOK, msg)
		return
	}

	year, week := time.Now().ISOWeek()
	id := fmt.Sprintf("%s-%d-%d", s.TeamID, year, week)
	fmt.Println(id)
	_, err = c.db.JackpotRow.UpsertOne(
		db.JackpotRow.ID.Equals(id),
	).Create(
		db.JackpotRow.SLAckWorkspaceID.Set(s.TeamID),
		db.JackpotRow.Date.Set(utils.StartOfWeek(time.Now())),
		db.JackpotRow.ID.Set(id),
	).Update().Exec(ctx)
	if err != nil {
		msg := &slack.Msg{
			Text: err.Error(),
		}
		ctx.JSON(http.StatusOK, msg)
		return
	}

	row, err := c.db.JackpotRow.FindUnique(
		db.JackpotRow.ID.Equals(id),
	).With(db.JackpotRow.JackpotNumbers.Fetch()).Exec(ctx)
	if err != nil {
		msg := &slack.Msg{
			Text: err.Error(),
		}
		ctx.JSON(http.StatusOK, msg)
		return
	}

	numbers := row.JackpotNumbers()
	if len(numbers) > TotalMaxNumbers {
		msg := &slack.Msg{
			Text: fmt.Sprintf("You have already chosen %d numbers this week", len(numbers)),
		}
		ctx.JSON(http.StatusOK, msg)
		return
	}
	if countPowerballs(numbers) > MaxPowerballs {
		msg := &slack.Msg{
			Text: fmt.Sprintf("You have already chosen %d powerballs this week", countPowerballs(numbers)),
		}
		ctx.JSON(http.StatusOK, msg)
		return
	}
	if countNumbers(numbers) > MaxNumbers {
		msg := &slack.Msg{
			Text: fmt.Sprintf("You have already chosen %d numbers this week", countNumbers(numbers)),
		}
		ctx.JSON(http.StatusOK, msg)
		return
	}
	if userAlreadyANumber(numbers, s.UserID) {
		msg := &slack.Msg{
			Text: fmt.Sprintf("You have already chosen a number this week"),
		}
		ctx.JSON(http.StatusOK, msg)
		return
	}
	if numberAlreadyUsed(numbers, number) {
		msg := &slack.Msg{
			Text: fmt.Sprintf("That number has already been chosen"),
		}
		ctx.JSON(http.StatusOK, msg)
		return
	}

	_, err = c.db.JackpotNumber.CreateOne(
		db.JackpotNumber.SLAckUserID.Set(s.UserID),
		db.JackpotNumber.JackpotRow.Link(db.JackpotRow.ID.Equals(id)),
		db.JackpotNumber.NumberType.Set(numberType),
		db.JackpotNumber.Number.Set(number),
	).Exec(ctx)

	if err != nil {
		msg := &slack.Msg{
			Text: err.Error(),
		}
		ctx.JSON(http.StatusOK, msg)
		return
	}
	rowMsg, err := c.rowMessage(s.TeamID)
	if err != nil {
		msg := &slack.Msg{
			Text: err.Error(),
		}
		ctx.JSON(http.StatusOK, msg)
		return
	}
	msg := &slack.Msg{
		Text: rowMsg,
	}
	ctx.JSON(http.StatusOK, msg)
}

func userAlreadyANumber(numbers []db.JackpotNumberModel, userId string) bool {
	for _, number := range numbers {
		if number.SLAckUserID == userId {
			return true
		}
	}
	return false
}

func numberAlreadyUsed(numbers []db.JackpotNumberModel, number int) bool {
	for _, n := range numbers {
		if n.Number == number {
			return true
		}
	}
	return false
}

func countPowerballs(numbers []db.JackpotNumberModel) int {
	count := 0
	for _, number := range numbers {
		if number.NumberType == db.JackpotNumberTypePOWERBALL {
			count++
		}
	}
	return count
}

func countNumbers(numbers []db.JackpotNumberModel) int {
	count := 0
	for _, number := range numbers {
		if number.NumberType == db.JackpotNumberTypeNUMBER {
			count++
		}
	}
	return count
}

func parseNumber(text string) (int, db.JackpotNumberType, error) {
	for _, trigger := range PowerballTriggers {
		if strings.Contains(text, trigger) {
			return parsePowerballNumber(text)
		}
	}

	parseInt, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return 0, db.JackpotNumberTypeNUMBER, err
	}
	if parseInt < 0 {
		return 0, db.JackpotNumberTypeNUMBER, fmt.Errorf("number must be positive")
	}
	return int(parseInt), db.JackpotNumberTypeNUMBER, nil
}

func parsePowerballNumber(text string) (int, db.JackpotNumberType, error) {
	for _, trigger := range PowerballTriggers {
		text = strings.ReplaceAll(text, trigger, "")
	}
	parseInt, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return 0, db.JackpotNumberTypePOWERBALL, err
	}
	if parseInt < 0 {
		return 0, db.JackpotNumberTypePOWERBALL, fmt.Errorf("number must be positive")
	}
	if parseInt > MaxPowerballNumber {
		return 0, db.JackpotNumberTypePOWERBALL, fmt.Errorf("number must be less than %d", MaxPowerballNumber)
	}
	return int(parseInt), db.JackpotNumberTypePOWERBALL, nil
}
