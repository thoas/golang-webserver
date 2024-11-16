package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func getBody(r *http.Request) (interface{}, []byte, error) {
	var (
		err  error
		save io.ReadCloser
	)
	save, r.Body, err = drainBody(r.Body)
	if err != nil {
		return nil, nil, err
	}
	if save != nil {
		if r.Header.Get("Content-Type") == "application/json" {
			out, err := dumpJson(save)
			if err != nil {
				return nil, nil, err
			}

			return out, nil, nil
		}

		body, err := io.ReadAll(save)
		if err != nil {
			return nil, nil, err
		}
		return nil, body, nil
	}

	return nil, nil, nil
}

func remoteIP(req *http.Request) string {
	url, err := url.Parse(req.Header.Get("Origin"))
	if err == nil {
		host := url.Host
		ip, _, err := net.SplitHostPort(host)
		if err == nil {
			return ip
		}
	}

	ip := getClientIPByRequestRemoteAddr(req)
	if ip != "" {
		return ip
	}

	ip = getClientIPByHeaders(req)
	if ip != "" {
		return ip
	}

	return ""
}

func getClientIPByHeaders(req *http.Request) string {
	ipSlice := make([]string, 0)
	ipSlice = append(ipSlice, req.Header.Get("X-Forwarded-For"))
	ipSlice = append(ipSlice, req.Header.Get("x-forwarded-for"))
	ipSlice = append(ipSlice, req.Header.Get("X-FORWARDED-FOR"))

	for _, v := range ipSlice {
		if v != "" {
			return v
		}
	}
	return ""
}

func getClientIPByRequestRemoteAddr(r *http.Request) string {
	ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err != nil {
		return ""
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return ""
	}
	return userIP.String()

}

func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == nil || b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return io.NopCloser(&buf), io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

func dumpJson(save io.ReadCloser) (interface{}, error) {
	var out interface{}
	body, err := io.ReadAll(save)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}

	return out, nil
}
