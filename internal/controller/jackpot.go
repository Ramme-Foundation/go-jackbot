package controller

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"github.com/ramme-foundation/go-jackbot/internal/models"
	"github.com/slack-go/slack"
	"io"
	"net/http"
	"strconv"
)

const JackpotResultUrl = "https://www.lottoland.com/api/drawings/euroJackpot"
const EuroExchangeRateUrl = "https://cdn.jsdelivr.net/gh/fawazahmed0/currency-api@1/latest/currencies/eur.json"

func (c *Controller) jackpot(ctx *gin.Context) {
	result, err := euroJackpotResult()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status:": "error", "message": err.Error()})
		return
	}
	msgText, err := displayJackpot(result)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"status:": "error", "message": err.Error()})
		return
	}
	msg := &slack.Msg{
		Text: msgText,
	}
	ctx.JSON(http.StatusOK, msg)
}

func euroJackpotResult() (*models.EuroJackpotResponse, error) {
	get, err := http.Get(JackpotResultUrl)
	if err != nil {
		return nil, err
	}
	defer get.Body.Close()
	stringData, err := io.ReadAll(get.Body)
	if err != nil {
		return nil, err
	}
	var result models.EuroJackpotResponse
	err = json.Unmarshal(stringData, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func displayJackpot(result *models.EuroJackpotResponse) (string, error) {
	jackpot, err := jackpotString(result, 1)
	if err != nil {
		return "", err
	}
	exchangeRate, err := getExchangeRate()
	if err != nil {
		return "", err
	}
	sweJackpot, err := jackpotString(result, exchangeRate)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("This weeks jackpot is €%s which is %s SEK (using EUR/SEK exchange rate of %.2f)", jackpot, sweJackpot, exchangeRate), nil
}

func jackpotString(result *models.EuroJackpotResponse, multiplier float64) (string, error) {
	jackpotInt, err := strconv.ParseFloat(result.Next.Jackpot, 64)
	if err != nil {
		return "", err
	}
	jackpot := jackpotInt * 1000000 * multiplier
	return humanize.Commaf(jackpot), nil
}

func getExchangeRate() (float64, error) {
	get, err := http.Get(EuroExchangeRateUrl)
	if err != nil {
		return 0, err
	}
	defer get.Body.Close()
	stringData, err := io.ReadAll(get.Body)
	if err != nil {
		return 0, err
	}
	var result models.EuroExchangeRateResponse
	err = json.Unmarshal(stringData, &result)
	if err != nil {
		return 0, err
	}
	return result.Eur.Sek, nil
}
