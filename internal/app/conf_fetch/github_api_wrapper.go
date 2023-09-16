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

func (w *Wrapper) FetchNames() ([]byte, error) {
	client := http.Client{}

	res, ok := client.Do(w.req_user)
	if ok != nil {
		return nil, ok
	}

	res_bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return res_bytes, nil
}

func (w *Wrapper) UnmarshalNames(nameBytes []byte) ([]RepoName, error) {
	var res []RepoName

	err := json.Unmarshal(nameBytes, &res)
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
