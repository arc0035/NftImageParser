package parser

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

//https://www.larvalabs.com/public/images/cryptopunks/punk0001.png
type PunkParser struct {
	Address string
	Name    string
}

func NewPunkParser(address string, name string) *PunkParser {
	parser := PunkParser{}
	parser.Address = address
	parser.Name = name
	return &parser
}

func (p *PunkParser) DownloadImages(folder string, startToken, endToken int, concurrency int) {
	downloadTo := folder
	_ = os.Mkdir(downloadTo, 0755)
	//punk0000.png ~ punk9999.png
	var wg sync.WaitGroup
	nfts := endToken - startToken + 1
	wg.Add(nfts)

	guard := make(chan struct{}, concurrency)
	for i := startToken; i <= endToken; i++ {
		guard <- struct{}{}
		go func(tokenId int) {
			defer func() {
				<-guard
				wg.Done()
			}()
			punkPng := fmt.Sprintf("punk%04d.png", tokenId)
			url := p.Address + "/" + punkPng
			resp, err := http.Get(url)
			if err != nil {
				log.Fatalln(err)
				return
			}
			fpath := filepath.Join(downloadTo, fmt.Sprintf("%v_%v.png", p.Name, tokenId))
			f, err := os.Create(fpath)
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = io.Copy(f, resp.Body)
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Printf("%v download complete\n", tokenId)
		}(i)

	}

	wg.Wait()
	fmt.Println("Cryptopunks completed")
}
