package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"proxy/client"
)

type userArr struct {
	elem map[string]bool
}

func (u *userArr) new()  {
	if len(u.elem) == 0 {
		originArr := []string{
			"class",
			"href",
			"id",
		}
		u.elem = make(map[string]bool)
		for _, t :=range originArr{
			u.elem[t] = true
		}
	}
}


func (u *userArr) includes(t string) bool{
	return u.elem[t]
}

func _check(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func _checkErr(err error) int {
	if err != nil {
		return 500
	}
	return 200
}

func changeQuery(query url.Values) string {
	//Временно уберем
	/*query.Set("tracking", "jobsgalore")
	query.Set("utm_source", "jobsgaloreeu")
	query.Set("utm_campaign", "jobsgaloreeu")
	query.Set("utm_medium", "organic")
	query.Set("source", "jobsgaloreeu")*/
	return query.Encode()
}

func updateUrl(urlPath string) string{
	log.Println( "Старый урл " + urlPath)
	result := urlPath
	u, err := url.Parse(urlPath)
	_check(err)

	if u != nil {
		u.RawQuery = changeQuery(u.Query())
		result = u.String()
	}

	log.Println( "Новый урл " + result)

	return result

}



func main() {

	httpClient := client.HttpClient{}

	m := gin.Default()
	m.GET("/", func(c *gin.Context)  {
		c.String(http.StatusOK, "Working")
	})

	m.GET("/open", func(c *gin.Context) {
		gottenUrl := c.Query("url")
		result :=""
		if gottenUrl !="" {
			response, err := httpClient.Get(gottenUrl)
			log.Println(fmt.Sprintf("url: \"%s  | agent: \"%s", gottenUrl, httpClient.Agent()))
			doc, err := goquery.NewDocumentFromResponse(response)
			_check(err)

			u := userArr{}
			u.new()

			body := doc.Find("div.search-results-container, div.job_info, td#resultsCol, " +
				"div.jobsearch-JobComponent, div#jobResults").Each(func(i int, selection *goquery.Selection) {
				n := selection.Nodes
				for _, elem := range n {
					for _, at :=range elem.Attr{
						if !u.includes(at.Key){
							selection.RemoveAttr(at.Key)
						}
					}
				}
			})
			body.Find("svg, link, script, style").Remove()

			//removeExcess(body)
			result,_ = body.Html()

		}
		//c.String(http.StatusOK, result)
		c.Data(http.StatusOK, "text/html; charset=utf-8",[]byte(result))

	})

	m.GET("/redirect", func(c *gin.Context) {

		originUrl := c.Query("url")
		//aplitrak  := "www.aplitrak.com"
		status :=200
		res := ""
		if originUrl !="" {

			//Получаем страницу по url
			response, err := httpClient.Get(originUrl)
			log.Println(fmt.Sprintf("url: \"%s  | agent: \"%s", originUrl, httpClient.Agent()))

			status =_checkErr(err)

			//Если получили страницу, то продолжаем иначе у оригинального URL меняем метки и выходим
			if status==200 {
				/*new_url := ""
				//Если ссылается на аплитрак, то получаем с него ссылку
				if aplitrak == response.Request.URL.Host {
					log.Println("Aplitrak url: " + response.Request.URL.String())
					_, err := goquery.NewDocumentFromResponse(response)
					flag := _checkErr(err)
					if flag == 200 {
						new_url = response.Request.URL.String()
						/*body := doc.Find("div.search-results-container, div.job_info, td#resultsCol, " +
							"div.jobsearch-JobComponent, div#jobResults")
						newUrl :=
					} else {
						new_url = response.Request.URL.String()
					}
				} else{
					log.Println("Не литрак url: " + response.Request.URL.String())
					new_url = response.Request.URL.String()
				}

				//Проверяем полученный URL
				respChecking, err := httpClient.Get(new_url)
				_check(err)

				if respChecking.StatusCode == 200{
					log.Println("Ссылка работает : " + new_url)
					res = new_url
				} else {
					log.Println(fmt.Sprintf("url: \"%s  | agent: \"%s", new_url, httpClient.Agent()))
					log.Println("Ссылка не аботает : " + new_url)
					log.Println(respChecking.StatusCode)
					res = originUrl
				}
				res = new_url
				*/

				res = response.Request.URL.String()
			} else {
				res = originUrl
			}

			res = updateUrl(res)
			//Вставляем в  URL свои метки


		}

		c.PureJSON(status, gin.H{
			"uri":res,
			"body": nil,
		})
	})





	m.Run()
}
