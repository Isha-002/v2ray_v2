package main

import (
    tele "gopkg.in/telebot.v3"
    "time"
    "fmt"
)

var b *tele.Bot  

func InitBot() error {
    b, err := tele.NewBot(tele.Settings{
        Token:  "7743565843:AAHjGGubY_yLb51cacm0T-gEMtHtHd-IYWQ",
        Poller: &tele.LongPoller{Timeout: 10 * time.Second},
    })

    if err != nil {
        return fmt.Errorf("failed to initialize bot: %w", err)
    }

    go b.Start()

    return nil
}
