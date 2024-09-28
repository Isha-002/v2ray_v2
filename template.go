package main

import (tele "gopkg.in/telebot.v3")

type UserState struct {
	HasSelectedPlan bool
	selectedPlan    string
	Referee         bool
	RefereeName     string
	newUser      bool
	//
	Renew        bool
	username     string
	HasPanelName bool
	PanelName    string
	SendReceipt  bool
	Receipt      *tele.Photo
}

func resetStruct(u *UserState) {
	*u = UserState{}
}