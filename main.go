package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
)


func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(envError)
	}

	receiptChannelIDstr := os.Getenv("RECEIPT_CHANNEL_ID")

	receiptChannelID, err := strconv.ParseInt(receiptChannelIDstr, 10, 64)
	if err != nil {
		fmt.Println(typeConvErr, err)
		return
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")

	b, err := tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	var menu = &tele.ReplyMarkup{}
	var returnBtn = &tele.ReplyMarkup{}
	var helpOptions = &tele.ReplyMarkup{}

	// initial inline menu
	menu.Inline(
		menu.Row(buyBtn, renewBtn),
	)
	returnBtn.Inline(
		returnBtn.Row(restartBtn, helpBtn),
	)
	helpOptions.Inline(
		helpOptions.Row(androidHelp, iosHelp),
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
	b.Handle("/restart", func(c tele.Context) error {
		userState := getUserState(c.Sender().ID)
		resetStruct(userState)
		return c.Send(chooseMenu, menu)
	})
	b.Handle(&restartBtn, func(c tele.Context) error {
		userState := getUserState(c.Sender().ID)
		resetStruct(userState)
		return c.Send(chooseMenu, menu)
	})

	b.Handle(&helpBtn, func(c tele.Context) error {
		return c.Edit(helpMsg, helpOptions)
	})
	b.Handle("/help", func(c tele.Context) error {

		return c.Send(helpMsg, helpOptions)
	})

	b.Handle(&androidHelp, func(c tele.Context) error {
		originalMessageID := 56
		sourceChatID := int64(receiptChannelID)
		_, err := b.Copy(tele.ChatID(c.Chat().ID), &tele.Message{
			ID: originalMessageID,
			Chat: &tele.Chat{
				ID: sourceChatID,
			},
		})

		if err != nil {
			return c.Send(errForward)
		}
		return nil
	})

	b.Handle(&iosHelp, func(c tele.Context) error {
		originalMessageID := 57
		videoMessageID := 87
		sourceChatID := int64(receiptChannelID)
		_, err := b.Copy(tele.ChatID(c.Chat().ID), &tele.Message{
			ID: originalMessageID,
			Chat: &tele.Chat{
				ID: sourceChatID,
			},
		})

		if err != nil {
			return c.Send(errForward)
		}

		_, err = b.Copy(tele.ChatID(c.Chat().ID), &tele.Message{
			ID: videoMessageID,
			Chat: &tele.Chat{
				ID: sourceChatID,
			},
		})

		if err != nil {
			return c.Send(errForward)
		}
		return nil
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
		if userState.done {
			return c.Send(doneChat, returnBtn)
		}
		if !userState.Renew {
			if userState.newUser && userState.HasSelectedPlan && !userState.Referee && !userState.hasphoneNumber {
				userState.Referee = true
				userState.username = c.Message().Text
				return c.Send(chooseReferee)
			}
			if userState.HasSelectedPlan && userState.Referee && !userState.hasphoneNumber {
				userState.RefereeName = c.Message().Text
			}

			user := c.Sender().Username
			if user != "" {
				sendToChannel := fmt.Sprintf(newRequestMsg, userState.username, "@"+user, userState.RefereeName, userState.selectedPlan)
				_, err := b.Send(tele.ChatID(receiptChannelID), sendToChannel)
				if err != nil {
					fmt.Println(err)
					return nil
				}
				userState.done = true
				return c.Send(successPurchase, returnBtn)
			}

			if user == "" && !userState.hasphoneNumber {
				userState.hasphoneNumber = true
				return c.Send(askPhoneNumber)
			}
			if userState.hasphoneNumber {
				noUserId := noId
				userState.phoneNumber = c.Message().Text
				sendToChannel := fmt.Sprintf(newRequestMsgWithPhone, userState.username, noUserId, userState.phoneNumber, userState.RefereeName, userState.selectedPlan)
				_, err := b.Send(tele.ChatID(receiptChannelID), sendToChannel)
				if err != nil {
					fmt.Println(err)
					return nil
				}
				userState.done = true
				return c.Send(successPurchase, returnBtn)
			}
			return c.Send(choosePlanError)
		}
		if userState.Renew {
			userId := c.Sender().Username
			if !userState.HasPanelName && userId != "" {
				userState.PanelName = c.Message().Text
				userState.HasPanelName = true
				return c.Send(sendReceiptMsg)
			}
			if userId == "" && !userState.hasphoneNumber {
				userState.HasPanelName = true
				userState.PanelName = c.Message().Text
				userState.hasphoneNumber = true
				return c.Send(askPhoneNumber)
			}
			if userState.hasphoneNumber {
				userState.phoneNumber = c.Message().Text
				return c.Send(sendReceiptMsg)
			}

		}
		return c.Send(chooseServiceError)
	})

	b.Handle(tele.OnPhoto, func(c tele.Context) error {
		userState := getUserState(c.Sender().ID)
		userId := c.Sender().Username
		if userState.HasPanelName {
			userState.Receipt = c.Message().Photo
			var plan string
			if userState.HasSelectedPlan {
				plan = userState.selectedPlan
			} else {
				plan = currentPlan
			}
			if userId == "" {
				userId = userState.phoneNumber
			} else {
				userId = "@" + c.Sender().Username
			}
			sendToChannel := fmt.Sprintf(renewRequestMsg, userState.PanelName, userId, plan)

			id, err := b.Send(tele.ChatID(receiptChannelID), userState.Receipt)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			b.EditCaption(id, sendToChannel)
			userState.done = true
			return c.Send(successRenew, returnBtn)
		}
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
