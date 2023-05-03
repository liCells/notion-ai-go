package main

import (
    "fmt"
    "github.com/jyz0309/notionAI-go/client"
    "os"
)

func main() {
    args := os.Args
    if len(args) != 6 {
        fmt.Println("Usage: <option> <token> <spaceId> <prompt> <proxy>")
        return
    }
    execute(args[1], args[2], args[3], args[4], args[5])
}

func execute(option string, token string, spaceId string, q string, proxy string) {
    cli := client.NewProxyClient(token, spaceId, proxy)
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
