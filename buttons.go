package main

import (
    tele "gopkg.in/telebot.v3"
)

// Create inline buttons
var buyBtn = tele.Btn{Unique: "buy_service", Text: "خرید سرویس"}
var renewBtn = tele.Btn{Unique: "renew_service", Text: "تمدید سرویس"}
var oneUserPlans = tele.Btn{Unique: "one_user_plans", Text: "تک کاربره"}
var twoUserPlans = tele.Btn{Unique: "two_user_plan", Text: "دو کاربره"}
var unlimitedUserPlans = tele.Btn{Unique: "unlimited_user_plan", Text: "کاربر نامحدود"}

// Renew buttons
var renewPlan = tele.Btn{Unique: "renew_same", Text: "تمدید سرویس فعلی"}
var renewAnotherPlan = tele.Btn{Unique: "renew_another", Text: "تغییر سرویس فعلی"}

// Back buttons
var backBtn = tele.Btn{Unique: "back_btn", Text: "↩️ بازگشت"}
var backtoMainBtn = tele.Btn{Unique: "backToMain_btn", Text: "↩️ بازگشت"}

