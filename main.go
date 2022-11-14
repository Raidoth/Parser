package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
)

type Purchase struct {
	Number string
	Name   string
}

var PurchasesSite1 = make([]Purchase, 0, 10)
var PurchasesSite2 = make([]Purchase, 0, 10)
var PurchasesSite3 = make([]Purchase, 0, 10)
var Client = &http.Client{}
var reqSite1, _ = http.NewRequest("GET", "https://www.fabrikant.ru/trades/procedure/search/", nil)
var reqSite2, _ = http.NewRequest("GET", "https://estp.ru/purchases", nil)
var reqSite3 *http.Request

func AddSite1(number, name string) {
	tmp := Purchase{Number: number, Name: name}
	PurchasesSite1 = append(PurchasesSite1, tmp)
}

func AddSite2(number, name string) {
	tmp := Purchase{Number: number, Name: name}
	PurchasesSite2 = append(PurchasesSite2, tmp)
}
func AddSite3(number, name string) {
	tmp := Purchase{Number: number, Name: name}
	PurchasesSite3 = append(PurchasesSite3, tmp)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Write find arg\n")
		os.Exit(1)
	} else {

		findName := os.Args[1]
		cntPages := getCountPageSite1(findName)
		for i := 2; i <= cntPages+1; i++ {
			getPurchasesSite1OnPage(i)
		}

		cntPages = getCountPageSite2(findName)
		for i := 2; i <= cntPages+1; i++ {
			getPurchasesSite2OnPage(i)
		}

		findName = "https://otc.ru/marketplace-b2b/query/" + os.Args[1]
		reqSite3, _ = http.NewRequest("GET", findName, nil)
		cntPages = getCountPageSite3(findName)
		for i := 2; i <= cntPages+1; i++ {
			getPurchasesSite3OnPage(i)
		}
		fmt.Printf("\nSearch %s on site https://www.fabrikant.ru/\n", os.Args[1])
		for i, val := range PurchasesSite1 {
			fmt.Println(i+1, "\tNumber: ", val.Number, "\n\tName: ", val.Name)
		}
		fmt.Printf("\nSearch %s on site https://estp.ru/purchases\n", os.Args[1])
		for i, val := range PurchasesSite2 {
			fmt.Println(i+1, "\tNumber: ", val.Number, "\n\tName: ", val.Name)
		}
		fmt.Printf("\nSearch %s on site https://otc.ru\n", os.Args[1])
		for i, val := range PurchasesSite3 {
			fmt.Println(i+1, "\tNumber: ", val.Number, "\n\tName: ", val.Name)
		}

	}

}

func getPurchasesSite1OnPage(page int) {
	q := reqSite1.URL.Query()
	q.Set("page", strconv.Itoa(page))
	reqSite1.URL.RawQuery = q.Encode()
	req2Site1, _ := http.NewRequest("GET", reqSite1.URL.String(), nil)
	req2Site1.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:106.0) Gecko/20100101 Firefox/106.0")
	resp, err := Client.Do(req2Site1)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer resp.Body.Close()
	htmlNods, _ := htmlquery.Parse(resp.Body)

	NumberPurchasesSite1 := htmlquery.Find(htmlNods, "//body//div[@class=\"col-xs-8\"]//div[@class=\"marketplace-unit__info__name\"]")
	NamePurchasesSite1 := htmlquery.Find(htmlNods, "//body//div[@class=\"col-xs-8\"]//div[@class=\"marketplace-unit__cut-wrap\"]//h4[@class=\"marketplace-unit__title\"]//a")
	var sb strings.Builder
	var sb2 strings.Builder
	cntSpace := 0
	for i := 0; i < len(NumberPurchasesSite1); i++ {
		for _, j := range htmlquery.InnerText(NumberPurchasesSite1[i]) {
			if string(j) == "\n" {
				continue
			}
			if string(j) == " " && cntSpace > 0 {
				continue
			}
			if string(j) == " " {
				sb.WriteRune(j)
				cntSpace++
				continue
			}
			sb.WriteRune(j)
			cntSpace = 0
		}
		cntSpace = 0
		for _, j := range htmlquery.InnerText(NamePurchasesSite1[i]) {
			if string(j) == "\n" {
				continue
			}
			if string(j) == " " && cntSpace > 0 {
				continue
			}
			if string(j) == " " {
				sb2.WriteRune(j)
				cntSpace++
				continue
			}
			sb2.WriteRune(j)
			cntSpace = 0
		}

		AddSite1(sb.String(), sb2.String())
		sb.Reset()
		sb2.Reset()
	}
	fmt.Printf("Page %d get site https://www.fabrikant.ru/\n", page)
}

func getCountPageSite1(find string) int {
	q := reqSite1.URL.Query()
	q.Add("types", "0")
	q.Add("procedure_stage", "0")
	q.Add("price_from", "")
	q.Add("price_to", "")
	q.Add("currency", "0")
	q.Add("date_type", "date_publication")
	q.Add("date_from", "")
	q.Add("date_to", "")
	q.Add("ensure", "date_publication")
	q.Add("count_on_page", "40")
	q.Add("query", find)
	q.Add("page", "1")
	reqSite1.URL.RawQuery = q.Encode()
	req2, _ := http.NewRequest("GET", reqSite1.URL.String(), nil)
	req2.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:106.0) Gecko/20100101 Firefox/106.0")
	resp, err := Client.Do(req2)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return 0
	}
	defer resp.Body.Close()
	htmlNods, _ := htmlquery.Parse(resp.Body)
	Pages := htmlquery.Find(htmlNods, "//body//div[@class=\"col-xs-8\"]//a[@class=\"pagination__lt__ref pagination__link\"]")

	Number := htmlquery.Find(htmlNods, "//body//div[@class=\"col-xs-8\"]//div[@class=\"marketplace-unit__info__name\"]")
	Name := htmlquery.Find(htmlNods, "//body//div[@class=\"col-xs-8\"]//div[@class=\"marketplace-unit__cut-wrap\"]//h4[@class=\"marketplace-unit__title\"]//a")
	var sb strings.Builder
	var sb2 strings.Builder
	cntSpace := 0
	for i := 0; i < len(Number); i++ {
		for _, j := range htmlquery.InnerText(Number[i]) {
			if string(j) == "\n" {
				continue
			}
			if string(j) == " " && cntSpace > 0 {
				continue
			}
			if string(j) == " " {
				sb.WriteRune(j)
				cntSpace++
				continue
			}
			sb.WriteRune(j)
			cntSpace = 0
		}
		cntSpace = 0
		for _, j := range htmlquery.InnerText(Name[i]) {
			if string(j) == "\n" {
				continue
			}
			if string(j) == " " && cntSpace > 0 {
				continue
			}
			if string(j) == " " {
				sb2.WriteRune(j)
				cntSpace++
				continue
			}
			sb2.WriteRune(j)
			cntSpace = 0
		}

		AddSite1(sb.String(), sb2.String())
		sb.Reset()
		sb2.Reset()
	}
	var cntPage int
	if Pages == nil {
		cntPage = 1
	} else {
		cntPage, _ = strconv.Atoi(htmlquery.InnerText(Pages[len(Pages)-1]))
	}

	fmt.Printf("Page %d get site https://www.fabrikant.ru/\n", 1)
	return cntPage

}

func getCountPageSite2(find string) int {
	q := reqSite2.URL.Query()
	q.Add("search", find)
	reqSite2.URL.RawQuery = q.Encode()
	req2Site2, _ := http.NewRequest("GET", reqSite2.URL.String(), nil)

	resp, err := Client.Do(req2Site2)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return 0
	}
	defer resp.Body.Close()
	htmlNods, _ := htmlquery.Parse(resp.Body)

	Pages := htmlquery.Find(htmlNods, "//body//div[@class=\"PurchasesContainer_pagination__EmeP5\"]//li[@class=\"Pagination_item__ETr8q\"]")
	Number := htmlquery.Find(htmlNods, "//body//div[@class=\"ListItem_number__423sR\"]//span[@class=\"AnnounceNumber_number__Ij3_7\"]")
	Name := htmlquery.Find(htmlNods, "//body//div[@class=\"ListItem_title__ktHjb\"]")

	for i := 0; i < len(Number); i++ {
		AddSite2(htmlquery.InnerText(Number[i]), htmlquery.InnerText(Name[i]))
	}

	var cntPage int
	if Pages == nil {
		cntPage = 1
	} else {
		cntPage, _ = strconv.Atoi(htmlquery.InnerText(Pages[len(Pages)-1]))
	}
	fmt.Printf("Page %d get site https://estp.ru/purchases\n", 1)
	return cntPage
}

var isSkip = false

func getPurchasesSite2OnPage(page int) {
	q := reqSite2.URL.Query()
	if isSkip {
		q.Set("page", strconv.Itoa(page))

	} else {
		q.Add("page", strconv.Itoa(page))
		isSkip = true
	}
	reqSite2.URL.RawQuery = q.Encode()
	req2Site2, _ := http.NewRequest("GET", reqSite2.URL.String(), nil)

	resp, err := Client.Do(req2Site2)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer resp.Body.Close()
	htmlNods, _ := htmlquery.Parse(resp.Body)
	Number := htmlquery.Find(htmlNods, "//body//div[@class=\"ListItem_number__423sR\"]//span[@class=\"AnnounceNumber_number__Ij3_7\"]")
	Name := htmlquery.Find(htmlNods, "//body//div[@class=\"ListItem_title__ktHjb\"]")

	for i := 0; i < len(Number); i++ {
		AddSite2(htmlquery.InnerText(Number[i]), htmlquery.InnerText(Name[i]))
	}
	fmt.Printf("Page %d get site https://estp.ru/purchases\n", page)
}

func getCountPageSite3(find string) int {
	resp, err := Client.Do(reqSite3)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return 0
	}
	defer resp.Body.Close()
	htmlNods, _ := htmlquery.Parse(resp.Body)
	Pages := htmlquery.Find(htmlNods, "//body//button[@class=\"Button ButtonSteel\"]")
	Number := htmlquery.Find(htmlNods, "//body//a[@class=\"name_2D3SN\"]/@href")
	Name := htmlquery.Find(htmlNods, "//body//a[@class=\"name_2D3SN\"]")

	for i := 0; i < len(Name); i++ {
		AddSite3(getId(htmlquery.InnerText(Number[i])), strings.TrimSpace(strings.TrimSuffix(htmlquery.InnerText(Name[i]), "\n")))
	}

	var cntPage int
	if Pages == nil {
		cntPage = 1
	} else {
		tmp := strings.TrimSpace(strings.TrimSuffix(htmlquery.InnerText(Pages[len(Pages)-2]), "\n"))

		cntPage, _ = strconv.Atoi(tmp)

	}
	fmt.Printf("Page %d get site https://otc.ru\n", 1)
	return cntPage
}

func getId(s string) string {
	var sb strings.Builder

	for _, val := range s {
		sb.WriteRune(val)
		if string(val) == "/" {
			sb.Reset()
		}
	}
	return sb.String()
}

func getPurchasesSite3OnPage(page int) {
	q := reqSite3.URL.Query()
	if isSkip {
		q.Set("p", strconv.Itoa(page))

	} else {
		q.Add("p", strconv.Itoa(page))
		isSkip = true
	}
	reqSite3.URL.RawQuery = q.Encode()
	req2Site3, _ := http.NewRequest("GET", reqSite3.URL.String(), nil)

	resp, err := Client.Do(req2Site3)
	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}
	defer resp.Body.Close()
	htmlNods, _ := htmlquery.Parse(resp.Body)
	Number := htmlquery.Find(htmlNods, "//body//a[@class=\"name_2D3SN\"]/@href")
	Name := htmlquery.Find(htmlNods, "//body//a[@class=\"name_2D3SN\"]")

	for i := 0; i < len(Name); i++ {
		AddSite3(getId(htmlquery.InnerText(Number[i])), strings.TrimSpace(strings.TrimSuffix(htmlquery.InnerText(Name[i]), "\n")))
	}
	fmt.Printf("Page %d get site https://otc.ru\n", page)
}
