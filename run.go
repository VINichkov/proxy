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

type user_arr struct {
	elem map[string]bool
}

func (u *user_arr) new()  {
	if len(u.elem) == 0 {
		origin_arr := []string{
			"class",
			"href",
			"id",
		}
		u.elem = make(map[string]bool)
		for _, t :=range origin_arr{
			u.elem[t] = true
		}
	}
}

var u user_arr

func (u *user_arr) includes(t string) bool{
	return u.elem[t]
}

func _check(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func _check_err(err error) int {
	if err != nil {
		return 500
	}
	return 200
}

func changeQuery(query url.Values) string {
	query.Set("tracking", "jobsgalore")
	query.Set("utm_source", "jobsgaloreeu")
	query.Set("utm_campaign", "jobsgaloreeu")
	query.Set("utm_medium", "organic")
	return query.Encode()
}

func updateUrl(url_path string) string{
	log.Println( "Старый урл " + url_path)
	result := url_path
	u, err := url.Parse(url_path)
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
		url := c.Query("url")
		result :=""
		if url !="" {
			response, err := httpClient.Get(url)
			log.Println(fmt.Sprintf("url: \"%s  | agent: \"%s", url, httpClient.Agent()))
			doc, err := goquery.NewDocumentFromResponse(response)
			_check(err)

			u := user_arr{}
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

		origin_url := c.Query("url")
		//aplitrak  := "www.aplitrak.com"
		status :=200
		res := ""
		if origin_url !="" {

			//Получаем страницу по url
			response, err := httpClient.Get(origin_url)
			log.Println(fmt.Sprintf("url: \"%s  | agent: \"%s", origin_url, httpClient.Agent()))

			status =_check_err(err)

			//Если получили страницу, то продолжаем иначе у оригинального URL меняем метки и выходим
			if status==200 {
				/*new_url := ""
				//Если ссылается на аплитрак, то получаем с него ссылку
				if aplitrak == response.Request.URL.Host {
					log.Println("Aplitrak url: " + response.Request.URL.String())
					_, err := goquery.NewDocumentFromResponse(response)
					flag := _check_err(err)
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
					res = origin_url
				}
				res = new_url
				*/

				res = response.Request.URL.String()
			} else {
				res = origin_url
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
