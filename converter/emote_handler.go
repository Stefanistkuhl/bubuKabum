package main

import (
	"sync"
)

type inputData struct {
	link         string
	is2framegif  bool
	desiredNamed string
	guildId      string
}

func process_emote_requests(request Request) Response {
	inputs := []inputData{}
	for i := range request.Links {
		var tmp inputData
		tmp.link = request.Links[i].Link
		tmp.is2framegif = request.Links[i].Is2FrameGif
		tmp.desiredNamed = request.Links[i].DesiredName
		tmp.guildId = request.Links[i].GuildID
		inputs = append(inputs, tmp)
	}

	response := Response{
		ResponseObject: make([]ResponseElements, len(inputs)),
	}
	var wg sync.WaitGroup

	for i, input := range inputs {
		wg.Add(1)
		go func(idx int, input inputData) {
			var responseData ResponseElements
			defer wg.Done()
			responseData = get_emote(input)
			response.ResponseObject[idx].Filename = responseData.Filename
			response.ResponseObject[idx].GuildId = responseData.GuildId
		}(i, input)
	}

	wg.Wait()
	return response
}
