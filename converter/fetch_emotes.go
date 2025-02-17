package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const CDN_URL = "https://cdn.7tv.app/emote/"
const API_URL = "https://7tv.io/v3/emotes/"

type FileMetadata struct {
	filename      string
	size          int64
	frame_count   int64
	make2framegif bool
	guildId       string
}

type Emote struct {
	emote_name string
	animated   bool
	id         string
	metadata   FileMetadata
}

func get_emote(input inputData) ResponseElements {
	data := make_request(input.link)
	emote := parse_data(data, input)
	download_emote(emote)
	response := compress_emote(emote)
	return response

}

func parse_data(data []byte, input inputData) Emote {
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
		emotetmp.metadata.guildId = input.guildId
		emotearr = append(emotearr, emotetmp)
	}
	for i := 0; i < len(emotearr); i++ {
		if i == 0 {
			if emotearr[0].metadata.size < TARGET_SIZE {
				emote.emote_name = emotearr[i].emote_name
				emote.animated = emotearr[i].animated
				emote.id = emotearr[i].id
				emote.metadata.filename = emotearr[i].metadata.filename
				emote.metadata.size = emotearr[i].metadata.size
				emote.metadata.frame_count = emotearr[i].metadata.frame_count
				emote.metadata.make2framegif = false
				emote.metadata.guildId = input.guildId
				if input.desiredNamed != "" {
					emote.emote_name = input.desiredNamed
				}
				if input.is2framegif {
					emote.metadata.make2framegif = true
				}
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
				emote.metadata.guildId = input.guildId
				if input.desiredNamed != "" {
					emote.emote_name = input.desiredNamed
				}
				if input.is2framegif {
					emote.metadata.make2framegif = true
				}
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
	emote.metadata.guildId = input.guildId
	if input.desiredNamed != "" {
		emote.emote_name = input.desiredNamed
	}
	if input.is2framegif {
		emote.metadata.make2framegif = true
	}
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
	os.MkdirAll(filepath.Join("to-convert", emote.metadata.guildId), 0700)
	fileName := filepath.Join("to-convert", string(emote.metadata.guildId), emote.emote_name+emote.metadata.filename)
	emote_url := CDN_URL + emote.id + "/" + emote.metadata.filename
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
	resp, err := client.Get(emote_url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)

	defer file.Close()

	fmt.Printf("Downloaded a file %s with size %d \n", fileName, size)
}
