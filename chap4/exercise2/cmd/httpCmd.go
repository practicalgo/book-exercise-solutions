package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type httpConfig struct {
	url             string
	postBody        string
	verb            string
	disableRedirect bool
	basicAuth       string
	headers         []string
}

func validateConfig(c httpConfig) error {
	var validMethod bool
	allowedVerbs := []string{http.MethodGet, http.MethodPost, http.MethodHead}
	for _, v := range allowedVerbs {
		if c.verb == v {
			validMethod = true
		}
	}
	if !validMethod {
		return ErrInvalidHTTPMethod
	}

	if c.verb == http.MethodPost && len(c.postBody) == 0 {
		return ErrInvalidHTTPPostRequest
	}

	if c.verb != http.MethodPost && len(c.postBody) != 0 {
		return ErrInvalidHTTPCommand
	}

	return nil
}

func addHeaders(c httpConfig, req *http.Request) {
	for _, h := range c.headers {
		kv := strings.Split(h, "=")
		req.Header.Add(kv[0], kv[1])
	}
}

func addBasicAuth(c httpConfig, req *http.Request) {
	if len(c.basicAuth) != 0 {
		up := strings.Split(c.basicAuth, "=")
		req.SetBasicAuth(up[0], up[1])
	}
}

func HandleHttp(w io.Writer, args []string) error {
	var outputFile string
	var postBodyFile string
	var responseBody []byte
	var redirectPolicyFunc func(req *http.Request, via []*http.Request) error
	var req *http.Request
	var httpClient http.Client

	c := httpConfig{}

	fs := flag.NewFlagSet("http", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&c.verb, "verb", "GET", "HTTP method")
	fs.StringVar(&c.postBody, "body", "", "JSON data for HTTP POST request")
	fs.StringVar(&postBodyFile, "body-file", "", "File containing JSON data for HTTP POST request")
	fs.StringVar(&outputFile, "output", "", "File path to write the response into")
	fs.BoolVar(&c.disableRedirect, "disable-redirect", false, "Do not follow redirection request")
	fs.StringVar(&c.basicAuth, "basicauth", "", "Add basic auth (username:password) credentials to the outgoing request")

	// we want the user to be able to specify the -header option one or more times
	// hence we use this method:
	// https://pkg.go.dev/flag#FlagSet.Func
	headerOptionFunc := func(v string) error {
		c.headers = append(c.headers, v)
		return nil
	}
	fs.Func("header", "Add one or more headers to the outgoing request (key=value)", headerOptionFunc)

	fs.Usage = func() {
		var usageString = `
http: A HTTP client.

http: <options> server`
		fmt.Fprintf(w, usageString)

		fmt.Fprintln(w)
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: ")
		fs.PrintDefaults()
	}

	err := fs.Parse(args)
	if err != nil {
		return FlagParsingError{err}
	}

	if fs.NArg() != 1 {
		return InvalidInputError{ErrNoServerSpecified}
	}

	if len(postBodyFile) != 0 && len(c.postBody) != 0 {
		return InvalidInputError{ErrInvalidHTTPPostCommand}
	}

	if c.verb == http.MethodPost && len(postBodyFile) != 0 {
		data, err := os.ReadFile(postBodyFile)
		if err != nil {
			return err
		}
		c.postBody = string(data)
	}

	err = validateConfig(c)
	if err != nil {
		return InvalidInputError{err}
	}

	c.url = fs.Arg(0)

	if c.disableRedirect {
		redirectPolicyFunc = func(req *http.Request, via []*http.Request) error {
			if len(via) >= 1 {
				return errors.New("stopped after 1 redirect")
			}
			return nil
		}
	}

	httpClient = http.Client{CheckRedirect: redirectPolicyFunc}
	switch c.verb {
	case http.MethodGet:
		req, err = http.NewRequestWithContext(
			// FIXME: Add a valid context
			context.Background(),
			http.MethodGet,
			c.url,
			nil,
		)
		if err != nil {
			return err
		}
	case http.MethodPost:
		postBodyReader := strings.NewReader(c.postBody)
		req, err = http.NewRequestWithContext(
			// FIXME: Add a valid context
			context.Background(),
			http.MethodPost,
			c.url,
			postBodyReader,
		)
		if err != nil {
			return err
		}
		c.headers = append(c.headers, "Content-Type=application/json")
	}

	addHeaders(c, req)
	addBasicAuth(c, req)

	r, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	responseBody, err = io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if len(outputFile) != 0 {
		f, err := os.Create(outputFile)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write(responseBody)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "Data saved to: %s\n", outputFile)
		return err
	}
	fmt.Fprintln(w, string(responseBody))
	return nil
}
