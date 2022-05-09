package main

import (
	"fmt"
	"strconv"
	"time"

	log "github.com/iliesh/go-templates/logger"
)

func newGamePermission(uid int64) (bool, error) {
	// TODO: Check if user is allowed to create new game
	return true, nil
}

func selectDate(id int64) error {
	currTime := time.Now()

	monthName := inlineKeyboardButtonT{Text: currTime.Month().String() + " " + fmt.Sprintf("%d", currTime.Year()), CallBackData: "month"}
	monthButton := []inlineKeyboardButtonT{monthName}

	weekDayMo := inlineKeyboardButtonT{Text: "Mo", CallBackData: "monday"}
	weekDayTu := inlineKeyboardButtonT{Text: "Tu", CallBackData: "tuesday"}
	weekDayWe := inlineKeyboardButtonT{Text: "We", CallBackData: "wednesday"}
	weekDayTh := inlineKeyboardButtonT{Text: "Th", CallBackData: "thursday"}
	weekDayFr := inlineKeyboardButtonT{Text: "Fr", CallBackData: "friday"}
	weekDaySa := inlineKeyboardButtonT{Text: "Sa", CallBackData: "saturday"}
	weekDaySu := inlineKeyboardButtonT{Text: "Su", CallBackData: "sunday"}
	weekDaysButton := []inlineKeyboardButtonT{weekDayMo, weekDayTu, weekDayWe, weekDayTh, weekDayFr, weekDaySa, weekDaySu}
	firstLineDate := []inlineKeyboardButtonT{}
	secondLineDate := []inlineKeyboardButtonT{}

	weekDayNumber := currTime.Weekday()
	if weekDayNumber == 0 {
		weekDayNumber = 7
	}
	for i := 1; i < 15; i++ {
		if i < int(weekDayNumber) {
			firstLineDate = append(firstLineDate, inlineKeyboardButtonT{Text: "\u2718", CallBackData: "Invalid Date"})
			continue
		}
		callBackDate := currTime.AddDate(0, 0, i-int(weekDayNumber))
		_, _, date := callBackDate.Date()
		if i < 8 {
			firstLineDate = append(firstLineDate, inlineKeyboardButtonT{Text: strconv.Itoa(date), CallBackData: callBackDate.Format("2006-01-02")})
			continue
		}
		secondLineDate = append(secondLineDate, inlineKeyboardButtonT{Text: strconv.Itoa(date), CallBackData: callBackDate.Format("2006-01-02")})
	}

	b := [][]inlineKeyboardButtonT{monthButton, weekDaysButton, firstLineDate, secondLineDate}
	c := inlineKeyboardMarkupT{b}
	m := sendMessageReqBodyT{ChatID: id, Text: "Please select the date:", ReplyMarkup: c}
	log.Info("Sending Message: %v", m)
	err := sendMessage(m)
	if err != nil {
		return err
	}
	return nil
}
