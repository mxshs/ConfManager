package confautocomplete

import (
	"bufio"
	"os"
	"strings"
	"time"
)

const TIME_FORMAT string = "Mon Jan _2 15:04:05 2006"

type Cache interface {
    Open(string) error
	ReadCache(func() ([]string, error)) ([]string, error)
	isValid() bool
	updateCache(func() ([]string, error)) ([]string, error)
}

type FileCache struct {
	path   string
	handle *os.File
	reader *bufio.Reader
}

func (fc *FileCache) Open(path string) error {
	handle, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(handle)

	fc.path = path
	fc.handle = handle
	fc.reader = reader

	return nil
}

func (fc *FileCache) ReadCache(fetchFn func() ([]string, error)) ([]string, error) {
	if !fc.isValid() {
		return fc.updateCache(fetchFn)
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

func (fc *FileCache) isValid() bool {
	cacheTime, err := fc.reader.ReadString('\n')
	if err != nil {
		return false
	} else {
		cacheTime = cacheTime[:len(cacheTime) - 1]
	}

	cacheTimeAsTime, err := time.Parse(TIME_FORMAT, cacheTime)
	if err != nil {
		return false
	}

	if time.Since(cacheTimeAsTime) > 15 * time.Second {
		return false
	}

	return true
}

func (fc *FileCache) updateCache(fetchFn func() ([]string, error)) ([]string, error) {
	err := os.Truncate(fc.path, 0)
	if err != nil {
		return nil, err
	}

	_, err = fc.handle.Seek(0, 0)
	if err != nil {
		return nil, err
	}

    res, err := fetchFn()
    if err != nil {
        return nil, err
    }

	to_write := time.Now().UTC().Format(TIME_FORMAT) + "\n" + strings.Join(res, "\n")

	fc.handle.Write([]byte(to_write))

	return res, nil
}

