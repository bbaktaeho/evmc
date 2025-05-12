package evmc

import "context"

type web3Namespace struct {
	c caller
}

func (w *web3Namespace) ClientVersion() (string, error) {
	result := new(string)
	if err := w.c.call(context.Background(), result, Web3ClientVersion); err != nil {
		return "", err
	}
	return *result, nil
}
