package client

import (
	"math/rand"
	"net/http"
	"time"
)

type UserSlice struct {
	arr []string
}

func (array *UserSlice)sample() string  {
	rand.Seed(time.Now().Unix())
	result :=array.arr[(rand.Int() % len(array.arr))]
	return result
}

type HttpClient struct {
	agent string
	lang string
	client http.Client
}

func (client *HttpClient) newRequest(url string) (*http.Request, error)   {
	client.selectAgent()
	client.selectLang()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", client.agent)
	req.Header.Set("Accept-Language", client.lang)

	return req, nil

}

func (client *HttpClient) selectAgent()  {
	if client.agent == "" {
		agents := UserSlice{[]string{
			"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:64.0) Gecko/20100101 Firefox/64.0",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/600.1.25 (KHTML, like Gecko) Version/8.0 Safari/600.1.25",
			"Mozilla/5.0 (X11; Ubuntu; Linux x86_64) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/11.0 Safari/605.1.15 Epiphany/605.1.15",
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) snap Chromium/70.0.3538.110 Chrome/70.0.3538.110 Safari/537.36",
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) snap Chromium/70.0.3538.110 Chrome/70.0.3538.110 Safari/537.36",
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/602.1 (KHTML, like Gecko) Arora/0.11.0 Version/9.0 Safari/602.1",
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) snap Chromium/70.0.3538.110 Chrome/70.0.3538.110 Safari/537.36",
			"Mozilla/5.0 (X11; U; Linux i686; ru; rv:1.9b5) Gecko/2008050509 Firefox/3.0b5",
			"Mozilla/5.0 (X11; U; Linux i686; ru; rv:1.9b5) Gecko/2008050509 Firefox/3.0b5",
			"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:63.0) Gecko/20100101 Firefox/63.0",
			"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:64.0) Gecko/20100101 Firefox/64.0",
		}}

		client.agent = agents.sample()
	}
}

func (client *HttpClient) selectLang()  {
	if client.lang =="" {
		langs := UserSlice{[]string{
			"de-CH",
			"en-US,en;q=0.5",
			"fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5",
			"ru-RU, ru;q=0.9, en-US;q=0.8, en;q=0.7, fr;q=0.6",
		}}
		client.lang = langs.sample()
	}
}

func (client *HttpClient) Agent() string  {
	return client.agent
}

func (client *HttpClient) Get(url string) (*http.Response,  error) {

	req, err := client.newRequest(url)
	if err != nil {
		panic(err)
	}

	if &client.client == nil {
		client.client = http.Client{}
	}
	resp, err := client.client.Do(req)
	if err != nil {
		panic(err)
	}

	//defer resp.Body.Close()

	/*responseData,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	res := string(responseData)*/

	return resp, nil

}
