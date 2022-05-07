package main

import (
	"fmt"
	"nftParser/parser"
	//"nftParser/remote"
)

func main() {
	fmt.Println("Hello world")

	baycParser := parser.NewErc721Parser(
		"0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D",
		"bayc",
	)

	baycParser.DownloadImages("./nfts", 1, 200, 100)

	// maycParser := parser.NewErc721Parser(
	// 	"0x60E4d786628Fea6478F785A6d7e704777c86a7c6",
	// 	"mayc",
	// )

	// maycParser.DownloadImages("./nfts", 143, 143, 50)

	// bakcParser := parser.NewErc721Parser(
	// 	"0xba30e5f9bb24caa003e9f2f0497ad287fdf95623",
	// 	"bakc",
	// )

	// bakcParser.DownloadImages("./nfts", 300, 310, 50)

	// meebitsParser := parser.NewErc721Parser(
	// 	"0x7Bd29408f11D2bFC23c34f18275bBf23bB716Bc7",
	// 	"meebits",
	// )

	// meebitsParser.DownloadImages("./nfts", 1, 400, 100)

	// punkParser := parser.NewPunkParser(
	// 	"http://www.larvalabs.com/public/images/cryptopunks",
	// 	"cryptopunks",
	// )

	// punkParser.DownloadImages("./nfts", 0, 399, 100)
}
