package communication

import (
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/evcc-io/eebus/spine/model"
)

// TODO heartBeatCounter should be global on CEM level, not on connection level
func (c *ConnectionController) heartBeatCounter() *uint64 {
	i := atomic.AddUint64(&c.heartBeatNum, 1)
	return &i
}

func (c *ConnectionController) sendHearbeat(stopC chan struct{}, d time.Duration) {
	ticker := time.NewTicker(d)
	for {
		select {
		case <-ticker.C:
			var heartBeatTimeout string = "PT4S"

			ctx := c.context(nil)

			var senderAddr, destinationAddr *model.FeatureAddressType

			for _, item := range c.subscriptionEntries {
				if *item.SubscriptionId == 2 { // TODO FIX THIS AS THIS IS DIRTY
					senderAddr = item.ServerAddress
					destinationAddr = item.ClientAddress
				}
			}

			if senderAddr == nil || destinationAddr == nil {
				break
			}

			res := []model.CmdType{{
				DeviceDiagnosisHeartbeatData: &model.DeviceDiagnosisHeartbeatDataType{
					HeartbeatCounter: c.heartBeatCounter(),
					HeartbeatTimeout: &heartBeatTimeout,
				},
			}}

			// err := ctx.Notify(lf.GetAddress(), rf.GetAddress(), res)
			err := ctx.Notify(senderAddr, destinationAddr, res)
			if err != nil {
				c.log.Println("ERROR sending heartbeat: ", err)
			}
		case <-stopC:
			return
		}
	}
}

func (c *ConnectionController) IsHeartbeatClosed() bool {
	select {
	case <-c.stopHeartbeatC:
		return true
	default:
	}

	return false
}

func (c *ConnectionController) stopHeartbeat() {
	c.stopMux.Lock()
	defer c.stopMux.Unlock()

	if c.stopHeartbeatC != nil && !c.IsHeartbeatClosed() {
		close(c.stopHeartbeatC)
	}
}

func (c *ConnectionController) startHeartBeatSend() {
	c.stopHeartbeatC = make(chan struct{})

	go func() {
		c.sendHearbeat(c.stopHeartbeatC, 800*time.Millisecond)
	}()

	// catch signals
	go func() {
		signalC := make(chan os.Signal, 1)
		signal.Notify(signalC, os.Interrupt, syscall.SIGTERM)

		<-signalC // wait for signal
		c.stopHeartbeat()
	}()
}
