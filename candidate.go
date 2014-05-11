package ice

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

type CandidatePair struct {
    Local      Candidate
    Remote     Candidate
    Default    bool
    Valid      bool
    Nominated  bool
    State      byte
    LocalCred  Credential
    RemoteCred Credential
}
