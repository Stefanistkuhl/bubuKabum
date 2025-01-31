package main

import (
	"fmt"
	"sync"
)

type inputData struct {
	link         string
	is2framegif  bool
	desiredNamed string
}

func processEmoteRequests(request Request) Response {
	var response Response
	inputs := []inputData{}
	for i := range request.Links {
		fmt.Println(request.Links[i].Link)
		var tmp inputData
		tmp.link = request.Links[i].Link
		tmp.is2framegif = request.Links[i].Is2FrameGif
		tmp.desiredNamed = request.Links[i].DesiredName
		inputs = append(inputs, tmp)
	}
	fmt.Println(inputs)

	var wg sync.WaitGroup

	for _, input := range inputs {
		wg.Add(1)
		go func(input inputData) {
			defer wg.Done()
			get_emote(input)
		}(input)
	}

	wg.Wait()
	fmt.Println("Done")
	return response
}
