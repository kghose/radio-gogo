package mpv

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"time"
)

const (
	mpvSocket = "/tmp/mpv.sock"
	mpvWait   = 1 * time.Second
)

type MpvRequest struct {
	Command    []string `json:"command"`
	Request_id int      `json:"request_id"`
}

type MpvResponse struct {
	Request_id int
	Data       json.RawMessage
	Error      string
}

type MpvMetadata struct {
	Name        string `json:"icy-name"`
	Description string `json:"icy-description"`
	Genre       string `json:"icy-genre"`
	Title       string `json:"icy-title"`
}

type Player struct {
	url        string
	playing    bool
	request_id int
}

func (p *Player) Start() {
	p.request_id = 0
	cmd := exec.Command(
		"mpv",
		fmt.Sprintf("--input-ipc-server=%s", mpvSocket),
		"--idle=yes",
		"--profile=low-latency",
	)
	cmd.Start()
	time.Sleep(mpvWait) // mpv takes a tiny bit to startup and have the socket ready
}

func (p *Player) Play(url string) MpvResponse {
	p.url = url
	p.playing = true
	return p.command([]string{"loadfile", url})
}

func (p *Player) TogglePause() MpvResponse {
	p.command([]string{"cycle", "pause"}) // Most reliable way to pause/unpause
	r := p.command([]string{"get_property", "pause"})
	p.playing = string(r.Data) == "true"
	return r
}

func (p *Player) Meta() MpvMetadata {
	response := p.command([]string{"get_property", "metadata"})
	meta := MpvMetadata{}
	json.Unmarshal(response.Data, &meta)
	return meta
}

func (p *Player) Quit() MpvResponse {
	return p.command([]string{"quit"})
}

func (p *Player) command(cmd []string) (resp MpvResponse) {
	resp = MpvResponse{}
	conn, err := net.Dial("unix", mpvSocket)
	if err != nil {
		resp.Error = err.Error()
		return resp
	}
	defer conn.Close()

	request := MpvRequest{Command: cmd, Request_id: p.request_id}
	p.request_id++

	req_str, err := json.Marshal(request)

	_, err = conn.Write(append(req_str, []byte("\n")...))
	if err != nil {
		resp.Error = err.Error()
		return resp
	}

	resp = MpvResponse{}
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	data := scanner.Bytes()
	err = json.Unmarshal(data, &resp)
	if err != nil {
		resp.Error = err.Error()
		return resp
	}
	return resp

}
