package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var directDonation = "http://cure.pcrf-kids.org/site/TR/Events/General?px=1071909&pg=personal&fr_id=1050"

var shopDonations = []string{"https://www.ownaj.com/products/2018-vinesauce-is-hope-magazine",
	"https://www.ownaj.com/collections/limitededition/products/2018-vinesauce-is-hope-charity-pin",
	"https://www.ownaj.com/collections/limitededition/products/2018-vinesauce-is-hope-charity-shirt?variant=12407106928663", "https://www.ownaj.com/collections/limitededition/products/pcrf-vineshroom-plushy"}

var lastYear = float32(137366)
var lastDirect = 0
var beat = false
var beat2 = false

func main() {
	for {
		directAmount := getDirect()
		shopAmount := getShop()
		if directAmount < lastDirect {
			directAmount = lastDirect
		} else {
			lastDirect = directAmount
		}
		shopAmountHalf := shopAmount / 2
		total := shopAmount + float32(directAmount)
		totalHalfed := shopAmountHalf + float32(directAmount)
		fmt.Println(fmt.Sprintf("the total amount of donations with the shop not halfed is $%.2f\n", total))
		fmt.Println(fmt.Sprintf("the total amount of donations with the shop halfed is $%.2f\n", totalHalfed))
		fmt.Println(fmt.Sprintf("$%d is from the direct donations (which is pretty much 100 percent accurate)\n", directAmount))
		fmt.Println(fmt.Sprintf("$%.2f is from the shop purchases (without factoring in production costs)\n", shopAmount))
		fmt.Println(fmt.Sprintf("$%.2f is from the shop purchases (which attempts to factor in the production costs by dividing the total amount in half, this might not be 100 percent accurate but I'm doing my best with what I got)\n", shopAmountHalf))
		howMuchTillWePassLastYear(totalHalfed, true)
		fmt.Println("If we didn't half the shop then:")
		howMuchTillWePassLastYear(total, false)
		fmt.Println(fmt.Sprintf("Last updated at %s", time.Now()))
		fmt.Println("Checking amount again in 10 seconds!\n")
		time.Sleep(10 * time.Second)
	}
}

func getDirect() int {
	// make a request to the donation page
	res, err := http.Get(directDonation)
	if err != nil {
		log.Fatal(err)
	}
	// defer closing the body
	defer res.Body.Close()

	// load the document reader
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	amountText := doc.Find(".amount-raised-value").First().Text()
	amount, _ := strconv.Atoi(strings.Replace(strings.Replace(strings.TrimSpace(amountText), "$", "", -1), ",", "", -1))
	res.Body.Close()
	return amount
}

func getShop() float32 {
	total := float32(0)
	for _, shop := range shopDonations {
		// make a request to the shop page
		res, err := http.Get(shop)
		if err != nil {
			log.Fatal(err)
		}
		// defer closing the body
		defer res.Body.Close()

		// load the document reader
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		// find the cost
		costText := doc.Find(".money").First().Text()
		cost, _ := strconv.ParseFloat(strings.TrimSpace(strings.Replace(strings.TrimSpace(costText), "$", "", -1)), 32)
		// now find the total bought
		totalText := doc.Find(".cf-backertotal").First().Text()
		totalItemSpent, _ := strconv.ParseFloat(strings.TrimSpace(strings.Replace(strings.TrimSpace(totalText), "$", "", -1)), 32)
		total += float32(cost * totalItemSpent)
		res.Body.Close()
	}
	return total
}

func howMuchTillWePassLastYear(total float32, half bool) {
	amountToGo := lastYear - total
	if half {
		if amountToGo < 0 && half {
			if !beat {
				fmt.Println("WE BEAT LAST YEAR WOOHOOO!!!!!")
				fmt.Println(fmt.Sprintf("We beat last year at about %s", time.Now()))
				beat = true
			}
		} else {
			fmt.Println(fmt.Sprintf("(Half) We got about $%.2f to go until we beat last year!", amountToGo))
		}
	} else {
		if amountToGo < 0 {
			if !beat2 {
				fmt.Println("WE BEAT LAST YEAR WOOHOOO!!!!!")
				fmt.Println(fmt.Sprintf("We beat last year at about %s", time.Now()))
				beat2 = true
			}
		} else {
			fmt.Println(fmt.Sprintf("(Full) We got about $%.2f to go until we beat last year!", amountToGo))
		}
	}
}
