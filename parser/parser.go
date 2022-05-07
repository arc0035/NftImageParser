package parser

type parser interface {
	DownloadImages(folder string, startToken int, endToken int, concurrency int)
}
