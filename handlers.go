package main

import (
tele "gopkg.in/telebot.v3"
)

func startBot() error {
	b.Handle("/start", func(c tele.Context) error {
	userState := getUserState(c.Sender().ID)
	resetStruct(userState)
	return c.Send("یک گزینه را انتخاب کنید:", menu)
})
	return nil
}


b.Handle(&buyBtn, func(c tele.Context) error {
	userState := getUserState(c.Sender().ID)
	userState.newUser = true
	return c.Edit("لطفا سرویس مورد نظر خود را انتخاب کنید:", userSelection_Menu)
})