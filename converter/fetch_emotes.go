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

func get_emote(emote_url string) {
	data := make_request(emote_url)
	emote := parse_data(data)
	fmt.Println(emote)
	download_emote(emote)
	compress_emote(emote)
}

func parse_data(data []byte) Emote {
	var emote Emote
	file := gjson.GetBytes(data, "host.files.11")
	emote.emote_name = gjson.GetBytes(data, "name").String()
	emote.animated = gjson.GetBytes(data, "animated").Bool()
	emote.id = gjson.GetBytes(data, "id").String()
	emote.metadata.filename = file.Get("name").String()
	emote.metadata.size = file.Get("size").Int()
	emote.metadata.frame_count = file.Get("frame_count").Int()
	return emote
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

func download_emote(emote Emote) {
	fileName := "./to-convert/" + emote.emote_name + emote.metadata.filename
	emote_url := CDN_URL + emote.id + "/" + emote.metadata.filename
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
