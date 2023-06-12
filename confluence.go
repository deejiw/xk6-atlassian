package atlassian

import (
	"encoding/json"
	"fmt"
	"os"

	confluence "github.com/ctreminiom/go-atlassian/confluence/v2"
	"github.com/dop251/goja"
	"go.k6.io/k6/js/common"
)

type (
	ConfluenceConfig struct {
		Host  string
		Email string
		Token string
	}

	ConfluenceOption func(*Confluence) error

	Confluence struct {
		host  string
		email string
		token string
	}
)

func (a *Atlassian) confluenceClass(call goja.ConstructorCall) *goja.Object {
	runtime := a.vu.Runtime()
	var confluenceConfig *ConfluenceConfig
	if len(call.Arguments) == 0 {
		common.Throw(runtime, ErrNotEnoughArguments)
	}

	if params, ok := call.Argument(0).Export().(map[string]interface{}); ok {
		if b, err := json.Marshal(params); err != nil {
			common.Throw(runtime, err)
		} else {
			if err = json.Unmarshal(b, &confluenceConfig); err != nil {
				common.Throw(runtime, err)
			}
		}
	}

	c := a.confluence(confluenceConfig)

	obj := runtime.NewObject()
	// This is the writer object itself.
	if err := obj.Set("This", c); err != nil {
		common.Throw(runtime, err)
	}

	return runtime.ToValue(obj).ToObject(runtime)
}

func (a *Atlassian) confluence(config *ConfluenceConfig) *confluence.Client {
	envHost := "CONFLUENCE_HOST"
	envEmail := "CONFLUENCE_EMAIL"
	envToken := "CONFLUENCE_TOKEN"

	c, err := newConfluenceConstructor(
		withConfluenceConstructorHost(config.Host, envHost),
		withConfluenceConstructorEmail(config.Email, envEmail),
		withConfluenceConstructorToken(config.Token, envToken),
	)

	if err != nil {
		common.Throw(a.vu.Runtime(), fmt.Errorf("cannot initialize Confluence constructor <%w>", err))
	}

	client, err := confluence.New(nil, c.host)
	if err != nil {
		common.Throw(a.vu.Runtime(), err)
	}
	client.Auth.SetBasicAuth(c.email, c.token)

	return client
}

func newConfluenceConstructor(opts ...ConfluenceOption) (Confluence, error) {
	c := Confluence{}

	for _, opt := range opts {
		if err := opt(&c); err != nil {
			return Confluence{}, fmt.Errorf("Confluence constructor fails to read options %w", err)
		}
	}

	return c, nil
}

func withConfluenceConstructorHost(host string, env string) func(*Confluence) error {
	return func(c *Confluence) error {
		if !isEmpty(host) {
			c.host = host

			return nil
		}

		if envString := os.Getenv(env); envString != "" {
			c.host = envString

			return nil
		}

		return fmt.Errorf("host not found. Please use %s or input 'host' parameter", env)
	}
}

func withConfluenceConstructorEmail(email string, env string) func(*Confluence) error {
	return func(c *Confluence) error {
		if !isEmpty(email) {
			c.token = email

			return nil
		}

		if envString := os.Getenv(env); envString != "" {
			c.token = envString

			return nil
		}

		return fmt.Errorf("email not found. Please use %s or input 'email' parameter", env)
	}
}

func withConfluenceConstructorToken(token string, env string) func(*Confluence) error {
	return func(c *Confluence) error {
		if !isEmpty(token) {
			c.token = token

			return nil
		}

		if envString := os.Getenv(env); envString != "" {
			c.token = envString

			return nil
		}

		return fmt.Errorf("token not found. Please use %s or input 'token' parameter", env)
	}
}
