package handlers_test

import (
	"testing"
	"time"

	"github.com/dedis/cothority/example/handlers"
	"github.com/dedis/kyber/suites"
	"github.com/dedis/onet"
	"github.com/dedis/onet/log"
	"github.com/dedis/onet/network"
)

func TestMain(m *testing.M) {
	log.MainTest(m)
}

// Tests a 2-node system
func TestNode(t *testing.T) {
	suite := suites.MustFind("Ed25519")
	local := onet.NewLocalTest(suite)
	nbrNodes := 2
	_, _, tree := local.GenTree(nbrNodes, true)
	//log.Lvl3(tree.Dump())
	defer local.CloseAll()

	pi, err := local.StartProtocol("ExampleHandlers", tree)
	if err != nil {
		t.Fatal("Couldn't start protocol:", err)
	}
	protocol := pi.(*handlers.ProtocolExampleHandlers)
	timeout := network.WaitRetry * time.Duration(network.MaxRetryConnect*nbrNodes*2) * time.Millisecond
	select {
	case children := <-protocol.ChildCount:
		log.Lvl2("Instance 1 is done")
		if children != nbrNodes {
			t.Fatal("Didn't get a child-cound of", nbrNodes)
		}
	case <-time.After(timeout):
		t.Fatal("Didn't finish in time")
	}
}
