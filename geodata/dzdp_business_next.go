// Crawl all pages of business URLs.

package main

import "crypto/md5"
import "flag"
import "fmt"
import "io/ioutil"
import "log"
import "math/rand"
import "os"
import "path"
import "sort"
import "strconv"
import "strings"
import "time"

import "github.com/lxhuang/geodata/crawl"
import "github.com/lxhuang/geodata/file"
import "golang.org/x/net/html"

var FLAGS_entrypage_folder string
var FLAGS_category_folder string
var FLAGS_nextpage_folder string
var FLAGS_crawl_city string

func find_total_page_num(node *html.Node, total_page *int) {
	if node == nil {
		return
	}
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attribute := range node.Attr {
			if attribute.Key == "class" && strings.EqualFold(attribute.Val, "PageLink") {
				page_index, err := strconv.Atoi(strings.TrimSpace(node.FirstChild.Data))
				if err != nil {
					log.Printf("Invalid page number: %s\n", node.FirstChild.Data)
					continue
				}
				if page_index > *total_page {
					*total_page = page_index
				}
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		find_total_page_num(c, total_page)
	}
}

func init() {
	flag.StringVar(&FLAGS_entrypage_folder, "entrypage_folder", "data/dianping/entrypage/",
		"Specify the folder of entry pages.")
	flag.StringVar(&FLAGS_category_folder, "category_folder", "data/dianping/category/",
		"Specify the folder of category files.")
	flag.StringVar(&FLAGS_nextpage_folder, "nextpage_folder", "data/dianping/nextpage/",
		"Specify the folder of nextpage files.")
	flag.StringVar(&FLAGS_crawl_city, "city", "beijing", "Specify which city to crawl.")
}

func category_file() string {
	category_file_list := make([]string, 0, 100)

	files, _ := ioutil.ReadDir(FLAGS_category_folder)
	for _, file := range files {
		if strings.Contains(file.Name(), FLAGS_crawl_city) {
			category_file_list = append(category_file_list, file.Name())
		}
	}

	// Find the latest one.
	sort.Strings(category_file_list)
	return category_file_list[len(category_file_list)-1]
}

func main() {
	flag.Parse()

	rand.Seed(time.Now().Unix())

	categories := file.ReadLines(path.Join(FLAGS_category_folder, category_file()))
	if len(categories) == 0 {
		log.Fatal("dzdp_business.go: No Categories Found.")
	} else {
		log.Printf("Have %d categories in total\n", len(categories))
	}

	// Find existing files under data/dianping/nextpage/...
	existed_files := make(map[string]time.Time)
	files, _ := ioutil.ReadDir(FLAGS_nextpage_folder)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		existed_files[file.Name()] = file.ModTime()
	}
	log.Printf("Existed Files: %d\n", len(existed_files))

	var category_progress = 0
	for _, category := range categories {
		category_progress++
		fmt.Printf("Category progress: %d\n", category_progress)
		tokens := strings.Split(category, "\t")
		if len(tokens) == 0 {
			log.Fatal("dzdp_business_next.go: Invalid line in Category.")
		}

		url := strings.TrimSpace(tokens[0])
		if !strings.Contains(url, "http://www.dianping.com") {
			url = "http://www.dianping.com" + url
		}
		url_enc := fmt.Sprintf("%x", md5.Sum([]byte(url)))
		// Check the existence of entrypage file. It should exist,
		// otherwise, something is wrong. and we should skip it in such case.
		entrypage_file := path.Join(FLAGS_entrypage_folder, url_enc+".txt")
		if _, err := os.Stat(entrypage_file); os.IsNotExist(err) {
			log.Printf("%s not exist\n", entrypage_file)
			continue
		}

		// Parse the entrypage file, and find the total number of pages.
		log.Printf("Parsing %s\n", entrypage_file)
		content := file.ReadFile(entrypage_file)
		root, err := html.Parse(strings.NewReader(content))
		if err != nil {
			log.Fatalf("%s is invalid HTML\n", entrypage_file)
		}
		var total_page_num int = 1
		find_total_page_num(root, &total_page_num)

		// Crawl each page.
		for i := 2; i <= total_page_num; i++ {
			next_page_url := fmt.Sprintf("%sp%d", url, i)
			dump_filename := fmt.Sprintf("%x.txt", md5.Sum([]byte(next_page_url)))
			// Whether we've crawled the page or not.
			if _, ok := existed_files[dump_filename]; ok {
				log.Printf("Skip File: %s\n", dump_filename)
				continue
			}

			headers := make(map[string]string)
			headers["User-Agent"] = "Chrome/48.0.2564.109"
			content = crawl.Get(next_page_url, headers)
			if len(content) == 0 {
				log.Printf("Error when crawling: %s\n", next_page_url)
				continue
			}
			file.WriteFile(path.Join(FLAGS_nextpage_folder, dump_filename), content)

			existed_files[dump_filename] = time.Now()

			time.Sleep(time.Duration(rand.Intn(5) + 5) * time.Second)
		}
	}
}

