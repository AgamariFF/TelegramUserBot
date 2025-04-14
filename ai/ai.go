package ai

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"telegram/internal"
	"telegram/tg"
	"time"

	"github.com/chromedp/chromedp"
)

const ProfilePathAi = "./chrome_profile_ai_"

func setupAi(login, password string, logInfo *log.Logger, logErr *log.Logger, ctx context.Context) {
	logInfo.Println("Chrome to logining in Pi started")

	err := chromedp.Run(ctx,
		chromedp.Navigate("https://pi.ai/profile/account"),
		chromedp.Sleep(100*time.Millisecond),
		chromedp.Click(`//*[@id="__next"]/main/div/div[3]/div[2]/div/div/div/div/div/div[1]/button[1]`, chromedp.NodeNotVisible),
		chromedp.Sleep(100*time.Millisecond),
		chromedp.SendKeys("input[type=email]", login, chromedp.NodeVisible),
		chromedp.Sleep(100*time.Millisecond),
		chromedp.Click(`#identifierNext > div > button`, chromedp.NodeVisible),
		chromedp.SendKeys("input[type=password]", password, chromedp.NodeVisible),
		chromedp.Sleep(100*time.Millisecond),
		chromedp.Click(`#passwordNext > div > button`, chromedp.NodeVisible),
		chromedp.Sleep(100*time.Millisecond),
		chromedp.WaitVisible("textarea"),
	)
	if err != nil {
		logErr.Println("Error while performing the automation logic:", err)
	}
}

func AiHeandler(login, pasword string, logInfo *log.Logger, logErr *log.Logger, lastDialog os.File, consoleMutex *sync.Mutex, incoming chan tg.Message, outcoming chan tg.Message, wg *sync.WaitGroup) {
	var profilePath string
	var reloadThreads bool
	for i := 0; i < 100; i++ {
		profilePath = ProfilePathAi + strconv.Itoa(i) + "/"
		if !internal.IsBrowserRunning(profilePath) {
			break
		}
		if i == 99 {
			panic("Не удалось найти свободный профиль")
		}
		logInfo.Println("Свободный путь для Ai: ", profilePath)
	}
	if err := os.MkdirAll(profilePath, 0755); err != nil {
		panic(err)
	}
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("disable-infobars", true),
		chromedp.Flag("user-data-dir", profilePath),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-default-browser-check", true),
	)
	var ctx context.Context

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(allocCtx,
		chromedp.WithLogf(func(string, ...interface{}) {}),
		chromedp.WithDebugf(func(string, ...interface{}) {}),
		chromedp.WithErrorf(func(string, ...interface{}) {}),
	)
	defer cancel()

	url := "https://pi.ai/threads"
	var exists bool

	err := chromedp.Run(ctx,
		chromedp.EmulateViewport(1570, 730),
		chromedp.ActionFunc(func(ctx context.Context) error {
			logInfo.Println("Chrome to Pi started")
			return nil
		}),
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second),
		chromedp.Evaluate(`
            (function() {
                const xpath = '//*[@id="__next"]/main/div/div/div[2]/div[3]/div/div/button[2]';
                const result = document.evaluate(
                    xpath,
                    document,
                    null,
                    XPathResult.FIRST_ORDERED_NODE_TYPE,
                    null
                );
                return !!result.singleNodeValue;
            })()
        `, &exists),
	)

	if err != nil {
		logErr.Println(err)
	}

	if !exists {
		logInfo.Println("login is required in AI, start setupAi")
		setupAi(login, pasword, logInfo, logErr, ctx)
	} else {
		logInfo.Println("the entry has already been completed in AI")
	}

	var xpath string
	var count int
	var screenBuffer []byte

	logInfo.Println("AiHeandler started")
	err = chromedp.Run(ctx,
		chromedp.Navigate("https://pi.ai/threads"),
		chromedp.Sleep(1*time.Second),
		chromedp.Evaluate(`
  document.querySelectorAll(
    '#__next > main > div > div > div:nth-child(2) > div:nth-child(3) > div > div button'
  ).length
`, &count),
		chromedp.Screenshot(`//*[@id="__next"]/main/div/div/div[2]/div[3]`, &screenBuffer, chromedp.NodeVisible),
		chromedp.ActionFunc(func(ctx context.Context) error {
			os.WriteFile("Chats.png", screenBuffer, 0644)
			return nil
		}),
	)
	if err != nil {
		logErr.Println(err)
	}
	consoleMutex.Lock()
	var b int

	for {
		fmt.Printf("Скриншот чатов сохранен в Chats.png\nДоступно %d чатов. Выбирите номер чата по счету или 0 если делаете это вручную\n", count)
		fmt.Scan(&b)
		os.Remove("Chats.png")
		if b == 0 {
			break
		} else if b > count {
			fmt.Println("Нет такого треда, повторите ввод")
			continue
		}
		xpath = fmt.Sprintf(`/html/body/div/main/div/div/div[2]/div[3]/div/div/button[%d]`, b)
		break
	}

	wg.Done()
	consoleMutex.Unlock()

	err = chromedp.Run(ctx,
		chromedp.Click(xpath, chromedp.NodeVisible),
	)
	if err != nil {
		logErr.Println(err)
	}

	logInfo.Println("Выбор треда сделан")

	var incomingMsg tg.Message
	var outcomingMsg tg.Message

	for {
		incomingMsg = <-incoming
		logInfo.Println("AiHeabdler считал входящее сообщение: ", incomingMsg)
		outcomingMsg.Id = incomingMsg.Id
		err = chromedp.Run(ctx,
			chromedp.SendKeys(`//*[@id="__next"]/main/div/div/div[3]/div[1]/div[4]/div/div/textarea`, incomingMsg.Text, chromedp.NodeVisible),
			chromedp.Click(`//*[@id="__next"]/main/div/div/div[3]/div[1]/div[4]/div/button`, chromedp.NodeVisible),
		)
		if b != 1 && !reloadThreads {
			reloadThreads = true
			err = chromedp.Run(ctx,
				chromedp.Sleep(3*time.Second),
				chromedp.Navigate(url),
				chromedp.Sleep(3*time.Second),
				chromedp.Click(`/html/body/div/main/div/div/div[2]/div[3]/div/div/button[2]`, chromedp.NodeVisible),
				chromedp.Sleep(4*time.Second),
			)
		} else {
			time.Sleep(10 * time.Second)
		}

		err = chromedp.Run(ctx,
			chromedp.Text(`//*[@id="__next"]/main/div/div/div[3]/div[1]/div[2]/div/div/div/div[3]/div/div/div[2]/div[1]`, &outcomingMsg.Text, chromedp.NodeVisible),
		)
		logInfo.Println("Овет от Ai: ", outcomingMsg.Text)
		if err != nil {
			logErr.Println(err)
		}
		if len(outcomingMsg.Text) > 0 {
			if outcomingMsg.Text[len(outcomingMsg.Text)-1] == '.' {
				outcomingMsg.Text = outcomingMsg.Text[:len(outcomingMsg.Text)-1]
			}
		}
		fmt.Println("Ответ от Ai: " + outcomingMsg.Text)
		_, err = lastDialog.WriteString("\nin: " + incomingMsg.Text + "\nout: " + outcomingMsg.Text)
		if err != nil {
			logErr.Println(err)
		}
		logInfo.Println("AiHeabdler обработал сообщение и пытается отправить ответ в oucoming: ", outcomingMsg)
		outcoming <- outcomingMsg
		logInfo.Println("AiHeabdler отправил ответ в outcoming: ", outcomingMsg)
	}
}
