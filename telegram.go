package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	os.Remove("screenshot.png")
	os.Remove("screenshotNumbEmpty.png")
	os.Remove("screenshotNumbEnter.png")
	os.Remove("screenshotCodeEmpty.png")
	os.Remove("screenshotCodeEnter.png")
	file, _ := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer file.Close()
	logInfo := log.New(file, "INFO\t", log.Ldate|log.Ltime)
	logErr := log.New(file, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	defer logInfo.Println("app exit")

	var screenshotBuffer []byte
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()
	var code, phone string
	phone = ""
	url := "https://web.telegram.org/a"
	err := chromedp.Run(ctx,
		chromedp.EmulateViewport(1920, 1080),
		chromedp.ActionFunc(func(ctx context.Context) error {
			logInfo.Println("Chrome started")
			fmt.Println("Chrome started")
			return nil
		}),
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			logInfo.Printf("Chrome visited %s\n", url)
			fmt.Printf("Chrome visited %s\n", url)
			return nil
		}),
		chromedp.Sleep(5000*time.Millisecond),
		chromedp.Screenshot("html", &screenshotBuffer),
		chromedp.ActionFunc(func(ctx context.Context) error {
			err := os.WriteFile("screenshot.png", screenshotBuffer, 02)
			return err
		}),
		chromedp.Click("button"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			logInfo.Println("Chrome clicked button to reg by phone number")
			fmt.Println("Chrome clicked button to reg by phone number")
			return nil
		}),
		chromedp.Sleep(4000*time.Millisecond),
		chromedp.Screenshot("html", &screenshotBuffer),
		chromedp.ActionFunc(func(ctx context.Context) error {
			err := os.WriteFile("screenshotNumbEmpty.png", screenshotBuffer, 02)
			return err
		}),
		chromedp.SendKeys("input[id=sign-in-phone-number]", phone),
		chromedp.ActionFunc(func(ctx context.Context) error {
			logInfo.Println("Chrome entered phone number")
			fmt.Println("Chrome entered phone number")
			return nil
		}),
		chromedp.Screenshot("html", &screenshotBuffer),
		chromedp.ActionFunc(func(ctx context.Context) error {
			err := os.WriteFile("screenshotNumbEnter.png", screenshotBuffer, 02)
			return err
		}),
		chromedp.Click("div[class=ripple-container]"),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.Screenshot("html", &screenshotBuffer),
		chromedp.ActionFunc(func(ctx context.Context) error {
			err := os.WriteFile("screenshotCodeEmpty.png", screenshotBuffer, 02)
			return err
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := fmt.Println(`Chrome clicked button "next". Enter your telegram code`)
			logInfo.Println(`Chrome clicked button "next".`)
			fmt.Scan(code)
			return err
		}),
		chromedp.SendKeys("input", code),
		chromedp.Sleep(1000*time.Millisecond),
		chromedp.Screenshot("html", &screenshotBuffer),
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("Chrome enter telegram code", code)
			logInfo.Println("Chrome enter telegram code", code)
			err := os.WriteFile("screenshotCodeEnter.png", screenshotBuffer, 02)
			return err
		}),
	)
	if err != nil {
		logErr.Panicln("Error while performing the automation logic:", err)
	}
}
