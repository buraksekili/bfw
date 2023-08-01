package main

type BFWConf struct {
	Proxy []Proxy `json:"proxy"`
}

type Target struct {
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
	Path   string `json:"path"`
}

type Proxy struct {
	ListenPath string `json:"listenPath"`
	Target     Target `json:"target"`
}
