package main

type Message struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"-"`
	Data     string `json:"data"`
}
