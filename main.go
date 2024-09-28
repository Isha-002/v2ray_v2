package main

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)


const receiptChannelID = -1002457603510

func main() {
	err := InitBot()
	if err != nil {
		fmt.Println(err)
		return
	}

  err = startBot()
	if err != nil {
		fmt.Println(err)
		return
	}

	// initial inline menu
	menu := &tele.ReplyMarkup{}
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

	b.Handle(&oneUserPlans, func(c tele.Context) error {

		plansWithBack := append(oneUser_plans, backBtn)
		plansMenu.Inline(plansMenu.Split(2, plansWithBack)...)

		return c.Edit("لطفا پلن مورد نظر خود را انتخاب کنید:", plansMenu)
	})

	b.Handle(&twoUserPlans, func(c tele.Context) error {
		plansWithBack := append(twoUser_plans, backBtn)
		plansMenu.Inline(plansMenu.Split(2, plansWithBack)...)
		return c.Edit("لطفا پلن مورد نظر خود را انتخاب کنید:", plansMenu)
	})

	b.Handle(&unlimitedUserPlans, func(c tele.Context) error {
		plansWithBack := append(unlimitedUser_plans, backBtn)
		plansMenu.Inline(plansMenu.Split(1, plansWithBack)...)
		return c.Edit("لطفا پلن مورد نظر خود را انتخاب کنید:", plansMenu)
	})

	b.Handle(&backBtn, func(c tele.Context) error {
		return c.Edit("لطفا سرویس مورد نظر خود را انتخاب کنید:", userSelection_Menu)
	})
	b.Handle(&backtoMainBtn, func(c tele.Context) error {
		return c.Edit("یک گزینه را انتخاب کنید:", menu)
	})

	// Handle plans

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
				return c.Send("لطفا نام معرف خود را وارد کنید: ")
			}
			if userState.HasSelectedPlan && userState.Referee {
				userName := c.Message().Text
				user := c.Sender().Username
				plan := userState.selectedPlan
				referee := userState.RefereeName
				sendToChannel := fmt.Sprintf("درخواست جدید:\n\nنام کاربری: %s\nآیدی: @%s\nمعرف: %s\nپلن: %s\nㅤ", userName, user, referee, plan)
				_, err := b.Send(tele.ChatID(receiptChannelID), sendToChannel)
				if err != nil {
					fmt.Println(err)
					return nil
				}

				return c.Send("درخواست شما با موفقیت ثبت گردید. \nدر اسرع وقت توسط پشتیبانی با شما تماس گرفته میشود.\n ㅤ")
			}
			return c.Send("لطفا ابتدا یک پلن انتخاب کنید.")
		}
		if userState.Renew {
			userState.username = c.Message().Text
			if !userState.HasPanelName {
				userState.PanelName = c.Message().Text
				userState.HasPanelName = true
				return c.Send("لطفاعکس رسید پرداختی خود را ارسال نمایید:")
			}
			return c.Send("لطفا نام پنل کاربری خود را وارد کنید:")
		}
		return c.Send("برای شروع لطفا از منو گزینه مورد نظر خود را انتخاب کنید!")
	})

	b.Handle(tele.OnPhoto, func(c tele.Context) error {
		userState := getUserState(c.Sender().ID)
		if userState.HasPanelName {
		userState.Receipt = c.Message().Photo
		var plan string
		if userState.HasSelectedPlan {
			plan = userState.selectedPlan
		} else {
			plan = "تمدید فعلی"
		}
		sendToChannel := fmt.Sprintf("درخواست تمدید \n\nنام کاربر: %s\nآیدی: @%s\nپلن: %s\nㅤ", userState.username, c.Sender().Username, plan)

		id , err := b.Send(tele.ChatID(receiptChannelID), userState.Receipt)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		b.EditCaption(id, sendToChannel)
		b.Send(tele.ChatID(receiptChannelID), fmt.Sprintf("user: %+v", userState))
		return c.Send("درخواست شما با موفقیت ثبت گردید و پنل کاربری شما به زودی شارژ میشود!")}
		return c.Send("Invalid response")
	})

	// renew logic
	b.Handle(&renewBtn, func(c tele.Context) error {
		userState := getUserState(c.Sender().ID)
		userState.Renew = true
		return c.Edit("برای تمدید کردن سرویس خود یکی از گزینه های زیر را انتخاب کنید:", renewMenu)
	})
	b.Handle(&renewPlan, func(c tele.Context) error {
		return c.Edit("لطفا نام پنل کاربری خود را وارد کنید:")
	})
	b.Handle(&renewAnotherPlan, func(c tele.Context) error {
		return c.Edit("یک گزینه را انتخاب کنید:", userSelection_Menu)
	})

	// Start the bot
	fmt.Println("application started")
	b.Start()
}
