package ice

import (
    "sort"
    )

var MaxPairsPerStream = 50

// collection of components
type MediaStream struct {
    Name        string
    parent      *Agent
    CredLocal   *Credential
    CredRemote  *Credential
    Description *sdp.MediaDescription
    Components  map[int]*Component
    CheckList   *PairList
    ValidList   *PairList
}

func (m *MediaStream) initCheckList() {
    m.createCheckList()
    m.orderCheckList()
    m.pruneCheckList()
}

func (m *MediaStream) createCheckList() {
    pl := &PairList{}
    for _, comp := range m.Components {
        pl = append(pl, comp.Pairs()...)
    }
    
    m.CheckList = pl
}

func (m *MediaStream) orderCheckList() {
    sort.Sort(PairList(m.CheckList))
}

func (m *MediaStream) pruneCheckList() {
    pl := &PairList{}
    pl = append(pl, m.CheckList[0])
    Out:
    for _, cp1 := range m.CheckList {
        for _, cp2 := range pl {
            if cp1 != cp2 {
                pl = append(pl, cp1)
                if len(pl) >= MaxPairs {
                    break Out
                }
            }
        }
    }
    m.CheckList = pl
}

func (m *MediaStream) getCredLocal() *Credential {
    if m.CredLocal.Username == "" {
        return m.parent.CredLocal
    } else {
        return m.CredLocal
    }
}

func (m *MediaStream) getCredRemote() *Credential {
    if m.CredRemote.Username == "" {
        return m.parent.CredRemote
    } else {
        return m.CredRemote
    }
}

// sets one pair per foundation (lowest ComponentId, highest Priority) to Waiting
func (m *MediaStream) setInitialCheckListStates() {
    foundations := make(map[int]*CandidatePair)
    // one pair per foundation, lowest ComponentId, highest Priority
    for _, cp := range m.CheckList {
        if pair, ok := foundations[cp.Foundation()]; ok {
            if cp.ComponentId > pair.ComponentId {
                foundations[cp.Foundation()] = cp
            } else if cp.ComponentId == pair.ComponentId {
                if cp.Priority > pair.Priority {
                    foundations[cp.Foundation()] = cp
                }
            }
        } else {
            foundations[cp.Foundation{}] = cp
        }
    }
    for _, cp := range foundations {
        cp.State = Waiting
    }
}
