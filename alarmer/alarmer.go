package alarmer

import (
	"fmt"
	"time"

	"github.com/AmogusAzul/AutoShutdown/config"
)

type Alarmer struct {
	TargetFunc func()
	AlertFunc  func()

	TargetTimer time.Timer
	AlertTimer  time.Timer

	Config *config.Config

	KillChan chan bool
}

func GetAlarmer(targetFunc func(), alertFunc func(), config *config.Config) *Alarmer {

	alertTime, targetTime := CalculateTimes(config.AlertAdvance, config.Target.AsDuration())

	a := &Alarmer{
		TargetFunc: targetFunc,
		AlertFunc:  alertFunc,
		Config:     config,

		TargetTimer: *time.NewTimer(time.Until(targetTime)),
		AlertTimer:  *time.NewTimer(time.Until(alertTime)),
		KillChan:    make(chan bool),
	}

	return a
}

func (a *Alarmer) Activate() {

	go func() {

		var killed bool

		for {

			if killed {
				break
			}

			select {
			case <-a.Config.UpdateChan:
				a.UpdateTimers()

			case <-a.AlertTimer.C:
				a.AlertFunc()
			case <-a.TargetTimer.C:
				a.TargetFunc()

			case killed = <-a.KillChan:
			}

		}

	}()

}

func (a *Alarmer) UpdateTimers() {

	a.AlertTimer.Stop()
	a.TargetTimer.Stop()

	alertTime, targetTime := CalculateTimes(a.Config.AlertAdvance, a.Config.Target.AsDuration())

	a.AlertTimer.Reset(time.Until(alertTime))
	a.TargetTimer.Reset(time.Until(targetTime))

	fmt.Println("Updated to:", a)

}

func CalculateTimes(alertAdvanceDuration, targetDuration time.Duration) (alertTime, targetTime time.Time) {

	now := time.Now()

	startOfTheDay := time.Date(now.Year(), now.Month(), time.Now().Day(), 0, 0, 0, 0, now.Location())

	targetTime = startOfTheDay.Add(targetDuration)

	alertTime = startOfTheDay.Add(targetDuration - alertAdvanceDuration)

	return
}
