package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"net/http"
	"log"
	"io/ioutil"
	"regexp"
	"net/url"
)

var TelegramToken = os.Getenv("telegram_token")

var bot, _ = tb.NewBot(tb.Settings{
	Token: TelegramToken,
})

var (
	receivedRe   = regexp.MustCompile(`"name": "Case Was Received[\S\s]+?barCategoryGap`)
	newCardRe    = regexp.MustCompile(`"name": "New Card Is Being[\S\s]+?barCategoryGap`)
	mailCardRe   = regexp.MustCompile(`"name": "Card Was Mai[\S\s]+?barCategoryGap`)
	pickCardRe   = regexp.MustCompile(`"name": "Card Was Pic[\S\s]+?barCategoryGap`)
	deliveCardRe = regexp.MustCompile(`"name": "Card Was Deli[\S\s]+?barCategoryGap`)
	yscRe        = regexp.MustCompile(`"YSC[\s\S]+?]`)
	statusRe     = regexp.MustCompile(`Your Current Status[\s\S]+?<span`)
)

func init() {
	//fmt.Println("token:", TelegramToken, "bot:", bot)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello from Jean</h1>")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	update := tb.Update{}
	if err := json.Unmarshal(body, &update); err != nil {
		fmt.Println("parse request error")
		return
	}
	message := update.Message
	if message == nil {
		fmt.Println("Message is nil, skip")
		return
	}
	fmt.Printf("request message: %#v\n", message)
	fmt.Println(message.Entities, len(message.Entities))

	if strings.Contains(message.Text, "opt") {
		optRes := opt()
		bot.Send(message.Chat, fmt.Sprintf("YSC1990018000:\nCase Was Received: %s\nNew Card Is Being Produced: %s\nCard Was Mailed To Me: %s\nCard Was Picked Up By The United States Postal Service: %s\nCard Was Delivered To Me By The Post Office: %s\n", optRes[0], optRes[1], optRes[2], optRes[3], optRes[4]))
		bot.Send(message.Chat, fmt.Sprintf("Your Current Status:\n%s", checkStatus("YSC1990018216")))
		return
	}

	if message.FromGroup() {
		if strings.Contains(message.Text, "小君") {
			fmt.Println("hit!!!!!!!!!")
			bot.Send(message.Chat, "小君君真好看！")
		}
	} else {
		bot.Send(message.Chat, "小君君真好看！")
	}

}

func opt() []string {

	resp, err := http.Get("http://www.checkuscis.com/en")
	//resp, err := http.Get("http://www.qq.com")
	if err != nil {
		fmt.Println("get url error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)

	rec := genList(receivedRe.FindAllString(bodyStr, -1)[1])
	newc := genList(newCardRe.FindAllString(bodyStr, -1)[1])
	mailc := genList(mailCardRe.FindAllString(bodyStr, -1)[1])
	pickc := genList(pickCardRe.FindAllString(bodyStr, -1)[1])
	deliv := genList(deliveCardRe.FindAllString(bodyStr, -1)[1])

	ysc := genList(yscRe.FindAllString(bodyStr, -1)[1])
	fmt.Println(rec, newc, mailc, pickc, deliv, ysc)
	fmt.Println(len(rec), len(newc), len(mailc), len(pickc), len(deliv), len(ysc))
	index := findIndex(ysc, "YSC1990018000")

	return []string{rec[index], newc[index], mailc[index], pickc[index], deliv[index]}
}

//func main() {
//	//optRes := opt()
//	//fmt.Println(fmt.Sprintf("YSC1990018000:\nCase Was Received: %s\nNew Card Is Being Produced: %s\nCard Was Mailed To Me: %s\nCard Was Picked Up By The United States Postal Service: %s\nCard Was Delivered To Me By The Post Office: %s\n", optRes[0], optRes[1], optRes[2], optRes[3], optRes[4]))
//
//	//fmt.Println(opt())
//	fmt.Println(checkStatus("YSC1990018216"))
//}

func genList(str string) []string {
	left := strings.Index(str, "[")
	right := strings.Index(str, "]")
	sp := strings.Split(str[left+1:right-1], ",")
	ans := make([]string, len(sp))
	for i, item := range sp {
		ans[i] = strings.TrimSpace(item)
	}
	fmt.Println(ans)
	return ans
}
func findIndex(ysc []string, name string) int {
	for i, item := range ysc {
		if strings.Contains(item, name) {
			return i
		}
	}
	return -1
}

func checkStatus(number string) string {

	form := url.Values{}
	form.Add("appReceiptNum", number)
	form.Add("initCaseSearch", "CHECK STATUS")
	form.Add("changeLocale", "")

	resp, err := http.Post("https://egov.uscis.gov/casestatus/mycasestatus.do", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return ""
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)

	res := statusRe.FindString(bodyStr)
	left := strings.Index(res, "strong")
	right := strings.Index(res, "span")
	fmt.Println(res[left+19 : right-5])
	return strings.TrimSpace(res[left+19 : right-5])
}
