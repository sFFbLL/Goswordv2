package utils

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func OperateLog(c *gin.Context) ([]byte, []byte) {
	var mBody []byte
	var query string
	var queryGet []byte
	if c.Request.Method != http.MethodGet {
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		var mPost map[string]string
		_ = json.Unmarshal(body, &mPost)
		for k, _ := range mPost {
			if k == "phone" {
				delete(mPost, k)
			}
		}
		mBody, _ = json.Marshal(mPost)
	} else {
		query = c.Request.URL.RawQuery
		query, _ = url.QueryUnescape(query)
		split := strings.Split(query, "&")
		mGet := make(map[string]string)
		for _, v := range split {
			kv := strings.Split(v, "=")
			if len(kv) == 2 {
				mGet[kv[0]] = kv[1]
			}
		}
		queryGet, _ = json.Marshal(mGet)
	}
	return mBody, queryGet
}
