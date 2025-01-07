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

const CDN_URL = "https://cdn.7tv.app/emote/"
const API_URL = "https://7tv.io/v3/emotes/"

type FileMetadata struct {
	filename    string
	size        int64
	frame_count int64
}

type Emote struct {
	emote_name string
	animated   bool
	id         string
	metadata   FileMetadata
}

func main() {
	get_emote("https://7tv.app/emotes/01F6NTBCYG000B70V1XA8K6W4P")
	get_emote("https://7tv.app/emotes/01F6NACCD80006SZ7ZW5FMWKWK")

}

func get_emote(emote_url string) {
	data := make_request(emote_url)
	emote := parse_data(data)
	fmt.Println(emote)
	download_emote(emote.id, emote.emote_name, emote.metadata.filename)
}

func parse_data(data []byte) *Emote {
	emote := new(Emote)
	emote_name := gjson.GetBytes(data, "name").String()
	id := gjson.GetBytes(data, "id").String()
	fmt.Println("emote name: ", emote_name)
	file := gjson.GetBytes(data, "host.files.11")
	file_extension := file.Get("name").String()
	state := gjson.GetBytes(data, "animated").Bool()
	if state == true {
		size := file.Get("size").Int()
		frame_count := file.Get("frame_count").Int()
		emote.emote_name = emote_name
		emote.animated = state
		emote.id = id
		emote.metadata.size = size
		emote.metadata.filename = file_extension
		emote.metadata.frame_count = frame_count
		return emote

	} else {
		size := file.Get("size").Int()
		frame_count := file.Get("frame_count").Int()
		fmt.Println("this emote is not animated")
		emote.emote_name = emote_name
		emote.animated = state
		emote.id = id
		emote.metadata.size = size
		emote.metadata.filename = file_extension
		emote.metadata.frame_count = frame_count
		return emote
	}
}

func make_request(url string) []byte {
	split := strings.Split(url, "/")
	url = API_URL + split[len(split)-1]
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

func download_emote(emote_id string, emote_name string, file_extension string) {
	fileName := emote_name + file_extension
	emote_url := CDN_URL + emote_id + "/" + file_extension
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

	fmt.Printf("Downloaded a file %s with size %d \n", fileName, size)
}
