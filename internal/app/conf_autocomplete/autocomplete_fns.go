package confautocomplete

import (
	"confmanager/internal/app/conf_fetch"
)

func GetRepoNames() ([]string, error) {
    cache := &FileCache{}
    err := cache.Open("names.cache")
    if err != nil {
        return nil, err
    }

    res, err := cache.ReadCache(getNewRepoNames)
    if err != nil {
        return nil, err
    }

    return res, nil
}

func GetFileNames(repo string) ([]string, error) {
    names, err := conf_fetch.W.FetchFileNames(repo)
    if err != nil {
        return nil, err
    }

    comp := []string{}

    for _, name := range names {
        comp = append(comp, name.Path)
    }

    return comp, nil
}

func getNewRepoNames() ([]string, error) {
	names, err := conf_fetch.W.FetchNames()
	if err != nil {
		return nil, err
	}

	comp := []string{}

	for _, name := range names {
		comp = append(comp, name.Name)
	}

    return comp, nil
}
