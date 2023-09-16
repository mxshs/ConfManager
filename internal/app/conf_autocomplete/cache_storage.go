package confautocomplete

import (
	"bufio"
	"os"
	"time"

	"confmanager/internal/app/conf_fetch"
)

const TIME_FORMAT string = "Mon Jan _2 15:04:05 2006"

type Cache interface {
	Open(string) (Cache, error)
	ReadCache() ([]string, error)
	checkValidCache() bool
	updateCache() ([]string, error)
}

type FileCache struct {
	path   string
	handle *os.File
	reader *bufio.Reader
}

func (fc *FileCache) Open(path string) (Cache, error) {
	handle, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(handle)

	fc.path = path
	fc.handle = handle
	fc.reader = reader

	return fc, nil
}

func (fc *FileCache) ReadCache() ([]string, error) {
	validity := fc.checkValidCache()

	if !validity {
		return fc.updateCache()
	} else {
		comp := []string{}

		for {
			name, err := fc.reader.ReadString('\n')
			comp = append(comp, name)
			if err != nil {
				break
			}
		}

		return comp, nil
	}
}

func (fc *FileCache) checkValidCache() bool {
	cacheTime, err := fc.reader.ReadString('\n')
	if err != nil {
		return false
	} else {
		cacheTime = cacheTime[:len(cacheTime)-1]
	}

	cacheTimeAsTime, err := time.Parse(TIME_FORMAT, cacheTime)
	if err != nil {
		return false
	}

	if time.Since(cacheTimeAsTime) > 15*time.Second {
		return false
	}

	return true
}

func (fc *FileCache) updateCache() ([]string, error) {
	err := os.Truncate(fc.path, 0)
	if err != nil {
		return nil, err
	}

	_, err = fc.handle.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	names, err := conf_fetch.W.FetchNames()
	if err != nil {
		return nil, err
	}

	res, err := conf_fetch.W.UnmarshalNames(names)
	if err != nil {
		return nil, err
	}

	comp := []string{}
	to_write := []byte(time.Now().UTC().Format(TIME_FORMAT) + "\n")

	for _, name := range res {
		temp := []byte(name.Name + "\n")
		to_write = append(to_write, temp...)
		comp = append(comp, name.Name)
	}

	fc.handle.Write(to_write)

	return comp, nil
}
