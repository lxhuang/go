package main

import "flag"
import "fmt"
import "log"
import "path"
import "strings"
import "time"

import "github.com/lxhuang/geodata/crawl"
import "github.com/lxhuang/geodata/file"
import "golang.org/x/net/html" // go get golang.org/x/net/html

var FLAGS_crawl_city string
var FLAGS_category_folder string

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func traverse(node *html.Node, link *map[string][]string) {
	if node == nil {
		return
	}
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attribute := range node.Attr {
			if attribute.Key == "href" && strings.Contains(attribute.Val, "/search/category/") {
				_, ok := (*link)[attribute.Val]
				if !ok {
					(*link)[attribute.Val] = make([]string, 0, 10)
				}
				var biz_type string = strings.TrimSpace(node.FirstChild.Data)
				if !contains((*link)[attribute.Val], biz_type) {
					(*link)[attribute.Val] = append((*link)[attribute.Val], biz_type)
				}
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		traverse(c, link)
	}
}

func init() {
	flag.StringVar(&FLAGS_crawl_city, "city", "beijing", "Specify which city to crawl.")
	flag.StringVar(&FLAGS_category_folder, "category_folder", "data/dianping/category/",
		"Specify the folder of category files.")
}

func main() {
	flag.Parse()

	content := crawl.Get("http://www.dianping.com/"+FLAGS_crawl_city, make(map[string]string))
	dump_filename := "tmp_dianping_" + FLAGS_crawl_city + ".txt"
	file.WriteFile(dump_filename, content)

	content = file.ReadFile(dump_filename)
	// Note: content must be UTF-8 encoded HTML.
	root, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return
	}
	// Find all <a> with href contains "/search/category/".
	seed_links := make(map[string][]string)
	traverse(root, &seed_links)

	// Print them out.
	fmt.Println(len(seed_links))
	if len(seed_links) < 150 {
		log.Fatal("Dianping Category Crawler is broken...")
	}

	// Find the current date.
	now := time.Now()
	now_string := path.Join(
		FLAGS_category_folder,
		fmt.Sprintf("%s_%04d%02d%02d%02d.txt",
			FLAGS_crawl_city, now.Year(), now.Month(), now.Day(), now.Hour()))

	for k, v := range seed_links {
		line := fmt.Sprintf("%s\t%s\n", k, v)
		file.AppendToFile(now_string, line)
	}
}
