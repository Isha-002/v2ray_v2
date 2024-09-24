package main

import (
	"fmt"
	"time"

	tele "gopkg.in/telebot.v3"
)

func main() {
	b, err := tele.NewBot(tele.Settings{
		Token:  "7743565843:AAHjGGubY_yLb51cacm0T-gEMtHtHd-IYWQ",
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create inline buttons
	buyBtn := tele.Btn{Unique: "buy_service", Text: "خرید سرویس"}
	renewBtn := tele.Btn{Unique: "renew_service", Text: "تمدید سرویس"}
	oneUserPlans := tele.Btn{Unique: "one_user_plans", Text: "تک کاربره"}
	twoUserPlans := tele.Btn{Unique: "two_user_plan", Text: "دو کاربره"}
	unlimitedUserPlans := tele.Btn{Unique: "unlimited_user_plan", Text: "کاربر نامحدود"}

	// one user plans
	oneUser_plans := []tele.Btn{
		{Unique: "oneuser_plan1", Text: "۴۰ گیگ ۱ ماهه: ۷۵ تومن"},
		{Unique: "oneuser_plan1", Text: "۴۰ گیگ ۱ ماهه: ۷۵ تومن"},
		{Unique: "oneuser_plan3", Text: "۷۵ گیگ ۱ ماهه: ۱۰۰ تومن"},
		{Unique: "oneuser_plan4", Text: "۱۰۰ گیگ ۱ ماهه: ۱۲۰ تومن"},
	}


	// initial inline keyboard
	menu := &tele.ReplyMarkup{}
	menu.Inline(
		menu.Row(buyBtn, renewBtn),
	)


	b.Handle("/start", func(c tele.Context) error {
		return c.Send("یک گزینه را انتخاب کنید:", menu)
	})

	b.Handle(&buyBtn, func(c tele.Context) error {

		user_Menu := &tele.ReplyMarkup{}
		user_Menu.Inline(user_Menu.Row(oneUserPlans, twoUserPlans, unlimitedUserPlans))

		return c.Edit("لطفا سرویس مورد نظر خود را انتخاب کنید:", user_Menu)
	})

	b.Handle(&oneUserPlans, func(c tele.Context) error{

		plansMenu := &tele.ReplyMarkup{}
		plansMenu.Inline(plansMenu.Split(2, oneUser_plans)...)

		return c.Edit("لطفا پلن مورد نظر خود را انتخاب کنید:", plansMenu)
	})










	// renew logic
	b.Handle(&renewBtn, func(c tele.Context) error {
		return c.Edit("برای تمدید کردن سرویس فعلی، نام پنل کاربری خود را وارد کنید:")
	})

	// Start the bot
	fmt.Println("application started")
	b.Start()
}
