package main

type BFWConf struct {
	Destinations []Destination `json:"destinations"`
}

type Destination struct {
	From string `json:"from"`
	To   string `json:"to"`
	Port int    `json:"port"`
}
