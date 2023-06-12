package atlassian

import (
	"encoding/json"
	"fmt"
	"os"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/dop251/goja"
	"go.k6.io/k6/js/common"
)

type (
	JiraConfig struct {
		Host  string
		Email string
		Token string
	}

	JiraOption func(*Jira) error

	Jira struct {
		host  string
		email string
		token string
	}
)

func (a *Atlassian) jiraClass(call goja.ConstructorCall) *goja.Object {
	runtime := a.vu.Runtime()
	var jiraConfig *JiraConfig
	if len(call.Arguments) == 0 {
		common.Throw(runtime, ErrNotEnoughArguments)
	}

	if params, ok := call.Argument(0).Export().(map[string]interface{}); ok {
		if b, err := json.Marshal(params); err != nil {
			common.Throw(runtime, err)
		} else {
			if err = json.Unmarshal(b, &jiraConfig); err != nil {
				common.Throw(runtime, err)
			}
		}
	}

	c := a.jira(jiraConfig)

	obj := runtime.NewObject()
	// This is the writer object itself.
	if err := obj.Set("This", c); err != nil {
		common.Throw(runtime, err)
	}

	return runtime.ToValue(obj).ToObject(runtime)
}

func (a *Atlassian) jira(config *JiraConfig) *jira.Client {
	envHost := "CONFLUENCE_HOST"
	envEmail := "CONFLUENCE_EMAIL"
	envToken := "CONFLUENCE_TOKEN"

	c, err := newJiraConstructor(
		withJiraConstructorHost(config.Host, envHost),
		withJiraConstructorEmail(config.Email, envEmail),
		withJiraConstructorToken(config.Token, envToken),
	)

	if err != nil {
		common.Throw(a.vu.Runtime(), fmt.Errorf("cannot initialize Jira constructor <%w>", err))
	}

	client, err := jira.New(nil, c.host)
	if err != nil {
		common.Throw(a.vu.Runtime(), err)
	}
	client.Auth.SetBasicAuth(c.email, c.token)

	return client
}

func newJiraConstructor(opts ...JiraOption) (Jira, error) {
	c := Jira{}

	for _, opt := range opts {
		if err := opt(&c); err != nil {
			return Jira{}, fmt.Errorf("Jira constructor fails to read options %w", err)
		}
	}

	return c, nil
}

func withJiraConstructorHost(host string, env string) func(*Jira) error {
	return func(c *Jira) error {
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

func withJiraConstructorEmail(email string, env string) func(*Jira) error {
	return func(c *Jira) error {
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

func withJiraConstructorToken(token string, env string) func(*Jira) error {
	return func(c *Jira) error {
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
