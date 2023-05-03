package main

import (
	"fmt"
	"github.com/jyz0309/notionAI-go/client"
	"os"
)

func main() {
	args := os.Args
	if len(args) == 5 {
		cli := client.NewClient(args[2], args[3])
		execute(args[1], args[4], cli)
		return
	}
	if len(args) == 6 {
		cli := client.NewProxyClient(args[2], args[3], args[5])
		execute(args[1], args[4], cli)
		return
	}
	fmt.Println("Usage: <option> <token> <spaceId> <prompt> <proxy>")
}

func execute(option string, q string, cli *client.NotionClient) {
	var val = ""
	var err error
	switch option {
	case "draft":
		val, err = cli.HelpMeDraft(q)
	case "translate_into_chinese":
		val, err = cli.HelpMeEdit(q, "Translate into Chinese")
	case "translate_into_english":
		val, err = cli.HelpMeEdit(q, "Translate into English")
	case "continue":
		val, err = cli.ContinueWriting(q)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(val)
}
