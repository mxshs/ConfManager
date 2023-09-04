package conf_fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Wrapper struct {
    token string
    username string
    req_repo *http.Request
    req_user *http.Request
}

type RepoName struct {
    Name string `json:"name"`
}

func (w *Wrapper) SetToken(token string) {
    w.token = token
}

func (w *Wrapper) SetUsername(username string) {
    w.username = username
}

func (w *Wrapper) ConstructUserUrl() {
    req, _ := http.NewRequest(
        "GET",
        "https://api.github.com/user/repos",
        nil,
    )

    req.Header.Set("Authorization", "token " + w.token)
    w.req_user = req
}

func (w *Wrapper) ConstructRepoUrl() {
    req, _ := http.NewRequest(
        "GET",
        fmt.Sprintf("https://api.github.com/repos/%s/", w.username),
        nil,
    )

    req.Header.Set("Authorization", "token " + w.token)
    w.req_repo = req
}

func (w *Wrapper) FetchNames() ([]RepoName, error) {
    client := http.Client{}

    res, ok := client.Do(w.req_user)
    if ok != nil {
        return nil, ok
    }

    res_bytes, err:= io.ReadAll(res.Body)
    if err != nil {
        return nil, err 
    }

    var resp []RepoName
    _ = json.Unmarshal(res_bytes, &resp)
    
    return resp, nil 
}

func (w *Wrapper) FetchConf(name string) error {
    client := http.Client{}

    req := w.req_repo
    req.URL = req.URL.JoinPath(name)

    res, ok := client.Do(req)
    if ok != nil {
        return ok
    }
    
    res_bytes, err:= io.ReadAll(res.Body)
    if err != nil {
        return err
    }

    fmt.Printf(string(res_bytes))

    return nil
}

