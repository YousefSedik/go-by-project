package main

import (
	"fmt"
	"github.com/playwright-community/playwright-go"
	"log"
	"net/http"
	"os"
	"strings"
)

// 0 - not in queue, not visited, 1 - in queue, not visited, 2 - visited, not in queue
var visited = make(map[string]int8)

type Queue []string

func (q *Queue) Enqueue(item string) {
	*q = append(*q, item)
}

func (q *Queue) Dequeue() (string, bool) {
	if len(*q) == 0 {
		return "", false
	}
	item := (*q)[0]
	*q = (*q)[1:]
	return item, true
}

func isURLDead(url string) (bool, error) {
	/* A dead link is defined as one that returns a status code in the range of 4xx or 5xx. */
	response, err := http.Get(url)
	if err != nil {
		return true, err
	}
	if 400 <= response.StatusCode && response.StatusCode <= 599 {
		return true, nil
	}
	return false, nil
}

func getDomain(url string) (string, error) {
	/* Get the domain of the URL */
	if !strings.HasPrefix(url, "http") {
		return "", fmt.Errorf("URL must start with http or https")
	}
	parts := strings.Split(url, "/")
	if len(parts) < 3 {
		return "", fmt.Errorf("URL is not valid")
	}
	return parts[2], nil
}

func buildURL(domain, url string) string {
	/* Build the URL from the domain and the URL */
	if strings.HasPrefix(url, "http") {
		return url
	}
	if strings.HasPrefix(url, "/") {
		return "http://" + domain + url
	}
	return "http://" + domain + "/" + url
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: main <url>")
		return
	}
	url := os.Args[1]
	domain, err := getDomain(url)
	if err != nil {
		log.Fatalf("could not get url prefix: %v", err)
	}
	bad_urls, not_bad_urls := make([]string, 0), make([]string, 0)
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch()
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	queue := Queue{}
	queue.Enqueue(url)
	for len(queue) != 0 {
		// time.Sleep(1 * time.Second) // Sleep for 1 second to avoid overwhelming the server
		url_to_check, _ := queue.Dequeue()
		visited[url_to_check] = 2
		fmt.Printf("Checking %s\n", url_to_check)
		if id_dead, _ := isURLDead(url_to_check); id_dead == true {
			fmt.Printf(" %s is bad \n", url_to_check)
			bad_urls = append(bad_urls, url_to_check)
			continue
		}
		not_bad_urls = append(not_bad_urls, url_to_check)
		if _, err = page.Goto(url_to_check); err != nil {
			log.Fatalf("could not goto: %v", err)
		}
		links_locators, err := page.Locator("a").All()
		if err != nil {
			log.Fatalf("could not get entries: %v", err)
		}
		links := []string{}
		for _, link := range links_locators {
			href, err := link.GetAttribute("href")
			if err == nil && href != "" {
				href = buildURL(domain, href)
				links = append(links, href)
			}
		}
		for _, link := range links {
			if link_domain, _ := getDomain(link); link_domain != domain {
				fmt.Printf("Skipping %s, not the same domain\n", link)
				continue
			} else if visited[link] == 0 {
				visited[link] = 1
				queue.Enqueue(link)
			} else {
				fmt.Printf("Skipping %s, already visited\n", link)
			}
		}
	}
	fmt.Println("Task Finished.")
	fmt.Println("Bad URLs:", bad_urls)
	fmt.Println("Not Bad URLs:", not_bad_urls)
	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
}
