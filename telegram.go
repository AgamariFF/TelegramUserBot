package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func telegram(dialogId string) {
	file, _ := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer file.Close()
	logInfo := log.New(file, "TG_INFO\t", log.Ldate|log.Ltime)
	logErr := log.New(file, "TG_ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	var screenshotBuffer []byte
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()
	url := "https://web.telegram.org/a"
	err := chromedp.Run(ctx,
		chromedp.EmulateViewport(1920, 1080),
		chromedp.ActionFunc(func(ctx context.Context) error {
			logInfo.Println("Chrome started")
			fmt.Println("Chrome started")
			return nil
		}),
		chromedp.Navigate(url),
		chromedp.Sleep(1000*time.Millisecond),
		chromedp.FullScreenshot(&screenshotBuffer, 100),
		chromedp.ActionFunc(func(ctx context.Context) error {
			err := os.WriteFile("TgScreen\\StartTelegram.png", screenshotBuffer, 02)
			logInfo.Printf("Chrome visited %s\n", url)
			fmt.Printf("Chrome visited %s\n", url)
			return err
		}),
		chromedp.Sleep(7000*time.Millisecond),
		chromedp.Screenshot("div[class=qr-outer]", &screenshotBuffer),
		chromedp.ActionFunc(func(ctx context.Context) error {
			err := os.WriteFile("TgScreen\\QR.png", screenshotBuffer, 02)
			fmt.Println("Scan QR code in QR.png")
			logInfo.Println("Made a screen of QR")
			return err
		}),

		chromedp.WaitVisible("div[id=peer-story"+dialogId+"]"),
		chromedp.FullScreenshot(&screenshotBuffer, 100),
		chromedp.ActionFunc(func(ctx context.Context) error {
			err := os.WriteFile("TgScreen\\UserTelegram.png", screenshotBuffer, 02)
			fmt.Println("Chrome etner user account by QR")
			logInfo.Println("Chrome etner user account by QR")
			return err
		}),
		chromedp.Click("div[id=peer-story"+dialogId+"]"),
		chromedp.Sleep(1000*time.Millisecond),
		chromedp.FullScreenshot(&screenshotBuffer, 100),
		chromedp.ActionFunc(func(ctx context.Context) error {
			err := os.WriteFile("TgScreen\\DialogTelegram.png", screenshotBuffer, 02)
			fmt.Println("Open Dialog in Telegram")
			logInfo.Println("Open Dialog in Telegram")
			return err
		}),
	)
	if err != nil {
		logErr.Panicln("Error while performing the automation logic:", err)
	}
}

func gpt(login, password string) {
	file, _ := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer file.Close()
	logInfo := log.New(file, "GPT_INFO\t", log.Ldate|log.Ltime)
	logErr := log.New(file, "GPT_ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	// var screenshotBuffer []byte
	options := append(
		chromedp.DefaultExecAllocatorOptions[:],
		// chromedp.ProxyServer("45.8.211.64:80"),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"),
		chromedp.Flag("headless", false),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-web-security", false),
		chromedp.Flag("allow-running-insecure-content", true),
	)
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()
	url := "https://pi.ai/"
	err := chromedp.Run(ctx,
		chromedp.EmulateViewport(1280, 720),
		chromedp.ActionFunc(func(ctx context.Context) error {
			logInfo.Println("Chrome started")
			fmt.Println("Chrome started")
			return nil
		}),
		chromedp.Navigate(url),
		chromedp.WaitVisible("input[type=email]"),
		chromedp.SendKeys("input[type=email]", login),
		chromedp.Sleep(505 * time.Millisecond),
		chromedp.MouseClickXY(1250, 740),
		chromedp.WaitVisible("input[type=password]"),
		chromedp.SendKeys("input[type=password]", password),
		// chromedp.ActionFunc(func(ctx context.Context) error {
		// 	fmt.Println("Вижу пароль")
		// 	return nil
		// }),
		// chromedp.SendKeys("input[name=password]",password),
		// chromedp.Click("button[type=submit]"),
		chromedp.Sleep(300000*time.Millisecond),
	)
	if err != nil {
		logErr.Panicln("Error while performing the automation logic:", err)
	}
}

func main() {
	os.Remove("TgScreen\\QR.png")
	os.Remove("TgScreen\\StartTelegram.png")
	os.Remove("TgScreen\\UserTelegram.png")
	os.Remove("TgScreen\\DialogTelegram.png")
	os.Remove("GPT\\1.png")
	os.Remove("GPT\\2.png")
	os.Remove("GPT\\3.png")
	os.Remove("GPT\\4.png")
	os.Remove("GPT\\5.png")
	os.Remove("GPT\\6.png")
	os.Remove("GPT\\action.png")

	file, _ := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer file.Close()
	logInfo := log.New(file, "APP_INFO\t", log.Ldate|log.Ltime)
	defer logInfo.Println("app exit")
	// dialogId := "1238372228"

	gpt("bot408916@gmail.com", "1818ASIUbf23")
	// telegram(dialogId)
}
