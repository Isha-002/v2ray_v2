package main

import (
		tele "gopkg.in/telebot.v3"
)

	// one user plans
	var oneUser_plans = []tele.Btn{
		{Unique: "oneuser_plan1", Text: "۴۰ گیگ ۱ ماهه: ۷۵ تومن"},
		{Unique: "oneuser_plan2", Text: "۶۰ گیگ ۱ ماهه: ۹۰ تومن"},
		{Unique: "oneuser_plan3", Text: "۷۵ گیگ ۱ ماهه: ۱۰۰ تومن"},
		{Unique: "oneuser_plan4", Text: "۱۰۰ گیگ ۱ ماهه: ۱۲۰ تومن"},
	}

	// two user plans
	var twoUser_plans = []tele.Btn{
		{Unique: "twouser_plan1", Text: "۷۰ گیگ ۱ ماهه ۱۲۰ تومن"},
		{Unique: "twouser_plan2", Text: "۹۰ گیگ ۱ ماهه ۱۴۰ تومن"},
		{Unique: "twouser_plan3", Text: "۱۲۰ گیگ ۱ ماهه ۱۶۰ تومن"},
		{Unique: "twouser_plan4", Text: "۲۰۰ گیگ ۱ ماهه ۲۲۰ تومن"},
	}

	// unlimited user plans
	var unlimitedUser_plans = []tele.Btn{
		{Unique: "unlimitedUser_plan1", Text: "۱ ماهه ۱۵۰ گیگ ۲۵۰ تومن"},
		{Unique: "unlimitedUser_plan2", Text: "۱ ماهه ۲۵۰ گیگ ۳۱۵ تومن"},
		{Unique: "unlimitedUser_plan3", Text: "۱ ماهه ۳۵۰ گیگ  ۴۰۰ تومن"},
	}