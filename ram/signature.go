package ram

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"net/url"
	"sort"
	"strings"
)

// CreateQueryStr create url query string
func CreateQueryStr(args map[string]string) string {
	keys := make([]string, 0, len(args))
	for k := range args {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	var buf bytes.Buffer
	for _, key := range keys {
		value := args[key]
		prefix := PercentEncode(key) + "="
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(prefix)
		buf.WriteString(PercentEncode(value))
	}
	return buf.String()
}

// PercentEncode encode string
func PercentEncode(s string) string {
	spec := map[string]string{"+": "%20", "*": "%2A", "%7E": "~"}
	s = url.QueryEscape(s)
	for key, value := range spec {
		s = strings.Replace(s, key, value, -1)
	}
	return s
}

// HmacSha1 is the signature algorithm
func HmacSha1(stringToSign, key string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(stringToSign))
	sha1Value := mac.Sum(nil)
	signature := base64.StdEncoding.EncodeToString(sha1Value)
	return signature
}

// CreateSignature  create signature
func CreateSignature(method, cquerystr, secret string) string {
	strToSign := method + "&" + PercentEncode("/") + "&" + PercentEncode(cquerystr)
	return HmacSha1(strToSign, secret+"&")
}
