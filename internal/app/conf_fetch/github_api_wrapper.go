package conf_fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

type Wrapper struct {
	token    string
	username string
	req_repo *http.Request
	req_user *http.Request
	req_q *http.Request
}

type Name struct {
	Name string `json:"name"`
}

type GitTree struct {
    Tree []Path `json:"tree"`
}

type Path struct {
    Path string `json:"path"`
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

func (w *Wrapper) ConstructSearchQuery() {
	req, err := http.NewRequest(
		"GET",
        "https://api.github.com/search/",
		nil,
	)

    if err != nil {
        panic(err)
    }

	req.Header.Set("Authorization", "token " + w.token)
	w.req_q = req
}

func (w *Wrapper) FetchNames() ([]Name, error) {
	client := http.Client{}

	response, ok := client.Do(w.req_user)
	if ok != nil {
		return nil, ok
	}

	response_bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var res []Name

	err = json.Unmarshal(response_bytes, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (w *Wrapper) FetchConf(name string) error {
	client := http.Client{}

	req := w.req_repo
	req.URL = req.URL.JoinPath(name + "/tarball")

	resp, ok := client.Do(req)
	if ok != nil {
		return ok
	}

	out, err := os.Create(name + ".tar.gz")
	if err != nil {
		return err
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	exec.Command("mkdir", name).Run()
	exec.Command("tar", "-xf", name+".tar.gz", "-C", name, "--strip-components", "1").Run()
	return nil
}

func (w *Wrapper) FetchFileNames(repo string) ([]Path, error) {
    client := http.Client{}

    req := w.req_repo

    req.URL = req.URL.JoinPath(repo + "/git/trees/main")
    req.URL.RawQuery = "recursive=1"

    resp, ok := client.Do(req)
    if ok != nil {
        return nil, ok 
    }

	resp_bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

    var tree GitTree

    err = json.Unmarshal(resp_bytes, &tree)
    if err != nil {
        return nil, err
    }

    return tree.Tree, nil
}
 
