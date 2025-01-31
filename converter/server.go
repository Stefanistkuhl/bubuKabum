package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmoteRequest struct {
	Link        string `json:"link"`
	Is2FrameGif bool   `json:"is_2_frame_gif"`
	DesiredName string `json:"desired_name"`
	GuildID     string `json:"guild_id"`
}

type Request struct {
	Links []EmoteRequest `json:"emotes"`
}
type Response struct {
	ResponseObject []ResponseElements `json:"emotes"`
}
type ResponseElements struct {
	Filename string `json:"filename"`
	GuildId  string `json:"guildId"`
}

func server() {
	r := gin.Default()
	r.POST("/api/emote", handleEmote)
	r.StaticFS("/converted", http.Dir("static/converted"))
	r.Run(":6999")
}
func handleEmote(c *gin.Context) {
	var request Request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}
	response := process_emote_requests(request)
	fmt.Println(response)
	// fmt.Printf("%+v\n", request)
	c.JSON(http.StatusOK, response)
}
