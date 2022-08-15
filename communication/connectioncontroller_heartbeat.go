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
			localEntity := c.localDevice.EntityByType(model.EntityTypeType(model.EntityTypeEnumTypeCEM))

			// we could have multiple subscriptions, e.g. if they are coming in for local client and server roles (which is wrong, but anyway)
			for _, item := range c.subscriptionEntries {
				// check if this is a subscription to a local devicediagnosis feature
				lfType, err := c.featureTypeForAddress(localEntity, item.ServerAddress)
				if err != nil {
					continue
				}

				if lfType != model.FeatureTypeEnumTypeDeviceDiagnosis {
					continue
				}

				senderAddr = item.ServerAddress
				destinationAddr = item.ClientAddress
				timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.9Z")

				res := []model.CmdType{{
					DeviceDiagnosisHeartbeatData: &model.DeviceDiagnosisHeartbeatDataType{
						Timestamp:        &timestamp,
						HeartbeatCounter: c.heartBeatCounter(),
						HeartbeatTimeout: &heartBeatTimeout,
					},
				}}

				// err := ctx.Notify(lf.GetAddress(), rf.GetAddress(), res)
				err = ctx.Notify(senderAddr, destinationAddr, res)
				if err != nil {
					c.log.Println("ERROR sending heartbeat: ", err)
					// TODO: when a connection is closed, we shouldn't get here
					return
				}
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
