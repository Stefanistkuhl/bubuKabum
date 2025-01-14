package main

import (
	"fmt"
	"sync"
)

func main() {
	urls := []string{
		// "https://7tv.app/emotes/01F6NTBCYG000B70V1XA8K6W4P",
		// "https://7tv.app/emotes/01F6NACCD80006SZ7ZW5FMWKWK",
		"https://7tv.app/emotes/01J22JMMVG000EE58Q2GN4M3ZW",
	}

	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			get_emote(url)
		}(url)
	}

	wg.Wait()
	fmt.Println("Done")
}
