package ice

import (
    "net"
    "github.com/ccding/go-stun/stun"
    "github.com/WeMeetAgain/go-sdp"
    "math"
    "strconv"
    "time"
    )

var MaxPairs = 100

const (
    IceWaiting byte = iota
    IceRunning
    IceCompleted
    IceFailed
    IceTerminated
    )

type IceServer struct {
    Urls        []string
    CredLocal   Credential
    CredRemote  Credential
    Description SessionDescription
}

type Credential struct {
    Username string
    Password string
}

type Agent struct {
    Server      IceServer
    Aggressive  bool
    Controlling bool
    State       byte
    Streams     []*MediaStream
    checkers    []*checker
}

func (a *Agent) GetOffer() sdp.SessionDescription {
    a.setCandidates()
    return a.formulateSDP()
}

func (a *Agent) SetRemoteOffer(offer sdp.SessionDescription) sdp.SessionDescription {
    a.Remote = decodeRemoteSDP(offer)
    a.setCandidates()
    a.StartConnectivityEstablishment()
    return a.formulateSDP()
}

func (a *Agent) SetRemoteAnswer(answer sdp.SessionDescription) {
    a.Remote = decodeRemoteSDP(answer)
    a.StartConnectivityEstablishment()
}

func (a *Agent) StartConnectivityEstablishment() {
    a.removeMismatchedStreams()
    a.initCheckLists()
    a.State = IceRunning
    // check candidatepairs who've already received connectivity checks
    //
    a.startChecks()
}

// remove dangling components and streams (ones with no remote candidates or components)
func (a *Agent) removeMismatchedStreams() {
    for i, s := range a.Streams {
        for j, comp := range s.Components {
            if len(comp.Remote) == 0 {
                delete(a.Streams[i].Component,j)
            }
        }
        if len(a.Streams[i]) == 0 {
            a.removeStream(s.Name)
        }
    }
}

func (a *Agent) removeStream(n string) {
    for i, s := range a.Streams {
        if s.Name == n {
            delete(a.Streams,i)
            return
        }
    }
}

func (a *Agent) initCheckLists() {
    MaxPairsPerStream = MaxPairs / len(a.Streams)
    for _, s := range a.Streams {
        s.initCheckList()
    }
    a.Streams[0].setInitialCheckListStates()
}

func (a *Agent) decodeRemoteSDP(s sdp.SessionDescription) {
    // addRemoteCandidate
    for _, m := range s.MediaDescriptions {
        
    }
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
    a.AddLocal(c)
}

func (a *Agent) AddLocal(c *Candidate) error {
    for _, m := range a.Streams {
        for _, comp := range m.Components {
            comp.AddLocal(c)
        }
    }
    return nil
}

func (a *Agent) SetDefaultCandidates() {
    for _, m := range a.Streams {
        for _, comp := range m.Components {
            comp.SetDefaultLocal()
        }
    }
}

func (a *Agent) pendingStreams() []*MediaStream {
    streams := make([]*MediaStream)
    for _, s := range a.Streams {
        if s.CheckList.State == IceWaiting {
            streams = append(streams,s)
        }
    }
    return streams
}

func (a *Agent) getWaitInterval() time.Duration {
    activeChecks := a.ActiveCheckListCount()
    if activeChecks < 1 {
        acativeChecks = 1
    }
    return time.Duration(a.calculateTa() * activeChecks)
}
