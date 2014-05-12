package ice

import (
    "sort"
    )

// collection of components
type MediaStream struct {
    parent      *Agent
    CredLocal   *Credential
    CredRemote  *Credential
    Description *sdp.MediaDescription
    Components  []*Component
    CheckList   []*CandidatePair
    ValidList   []*CandidatePair
}

func (m *MediaStream) initCheckList() {
    m.createCheckList()
    m.orderCheckList()
    m.pruneCheckList()
}

func (m *MediaStream) createCheckList() {
    cl := make([]*CandidatePair)
    for _, comp := range m.Components {
        cl = append(cl, comp.Pairs()...)
    }
    m.CheckList = cl
}

func (m *MediaStream) orderCheckList() {
    sort.Sort(PairList(m.CheckList))
}

func (m *MediaStream) pruneCheckList() {
    
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
