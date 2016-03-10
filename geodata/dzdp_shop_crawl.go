// Crawl business information.

package main

import "crypto/md5"
import "flag"
import "fmt"
import "log"
import "math/rand"
import "os"
import "path"
import "strings"
import "time"

import "github.com/lxhuang/geodata/crawl"
import "github.com/lxhuang/geodata/file"

var FLAGS_parallelism int
var FLAGS_shopurl_file string
var FLAGS_shop_folder string

func init() {
	flag.StringVar(&FLAGS_shopurl_file, "shopurl_file", "data/dianping/shopurl/beijing_2016030611.txt",
		"The path of file containing shop URLs.")
	flag.StringVar(&FLAGS_shop_folder, "shop_folder", "data/dianping/shop/",
		"Specify the folder to dump shop details.")
	flag.IntVar(&FLAGS_parallelism, "parallelism", 2, "Specify the parallelism of crawling webpages.")
}

func main() {
	flag.Parse()

	urls := file.ReadLines(FLAGS_shopurl_file)
	log.Printf("Number of URLs: %d\n", len(urls))

	quit := make(chan string)
	queue := make(chan string, FLAGS_parallelism)
	go func() {
		for _, v := range urls {
			tokens := strings.Split(v, "\t")
			if len(tokens) < 1 {
				continue
			}
			log.Printf("Queuing %s\n", strings.TrimSpace(tokens[0]))
			queue <- strings.TrimSpace(tokens[0])
		}
		quit <- "done"
	}()

	var url string
	for {
		select {
		case url = <- queue:
			fetch(url)
		case <- quit:
			log.Printf("Crawling done")
			return
		}
	}
}

func fetch(uri string) {
	uri = "http://www.dianping.com" + uri

	dump_filename := fmt.Sprintf("%x.txt", md5.Sum([]byte(uri)))
	dump_filepath := path.Join(FLAGS_shop_folder, dump_filename)

	if _, err := os.Stat(dump_filepath); os.IsNotExist(err) {
		headers := make(map[string]string)
		headers["User-Agent"] = "Chrome/48.0.2564.109"
		log.Printf("Crawling %s\n", uri)
		content := crawl.Get(uri, headers)
		if len(content) == 0 {
			log.Printf("Error when crawling: %s\n", uri)
			return
		} else {
			file.WriteFile(dump_filepath, content)
		}
		time.Sleep(time.Duration(rand.Intn(5) + 5) * time.Second)
	} else {
		log.Printf("Skip %s\n", dump_filename)
	}
}
