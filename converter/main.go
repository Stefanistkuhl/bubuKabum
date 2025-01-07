package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	get_emote("https://7tv.app/emotes/01F6NTBCYG000B70V1XA8K6W4P")
	get_emote("https://7tv.app/emotes/01F6NACCD80006SZ7ZW5FMWKWK")

}

func get_emote(emote_url string) {
	emote_api_url := change_url(emote_url)
	data := make_request(emote_api_url)
	emote_name, size, frame_count, is_animated := parse_data(data)
	fmt.Println(size, frame_count)
	download_url := make_download_url(emote_url, is_animated)
	fmt.Println(download_url)
	download_emote(download_url, emote_name, is_animated)
}
func parse_data(data []byte) (string, int64, int64, bool) {
	emote_name := gjson.GetBytes(data, "name")
	fmt.Println("emote name: ", emote_name)
	state := gjson.GetBytes(data, "animated")
	if state.Bool() == true {
		fmt.Println("this emote is animated")
		size := gjson.GetBytes(data, "host.files.11.size")
		frame_count := gjson.GetBytes(data, "host.files.11.frame_count")
		fmt.Println("size:", size)
		fmt.Println("frame_count:", frame_count)
		return emote_name.String(), size.Int(), frame_count.Int(), state.Bool()
	} else {
		fmt.Println("this emote is not animated")
		return emote_name.String(), 0, 0, state.Bool()
	}
}

func change_url(input string) string {
	split := strings.Split(input, "/")
	base_url := "https://7tv.io/v3/emotes/"
	url := base_url + split[len(split)-1]
	return url
}

func make_request(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func download_emote(emote_url string, emote_name string, is_animated bool) {
	// Build fileName from fullPath
	fileName := ""
	if is_animated {
		fileName = emote_name + ".gif"
	} else {
		fileName = emote_name + ".png"
	}

	// Create blank file
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	// Put content on file
	resp, err := client.Get(emote_url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)

	defer file.Close()

	fmt.Printf("Downloaded a file %s with size %d", fileName, size)
}
func make_download_url(emote_url string, is_animated bool) string {
	if is_animated {
		split := strings.Split(emote_url, "/")
		base_url := "https://cdn.7tv.app/emote/"
		download_url := base_url + split[len(split)-1] + "/4x.gif"
		return download_url
	} else {
		split := strings.Split(emote_url, "/")
		base_url := "https://cdn.7tv.app/emote/"
		download_url := base_url + split[len(split)-1] + "/4x.png"
		return download_url
	}
}
