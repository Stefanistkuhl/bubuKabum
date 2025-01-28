package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
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
	var emotetmp Emote
	var emotearr []Emote
	host := "host.files"
	for i := 11; i >= 8; i-- {
		host = "host.files." + strconv.Itoa(i)
		file := gjson.GetBytes(data, host)
		emotetmp.emote_name = gjson.GetBytes(data, "name").String()
		emotetmp.animated = gjson.GetBytes(data, "animated").Bool()
		emotetmp.id = gjson.GetBytes(data, "id").String()
		emotetmp.metadata.filename = file.Get("name").String()
		emotetmp.metadata.size = file.Get("size").Int()
		emotetmp.metadata.frame_count = file.Get("frame_count").Int()
		emotearr = append(emotearr, emotetmp)
	}
	fmt.Println(len(emotearr))
	for i := 0; i < len(emotearr); i++ {
		if i == 0 {
			if emotearr[0].metadata.size < TARGET_SIZE {
				emote.emote_name = emotearr[i].emote_name
				emote.animated = emotearr[i].animated
				emote.id = emotearr[i].id
				emote.metadata.filename = emotearr[i].metadata.filename
				emote.metadata.size = emotearr[i].metadata.size
				emote.metadata.frame_count = emotearr[i].metadata.frame_count
				return emote
			}
		} else {
			if emotearr[i].metadata.size > emotearr[i-1].metadata.size && emotearr[i].metadata.size < TARGET_SIZE {
				emote.emote_name = emotearr[i].emote_name
				emote.animated = emotearr[i].animated
				emote.id = emotearr[i].id
				emote.metadata.filename = emotearr[i].metadata.filename
				emote.metadata.size = emotearr[i].metadata.size
				emote.metadata.frame_count = emotearr[i].metadata.frame_count
				return emote
			}
		}
	}
	emote.emote_name = emotearr[len(emotearr)-1].emote_name
	emote.animated = emotearr[len(emotearr)-1].animated
	emote.id = emotearr[len(emotearr)-1].id
	emote.metadata.filename = emotearr[len(emotearr)-1].metadata.filename
	emote.metadata.size = emotearr[len(emotearr)-1].metadata.size
	emote.metadata.frame_count = emotearr[len(emotearr)-1].metadata.frame_count
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
