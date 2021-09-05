/*
 * app.go --- Application hacks.
 *
 * Copyright (c) 2021 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 */

package app

import (
	"github.com/Asmodai/gohacks/config"
	"github.com/Asmodai/gohacks/di"
	"github.com/Asmodai/gohacks/process"
	"github.com/Asmodai/gohacks/types"

	"log"
	"os"
	"os/signal"
	"syscall"
)

type OnExitFn func(os.Signal)
type MainLoopFn func()

type Application struct {
	Name    string
	Version string

	OnExit   OnExitFn
	MainLoop MainLoopFn

	running      bool
	exitMainLoop chan os.Signal
	dism         *di.Service
}

func DefaultOnExit(junk os.Signal) {
}

func DefaultMainLoop() {
}

func NewApplication(name string, version string) *Application {
	if name == "" {
		name = "<anonymous>"
	}

	if version == "" {
		version = "<local>"
	}

	a := &Application{
		Name:         name,
		Version:      version,
		OnExit:       DefaultOnExit,
		MainLoop:     DefaultMainLoop,
		exitMainLoop: make(chan os.Signal),
	}

	if err := a.init(); err != nil {
		log.Panic(err.Error())
	}

	return a
}

func (app *Application) InitConfig(confname string, confobj interface{}, fns config.ValidatorsMap) error {
	if app.dism == nil {
		return types.NewError(
			"APPLICATION",
			"Unable to find service manager.",
		)
	}

	_, found := app.dism.Get(confname)
	if found {
		return types.NewError(
			"APPLICATION",
			"Application configuration is already registered.",
		)
	}

	app.dism.Add(confname, confobj)

	cnf, err := config.InitWithDI(app.Name, app.Version, confname, fns)
	if err != nil {
		return err
	}
	cnf.Parse()

	return nil
}

func (app *Application) SetOnExit(fn OnExitFn) {
	app.OnExit = fn
}

func (app *Application) SetMainLoop(fn MainLoopFn) {
	app.MainLoop = fn
}

func (app *Application) IsRunning() bool {
	return app.running == true
}

func (app *Application) init() error {
	app.dism = di.GetInstance()
	if app.dism == nil {
		return types.NewErrorAndLog(
			"APPLICATION",
			"Could not locate DI service manager.",
		)
	}

	_, found := app.dism.Get("ProcMgr")
	if !found {
		app.dism.Add("ProcMgr", process.NewManager())
	}

	app.installSignals()

	return nil
}

func (app *Application) loop() {
	var sig os.Signal

	pmgr, found := app.dism.Get("ProcMgr")
	if !found {
		log.Panic("Could not locate process manager!")
	}

	app.running = true
	for app.running == true {
		select {
		case sig = <-app.exitMainLoop:
			app.running = false

		default:
		}

		app.MainLoop()
	}

	pmgr.(*process.Manager).StopAll()
	app.OnExit(sig)
	log.Printf("APPLICATION: Terminating.")
}

func (app *Application) installSignals() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Printf("APPLICATION: Received '%v' signal.", sig)
		app.exitMainLoop <- sig
	}()
}

func (app *Application) Run() {
	log.Printf("APPLICATION: Starting %s (%v).", app.Name, app.Version)

	app.loop()
}

/* app.go ends here. */
