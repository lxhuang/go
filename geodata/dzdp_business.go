// Crawl the entry page to businesses.

package main

import "crypto/md5"
import "flag"
import "fmt"
import "io/ioutil"
import "log"
import "path"
import "sort"
import "strings"
import "time"

import "github.com/lxhuang/geodata/crawl"
import "github.com/lxhuang/geodata/file"

var FLAGS_crawl_city string
var FLAGS_category_folder string
var FLAGS_entrypage_folder string

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

func init() {
	flag.StringVar(&FLAGS_crawl_city, "city", "beijing", "Specify which city to crawl.")
	flag.StringVar(&FLAGS_category_folder, "category_folder", "data/dianping/category/",
		"Specify the folder of category files.")
	flag.StringVar(&FLAGS_entrypage_folder, "entrypage_folder", "data/dianping/entrypage/",
		"Specify the folder of entry pages.")
}

func main() {
	flag.Parse()

	categories := file.ReadLines(path.Join(FLAGS_category_folder, category_file()))
	if len(categories) == 0 {
		log.Fatal("dzdp_business.go: No Categories Found.")
	} else {
		log.Printf("Have %d categories in total\n", len(categories))
	}

	// Find existing files under data/dianping/entrypage/...
	existed_files := make(map[string]time.Time)
	files, _ := ioutil.ReadDir(FLAGS_entrypage_folder)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		existed_files[file.Name()] = file.ModTime()
	}
	log.Printf("Existed Files: %d\n", len(existed_files))

	for _, category := range categories {
		tokens := strings.Split(category, "\t")
		if len(tokens) == 0 {
			log.Fatal("dzdp_business.go: Invalid line in Category.")
		}

		url := strings.TrimSpace(tokens[0])
		if !strings.Contains(url, "http://www.dianping.com") {
			url = "http://www.dianping.com" + url
		}
		dump_filename = fmt.Sprintf("%x.txt", md5.Sum([]byte(url)))

		// Whether we've crawled the page or not.
		if _, ok := existed_files[dump_filename]; ok {
			log.Printf("Skip File: %s\n", dump_filename)
			continue
		}

		headers := make(map[string]string)
		headers["User-Agent"] = "Chrome/48.0.2564.109"
		content := crawl.Get(url, headers)
		if len(content) == 0 {
			log.Printf("Error when crawling: %s\n", url)
			continue
		}

		file.WriteFile(path.Join(FLAGS_entrypage_folder, dump_filename), content)
		existed_files[dump_filename] = time.Now()

		time.Sleep(5 * time.Second) // Friendly crawler.
	}
}
