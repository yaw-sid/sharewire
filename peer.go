package main

type Peer struct {
	Address string `json:"address"`
	Name    string `json:"name"`
	ID      int    `json:"id"`
}

func NewPeer(addr, name string, id int) *Peer {
	return &Peer{
		Address: addr,
		Name:    name,
		ID:      id,
	}
}
