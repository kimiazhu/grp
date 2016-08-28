// Author: ZHU HAIHUA
// Date: 8/15/16
package filter

import (
	"bytes"
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"github.com/kimiazhu/log4go"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"github.com/kimiazhu/grp/util/io"
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

func isZipped(encode string) bool {
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
		header.Set("Content-Length", strconv.Itoa(length))
	}
}

func SmartRead(resp *http.Response, replaceHost bool) (body []byte, unzipped bool, err error) {
	encoding := resp.Header.Get("Content-Encoding")
	contentType := resp.Header.Get("Content-Type")
	if isText(contentType) && replaceHost {
		if isZipped(encoding) {
			// unzip the content
			var reader io.ReadCloser
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				log4go.Error("new zip reader failed: %v", err)
				DumpResponse(resp)
				return
			}

			unzipped = true
			// read to byte array then replace all
			body, err = ioutil.ReadAll(reader)
			reader.Close()
			if err != nil {
				log4go.Error("read from zip reader failed: %v", err)
				return
			}
		} else {
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				log4go.Error("read from formal http body failed: %v", err)
				return
			}
		}

		data := ioutils.ReplaceHost(string(body), false)
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

// resp是远程被代理的服务器返回的http应答。
// body是我们处理过后的远程服务器应答Body数据。
// unzipped 表示body数据是否经过解压,如果unzipped=true, 则表明解压过,
// 那么在反馈给客户端之前需要重新压缩。
func SmartWrite(c *gin.Context, resp *http.Response, body []byte, unzipped bool) {
	FilterHeader(resp.Header, c.Writer.Header(), "", false)

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
	//writeCookies(c, resp, true)
	renewContentLength(c.Writer.Header(), len(newBody))
	c.Writer.WriteHeader(resp.StatusCode)
	c.Writer.Write(newBody)
	c.Writer.Flush()
}

func writeCookies(c *gin.Context, resp *http.Response, reverse bool) {
	cookies := FilterCookie(resp.Cookies(), reverse)
	for _, cookie := range cookies {
		log4go.Debug("&&&&&&&&&&&&&&%s", cookie.Domain)
		log4go.Debug("&&&&&&&&&&&&&&%s", cookie.String())
		//c.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
		http.SetCookie(c.Writer, cookie)
	}
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
