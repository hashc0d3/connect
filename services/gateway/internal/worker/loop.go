package worker

import (
	"context"
	"log"
	"time"

	"github.com/connect/gateway/internal/sms"
)

func Run(ctx context.Context, service *sms.Service) {
	tick := time.NewTicker(15 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			if err := service.DispatchDue(ctx); err != nil {
				log.Printf("очередь: %v", err)
			}
			service.RunBirthdays(ctx)
		}
	}
}
