package ice

type Component struct {
    Id            int
    ParentStream  *MediaStream
    Conn          *net.Conn
    DefaultLocal  *Candidate
    DefaultRemote *Candidate
    Local         []*Candidate
    Remote        []*Candidate
}

// add candidate to local candidate list, prioritizing and not adding if redundant
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

// returns all possible pairs for this component (priorities already calculated)
func (comp *Component) Pairs() []*CandidatePair {
    cp := make([]*CandidatePair)
    for _, l := range comp.Local {
        for _, r := range comp.Remote {
            // create pair if same IP version (and same component ID)
            if l.Type == r.Type {
                var p int
                if comp.ParentStream.parent.Controlling {
                    p = pairPriority(l,r)
                } else {
                    p = pairPriority(r,l)
                }
                cp = append(cp, &CandidatePair{
                    ComponentId: comp.Id,
                    Priority: p,
                    State: Frozen,
                    Local: l,
                    Remote: r,
                    CredLocal: l.ParentStream.getCredLocal(),
                    CredRemote: l.ParentStream.getCredRemote(),
                })
            }
        }
    }
    return cp
}
