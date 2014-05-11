package ice

import (
    "net"
    "github.com/ccding/go-stun/stun"
    "github.com/WeMeetAgain/go-sdp"
    "math"
    "strconv"
    )



type IceServer struct {
    Urls []string
    LocalCred Credential
    RemoteCred Credential
    Description SessionDescription
}

type Credential struct {
    Username string
    Password string
}

type Media struct {
    Conn *net.Conn
    Credential Credential
    Description sdp.MediaDescription
    DefaultCandidate *Candidate
}

type Agent struct {
    Server IceServer
    Aggressive bool
    State int
    Streams []Media
    Local []*Candidate
    Remote []*Candidate
    Pairs []*CandidatePair
    mu sync.*RWMutex
}



func (a *Agent) GetOffer() sdp.SessionDescription {
    a.setCandidates()
    return a.formulateSDP()
}

func (a *Agent) SetRemoteOffer(offer sdp.SessionDescription) sdp.SessionDescription {
    a.Remote = decodeRemoteSDP(offer)
    a.setCandidates()
    a.Check()
    return a.formulateSDP()
}

func (a *Agent) SetRemoteAnswer(answer sdp.SessionDescription) {
    a.Remote = decodeRemoteSDP(answer)
    a.Check()
}

func (a *Agent) Check() {
    
}

func (a *Agent) decodeRemoteSDP(s sdp.SessionDescription) {
    // addRemoteCandidate
}

func (a *Agent) formulateSDP() sdp.SessionDescription {
    s := a.Description
    // Media Description for each Media stream
    for _, stream := range a.Streams {
        media := stream.Description
        if len(a.Streams) == 1 {
            s.Connection = sdp.Connection{"IN", "IP4",stream.DefaultCandidate.Address.IPAddr}
        } else {
            media.Connection = sdp.Connection{"IN", "IP4",stream.DefaultCandidate.Address.IPAddr}
        }
        // candidate attribute 
        for _, cand := a.Local {
            attrVal := cand.Foundation + " " + string(cand.ComponentId)
            + " " + cand.Address.Transport + " " + strconv.FormatUint(uint64(cand.Priority),10)
            + " " + cand.Address.IPAddr + " " + cand.Address.Port
            + " typ " + cand.Type
            if cand.RelatedAddr.IPAddr != "" {
                attrVal += " raddr " + cand.RelatedAddr.IPAddr
            }
            if cand.RelatedAddr.Port != "" {
                attrVal += " rport " + cand.RelatedAddr.Port
            }
            attr := sdp.Attribute{Key:"candidate", Value:attrVal}
            media.Attributes = append(media.Attributes, attr)
        }
        // media-specific credential attribute
        if stream.Credential.Username != "" && stream.Credential.Password != "" {
            media.Attributes = append(media.Attributes, sdp.Attribute{"ice-pwd",stream.Credential.Password})
            media.Attributes = append(media.Attributes, sdp.Attribute{"ice-ufrag",stream.Credential.Username})
        }
        s.MediaDescriptions = append(s.MediaDescriptions,media)
    }
    return s
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
    a.mu.Lock()
    defer a.mu.Unlock()
    for key, cand := range a.Local {
        // check for redundant
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
    a.mu.Lock()
    defer mu.Unlock()
    for key, cand := range a.Local {
        if a.Local[index].Type != cand.Type {
            index = key
        }
    }
    for _, conn := range a.Streams {
        conn.DefaultCandidate = a.Local[index]
    }
}
