// Extract shop URLs from business pages.

package main

import "flag"
import "fmt"
import "io/ioutil"
import "log"
import "path"
import "regexp"
import "strings"
import "time"

import "github.com/lxhuang/geodata/file"
import "golang.org/x/net/html"

var FLAGS_crawl_city string
var FLAGS_entrypage_folder string
var FLAGS_nextpage_folder string
var FLAGS_shopurl_folder string

var validShopUrl = regexp.MustCompile("^/shop/[0-9]+$")

func match(url string) bool { return validShopUrl.MatchString(url) }

func traverse(node *html.Node, shop *map[string]string) {
	if node == nil {
		return
	}
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attribute := range node.Attr {
			if attribute.Key == "href" && match(attribute.Val) {
				(*shop)[attribute.Val] = ""
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		traverse(c, shop)
	}
}

func ExtractShopUrl(filepath string) []string {
	var result []string = make([]string, 0, 100)
	var shops map[string]string = make(map[string]string)

	content := file.ReadFile(filepath)
	// Note: content must be UTF-8 encoded HTML
	root, err := html.Parse(strings.NewReader(content))
	if err != nil {
		log.Printf("%s: invalid UTF8 HTML", filepath)
		return result
	}

	traverse(root, &shops)

	for k, _ := range shops {
		result = append(result, k)
	}

	return result
}

func init() {
	flag.StringVar(&FLAGS_crawl_city, "city", "beijing", "Specify which city to crawl.")
	flag.StringVar(&FLAGS_entrypage_folder, "entrypage_folder", "data/dianping/entrypage/",
		"Specify the folder of entry pages.")
	flag.StringVar(&FLAGS_nextpage_folder, "nextpage_folder", "data/dianping/nextpage/",
		"Specify the folder of nextpage files.")
	flag.StringVar(&FLAGS_shopurl_folder, "shopurl_folder", "data/dianping/shopurl/",
		"Specify the folder of shopurl files.")
}

func main() {
	flag.Parse()

	shopurl_fileset := make([]string, 0, 100)
	files, _ := ioutil.ReadDir(FLAGS_entrypage_folder)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		shopurl_fileset = append(shopurl_fileset, path.Join(FLAGS_entrypage_folder, file.Name()))
	}
	files, _ = ioutil.ReadDir(FLAGS_nextpage_folder)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		shopurl_fileset = append(shopurl_fileset, path.Join(FLAGS_nextpage_folder, file.Name()))
	}

	log.Printf("Total file number: %d\n", len(shopurl_fileset))
	// key is the shop url, and value is the number of occurrence in the dataset.
	shopurl_set := make(map[string]int)
	for _, filepath := range shopurl_fileset {
		log.Printf("Parsing %s\n", filepath)
		urls := ExtractShopUrl(filepath)
		for _, v := range urls {
			_, ok := shopurl_set[v]
			if !ok { shopurl_set[v] = 1 }
			shopurl_set[v]++
		}
	}

	log.Printf("Total shop number: %d\n", len(shopurl_set))
	// Output the result to file.
	shopurl_result := make([]string, 0, 100)
	for k, v := range shopurl_set {
		shopurl_result = append(shopurl_result, fmt.Sprintf("%s\t%d", k, v))
	}
	now := time.Now()
	shopurl_result_filepath := path.Join(
		FLAGS_shopurl_folder,
		fmt.Sprintf("%s_%04d%02d%02d%02d.txt",
			FLAGS_crawl_city, now.Year(), now.Month(), now.Day(), now.Hour()))
	file.WriteFile(shopurl_result_filepath, strings.Join(shopurl_result, "\n"))
}


