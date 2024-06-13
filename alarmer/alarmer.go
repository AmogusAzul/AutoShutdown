package alarmer

import (
	"time"

	"github.com/AmogusAzul/AutoShutdown/config"
)

type Alarmer struct {
	TargetFunc *func()
	AlertFunc  *func()

	TargetTimer *time.Timer
	AlertTimer  *time.Timer

	Config *config.Config

	KillChan chan bool
}

func GetAlarmer(targetFunc *func(), alertFunc *func(), config *config.Config) *Alarmer {

	alertTime := CalculateTimes(config.AlertAdvance, config.Target.AsDuration())

	a := &Alarmer{
		TargetFunc: targetFunc,
		AlertFunc:  alertFunc,
		Config:     config,

		TargetTimer: time.NewTimer(time.Duration(time.Hour * 1)),
		AlertTimer:  time.NewTimer(time.Until(alertTime)),
		KillChan:    make(chan bool),
	}

	a.TargetTimer.Stop()

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
				a.TargetTimer.Stop()
				a.TargetTimer.Reset(a.Config.AlertAdvance)
				(*a.AlertFunc)()
			case <-a.TargetTimer.C:
				(*a.TargetFunc)()

			case killed = <-a.KillChan:
			}

		}

	}()

}

func (a *Alarmer) UpdateTimers() {

	a.AlertTimer.Stop()
	a.TargetTimer.Stop()

	alertTime := CalculateTimes(a.Config.AlertAdvance, a.Config.Target.AsDuration())

	a.AlertTimer.Reset(time.Until(alertTime))

}

func CalculateTimes(alertAdvanceDuration, targetDuration time.Duration) (alertTime time.Time) {

	now := time.Now()

	startOfTheDay := time.Date(now.Year(), now.Month(), time.Now().Day(), 0, 0, 0, 0, now.Location())

	alertTime = startOfTheDay.Add(targetDuration - alertAdvanceDuration)

	return
}
