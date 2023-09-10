package conf_fetch

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/joho/godotenv"
)

type cache struct {
    arr []string
    time time.Time
}

const CACHE_INVAL = 5 * time.Minute

var GITHUB_TOKEN string
var GITHUB_USERNAME string
var W *Wrapper
var Cache *cache

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

    fmt.Println(GITHUB_TOKEN)
    W = &Wrapper{}
    W.SetToken(GITHUB_TOKEN)
    W.SetUsername(GITHUB_USERNAME)
    W.ConstructRepoUrl()
    W.ConstructUserUrl()

    Cache = &cache{}
    updateCache()
}

func GetNamesForAutocomp() ([]string, error) {
    if time.Since(Cache.time) < CACHE_INVAL {
        fmt.Println("cache hit")
        return Cache.arr, nil
    }

    fmt.Println("cache outdated")
    err := updateCache()

    return Cache.arr, err 
}

func FetchRepo(name string) error {
    err := W.FetchConf(name)
    if err != nil {
        return err
    }

    exec.Command("mkdir", name).Run()
    exec.Command("tar", "-xf", name + ".tar.gz", "-C", name, "--strip-components", "1").Run()
    return nil
}

func updateCache() error {
    Cache.time = time.Now()

    names, err := W.FetchNames()
    if err != nil {
        return err
    }

    var name_strings []string

    for _, name := range names {
        name_strings = append(name_strings, name.Name)
    }

    Cache.arr = name_strings
    
    return nil
}
