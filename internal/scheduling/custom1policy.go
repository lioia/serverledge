package scheduling

import (
	"github.com/grussorusso/serverledge/internal/function"
	"github.com/grussorusso/serverledge/internal/node"
)

type Custom1Policy struct {
}

func (p *Custom1Policy) Init() {
}

func (p *Custom1Policy) OnCompletion(r *scheduledRequest) {

}

func (p *Custom1Policy) OnArrival(r *scheduledRequest) {

	containerID, err := node.AcquireWarmContainer(r.Fun)
	if err == nil {
		execLocally(r, containerID, true)
	} else if handleColdStart(r) {
		return
	} else if r.CanDoOffloading && r.RequestQoS.Class == function.HIGH_PERFORMANCE {
		handleCloudOffload(r)
	} else if r.CanDoOffloading {
		url := handleEdgeOffloading(r)
		if url != "" {
			handleOffload(r, url)
		} else {
			dropRequest(r)
		}
	} else {
		dropRequest(r)
	}
}
