package tg

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"telegram/internal"
	"time"

	"github.com/chromedp/chromedp"
)

const ProfilePathTg = "./chrome_profile_tg_"

// Структура для хранения сообщений
type Message struct {
	Text string
	Id   int
}

func SetupTg(logInfo *log.Logger, logErr *log.Logger, consoleMutex *sync.Mutex, wg *sync.WaitGroup) (context.Context, context.CancelFunc) {
	var profilePath string
	for i := 0; i < 100; i++ {
		profilePath = ProfilePathTg + strconv.Itoa(i) + "/"
		if !internal.IsBrowserRunning(profilePath) {
			break
		}
		if i == 99 {
			panic("Не удалось найти свободный профиль")
		}
		logInfo.Println("Свободный путь для Tg: ", profilePath)
	}
	if err := os.MkdirAll(profilePath, 0755); err != nil {
		panic(err)
	}
	optionsTg := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("user-data-dir", profilePath),
		// chromedp.ProxyServer("45.8.211.64:80"),
		chromedp.Flag("headless", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"),
		chromedp.Flag("enable-automation", false),
		// chromedp.Flag("disable-web-security", true),
		// chromedp.Flag("allow-running-insecure-content", true),
	)

	var dialogId string
	consoleMutex.Lock()
	fmt.Println("Change dialog Mis_Kitsune - \"1238372228\"   Inal - \"833591886\" Blodos_Dodos_Bot - \"5650924958\" Арт - \"6133569386\" Юра - \"871044396\"\nПроектик - \"-4081628480\" Ермолов - \"1498293686\" Саша ДВ(^^) -\"1891226386\" Настя ДВ(milk_catt) - 1180819964")
	fmt.Scan(&dialogId)
	logInfo.Println("Выбран чат с id = ", dialogId)
	consoleMutex.Unlock()
	wg.Done()

	var screenBuffer []byte
	var exists bool

	ctx, cancel := chromedp.NewContext(func() context.Context {
		ctx, _ := chromedp.NewExecAllocator(context.Background(), optionsTg...)
		return ctx
	}(),
		chromedp.WithLogf(func(string, ...interface{}) {}),
		chromedp.WithDebugf(func(string, ...interface{}) {}),
		chromedp.WithErrorf(func(string, ...interface{}) {}),
	)

	urlTg := "https://web.telegram.org/a/#" + dialogId
	err := chromedp.Run(ctx,
		chromedp.EmulateViewport(1500, 800),
		chromedp.ActionFunc(func(ctx context.Context) error {
			logInfo.Println("Chrome to TG started")
			return nil
		}),
		chromedp.Navigate(urlTg),
		chromedp.Sleep(3*time.Second),
		chromedp.Evaluate(`
            document.body.innerText.includes("Log in to Telegram by QR Code")
        `, &exists),
	)
	if err != nil {
		logErr.Println(err)
	}

	if !exists {
		logInfo.Println("the entry has already been completed in Tg")
		fmt.Println("Вход в Tg не требуется")
		return ctx, cancel
	}
	logInfo.Println("login is required in Tg")

	err = chromedp.Run(ctx,
		chromedp.EmulateViewport(1500, 800),
		chromedp.WaitVisible(`//*[@id="auth-qr-form"]/div/button[1]`),
		chromedp.Sleep(2*time.Second),
		chromedp.Screenshot("#auth-qr-form > div > div", &screenBuffer, chromedp.NodeVisible),
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("Scan QR in QR.png")
			os.WriteFile("QR.png", screenBuffer, 0644)
			logInfo.Println("Made the ScreenShot with QR")
			return nil
		}),
		chromedp.WaitVisible("#editable-message-text"),
	)
	if err != nil {
		logErr.Panicln("Error while performing the automation logic:", err)
	}
	os.Remove("QR.png")
	time.Sleep(6 * time.Second)

	return ctx, cancel
}

func MonitorPartnerMessagesAndSend(ctx context.Context, logInfo *log.Logger, logErr *log.Logger, lastDialog *os.File, incoming chan Message, outcoming chan Message, timeSleep time.Duration) {
	logInfo.Println("MonitorPartnerMessages запущен")
	var outcomingMsg Message
	var maxN *int
	var n int
	lastMsg := Message{"", 0}
	time.Sleep(time.Second)
	err := chromedp.Run(ctx,
		chromedp.EmulateViewport(1500, 800),
		chromedp.Evaluate(`
            (() => {
                const elements = Array.from(document.querySelectorAll('[id^="message-"]'));
                if (elements.length === 0) return null;
                
                const nums = elements
                    .map(el => {
                        const match = el.id.match(/message-(\d+)/);
                        return match ? parseInt(match[1], 10) : null;
                    })
                    .filter(n => n !== null);
                
                return nums.length > 0 ? Math.max(...nums) : null;
            })()
        `, &maxN),
	)
	logInfo.Println("Максимальноу Id сообщения: ", *maxN)

	if err != nil {
		logErr.Fatalln("Ошибка поиска максимального n:", err)
	}

	if maxN == nil {
		logErr.Fatalln("Элементы не найдены или некорректные ID")
	}
	// if (lastMsg.Id - 10000) >
	for n = *maxN; n >= lastMsg.Id; n-- {
		var outMsg bool
		var exists bool
		var classAttr string
		xpath := fmt.Sprintf(`//*[@id="message-%d"]/div[3]/div/div[1]/div`, n)

		// Проверка существования элемента
		err := chromedp.Run(ctx,
			chromedp.Evaluate(
				fmt.Sprintf(`
                !!document.evaluate('%s', document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue
            `, xpath),
				&exists,
			),
		)
		if err != nil || !exists {
			continue
		} else {
			err = chromedp.Run(ctx,
				chromedp.AttributeValue(xpath, "class", &classAttr, &exists),
			)
			if err != nil {
				logErr.Println(err)
				continue
			}
		}

		logInfo.Println("Проверяю сообщение с id=", n, "\tАттрибуты: ", classAttr)

		substrings := []string{"with-outgoing-icon", "own"}
		for _, substing := range substrings {
			if strings.Contains(classAttr, substing) {
				outMsg = true
			}
		}
		if outMsg {
			logInfo.Println("Это исходящее сообщение")
			outMsg = false
			continue
		}
		lastMsg.Id = n
		logInfo.Println("Последнее входящее сообщение в этом чате имеет id = ", n)
		break
	}
	fmt.Println("Начал считывать сообщения в Tg")

	var skipMsg bool

	for {
		select {
		case outcomingMsg = <-outcoming:
			logInfo.Println("Считано сообщение из outcoming: ", outcomingMsg)
			err := chromedp.Run(ctx,
				chromedp.SendKeys(`//*[@id="editable-message-text"]`, outcomingMsg.Text, chromedp.NodeVisible),
				chromedp.Sleep(500*time.Millisecond),
				chromedp.Click(`//*[@id="MiddleColumn"]/div[4]/div[3]/div/div[2]/div[1]/button`, chromedp.NodeVisible),
			)
			if err != nil {
				logErr.Println(err)
			}
			logInfo.Println("Outcoming message: ", outcomingMsg.Text, "has been sended")
		default:
			err := chromedp.Run(ctx,
				chromedp.Evaluate(`
            (() => {
                const elements = Array.from(document.querySelectorAll('[id^="message-"]'));
                if (elements.length === 0) return null;
                
                const nums = elements
                    .map(el => {
                        const match = el.id.match(/message-(\d+)/);
                        return match ? parseInt(match[1], 10) : null;
                    })
                    .filter(n => n !== null);
                
                return nums.length > 0 ? Math.max(...nums) : null;
            })()
        `, &maxN),
			)

			if err != nil {
				logErr.Fatalln("Ошибка поиска максимального n:", err)
			}

			if maxN == nil {
				logErr.Fatalln("Элементы не найдены или некорректные ID")
			}

			// Шаг 2: Собрать текст из элементов без with-outgoing-icon
			var message string
			for n := *maxN; n >= lastMsg.Id; n-- {
				var outMsg bool
				var exists bool
				var classAttr string
				xpath := fmt.Sprintf(`//*[@id="message-%d"]/div[3]/div/div[1]/div`, n)

				err = chromedp.Run(ctx,
					chromedp.Evaluate(
						fmt.Sprintf(`
                !!document.evaluate('%s', document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue
            `, xpath),
						&exists,
					),
				)

				if err != nil || !exists {
					continue
				}

				err = chromedp.Run(ctx,
					chromedp.AttributeValue(xpath, "class", &classAttr, &exists),
				)

				if err != nil {
					logErr.Println(err)
					continue
				}

				substrings := []string{"with-outgoing-icon", "own"}
				for _, substing := range substrings {
					if strings.Contains(classAttr, substing) {
						outMsg = true
					}
				}

				if outMsg {
					outMsg = false
					continue
				}

				if lastMsg.Id == n {
					skipMsg = true
					break
				}
				lastMsg.Id = n

				logInfo.Println("Обнаружено новое входящее сообщение")

				logInfo.Println("Сообщение входящее, его аттрибуты: " + classAttr)

				if strings.Contains(classAttr, "message-subheader") { // это ответ на сообщение
					xpath2 := fmt.Sprintf(`//*[@id="message-%d"]/div[3]/div/div[1]/div[2]`, n)
					err = chromedp.Run(ctx,
						chromedp.Text(xpath2, &message, chromedp.BySearch),
					)
					if err != nil {
						logErr.Fatalln(err)
					}
					if strings.Contains(classAttr, "Audio") { // Пришло гс ответом на сообщение
						message = convertVoice(ctx, n, logInfo, logErr, message, 1)
					}
					if strings.Contains(classAttr, "RoundVideo") { // Пришел кружок ответом на сообщение
						message = convertVoice(ctx, n, logInfo, logErr, message, 2)
					} else { // Текстовое сообщение ответом на сообщение
						message = message[:len(message)-6]
					}
					break
				} else { // это не ответ на сообщение
					err = chromedp.Run(ctx,
						chromedp.Text(xpath, &message, chromedp.BySearch),
					)
					if strings.Contains(classAttr, "Audio") { // Пришло гс
						message = convertVoice(ctx, n, logInfo, logErr, message, 1)
					}
					if strings.Contains(classAttr, "RoundVideo") { // Пришел кружок
						message = convertVoice(ctx, n, logInfo, logErr, message, 2)
					} else { // Текстовое сообщение
						message = message[:len(message)-6]
					}
					if err != nil {
						logErr.Fatalln(err)
					}
					break
				}
			}
			if skipMsg {
				skipMsg = false
				continue
			}
			fmt.Println("Обнаружено новое сообщение: ", message)
			lastMsg.Text = message
			incoming <- lastMsg
			logInfo.Println("В incoming отправлено: ", lastMsg)
			time.Sleep(timeSleep)
		}
	}
}

// typeMsg: 1 - гс, 2 - кружок
func convertVoice(ctx context.Context, id int, logInfo *log.Logger, logErr *log.Logger, msgTime string, typeMsg int) string {
	sleepTime, err := parseDuration(msgTime)
	var xpath string
	if err != nil {
		logErr.Println(err)
		sleepTime = time.Minute
	}
	switch typeMsg {
	case 1:
		xpath = fmt.Sprintf(`//*[@id="message-%d"]/div[3]/div/div[1]/div/div[2]/div/button`, id)
	case 2:
		xpath = fmt.Sprintf(`//*[@id="message-%d"]/div[3]/div/div[1]/div//button`, id)
	}
	xpath2 := fmt.Sprintf(`//*[@id="message-%d"]/div[3]/div/div[1]/p`, id)
	var message string
	err = chromedp.Run(ctx,
		chromedp.Click(xpath, chromedp.NodeVisible),
		chromedp.Sleep(sleepTime),
		chromedp.Text(xpath2, &message, chromedp.NodeVisible, chromedp.BySearch),
	)
	if err != nil {
		logErr.Println(err)
		return ""
	}
	return message
}

func parseDuration(s string) (time.Duration, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return 0, errors.New("некорректный формат, требуется 'мин:сек'")
	}

	min, err := strconv.Atoi(parts[0])
	if err != nil || min < 0 {
		return 0, errors.New("некорректные минуты")
	}

	sec, err := strconv.Atoi(parts[1])
	if err != nil || sec < 0 || sec >= 60 {
		return 0, errors.New("некорректные секунды")
	}

	return time.Duration(min)*time.Minute + time.Duration(sec)*time.Second, nil
}
