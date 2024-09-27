package main

import (
	"fmt"
	"sync"
	"time"

	"gopkg.in/telebot.v3"
	tele "gopkg.in/telebot.v3"
)

type UserState struct {
	HasSelectedPlan bool
	selectedPlan string
	Referee bool
	RefereeName string
}

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

	// Create inline buttons
	buyBtn := tele.Btn{Unique: "buy_service", Text: "خرید سرویس"}
	renewBtn := tele.Btn{Unique: "renew_service", Text: "تمدید سرویس"}
	oneUserPlans := tele.Btn{Unique: "one_user_plans", Text: "تک کاربره"}
	twoUserPlans := tele.Btn{Unique: "two_user_plan", Text: "دو کاربره"}
	unlimitedUserPlans := tele.Btn{Unique: "unlimited_user_plan", Text: "کاربر نامحدود"}

	// back buttons
	backBtn := tele.Btn{Unique: "back_btn", Text: "↩️ بازگشت"}
	backtoMainBtn := tele.Btn{Unique: "backToMain_btn", Text: "↩️ بازگشت"}

	// one user plans
	oneUser_plans := []tele.Btn{
		{Unique: "oneuser_plan1", Text: "۴۰ گیگ ۱ ماهه: ۷۵ تومن"},
		{Unique: "oneuser_plan2", Text: "۶۰ گیگ ۱ ماهه: ۹۰ تومن"},
		{Unique: "oneuser_plan3", Text: "۷۵ گیگ ۱ ماهه: ۱۰۰ تومن"},
		{Unique: "oneuser_plan4", Text: "۱۰۰ گیگ ۱ ماهه: ۱۲۰ تومن"},
	}

	// two user plans
	twoUser_plans := []tele.Btn{
		{Unique: "twouser_plan1", Text: "۷۰ گیگ ۱ ماهه ۱۲۰ تومن"},
		{Unique: "twouser_plan1", Text: "۹۰ گیگ ۱ ماهه ۱۴۰ تومن"},
		{Unique: "twouser_plan3", Text: "۱۲۰ گیگ ۱ ماهه ۱۶۰ تومن"},
		{Unique: "twouser_plan4", Text: "۲۰۰ گیگ ۱ ماهه ۲۲۰ تومن"},
	}

	// unlimited user plans
	unlimitedUser_plans := []tele.Btn{
		{Unique: "unlimitedUser_plan1", Text: "۱ ماهه ۱۵۰ گیگ ۲۵۰ تومن"},
		{Unique: "unlimitedUser_plan2", Text: "۱ ماهه ۲۵۰ گیگ ۳۱۵ تومن"},
		{Unique: "unlimitedUser_plan3", Text: "۱ ماهه ۳۵۰ گیگ  ۴۰۰ تومن"},
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

	// plans menu
	plansMenu := &tele.ReplyMarkup{}

	// Mutex to ensure safe access to shared state
	var mu sync.Mutex

	// Track selected plans for each user
	userStates := make(map[int64]*UserState)

	getUserState := func(userID int64) *UserState {
		mu.Lock()
		defer mu.Unlock()
		if state, exists := userStates[userID]; exists {
			return state
		}
		newState := &UserState{HasSelectedPlan: false}
		userStates[userID] = newState
		return newState
	}

	b.Handle("/start", func(c tele.Context) error {
		userState := getUserState(c.Sender().ID)
		userState.HasSelectedPlan = false
		userState.selectedPlan = ""
		return c.Send("یک گزینه را انتخاب کنید:", menu)
	})

	b.Handle(&buyBtn, func(c tele.Context) error {
		return c.Edit("لطفا سرویس مورد نظر خود را انتخاب کنید:", userSelection_Menu)
	})

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

	// Handle plan selection for one-user plans
	for _, p := range oneUser_plans {
		plan := p 
		b.Handle(&plan, func(c tele.Context) error {
			userState := getUserState(c.Sender().ID)
			userState.HasSelectedPlan = true
			userState.selectedPlan = plan.Text
			return c.Edit(fmt.Sprintf("پلن \"%s\" با موفقیت برای شما ثبت شد! \nلطفا برای ادامه نام خود را وارد کنید:", plan.Text), tele.ModeHTML)
		})
	}

	// Handle plan selection for two-user plans
	for _, p := range twoUser_plans {
		plan := p
		b.Handle(&plan, func(c tele.Context) error {
			userState := getUserState(c.Sender().ID)
			userState.HasSelectedPlan = true
			userState.Referee = false
			userState.selectedPlan = plan.Text
			return c.Edit(fmt.Sprintf("پلن \"%s\" با موفقیت برای شما ثبت شد! \nلطفا برای ادامه نام خود را وارد کنید:", plan.Text), tele.ModeHTML)
		})
	}

	// Handle plan selection for unlimited-user plans
	for _, p := range unlimitedUser_plans {
		plan := p
		b.Handle(&plan, func(c tele.Context) error {
			userState := getUserState(c.Sender().ID)
			userState.HasSelectedPlan = true
			userState.Referee = false
			userState.selectedPlan = plan.Text
			return c.Edit(fmt.Sprintf("پلن \"%s\" با موفقیت برای شما ثبت شد! \nلطفا برای ادامه نام خود را وارد کنید:", plan.Text), tele.ModeHTML)

		})
	}


	// Handle user text input only if they selected a plan
	b.Handle(telebot.OnText, func(c tele.Context) error {
		userState := getUserState(c.Sender().ID)
		if userState.HasSelectedPlan && !userState.Referee {
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
			_, err = b.Send(tele.ChatID(receiptChannelID), sendToChannel)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			return c.Send("درخواست شما با موفقیت ثبت گردید. \nدر اسرع وقت توسط ادمین با شما تماس گرفته میشود.\n ㅤ")
		}
		return c.Send("لطفا ابتدا یک پلن انتخاب کنید.")
	})

	// renew logic
	b.Handle(&renewBtn, func(c tele.Context) error {
		return c.Edit("برای تمدید کردن سرویس فعلی، نام پنل کاربری خود را وارد کنید:")
	})

	// Start the bot
	fmt.Println("application started")
	b.Start()
}
