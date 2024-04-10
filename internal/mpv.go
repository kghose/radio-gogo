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
	Data       any
	Error      string
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

	_, err = conn.Write(req_str)
	if err != nil {
		resp.Error = err.Error()
		return resp
	}
	_, err = conn.Write([]byte("\n"))
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
