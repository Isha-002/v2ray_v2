package main

import (
	"fmt"
	"time"

	tele "gopkg.in/telebot.v3"
)


const receiptChannelID = -1002457603510

func main() {

	b, err := tele.NewBot(tele.Settings{
		Token:  "7743565843:AAHjGGubY_yLb51cacm0T-gEMtHtHd-IYWQ",
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	var menu = &tele.ReplyMarkup{}

	// initial inline menu
	menu.Inline(
		menu.Row(buyBtn, renewBtn),
	)

	// user selection menu
	userSelection_Menu := &tele.ReplyMarkup{}
	userSelectSlice := []tele.Btn{oneUserPlans, twoUserPlans, unlimitedUserPlans, backtoMainBtn}
	userSelection_Menu.Inline(userSelection_Menu.Split(3, userSelectSlice)...)

	renewMenu := &tele.ReplyMarkup{}
	renewMenuSlice := []tele.Btn{renewPlan, renewAnotherPlan, backtoMainBtn}
	renewMenu.Inline(renewMenu.Split(2, renewMenuSlice)...)

	// plans menu
	plansMenu := &tele.ReplyMarkup{}

	b.Handle("/start", func(c tele.Context) error {
		userState := getUserState(c.Sender().ID)
		resetStruct(userState)
		return c.Send(chooseMenu, menu)
	})

	b.Handle(&buyBtn, func(c tele.Context) error {
		userState := getUserState(c.Sender().ID)
		userState.newUser = true
		return c.Edit(chooseService, userSelection_Menu)
	})

	b.Handle(&oneUserPlans, func(c tele.Context) error {

		plansWithBack := append(oneUser_plans, backBtn)
		plansMenu.Inline(plansMenu.Split(2, plansWithBack)...)

		return c.Edit(choosePlan, plansMenu)
	})

	b.Handle(&twoUserPlans, func(c tele.Context) error {
		plansWithBack := append(twoUser_plans, backBtn)
		plansMenu.Inline(plansMenu.Split(2, plansWithBack)...)
		return c.Edit(choosePlan, plansMenu)
	})

	b.Handle(&unlimitedUserPlans, func(c tele.Context) error {
		plansWithBack := append(unlimitedUser_plans, backBtn)
		plansMenu.Inline(plansMenu.Split(1, plansWithBack)...)
		return c.Edit(choosePlan, plansMenu)
	})

	b.Handle(&backBtn, func(c tele.Context) error {
		return c.Edit(chooseService, userSelection_Menu)
	})
	b.Handle(&backtoMainBtn, func(c tele.Context) error {
		return c.Edit(chooseMenu, menu)
	})

	// range over plans 
	rangePlans := func(plans []tele.Btn) {
    for _, p := range plans {
        plan := p
        b.Handle(&plan, func(c tele.Context) error {
            userState := getUserState(c.Sender().ID)
            userState.HasSelectedPlan = true
            userState.Referee = false
            userState.selectedPlan = plan.Text
						var planMsg string
						if userState.Renew {
							planMsg = reNewrangePlanMsg
						} else {
							planMsg = rangePlanMsg
						}
            return c.Edit(fmt.Sprintf(planMsg, plan.Text), tele.ModeHTML)
        })
    }
}

	rangePlans(oneUser_plans)
	rangePlans(twoUser_plans)
	rangePlans(unlimitedUser_plans)

	// Handle user text input
	b.Handle(tele.OnText, func(c tele.Context) error {
		userState := getUserState(c.Sender().ID)
		if !userState.Renew {
			if userState.newUser && userState.HasSelectedPlan && !userState.Referee {
				userState.Referee = true
				userState.RefereeName = c.Message().Text
				return c.Send(chooseReferee)
			}
			if userState.HasSelectedPlan && userState.Referee {
				userName := c.Message().Text
				user := c.Sender().Username
				plan := userState.selectedPlan
				referee := userState.RefereeName
				sendToChannel := fmt.Sprintf(newRequestMsg, userName, user, referee, plan)
				_, err := b.Send(tele.ChatID(receiptChannelID), sendToChannel)
				if err != nil {
					fmt.Println(err)
					return nil
				}

				return c.Send(successPurchase)
			}
			return c.Send(choosePlanError)
		}
		if userState.Renew {
			fmt.Println(userState.HasPanelName)
			fmt.Println(userState.username)
			if !userState.HasPanelName {
				userState.PanelName = c.Message().Text
				userState.HasPanelName = true
				return c.Send(sendReceiptMsg)
			}
			return c.Send(sendPanelName)
		}
		return c.Send(chooseServiceError)
	})

	b.Handle(tele.OnPhoto, func(c tele.Context) error {
		userState := getUserState(c.Sender().ID)
		if userState.HasPanelName {
		userState.Receipt = c.Message().Photo
		var plan string
		if userState.HasSelectedPlan {
			plan = userState.selectedPlan
		} else {
			plan = currentPlan
		}
		sendToChannel := fmt.Sprintf(renewRequestMsg, userState.PanelName, c.Sender().Username, plan)

		id , err := b.Send(tele.ChatID(receiptChannelID), userState.Receipt)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		b.EditCaption(id, sendToChannel)
		return c.Send(successRenew)}
		return c.Send(photoError)
	})

	// renew logic
	b.Handle(&renewBtn, func(c tele.Context) error {
		userState := getUserState(c.Sender().ID)
		userState.Renew = true
		return c.Edit(renewFirstMsg, renewMenu)
	})
	b.Handle(&renewPlan, func(c tele.Context) error {
		return c.Edit(sendPanelName)
	})
	b.Handle(&renewAnotherPlan, func(c tele.Context) error {
		return c.Edit(chooseMenu, userSelection_Menu)
	})

	// Start the bot
	fmt.Println(startMsg)
	b.Start()
}


