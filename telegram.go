package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"io"
	"time"

	// "github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

func dialog(login, password, dialogId string) {
	file, _ := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer file.Close()
	lastDialog, _ := os.OpenFile("lastDialog.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer lastDialog.Close()
	logInfo := log.New(file, "INFO\t", log.Ldate|log.Ltime)
	logErr := log.New(file, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	optionsAi := append(
		chromedp.DefaultExecAllocatorOptions[:],
		// chromedp.ProxyServer("45.8.211.64:80"),
		chromedp.Flag("headless", false),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"),
		chromedp.Flag("enable-automation", false),
		// chromedp.Flag("disable-web-security", false),
		chromedp.Flag("disable-web-security", false),
		chromedp.Flag("allow-running-insecure-content", true),
	)
	optionsTg := append(
		chromedp.DefaultExecAllocatorOptions[:],
		// chromedp.ProxyServer("45.8.211.64:80"),
		chromedp.Flag("headless", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"),
		chromedp.Flag("enable-automation", false),
		// chromedp.Flag("disable-web-security", false),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("allow-running-insecure-content", true),
	)
	allocCtxAi, cancel := chromedp.NewExecAllocator(context.Background(), optionsAi...)
	defer cancel()

	allocCtxTg, cancel := chromedp.NewExecAllocator(context.Background(), optionsTg...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtxAi)
	defer cancel()
	var flag bool
	var screenBuffer []byte
	url := "https://pi.ai/"
	err := chromedp.Run(ctx,
		chromedp.EmulateViewport(1570, 730),
		chromedp.ActionFunc(func(ctx context.Context) error {
			logInfo.Println("Chrom to Pi started")
			fmt.Println("Chrome to Pi started")
			return nil
		}),
		chromedp.Navigate(url),
		chromedp.Sleep(100*time.Millisecond),
		chromedp.Click(`#__next > main > div > div > div.relative.flex.h-full.flex-col.items-center > div.flex.w-full.flex-col.items-center.z-20 > button`, chromedp.NodeNotVisible),
		chromedp.Click(`#__next > main > div > div.flex.grow.flex-col.overflow-y-auto.px-6.pb-5.z-70 > div > div > div > div > div > div.space-y-4 > button.flex.items-center.justify-center.whitespace-nowrap.t-action-m.h-14.w-full.max-w-\[353px\].rounded-full.p-4.border-\[1\.5px\].border-neutral-500.bg-\[\#FFF\].text-primary-600`, chromedp.NodeVisible),
		chromedp.Sleep(100*time.Millisecond),
		chromedp.SendKeys("input[type=email]", login, chromedp.NodeVisible),
		chromedp.Sleep(100*time.Millisecond),
		chromedp.Click(`#identifierNext > div > button`, chromedp.NodeVisible),
		chromedp.SendKeys("input[type=password]", password, chromedp.NodeVisible),
		chromedp.Sleep(100*time.Millisecond),
		chromedp.Click(`#passwordNext > div > button`, chromedp.NodeVisible),
		chromedp.Sleep(100*time.Millisecond),
		chromedp.WaitVisible("textarea"),
		chromedp.Navigate("https://pi.ai/threads"),
		chromedp.Sleep(100*time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("Open the required dialog, then send a message to the console")
			var b string
			fmt.Scan(&b)
			return nil
		}),
	)
	if err != nil {
		logErr.Panicln("Error while performing the automation logic:", err)
	}

	var text0, text1 string
	ctxTg, cancel := chromedp.NewContext(allocCtxTg)
	defer cancel()
	urlTg := "https://web.telegram.org/a/#" + dialogId
	err = chromedp.Run(ctxTg,
		chromedp.EmulateViewport(1920, 1080),
		chromedp.ActionFunc(func(ctx context.Context) error {
			logInfo.Println("Chrome to TG started")
			return nil
		}),
		chromedp.Navigate(urlTg),
		chromedp.Sleep(10000*time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			logInfo.Printf("Chrome visited %s\n", urlTg)
			return err
		}),
		chromedp.Sleep(7000*time.Millisecond),
		chromedp.Screenshot("#auth-qr-form > div > div", &screenBuffer, chromedp.NodeVisible),
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("Scan QR in QR.png")
			os.WriteFile("QR.png", screenBuffer, 02)
			logInfo.Println("Made the ScreenShot with QR")
			return err
		}),
		chromedp.WaitVisible("#editable-message-text"),
	)
	fmt.Println("Wait")
	os.Remove("QR.png")
	time.Sleep(12000 * time.Millisecond)
	bufferedChan := make(chan string, 1)
	// var wg sync.WaitGroup
	// wg.Add(2)
	go func() {
		var command string
		for {
			if command == "a" && flag == false {
				time.Sleep(400 * time.Millisecond)
				continue
			}
			fmt.Scan(&command)
			bufferedChan <- command
		}
	}()
	func() {
		for {
			chromedp.Run(ctxTg,
				chromedp.ActionFunc(func(ctxTg context.Context) error {
					fmt.Println("Ready to process messages or your command(q - quit, s - start dialog request, c - continue dialog reauest, a - custom dialog request)")
					logInfo.Println("Ready to process messages or command")
					return err
				}),

				chromedp.Text(`div[class="messages-container"]`, &text0),
				chromedp.ActionFunc(func(ctxTg context.Context) error {
					var answer string
					if len(text0) > 1000 {
						text0 = text0[len(text0)-999:]
					}
					lastDialog.WriteString(text0)
					text1 = text0
					for text0 == text1 {
						select {
						case msg := <-bufferedChan:
							{
								logInfo.Println("The command has been received", msg)
								switch msg {
								case "q":
									logInfo.Panicln("App exit by User")
								case "s":
									{
										logInfo.Println("The command to start the dialogue has been received", msg)
										fmt.Println("The command to start the dialogue has been received")
										err = chromedp.Run(ctx,
											chromedp.SendKeys(`#__next > main > div > div > div.relative.grow.overflow-x-auto.hidden.lg\:flex.lg\:flex-col > div.relative.flex.flex-col.overflow-hidden.sm\:overflow-x-visible.h-full.pt-8.grow > div.max-h-\[40\%\].px-5.sm\:px-0.z-15.w-full.mx-auto.max-w-1\.5xl.\32 xl\:max-w-\[47rem\] > div > div > textarea`, "Представь, что этого сообщения нет, не отвечай на него, просто начни диалог"),
											chromedp.MouseClickXY(1360, 630),
											chromedp.Sleep(10000*time.Millisecond),
											chromedp.Text(`body > #__next > main > div > div > .relative > .relative > .grow > div > div > div > .pb-6 > div > div > .break-anywhere > .flex > div > div`, &answer),
											chromedp.ActionFunc(func(ctx context.Context) error {
												fmt.Println(`the response has been received and contains "` + answer + `"`)
												logInfo.Println(`the response has been received and contains "` + answer + `"`)
												return err
											}),
										)
										if err != nil {
											logErr.Panicln("Error while performing the automation logic:", err)
										}
										err = chromedp.Run(ctxTg,
											chromedp.SendKeys(`#editable-message-text`, answer),
											chromedp.Click(`#MiddleColumn > div.messages-layout > div.Transition > div > div.middle-column-footer > div.Composer.shown.mounted > button`),
											chromedp.Sleep(500*time.Millisecond),
										)
										logInfo.Println("The response has been sent(start dialog)")
										flag = true
										fmt.Println("Ready to process messages or your command(q - quit, s - start dialog request, c - continue dialog reauest, a - custom dialog request)")
										logInfo.Println("Ready to process messages or command")
									}
								case "c":
									{
										logInfo.Println("The command to continue the dialogue has been received", msg)
										fmt.Println("The command to continue the dialogue has been received")
										err = chromedp.Run(ctx,
											chromedp.SendKeys(`#__next > main > div > div > div.relative.grow.overflow-x-auto.hidden.lg\:flex.lg\:flex-col > div.relative.flex.flex-col.overflow-hidden.sm\:overflow-x-visible.h-full.pt-8.grow > div.max-h-\[40\%\].px-5.sm\:px-0.z-15.w-full.mx-auto.max-w-1\.5xl.\32 xl\:max-w-\[47rem\] > div > div > textarea`, "Представь, что этого сообщения нет, не отвечай на него, просто продолжи диалог на любую тему"),
											chromedp.MouseClickXY(1360, 630),
											chromedp.Sleep(10000*time.Millisecond),
											chromedp.Text(`body > #__next > main > div > div > .relative > .relative > .grow > div > div > div > .pb-6 > div > div > .break-anywhere > .flex > div > div`, &answer),
											chromedp.ActionFunc(func(ctx context.Context) error {
												fmt.Println(`the response has been received and contains "` + answer + `"`)
												logInfo.Println(`the response has been received and contains "` + answer + `"`)
												return err
											}),
										)
										if err != nil {
											logErr.Panicln("Error while performing the automation logic:", err)
										}
										err = chromedp.Run(ctxTg,
											chromedp.SendKeys(`#editable-message-text`, answer),
											chromedp.Click(`#MiddleColumn > div.messages-layout > div.Transition > div > div.middle-column-footer > div.Composer.shown.mounted > button`),
											chromedp.Sleep(500*time.Millisecond),
										)
										logInfo.Println("The response has been sent(start dialog)")
										flag = true
										fmt.Println("Ready to process messages or your command(q - quit, s - start dialog request, c - continue dialog reauest, a - custom dialog request)")
										logInfo.Println("Ready to process messages or command")
									}
								case "a":
									{
										logInfo.Println("The command to custom request has been received", msg)
										fmt.Println(`Enter your request`)
										in := bufio.NewReader(os.Stdin)
										_, _ = in.ReadString('\n')
										request, _ := in.ReadString('\n')
										logInfo.Println(`Custom request contained "` + request + `"`)
										err = chromedp.Run(ctx,
											chromedp.SendKeys(`#__next > main > div > div > div.relative.grow.overflow-x-auto.hidden.lg\:flex.lg\:flex-col > div.relative.flex.flex-col.overflow-hidden.sm\:overflow-x-visible.h-full.pt-8.grow > div.max-h-\[40\%\].px-5.sm\:px-0.z-15.w-full.mx-auto.max-w-1\.5xl.\32 xl\:max-w-\[47rem\] > div > div > textarea`, request),
											chromedp.MouseClickXY(1360, 630),
											chromedp.Sleep(10000*time.Millisecond),
											chromedp.Text(`body > #__next > main > div > div > .relative > .relative > .grow > div > div > div > .pb-6 > div > div > .break-anywhere > .flex > div > div`, &answer),
											chromedp.ActionFunc(func(ctx context.Context) error {
												fmt.Println(`the response has been received and contains "` + answer + `"`)
												logInfo.Println(`the response has been received and contains "` + answer + `"`)
												return err
											}),
										)
										if err != nil {
											logErr.Panicln("Error while performing the automation logic:", err)
										}
										err = chromedp.Run(ctxTg,
											chromedp.SendKeys(`#editable-message-text`, answer),
											chromedp.Click(`#MiddleColumn > div.messages-layout > div.Transition > div > div.middle-column-footer > div.Composer.shown.mounted > button`),
											chromedp.Sleep(500*time.Millisecond),
										)
										logInfo.Println("The response has been sent(start dialog)")
										flag = true
										fmt.Println("Ready to process messages or your command(q - quit, s - start dialog request, c - continue dialog reauest, a - custom dialog request)")
										logInfo.Println("Ready to process messages or command")
									}
								}
							}
						default:
							{
								err := chromedp.Run(ctxTg,
									chromedp.Sleep(10000*time.Millisecond),
									chromedp.Text(`div[class="messages-container"]`, &text0),
									chromedp.ActionFunc(func(ctxTg context.Context) error {
										if len(text0) > 1000 {
											text0 = text0[len(text0)-999:]
										}
										if flag {
											text1 = text0
											flag = false
										}
										return nil
									}),
								)
								if err != nil {
									logErr.Panicln("Error while performing the automation logic:", err)
								}
								if text0 == text1 {
									lastDialog.Seek(0, io.SeekStart)
									lastDialog.Truncate(0)
									lastDialog.WriteString(text0)
									text1 = text0
								} else {
									lastDialog.Seek(0, io.SeekStart)
									lastDialog.Truncate(0)
									lastDialog.WriteString(text0)
									fmt.Println("A new message has been received from telegram")
									logInfo.Println("A new message has been received from telegram")
									break
								}
							}
						}
					}
					return nil
				}),
			)
			if err != nil {
				logErr.Panicln("Error while performing the automation logic:", err)
			}
			str := string(text0[:len(text0)-5])
			for i := len(str); true; i-- {
				if string(str[i-1]) >= "0" && string(str[i-1]) <= "9" && string(str[i-2]) >= "0" && string(str[i-2]) <= "9" && string(str[i-3]) == ":" && string(str[i-4]) >= "0" && string(str[i-4]) <= "9" && string(str[i-5]) >= "0" && string(str[i-5]) <= "9" {
					str = str[i:]
					logInfo.Println(`New message contain "` + str + `"`)
					if len(str) == 0 {
						continue
					}
					break
				} else if str[i-20:i-2] == "Владислав" {
					str = str[i-2:]
					logInfo.Println(`New message contain "` + str + `"`)
					if len(str) == 0 {
						continue
					}
					break
				}
			}
			var answer string
			err = chromedp.Run(ctx,
				chromedp.SendKeys(`textarea[role="textbox"]`, str),
				// chromedp.MouseClickXY(1360, 630),
				chromedp.Sleep(7000*time.Millisecond),
				chromedp.Text(`body > #__next > main > div > div > .relative > .relative > .grow > div > div > div > .pb-6 > div > div > .break-anywhere > .flex > div > div`, &answer),
				chromedp.ActionFunc(func(ctx context.Context) error {
					fmt.Println(`the response has been received and contains "` + answer + `"`)
					logInfo.Println(`the response has been received and contains "` + answer + `"`)
					return err
				}),
			)
			if err != nil {
				logErr.Panicln("Error while performing the automation logic:", err)
			}
			err = chromedp.Run(ctxTg,
				chromedp.SendKeys(`#editable-message-text`, answer),
				chromedp.Click(`#MiddleColumn > div.messages-layout > div.Transition > div > div.middle-column-footer > div.Composer.shown.mounted > button`),
				chromedp.Sleep(500*time.Millisecond),
			)
			logInfo.Println("The response has been sent")
		}
	}()
}

func main() {
	os.Remove("QR.png")
	file, _ := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer file.Close()
	logInfo := log.New(file, "APP_INFO\t", log.Ldate|log.Ltime)
	defer logInfo.Println("app exit")
	var dialogId string
	fmt.Println("Change dialog Mis_Kitsune - \"1238372228\"   Inal - \"833591886\" Blodos_Dodos_Bot - \"5650924958\" Арт - \"6133569386\" Юра - \"871044396\"\nПроектик - \"-4081628480\" Ермолов - \"1498293686\" Саша ДВ(^^) -\"1891226386\" Настя ДВ(milk_catt) - 1180819964")
	fmt.Scan(&dialogId)

	dialog("bot408916@gmail.com", "1818ASIUbf23", dialogId)
	// telegram(dialogId)
}
