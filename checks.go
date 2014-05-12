package ice

import (
    "time"
    )

func (a *Agent) startChecks() {
    a.StartChecks(a.pendingStreams()[0].CheckList)
}

func (a *Agent) StartChecks(list *PairList) {
    c := &checker{a,list}
    a.checkers = append(a.checkers,c)
    go c.start()
}

type checker struct {
    parent  *Agent
    list    *PairList
    running bool
}

func (c *checker) start() {
    c.running = true
    while c.running {
        wait := getWaitInterval()
        if wait > 0 {
            time.Sleep(wait)
        }
        //
    }
}
