package ice

// collection of components
type MediaStream struct {
    CredLocal   Credential
    CredRemote  Credential
    Description sdp.MediaDescription
    Components  []*Component
}

type Component struct {
    Id            int
    Conn        *net.Conn
    DefaultLocal  *Candidate
    DefaultRemote *Candidate
    Local         []*Candidate
    Remote        []*Candidate
    Pairs         []*CandidatePair
}

func (comp *Component) AddLocal(c *Candidate) error {
    if len(comp.Local) == 100 {
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
    localPref := 65535 - len(comp.Local)
    priority = (math.Exp2(24) * typPref) + (math.Exp2(8) * localPref) + (256)
    c.Priority = priority
    var redudant int
    var added bool
    comp.mu.Lock()
    defer comp.mu.Unlock()
    for key, cand := range comp.Local {
        // check for redundant
        if c.TransportAddr.Transport == cand.TransportAddr.Transport && c.Base == cand.Base {
            if c.Priority > cand.Priority {
                redundant = key
            }
        }
        // add in proper place in list
        if !added {
            if c.Priority > cand.Priority {
                comp.Local = append([c],comp.Local...)
                added = true
            } else {
                if key > len(comp.Local) - 1 {
                    if c.Priority > comp.Local[key+1].Priority {
                        comp.Local = append(append(comp.Local[0:key],c)comp.Local[key+1:len(comp.Local)-1]...)
                        added = true
                    }
                } else {
                    comp.Local = append(comp.Local,c)
                    added = true
                }
            }
        }
    }
    // get rid of redundant, if it exists
    comp.Local = append(comp.Local[0:redundant],comp.Local[redundant+1,len(comp.Local)-1]
    return nil
}

func (comp *Component) setDefaultLocal() {
    index  := 0
    comp.mu.Lock()
    defer comp.mu.Unlock()
    for key, cand := range comp.Local {
        if comp.Local[index].Type != cand.Type {
            index = key
        }
    }
    comp.DefaultLocal = comp.Local[index]
}
