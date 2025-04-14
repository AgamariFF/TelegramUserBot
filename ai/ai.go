package ai

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"telegram/internal"
	"telegram/tg"
	"time"

	"github.com/chromedp/chromedp"
)

const ProfilePathAi = "./chrome_profile_ai_"

func setupAi(login, password string, logInfo *log.Logger, logErr *log.Logger, profilePath string, consoleMutex *sync.Mutex) {
	logInfo.Println("Chrome to logining in Ai started")
	var wait string

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"),
		chromedp.Flag("enable-automation", false),
		// chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("disable-infobars", true),
		chromedp.Flag("user-data-dir", profilePath),
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-default-browser-check", true),
	)

	allocCtx, cancel0 := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel0()

	ctx, cancel1 := chromedp.NewContext(allocCtx,
		chromedp.WithLogf(func(string, ...interface{}) {}),
		chromedp.WithDebugf(func(string, ...interface{}) {}),
		chromedp.WithErrorf(func(string, ...interface{}) {}),
	)
	defer cancel1()

	err := chromedp.Run(ctx,
		chromedp.Navigate("https://character.ai/"),
	)
	if err != nil {
		logErr.Println("Error while performing the automation logic:", err)
	}
	consoleMutex.Lock()
	fmt.Println("Отправьте сообщение в чат после авторизации")
	fmt.Scan(&wait)
	consoleMutex.Unlock()
	cancel0()
	cancel1()
}

func AiHeandler(login, pasword string, logInfo *log.Logger, logErr *log.Logger, lastDialog os.File, consoleMutex *sync.Mutex, incoming chan tg.Message, outcoming chan tg.Message, wg *sync.WaitGroup) {
	var profilePath string
	var incomingMsg, outcomingMsg tg.Message
	var auth bool
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
		// chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("disable-infobars", true),
		chromedp.Flag("user-data-dir", profilePath),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-default-browser-check", true),
	)
	var ctx context.Context

	allocCtx, cancel0 := chromedp.NewExecAllocator(context.Background(), opts...)

	ctx, cancel1 := chromedp.NewContext(allocCtx,
		chromedp.WithLogf(func(string, ...interface{}) {}),
		chromedp.WithDebugf(func(string, ...interface{}) {}),
		chromedp.WithErrorf(func(string, ...interface{}) {}),
	)

	defer cancel0()
	defer cancel1()

	url := "https://character.ai/"
	var text string

	err := chromedp.Run(ctx,
		chromedp.EmulateViewport(1570, 730),
		chromedp.ActionFunc(func(ctx context.Context) error {
			logInfo.Println("Chrome to Ai started")
			return nil
		}),
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second),
		chromedp.Text(`/html`, &text),
	)

	if err != nil {
		logErr.Println(err)
	}

	if strings.Contains(text, "Продолжить с Google") {
		auth = true
		fmt.Println("Необходима авторизация в AI для профиля: ", profilePath)
		logInfo.Println("Необходима авторизация в AI для профиля: ", profilePath)
		cancel0()
		cancel1()
		setupAi(login, pasword, logInfo, logErr, profilePath, consoleMutex)

		allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancel()

		ctx, cancel = chromedp.NewContext(allocCtx,
			chromedp.WithLogf(func(string, ...interface{}) {}),
			chromedp.WithDebugf(func(string, ...interface{}) {}),
			chromedp.WithErrorf(func(string, ...interface{}) {}),
		)
		defer cancel()

	} else {
		fmt.Println("Авторизация в Ai не требуйется")
		logInfo.Println("Авторизация в Ai не требуйется")
	}

	var exit bool
	var chat string
	consoleMutex.Lock()
	for !exit {
		fmt.Println("Укажите номер чата (хозяин - 7) или введите свою ссылку для нового чата")
		fmt.Scan(&chat)
		logInfo.Println("Пользователь выбрал чат в Ai: ", chat)
		if strings.Contains(chat, "http") {
			logInfo.Println("Пользователь указал свою ссылку на чат")
			url = chat
			break
		} else {
			switch chat {
			case "1":
				{
					url = "https://character.ai/chat/S8NFtpIzsPYAESWca80JCp-1U8aefZXg9ERgkE4UqW0"
					exit = true
					break
				}
			case "2":
				{
					url = "https://character.ai/chat/qBFHJFQHlcxUDHTiyT6qGdUkR-_zHA96yhJrMTO8hnc"
					exit = true
					break
				}
			case "3":
				{
					url = "https://character.ai/chat/q9wSQFshuaSZ6f66SiPQDDo84AAxY3BTmxC-90WeZ04"
					exit = true
					break
				}
			case "4":
				{
					url = "https://character.ai/chat/Ac3HPkMB9gUwl4HrpjP0wVza_Cc6ec11Xw5EvuZcgWg"
					exit = true
					break
				}
			case "5":
				{
					url = "https://character.ai/chat/tIo-MLdv9NeKVXtqsWRhLcXLVt4cZhlQ4NGGpdHdI8o"
					exit = true
					break
				}
			case "6":
				{
					url = "https://character.ai/chat/jKzVHisMWlU4oIQsETo2DTNiOD2NNyFBpp5mnUbT4V8"
					exit = true
					break
				}
			case "7":
				{
					url = "https://character.ai/chat/VnitHr6Vf1qStUsBKQjGV0gvJm9jIO_u-OY4qxCNFq8"
					exit = true
					break
				}
			default:
				fmt.Println("Неверный ввод, повторите попытку")
			}
		}
	}
	consoleMutex.Unlock()

	if auth {
		err = chromedp.Run(ctx,
			chromedp.EmulateViewport(1500, 800),
			chromedp.Navigate(url),
			chromedp.Sleep(5*time.Second),
			chromedp.MouseClickXY(1123, 193),
			chromedp.Sleep(5*time.Second),
			chromedp.Navigate(url),
		)
	} else {
		err = chromedp.Run(ctx,
			chromedp.Navigate(url),
		)
		if err != nil {
			logErr.Panicln(err)
		}
	}

	wg.Done()

	for {
		incomingMsg = <-incoming
		logInfo.Println("AiHeabdler считал входящее сообщение: ", incomingMsg)
		outcomingMsg.Id = incomingMsg.Id
		err = chromedp.Run(ctx,
			chromedp.SendKeys(`//*[@id="chat-body"]/div[2]/div/div/div/div[1]/textarea`, incomingMsg.Text, chromedp.NodeVisible),
			chromedp.Click(`//*[@id="chat-body"]/div[2]/div/div/div/div[2]/button`, chromedp.NodeVisible),
			chromedp.Sleep(10*time.Second),
			chromedp.Text(`//*[@id="chat-messages"]/div[1]/div[1]/div/div/div[1]/div/div[1]/div[1]/div[2]/div[2]/div/div[1]`, &outcomingMsg.Text, chromedp.NodeVisible),
		)
		if err != nil {
			logErr.Println(err)
			continue
		}
		if strings.HasSuffix(outcomingMsg.Text, ".") {
			outcomingMsg.Text = outcomingMsg.Text[:len(outcomingMsg.Text)-1]
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
