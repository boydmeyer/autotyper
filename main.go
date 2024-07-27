package main

import (
	"fmt"
	"strings"

	g "xabbo.b7c.io/goearth"
	"xabbo.b7c.io/goearth/shockwave/out"
)

var ext = g.NewExt(g.ExtInfo{
	Title:       "Autotyper",
	Description: "Automatically types a stored message by simply typing '!' in chat.",
	Version:     "1.0.0",
	Author:      "Nanobyte",
})

var storedMessage string

func main() {
	ext.Intercept(out.CHAT, out.SHOUT, out.WHISPER).With(handleChat)
	ext.Run()
}

func handleChat(e *g.Intercept) {
	msg := e.Packet.ReadString()
	fmt.Println("Received message:", msg)  // Debug print to see the raw message

	// Check if the message is the trigger to send the stored message
	if msg == "!" && storedMessage != "" && storedMessage != "!" {
		fmt.Println("Sending stored message:", storedMessage)  // Debug print to see the message being sent
		ext.Send(out.SHOUT, storedMessage)
		return
	}

	// Handle command messages
	if strings.HasPrefix(msg, ":") {
		// Remove the leading ':' from the message
		command := strings.TrimPrefix(msg, ":")
		
		// Split the command into the action and the content
		parts := strings.SplitN(command, " ", 2)
		if len(parts) < 2 {
			fmt.Println("Invalid command format:", command)
			return
		}
	
		// Determine the command action
		action := parts[0]
		content := parts[1]
		
		switch action {
		case "setmsg":
			e.Block()
			fmt.Println("Message updated to:", content)  // Debug print to see the updated message
			storedMessage = content
		default:
			fmt.Println("Unknown command:", action)
		}
	}
}
