package ice

import (
    "net"
    "github.com/ccding/go-stun/stun"
    "github.com/WeMeetAgain/go-sdp"
    "math"
    )

type IceServer struct {
    urls []string
    username string
    credential string
}

type Media struct {
    Conn *net.Conn
    Description sdp.MediaDescription
    DefaultCandidate *Candidate
}

type Agent struct {
    Server IceServer
    Aggressive bool
    Streams []Media
    Local []*Candidate
    Remote []*Candidate
    Pairs []*CandidatePair
    candidateMu sync.*RWMutex
}

func (a *Agent) GetOffer() sdp.SessionDescription {
    a.setCandidates()
    return a.formulateSDP()
}

func (a *Agent) SetRemoteOffer(offer sdp.SessionDescription) sdp.SessionDescription {
    a.Remote = decodeSDP(offer)
    a.setCandidates()
    a.Check()
    return a.formulateSDP()
}

func (a *Agent) SetRemoteAnswer(answer sdp.SessionDescription) {
    a.Remote = decodeSDP(answer)
    a.Check()
}

func (a *Agent) Check() {
    
}

func (a *Agent) decodeSDP(s sdp.SessionDescription) {
    
}

func (a *Agent) formulateSDP() sdp.SessionDescription {
    
}

func (a *Agent) setCandidates() {
    a.gatherCandidates()
    a.setDefaultCandidates()
}

func (a *Agent) gatherCandidates() {
    // gather using different methods
    localhosts := net.LookupHost("localhost")
    natInfo, stunHost, err := stun.Discover()
    // append
    var c Candidate
    for addr := range localhosts {
        c = &Canididate{
            Address: TransportAddr{IPAddr:addr,Port:NewPort(),Transport:"udp"},
            Type: host
        }
        a.addCandidate(c)
    }
    c = &Canididate{
        Address: TransportAddr{IPAddr:stunHost.Ip(),Port:stunHost.Port(),Transport:stunHost.Transport()},
        Type: serverReflex
    }
    a.addCandidate(c)
}

func (a *Agent) addCandidate(c *Candidate) error {
    if len(a.Local) == 100 {
        return errors.New("too many candidates")
    }
    // set priority
    var typPref int
    switch c.Type {
    case host:
        typPref = 126
    case serverReflex:
        typPref = 100
    case peerReflex:
        typPref = 110
    case relay:
        typPref = 0
    }
    localPref := 65535 - len(a.Local)
    priority = (math.Exp2(24) * typPref) + (math.Exp2(8) * localPref) + (256)
    c.Priority = priority
    var redudant int
    var added bool
    a.candidateMu.Lock()
    defer a.candidateMu.Unlock()
    for key, cand := range a.Local {
        // check for redundanct
        if c.TransportAddr.Transport == cand.TransportAddr.Transport && c.Base == cand.Base {
            if c.Priority > cand.Priority {
                redundant = key
            }
        }
        // add in proper place in list
        if !added {
            if c.Priority > cand.Priority {
                a.Local = append([c],a.Local...)
                added = true
            } else {
                if key > len(a.Local) - 1 {
                    if c.Priority > a.Local[key+1].Priority {
                        a.Local = append(append(a.Local[0:key],c)a.Local[key+1:len(a.Local)-1]...)
                        added = true
                    }
                } else {
                    a.Local = append(a.Local,c)
                    added = true
                }
            }
        }
    }
    // get rid of redundant, if it exists
    a.Local = append(a.Local[0:redundant],a.Local[redundant+1,len(a.Local)-1]
    return nil
}

func (a *Agent) setDefaultCandidates() {
    index  := 0
    a.candidateMu.Lock()
    defer candidateMu.Unlock()
    for key, cand := range a.Local {
        if a.Local[index].Type != cand.Type {
            index = key
        }
    }
    for _, conn := range a.Streams {
        conn.DefaultCandidate = a.Local[index]
    }
}
