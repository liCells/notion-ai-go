package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/google/uuid"
	"github.com/jyz0309/notionAI-go/client"
	"github.com/jyz0309/notionAI-go/common"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	args := os.Args
	if len(args) == 5 {
		req := buildClient(args[1], args[4])
		cli := http.Client{}
		executeStream(args[3], args[2], cli, req)
		return
	}
	if len(args) == 6 {
		proxyAddress, err := url.Parse(args[5])
		if err != nil {
			panic(err)
		}
		cli := http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyAddress)}}
		req := buildClient(args[1], args[4])
		executeStream(args[3], args[2], cli, req)
		return
	}
	fmt.Println("Usage: <option> <token> <spaceId> <prompt> <proxy>")
}

// Deprecated: Use executeStream instead.
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

func buildClient(symbol string, q string) *client.NotionRequest {
	switch symbol {
	case "draft":
		return &client.NotionRequest{
			Context: client.HelpDraftContext{
				Prompt: q,
				Type:   common.HelpMeDraft,
			},
			Model: client.OpenAIModel,
		}
	case "translate_into_chinese":
		return &client.NotionRequest{
			Context: client.HelpEditContext{
				Prompt:       q,
				Type:         common.HelpMeEdit,
				SelectedText: "Translate into Chinese",
			},
			Model: client.OpenAIModel,
		}
	case "translate_into_english":
		return &client.NotionRequest{
			Context: client.HelpEditContext{
				Prompt:       q,
				Type:         common.HelpMeEdit,
				SelectedText: "Translate into English",
			},
			Model: client.OpenAIModel,
		}
	case "continue":
		return &client.NotionRequest{
			Context: client.ContinueWritingContext{
				PreviousContent: q,
				Type:            client.ContinueWriting,
			},
			Model: client.OpenAIModel,
		}
	default:
		return &client.NotionRequest{
			Context: client.HelpDraftContext{
				Prompt: q,
				Type:   common.HelpMeDraft,
			},
			Model: client.OpenAIModel,
		}
	}
}

func executeStream(spaceId string, token string, c http.Client, req *client.NotionRequest) error {
	req.ID = uuid.NewString()
	req.SpaceId = spaceId
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}
	r, err := http.NewRequest("POST", client.API, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	r.Header.Set("accept", "application/json")
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("User-Agent", client.UserAgent)
	r.Header.Set("cookie", fmt.Sprintf("token_v2=%v", token))
	resp, err := c.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			var aiResp *client.NotionAIResp
			err := json.Unmarshal([]byte(scanner.Text()), &aiResp)
			if err != nil {
				return err
			}
			if strings.Contains(aiResp.Completion, "\n") {
				robotgo.TypeStr(strings.Replace(aiResp.Completion, "\n", "\r", -1))
			} else {
				robotgo.TypeStr(aiResp.Completion)
			}
		}
	}
	return nil
}
