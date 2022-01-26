package sip

import (
	"net"
	"net/url"
	"strconv"
	"strings"
)

type (
	Uri struct {
		IsEncrypted bool
		HasProtocol bool
		User        string
		Password    string
		Host        string
		Port        int
		Params      Map //params
		Queries     Map //queries
	}

	Map map[string]string
)

func NewUri(user, addr string, ps Map) *Uri {
	uri := &Uri{
		IsEncrypted: false,
		User:        user,
		Params:      ps,
		Queries:     Map{},
	}
	uri.SetAddress(addr)
	return uri
}

func (uri *Uri) EnableProtocol() *Uri {
	uri.HasProtocol = true
	return uri
}

func (uri *Uri) SetParams(m Map) *Uri {
	uri.Params = m
	return uri
}

func (uri *Uri) SetAddress(addr string) *Uri {
	if strings.Index(addr, ":") > -1 {
		if host, port, err := net.SplitHostPort(addr); err == nil {
			uri.Host = host
			uri.Port, _ = strconv.Atoi(port)
		}
	} else {
		uri.Host = addr
	}
	return uri
}

func (uri *Uri) Address() string {
	return net.JoinHostPort(uri.Host, strconv.Itoa(uri.Port))
}

func (m *Map) Set(k, v string) {
	if m == nil || len(*m) == 0 {
		*m = make(map[string]string)
	}
	mp := *m
	mp[k] = v
}

func (m Map) Get(k string) string {
	if m == nil {
		return ""
	}
	return m[k]
}

func (m Map) ToString(sep string) string {
	var (
		i      int
		length int
		sb     strings.Builder
	)
	length = len(m)
	for k, v := range m {
		sb.WriteString(k)
		if len(v) > 0 {
			sb.WriteString("=")
			if strings.ContainsAny(v, " \t") {
				sb.WriteString("\"")
				sb.WriteString(v)
				sb.WriteString("\"")
			} else {
				sb.WriteString(v)
			}
		}
		i++
		if i < length {
			sb.WriteString(sep)
		}
	}
	return sb.String()
}

func (m Map) String() string {
	return m.ToString(";")
}

func (m Map) Clone() Map {
	mm := make(map[string]string)
	for k, v := range m {
		mm[k] = v
	}
	return mm
}

func (uri *Uri) ParamValue(name string) string {
	return uri.Params.Get(name)
}

func (uri *Uri) FormValue(name string) string {
	return uri.Queries.Get(name)
}

func (uri *Uri) Clone() *Uri {
	u := &Uri{
		HasProtocol: uri.HasProtocol,
		IsEncrypted: uri.IsEncrypted,
		User:        uri.User,
		Password:    uri.Password,
		Host:        uri.Host,
		Port:        uri.Port,
	}
	if uri.Params != nil {
		u.Params = uri.Params.Clone()
	}
	if uri.Queries != nil {
		u.Queries = uri.Queries.Clone()
	}
	return u
}

func (uri *Uri) String() string {
	var sb strings.Builder
	// Compulsory protocol identifier.
	if uri.HasProtocol {
		if uri.IsEncrypted {
			sb.WriteString("sips")
			sb.WriteString(":")
		} else {
			sb.WriteString("sip")
			sb.WriteString(":")
		}
	}
	if uri.User != "" {
		sb.WriteString(uri.User)
		if uri.Password != "" {
			sb.WriteString(":" + uri.Password)
		}
		sb.WriteString("@")
	}
	// Compulsory hostname.
	sb.WriteString(uri.Host)
	// Optional port number.
	if uri.Port != 0 {
		sb.WriteString(":")
		sb.WriteString(strconv.Itoa(int(uri.Port)))
	}
	if (uri.Params != nil) && len(uri.Params) > 0 {
		sb.WriteString(";")
		sb.WriteString(uri.Params.String())
	}
	if (uri.Queries != nil) && len(uri.Queries) > 0 {
		sb.WriteString("?")
		sb.WriteString(uri.Queries.ToString("&"))
	}
	return sb.String()
}

func parseMap(s string) (m Map, err error) {
	for s != "" {
		key := s
		if i := strings.IndexAny(key, "&;"); i >= 0 {
			key, s = key[:i], key[i+1:]
		} else {
			s = ""
		}
		if key == "" {
			continue
		}
		value := ""
		if i := strings.Index(key, "="); i >= 0 {
			key, value = key[:i], key[i+1:]
		}
		key, err1 := url.QueryUnescape(key)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		value, err1 = url.QueryUnescape(value)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		if len(value) > 0 {
			value = strings.Trim(value, "\"")
			value = strings.TrimSpace(value)
		}
		m.Set(key, value)
	}
	return
}

//parseUri parse uri from string
func parseUri(s string) (uri *Uri, err error) {
	var (
		pos          int
		endOfUserPos int
		netAddrStr   string
		netSplitPos  int
		p            string
	)
	uri = &Uri{}
	p = s
	//sip or sips
	if pos = strings.Index(p, ":"); pos > -1 && pos < 4 {
		uri.HasProtocol = true
		if strings.ToLower(p[:pos]) == "sips" {
			uri.IsEncrypted = true
		}
		p = p[pos+1:]
	}
	if pos = strings.Index(p, "@"); pos != -1 {
		if endOfUserPos = strings.Index(p[:pos], ":"); endOfUserPos == -1 {
			uri.User = p[:pos]
		} else {
			uri.User = p[:endOfUserPos]
			uri.Password = p[endOfUserPos+1 : pos]
		}
		p = p[pos+1:]
	}
	if pos = strings.IndexAny(p, ";?"); pos == -1 {
		netAddrStr = p
	} else {
		netAddrStr = p[:pos]
	}
	if netSplitPos = strings.Index(netAddrStr, ":"); netSplitPos == -1 {
		uri.Host = netAddrStr
	} else {
		uri.Host = netAddrStr[:netSplitPos]
		uri.Port, _ = strconv.Atoi(netAddrStr[netSplitPos+1:])
	}
	if pos > -1 {
		p = p[pos:]
		if p[0] == '?' { //queries
			uri.Queries, err = parseMap(p[1:])
		} else { //params
			if pos = strings.Index(p, "?"); pos == -1 {
				uri.Params, err = parseMap(p[1:])
			} else {
				uri.Params, err = parseMap(p[1:pos])
				uri.Queries, err = parseMap(p[pos+1:])
			}
		}
	}
	return
}
