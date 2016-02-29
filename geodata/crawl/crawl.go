package crawl

import "fmt"
import "io/ioutil"
import "log"
import "net/http"
import "net/url"

func Get(uri string, headers map[string]string) string {
	parsedUrl, err := url.Parse(uri)
	if err != nil {
		log.Fatal(err)
	}

	if parsedUrl.Scheme == "" { parsedUrl.Scheme = "http" }

	client := &http.Client{}
	req, err := http.NewRequest("GET", parsedUrl.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return ""
	}
	log.Printf("%s\n", uri)
	return fmt.Sprintf("%s", body)
}


