package remote

import (
	"context"
	"net/http"

	ipfshttp "github.com/ipfs/go-ipfs-http-client"
	iface "github.com/ipfs/interface-go-ipfs-core"
)

var Client = &http.Client{}

func Http(ctx context.Context) (iface.CoreAPI, error) {

	//https://docs.ipfs.io/concepts/ipfs-gateway/#gateway-providers
	httpApi, err := ipfshttp.NewURLApiWithClient("https://cloudflare-ipfs.com", Client)
	if err != nil {
		return nil, err
	}

	err = httpApi.Request("version").Exec(ctx, nil)
	if err != nil {
		return nil, err
	}
	return httpApi, nil
}
