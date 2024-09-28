package main

import (
	"fmt"
	"sync"

	tele "gopkg.in/telebot.v3"
)

// Mutex to ensure safe access to shared state
	var mu sync.Mutex

	// Track selected plans for each user
	var userStates = make(map[int64]*UserState)

	var getUserState = func(userID int64) *UserState {
		mu.Lock()
		defer mu.Unlock()
		if state, exists := userStates[userID]; exists {
			return state
		}
		newState := &UserState{HasSelectedPlan: false}
		userStates[userID] = newState
		return newState
	}

	// range over plans 
	func rangePlans(plans []tele.Btn, ) {

		for _, p := range plans {
			plan := p
			b.Handle(&plan, func(c tele.Context) error {
				userState := getUserState(c.Sender().ID)
				userState.HasSelectedPlan = true
				userState.Referee = false
				userState.selectedPlan = plan.Text
				return c.Edit(fmt.Sprintf("پلن \"%s\" با موفقیت برای شما ثبت شد! \nلطفا برای ادامه نام خود را وارد کنید:", plan.Text), tele.ModeHTML)
			})
		}
	}