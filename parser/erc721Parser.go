package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"nftParser/remote"
	"nftParser/token/erc721"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	//"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Erc721Parser struct {
	Address string
	Name    string
}

func NewErc721Parser(address string, name string) *Erc721Parser {
	parser := Erc721Parser{}
	parser.Address = address
	parser.Name = name
	return &parser
}

func (p *Erc721Parser) DownloadImages(folder string, startToken, endToken, concurrency int) {
	defer remote.DoCleanup()
	downloadTo := folder
	_ = os.Mkdir(downloadTo, 0755)

	client, err := ethclient.Dial("https://mainnet.infura.io/v3/9b71514d49094c61b53c85287ebb4f57")

	if err != nil {
		log.Fatalf("Oops! There was a problem", err)
	} else {
		fmt.Println("Success! you are connected to the Ethereum Network")
	}

	address := common.HexToAddress(p.Address)
	instance, _ := erc721.NewToken(address, client)

	// Accesserc721 contract to find tokenUris
	var wg sync.WaitGroup
	wg.Add(endToken - startToken + 1)
	guard := make(chan struct{}, concurrency)
	for i := startToken; i <= endToken; i++ {
		guard <- struct{}{}
		fmt.Println("process", i)
		go func(tokenId int) {
			defer wg.Done()
			defer func() {
				<-guard
			}()

			uri, err := instance.TokenURI(nil, big.NewInt(int64(tokenId)))
			if err != nil {
				fmt.Println(err, tokenId)
				return
			}
			err = p.process(uri, tokenId, downloadTo)
			if err != nil {
				fmt.Println(err)
				return
			}
		}(i)
	}
	wg.Wait()
	fmt.Printf("%v fetch complete\n", p.Name)
}

func (p *Erc721Parser) process(url string, token int, downloadTo string) error {
	//fetch metadata
	var reference string
	var err error
	switch scheme := strings.Split(url, "://")[0]; scheme {
	case "ipfs":
		reference, err = fetchReferenceIpfs(url, token, downloadTo)
	case "http", "https":
		reference, err = fetchReferenceHttp(url, token, downloadTo)
	default:
		err = errors.New("unsupported scheme")
	}
	if err != nil {
		return err
	}
	//download image
	fpath := filepath.Join(downloadTo, fmt.Sprintf("%v_%v", p.Name, token))
	switch scheme := strings.Split(reference, "://")[0]; scheme {
	case "ipfs":
		err = downloadIpfs(reference, token, fpath)
	case "http", "https":
		err = downloadHttp(reference, token, fpath)
	default:
		err = errors.New("unsupported scheme")
	}
	if err != nil {
		return err
	}
	return nil
}

func fetchReferenceIpfs(url string, token int, downloadTo string) (string, error) {

	meta, err := remote.ReadToJson(url)
	if err != nil {
		return "", err
	}
	realImgAddr, ok := meta["image"]
	if !ok {
		return "", errors.New("no image arribute")
	}

	return realImgAddr.(string), nil
}

func fetchReferenceHttp(url string, token int, downloadTo string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	m := make(map[string]interface{})
	json.Unmarshal(bytes, &m)

	realImgAddr, ok := m["image"]
	if !ok {
		return "", errors.New("no image arribute")
	}
	return realImgAddr.(string), nil
}

func downloadIpfs(reference string, token int, fpath string) error {
	attempts := 0
	for {
		err := remote.ReadToFile(reference, fpath)
		if err == nil {
			return nil
		}
		if attempts >= 5 {
			return err
		}
		fmt.Println("retrying")
		time.Sleep(2 * time.Second)
		attempts++

	}

	fmt.Printf("%v saved \n", token)
	return nil
}

func downloadHttp(reference string, token int, fpath string) error {
	url := reference
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	f, err := os.Create(fpath)
	defer f.Close()
	if err != nil {
		return err
	}
	var r io.Reader = resp.Body
	fmt.Println("start copying")
	_, err = io.Copy(f, r)
	fmt.Printf("%v saved \n", token)
	return nil
}
