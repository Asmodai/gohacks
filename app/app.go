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
	"github.com/Asmodai/gohacks/logger"
	"github.com/Asmodai/gohacks/process"
	"github.com/Asmodai/gohacks/semver"
	"github.com/Asmodai/gohacks/types"

	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	EventLoopSleep time.Duration = 250 * time.Millisecond
)

type OnSignalFn func(*Application)
type MainLoopFn func(*Application)

type Application struct {
	Name    string         // Application name.
	Version *semver.SemVer // Version string.

	OnExit   OnSignalFn // Function called on app exit.
	OnHup    OnSignalFn // Function called when SIGHUP received.
	OnUsr1   OnSignalFn // Function called when SIGUSR1 received.
	OnUsr2   OnSignalFn // Function called when SIGUSR2 received.
	OnWinch  OnSignalFn // Function used when SIGWINCH received.
	MainLoop MainLoopFn // Application main loop function.

	running bool               // Is the app running?
	ctx     context.Context    // Main context.
	cancel  context.CancelFunc // Context cancellation function.

	config *config.Config // Application configuration.

	dism    *di.Service      // DI service manager.
	procmgr *process.Manager // Process manager.
	logger  *logger.Logger
}

// Default `OnExit` calllback.
func DefaultOnExit(*Application) {
}

// Default `OnHup` callback.
func DefaultOnHup(app *Application) {
	app.logger.Info(
		"Default SIGHUP handled invoked.",
	)
}

// Default callback for USR signals.
func DefaultOnUsr(*Application) {
}

// Default callback for WINCH signals.
func DefaultOnWinch(*Application) {
}

// Default main loop callback.
func DefaultMainLoop(*Application) {
}

// Create a new application.
func NewApplication(name string, version *semver.SemVer) *Application {
	if name == "" {
		name = "<anonymous>"
	}

	// Set up a new context for the application here.
	ctx, cancelfn := context.WithCancel(context.Background())

	a := &Application{
		Name:     name,
		Version:  version,
		OnExit:   DefaultOnExit,
		OnHup:    DefaultOnHup,
		OnUsr1:   DefaultOnUsr,
		OnUsr2:   DefaultOnUsr,
		OnWinch:  DefaultOnWinch,
		MainLoop: DefaultMainLoop,
		ctx:      ctx,
		cancel:   cancelfn,
		logger:   nil,
	}

	if err := a.init(); err != nil {
		log.Panic(err.Error())
	}

	return a
}

// Initialise the application's runtime.
func (app *Application) Init() error {
	app.logger.Info(
		"Application is starting.",
		"type", "start",
		"name", app.Name,
		"version", app.Version,
		"commit", app.Version.Commit,
	)

	pm, found := app.dism.Get("ProcMgr")
	if !found {
		pm = process.NewManager()
		pm.(*process.Manager).SetLogger(app.logger)
		pm.(*process.Manager).SetContext(app.ctx)

		app.dism.Add("ProcMgr", pm)
	}
	app.procmgr = pm.(*process.Manager)

	app.installSignals()

	return nil
}

// Initialise the application's configuration.
func (app *Application) InitConfig(confname string, confobj interface{}, fns config.ValidatorsMap) error {
	var err error

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

	app.config, err = config.InitWithDI(app.Name, app.Version, confname, fns)
	if err != nil {
		return err
	}
	app.config.Parse()

	// Set log file options here.
	app.logger.SetLogFile(app.config.LogFile())
	app.logger.SetDebug(app.config.IsDebug())

	return nil
}

// Return the application's context.
func (app *Application) Context() context.Context {
	return app.ctx
}

// Return the application's process manager.
func (app *Application) ProcessManager() *process.Manager {
	return app.procmgr
}

// Set the `OnExit` callback.
func (app *Application) SetOnExit(fn OnSignalFn) {
	app.OnExit = fn
}

// Set the `OnHup` callback.
func (app *Application) SetOnHup(fn OnSignalFn) {
	app.OnHup = fn
}

// Set the `OnUsr1` callback.
func (app *Application) SetOnUsr1(fn OnSignalFn) {
	app.OnUsr1 = fn
}

// Set the `OnUsr2` callback.
func (app *Application) SetOnUsr2(fn OnSignalFn) {
	app.OnUsr2 = fn
}

// Set the main loop callback.
func (app *Application) SetMainLoop(fn MainLoopFn) {
	app.MainLoop = fn
}

// Is the application running?
func (app *Application) IsRunning() bool {
	return app.running == true
}

// Is the application running in debug mode?
func (app *Application) IsDebug() bool {
	return app.config.IsDebug()
}

// Initialise the application.
func (app *Application) init() error {
	app.dism = di.GetInstance()
	if app.dism == nil {
		return types.NewErrorAndLog(
			"APPLICATION",
			"Could not locate DI service manager.",
		)
	}

	alogger, found := app.dism.Get("Logger")
	if !found {
		alogger = logger.NewLogger("")
		app.dism.Add("Logger", alogger)
	}
	app.logger = alogger.(*logger.Logger)

	return nil
}

// Main loop.
func (app *Application) loop() {
	pmgr, found := app.dism.Get("ProcMgr")
	if !found {
		app.logger.Fatal("Could not locate process manager!")
	}

	app.running = true
	for app.running == true {
		select {
		case <-app.ctx.Done():
			// Application context was cancelled.
			app.running = false

		default:
		}

		app.MainLoop(app)

		time.Sleep(EventLoopSleep)
	}

	app.OnExit(app)
	pmgr.(*process.Manager).StopAll()
	app.logger.Info(
		"Application is terminating.",
		"type", "stop",
	)
}

// Install signal handler.
func (app *Application) installSignals() {
	sigs := make(chan os.Signal, 1)

	// We don't care for the following signals:
	signal.Ignore(syscall.SIGURG)

	// Notify when a signal we care for is received.
	signal.Notify(sigs)

	go func() {
		for {
			sig := <-sigs

			switch sig {
			case syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM:
				// Handle termination.
				app.logger.Info(
					"Received signal",
					"signal", sig.String(),
				)
				app.Terminate()
				return

			case syscall.SIGHUP:
				// Handle SIGHUP.
				app.logger.Info(
					"Received signal",
					"signal", sig.String(),
				)
				app.OnHup(app)

			case syscall.SIGWINCH:
				// Handle WINCH.
				// Note: Do not bother logging this one.
				app.OnWinch(app)

			case syscall.SIGUSR1:
				// Handle user-defined signal #1.
				app.logger.Info(
					"Received signal",
					"signal", sig.String(),
				)
				app.OnUsr1(app)

			case syscall.SIGUSR2:
				// Handle user-defined signal #2.
				app.logger.Info(
					"Received signal",
					"signal", sig.String(),
				)
				app.OnUsr2(app)

			default:
				if sig == syscall.SIGURG {
					// This signal is noise, generated by the Go runtime.
					break
				}
				app.logger.Info(
					"Unhandled signal",
					"signal", sig.String(),
				)
			}
		}
	}()
}

// Run the application.
func (app *Application) Run() {
	app.logger.Info(
		"Application is running.",
		"type", "run",
	)

	app.loop()
}

// Terminate the application
func (app *Application) Terminate() {
	if !app.running {
		return
	}

	app.running = false
	app.cancel()
}

/* app.go ends here. */
