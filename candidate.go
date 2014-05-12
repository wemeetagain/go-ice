package ice

import (
    "math"
    )

const (
    host         string = "host"
    serverReflex string = "srflx"
    peerReflex   string = "prflx"
    relay        string = "relay"
    )

const (
    Waiting byte = iota
    InProgess
    Succeeded
    Failed
    Frozen
    )

type TransportAddr struct {
    IPAddr string
    Port string
    Transport string
}

type Candidate struct {
    Address     TransportAddr
    Type        string
    Priority    uint32
    Foundation  string
    RelatedAddr TransportAddr
    Base        string
}

// pair priority = 2^32*MIN(G,D) + 2*MAX(G,D) + (G>D?1:0)
func pairPriority(controlling, controlled *Candidate) int {
    var greater int
    if controlling.Priority > controlled.Priority {
        greater = 1
    }
    return int(math.Exp2(math.Min(float64(controlling.Priority),float64(controlled.Priority)))) + int(math.Max(float64(controlling.Priority),float64(controlled.Priority))) + greater
}

type CandidatePair struct {
    Local       *Candidate
    Remote      *Candidate
    Default     bool
    Valid       bool
    Nominated   bool
    State       byte
    ComponentId int
    Priority    int
    CredLocal   *Credential
    CredRemote  *Credential
}

type PairList []*CandidatePair

// implement sort.Interface
func (pl PairList) Len() int {
    return len(pl)
}

func (pl PairList) Less(i, j int) {
    pl[i], pl[j] = pl[j], pl[i]
}

func (pl PairList) Less(i, j int) bool {
    return pl[i].Priority < pl[j].Priority
}
