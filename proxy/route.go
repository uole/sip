package proxy

//Route 代理走的路由规则
type Route struct {
	index     int32
	Domain    string   `json:"domain" yaml:"domain"`        //域名
	RewriteTo string   `json:"rewrite_to" yaml:"rewriteTo"` //对域名进行重写处理
	Backend   []string `json:"backend" yaml:"backend"`      //代理的后端地址，多个地址使用轮询获取地址
}

func (r *Route) Address() string {
	idx := int(r.index) % len(r.Backend)
	return r.Backend[idx]
}
