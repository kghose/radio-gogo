package radio

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
)

const MPV_SOCKET = "/tmp/mpv.sock"

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
	cmd := exec.Command(
		"mpv",
		fmt.Sprintf("--input-ipc-server=%s", MPV_SOCKET),
		"--idle=yes",
		"--profile=low-latency",
	)
	cmd.Start()
}

func (p *Player) Play(url string) MpvResponse {
	p.url = url
	p.playing = true
	return p.command([]string{"loadfile", url})
}

func (p *Player) Pause() MpvResponse {
	if p.playing {
		p.playing = false
		return p.command([]string{"stop"})
	} else {
		return p.Play(p.url)
	}
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
	conn, err := net.Dial("unix", MPV_SOCKET)
	if err != nil {
		resp.Error = err.Error()
		return resp
	}

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
