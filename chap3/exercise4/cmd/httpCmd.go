package cmd

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type uploadFileData struct {
	name string
	data io.Reader
}

type httpConfig struct {
	url        string
	postBody   string
	verb       string
	uploadFile *uploadFileData
	formDataKv []string
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

	if c.verb == http.MethodPost && len(c.postBody) == 0 && c.uploadFile == nil && len(c.formDataKv) == 0 {
		return ErrInvalidHTTPPostRequest
	}

	if c.verb != http.MethodPost && len(c.postBody) != 0 {
		return ErrInvalidHTTPCommand
	}

	return nil
}

func createMultiPartBody(c httpConfig) (io.Reader, string, error) {
	var b bytes.Buffer
	var err error
	var fw io.Writer

	mw := multipart.NewWriter(&b)

	// each value is a key value pair, name=value
	for _, kv := range c.formDataKv {
		f := strings.Split(kv, "=")
		fw, err = mw.CreateFormField(f[0])
		if err != nil {
			return nil, "", err
		}
		fmt.Fprintf(fw, f[1])
	}

	// we create a new form field for any JSON data that we also need
	// to send
	if len(c.postBody) != 0 {
		fw, err = mw.CreateFormField("jsondata")
		if err != nil {
			return nil, "", err
		}
		fmt.Fprintf(fw, c.postBody)
	}

	// add file if any
	if c.uploadFile != nil {
		fw, err = mw.CreateFormFile("filedata", c.uploadFile.name)
		if err != nil {
			return nil, "", err
		}
		_, err = io.Copy(fw, c.uploadFile.data)
		if err != nil {
			return nil, "", err
		}
	}

	err = mw.Close()
	if err != nil {
		return nil, "", err
	}

	return bytes.NewReader(b.Bytes()), mw.FormDataContentType(), nil
}

func createPostBody(c httpConfig) (io.Reader, string, error) {
	// if we have both json data and form data to send, we will create a new special form field
	// to send the JSON data. This obviously requires that the server is aware of this
	// "protocol". hence, for our test server, which we will use to test the implementation, we
	// will expect the form field, "jsondata" to contain any JSON data if we get a
	// multipart/form-data request

	if c.uploadFile != nil || len(c.formDataKv) != 0 {
		return createMultiPartBody(c)
	}
	return strings.NewReader(c.postBody), "application/json", nil
}

func fetchRemoteResource(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

func createRemoteResource(url string, body io.Reader, contentType string) ([]byte, error) {
	r, err := http.Post(url, contentType, body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

func HandleHttp(w io.Writer, args []string) error {
	c := httpConfig{}
	var outputFile string
	var postBodyFile string
	var uploadFile string
	var responseBody []byte

	fs := flag.NewFlagSet("http", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.StringVar(&c.verb, "verb", "GET", "HTTP method")
	fs.StringVar(&c.postBody, "body", "", "JSON data for HTTP POST request")
	fs.StringVar(&postBodyFile, "body-file", "", "File containing JSON data for HTTP POST request")
	fs.StringVar(&outputFile, "output", "", "File path to write the response into")
	fs.StringVar(&uploadFile, "upload", "", "Path of file to upload")
	formDataKvOptionFunc := func(v string) error {
		c.formDataKv = append(c.formDataKv, v)
		return nil
	}
	fs.Func("form-data", "Add one or more key value pairs (key=value) to send as form data", formDataKvOptionFunc)

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
		return err
	}

	if fs.NArg() != 1 {
		return ErrNoServerSpecified
	}

	if len(postBodyFile) != 0 && len(c.postBody) != 0 {
		return ErrInvalidHTTPPostCommand
	}

	if c.verb == http.MethodPost {
		if len(postBodyFile) != 0 {
			data, err := os.ReadFile(postBodyFile)
			if err != nil {
				return err
			}
			c.postBody = string(data)
		}

		if len(uploadFile) != 0 {
			c.uploadFile = &uploadFileData{}
			f, err := os.Open(uploadFile)
			if err != nil {
				return err
			}
			c.uploadFile.data = f
			c.uploadFile.name = filepath.Base(uploadFile)

		}
	}

	err = validateConfig(c)
	if err != nil {
		if errors.Is(err, ErrInvalidHTTPMethod) || errors.Is(err, ErrInvalidHTTPPostRequest) {
			fmt.Fprintln(w, err.Error())
		}
		return err
	}

	c.url = fs.Arg(0)

	switch c.verb {
	case http.MethodGet:
		responseBody, err = fetchRemoteResource(c.url)
		if err != nil {
			return err
		}
	case http.MethodPost:
		postBody, contentType, err := createPostBody(c)
		if err != nil {
			return err
		}
		responseBody, err = createRemoteResource(c.url, postBody, contentType)
		if err != nil {
			return err
		}
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
