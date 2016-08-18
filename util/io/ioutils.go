// Author: ZHU HAIHUA
// Date: 8/15/16
package ioutils

import (
	"bytes"
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"github.com/kimiazhu/grp/model"
	"github.com/kimiazhu/log4go"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
)

var TEXT_CONTENT []string = []string{"text/plain", "text/html", "text/xml", "text/javascript", "text/json",
	"application/json", "application/xml", "application/javascript"}

func isText(contentType string) bool {
	for _, typ := range TEXT_CONTENT {
		if strings.Contains(contentType, typ) {
			return true
		}
	}
	return false
}

func isZiped(encode string) bool {
	return strings.Contains(encode, "gzip")
}

func DumpResponse(resp *http.Response) {
	dat, e := httputil.DumpResponse(resp, true)
	if e != nil {
		log4go.Error("dump response failed: %v", e)
	} else {
		log4go.Error("dumped response body: %s", string(dat))
	}
}

func renewContentLength(header http.Header, length int) {
	if header.Get("Content-Length") != "" {
		header.Del("Content-Length")
		header.Add("Content-Length", strconv.Itoa(length))
	}
}

func SmartRead(resp *http.Response, p model.Proxies, replaceHost bool) (body []byte, unzipped bool, err error) {
	encoding := resp.Header.Get("Content-Encoding")
	contentType := resp.Header.Get("Content-Type")
	if resp.StatusCode == http.StatusOK && isZiped(encoding) && isText(contentType) && replaceHost {
		// unzip the content
		var reader io.ReadCloser
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			log4go.Error("new zip reader failed: %v", err)
			DumpResponse(resp)
			return
		}
		defer reader.Close()

		unzipped = true
		// read to byte array then replace all
		body, err = ioutil.ReadAll(reader)
		if err != nil {
			log4go.Error("read from zip reader failed: %v", err)
			return
		}

		data := string(body)
		for _local, _remote := range p {
			data = strings.Replace(data, "https://"+_remote, "http://"+_local, -1)
			data = strings.Replace(data, "http://"+_remote, "http://"+_local, -1)
		}
		body = []byte(data)
	} else {
		// not compressed or not text(which can pass to client directly)
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log4go.Error("read from reader failed: %v", err)
			return
		}
	}
	return
}

func SmartWrite(c *gin.Context, resp *http.Response, body []byte, unzipped bool) {
	for k, v := range resp.Header {
		for _, vv := range v {
			c.Writer.Header().Add(k, vv)
		}
	}

	newBody := body
	log4go.Fine("origin body length: %v, need compress: %v", len(body), unzipped)
	if len(body) > 0 && unzipped {
		var buf bytes.Buffer
		gz := gzip.NewWriter(&buf)
		_, err := gz.Write(body)
		gz.Close()
		if err != nil {
			log4go.Debug("zip body failed: %v", err)
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		newBody = buf.Bytes()
		log4go.Fine("complete compress body, new length: %v", len(newBody))
	}
	renewContentLength(c.Writer.Header(), len(newBody))
	c.Writer.WriteHeader(resp.StatusCode)
	c.Writer.Write(newBody)
	c.Writer.Flush()
}

func ZipReader(encoding string, reader io.ReadCloser) (io.ReadCloser, error) {
	var r io.ReadCloser
	var err error
	switch encoding {
	case "gzip":
		r, err = gzip.NewReader(reader)
		if err != nil {
			return nil, err
		}
	default:
		r = reader
	}
	return r, nil
}

func ZipWriter(encoding string, writer io.Writer) io.Writer {
	var w io.Writer
	switch encoding {
	case "gzip":
		w = gzip.NewWriter(writer)
	default:
		w = writer
	}
	return w
}
