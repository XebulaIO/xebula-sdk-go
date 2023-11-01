package xebula

import (
	"fmt"
	"strings"
	"xebula/network"
	"xebula/utils"
)

type Entity struct {
	Config network.Config
}

func NewEntity(cfg network.Config) *Entity {
	return &Entity{
		Config: cfg,
	}
}

func (e *Entity) BuildURL(endpoint string) (string, error) {
	if endpoint == "" || endpoint == "/" {
		return "", NewError(INVALID_URL_ERROR)
	}
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = fmt.Sprintf("/%s", endpoint)
	}
	return fmt.Sprintf("%s/%s%s", e.Config.BaseURL, VERSION, endpoint), nil
}

func (e *Entity) ApiCall(endpoint string, method utils.HttpMethod, in, out interface{}, headers ...utils.Header) error {
	url, err := e.BuildURL(endpoint)
	if err != nil {
		return err
	}
	cli := utils.NewHttpClient()

	return cli.Do(url, method, in, out, headers...)

}
