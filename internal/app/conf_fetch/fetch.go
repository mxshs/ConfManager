package conf_fetch

import (
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

var GITHUB_TOKEN string
var GITHUB_USERNAME string
var W *Wrapper

func SetEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	GITHUB_TOKEN = os.Getenv("TOKEN")
	GITHUB_USERNAME = os.Getenv("USERNAME")
}

func init() {
	SetEnv()

	W = &Wrapper{}
	W.SetToken(GITHUB_TOKEN)
	W.SetUsername(GITHUB_USERNAME)
	W.ConstructRepoUrl()
	W.ConstructUserUrl()
}

func FetchRepo(name string) error {
	err := W.FetchConf(name)
	if err != nil {
		return err
	}

	exec.Command("mkdir", name).Run()
	exec.Command("tar", "-xf", name+".tar.gz", "-C", name, "--strip-components", "1").Run()
	return nil
}
