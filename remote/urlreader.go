package remote

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	//"sync/atomic"
)

// var flags int32 = 0

// Read ipfs/http as json
func ReadToJson(path string) (map[string]interface{}, error) {
	//schema analysis
	r, err := toReader(path)
	if err != nil {
		return nil, err

	}

	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	result := map[string]interface{}{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return result, nil

}

// Read and download to directory
func ReadToFile(path string, file string) error {
	r, err := toReader(path)
	if err != nil {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	if err != nil {
		return err
	}
	return nil
}

func toReader(path string) (io.Reader, error) {
	gateway := "https://ipfs.io/ipfs/"
	path = path[len("ipfs://"):]
	url := gateway + path
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
