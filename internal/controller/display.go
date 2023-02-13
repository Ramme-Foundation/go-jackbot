package controller

import (
	"context"
	"fmt"
	"github.com/ramme-foundation/go-jackbot/internal/utils"
	"github.com/ramme-foundation/go-jackbot/prisma/db"
	"sort"
	"time"
)

func (c *Controller) rowMessage(teamID string) (string, error) {
	row, err := c.db.JackpotRow.FindFirst(
		db.JackpotRow.SLAckWorkspaceID.Equals(teamID),
		db.JackpotRow.Date.Equals(utils.StartOfWeek(time.Now())),
	).With(db.JackpotRow.JackpotNumbers.Fetch()).Exec(context.Background())
	if err != nil {
		return "", err
	}
	SortJackpotNumbers(row.JackpotNumbers())

	displayRow := ""

	for i := 0; i < TotalMaxNumbers; i++ {
		if i == MaxNumbers {
			displayRow += " | "
		}
		if i < len(row.JackpotNumbers()) {
			displayRow += fmt.Sprintf("%d ", row.JackpotNumbers()[i].Number)
			if i != TotalMaxNumbers-1 {
				displayRow += " -"
			}
		} else {
			displayRow += "_"
			if i != MaxNumbers-1 {
				displayRow += " -"
			}
		}
	}
	return displayRow, err
}

func SortJackpotNumbers(numbers []db.JackpotNumberModel) {
	sort.Slice(numbers, func(i, j int) bool {
		if numbers[i].NumberType == numbers[j].NumberType {
			return numbers[i].NumberType > numbers[j].NumberType
		}
		return numbers[i].Number < numbers[j].Number
	})

}
