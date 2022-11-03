/*
 * app.go --- Application hacks.
 *
 * Copyright (c) 2021-2022 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person
 * obtaining a copy of this software and associated documentation files
 * (the "Software"), to deal in the Software without restriction,
 * including without limitation the rights to use, copy, modify, merge,
 * publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be
 * included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
 * BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
 * ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package app

import (
	"github.com/Asmodai/gohacks/config"
	"github.com/Asmodai/gohacks/logger"
	"github.com/Asmodai/gohacks/process"
	"github.com/Asmodai/gohacks/semver"

	"context"
	"time"
)

const (
	// Time to sleep during main loop so we're a nice neighbour.
	EventLoopSleep time.Duration = 250 * time.Millisecond
)

type OnSignalFn func(*Application) // Signal callback function.
type MainLoopFn func(*Application) // Main loop callback function.

type Application struct {
	config    *Config        // App object configuration.
	appconfig *config.Config // User's app configuration.

	OnStart  OnSignalFn // Function called on app startup.
	OnExit   OnSignalFn // Function called on app exit.
	OnHUP    OnSignalFn // Function called when SIGHUP received.
	OnUSR1   OnSignalFn // Function called when SIGUSR1 received.
	OnUSR2   OnSignalFn // Function called when SIGUSR2 received.
	OnWINCH  OnSignalFn // Function used when SIGWINCH received.
	OnCHLD   OnSignalFn // Function used when SIGCHLD received.
	MainLoop MainLoopFn // Application main loop function.

	running bool               // Is the app running?
	ctx     context.Context    // Main context.
	cancel  context.CancelFunc // Context cancellation function.
}

// Create a new application.
func NewApplication(cnf *Config) *Application {
	cnf.validate()

	// Set up a new parent context for the whole application.
	ctx, cancelFn := context.WithCancel(context.Background())

	a := &Application{
		config:   cnf,
		OnStart:  defaultHandler,
		OnExit:   defaultHandler,
		OnHUP:    defaultOnHUP,
		OnUSR1:   defaultHandler,
		OnUSR2:   defaultHandler,
		OnWINCH:  defaultHandler,
		OnCHLD:   defaultHandler,
		MainLoop: defaultMainLoop,
		ctx:      ctx,
		cancel:   cancelFn,
	}

	if cnf.AppConfig != nil {
		a.appconfig = config.Init(
			cnf.Name,
			cnf.Version,
			cnf.AppConfig,
			cnf.Validators,
		)
	}

	return a
}

func (app *Application) Init() {
	if app.config != nil {
		app.appconfig.Parse()
		app.config.Logger.SetLogFile(app.appconfig.LogFile())
		app.config.Logger.SetDebug(app.appconfig.IsDebug())
	}

	pm := app.ProcessManager()
	if pm != nil {
		app.ProcessManager().SetLogger(app.Logger())
		app.ProcessManager().SetContext(app.Context())
	}

	app.installSignals()

	app.Logger().Info(
		"Application initialised.",
		"type", "init",
		"name", app.config.Name,
		"version", app.config.Version,
		"commit", app.config.Version.Commit,
	)
}

func (app *Application) Name() string {
	return app.config.Name
}

func (app *Application) Version() *semver.SemVer {
	return app.config.Version
}

func (app *Application) Commit() string {
	return app.config.Version.Commit
}

// Return the application's context.
func (app *Application) Context() context.Context {
	return app.ctx
}

func (app *Application) ProcessManager() process.IManager {
	return app.config.ProcessManager
}

func (app *Application) Logger() logger.ILogger {
	return app.config.Logger
}

func (app *Application) Configuration() *config.Config {
	return app.appconfig
}

// Set the `OnStart` callback.
func (app *Application) SetOnStart(fn OnSignalFn) {
	app.OnStart = fn
}

// Set the `OnExit` callback.
func (app *Application) SetOnExit(fn OnSignalFn) {
	app.OnExit = fn
}

// Set the `OnHUP` callback.
func (app *Application) SetOnHUP(fn OnSignalFn) {
	app.OnHUP = fn
}

// Set the `OnUSR1` callback.
func (app *Application) SetOnUSR1(fn OnSignalFn) {
	app.OnUSR1 = fn
}

// Set the `OnUSR2` callback.
func (app *Application) SetOnUSR2(fn OnSignalFn) {
	app.OnUSR2 = fn
}

// Set the `OnWINCH` callback.
func (app *Application) SetOnWINCH(fn OnSignalFn) {
	app.OnWINCH = fn
}

// Set the `OnCHLD` callback.
func (app *Application) SetOnCHLD(fn OnSignalFn) {
	app.OnCHLD = fn
}

// Set the main loop callback.
func (app *Application) SetMainLoop(fn MainLoopFn) {
	app.MainLoop = fn
}

// Is the application running?
func (app *Application) IsRunning() bool {
	return app.running
}

// Is the application using debug mode?
func (app *Application) IsDebug() bool {
	return app.appconfig.IsDebug()
}

// Start the application.
func (app *Application) Run() {
	if app.running {
		return
	}

	app.Logger().Info(
		"Application is running.",
		"type", "run",
	)

	app.loop()
}

// Stop the application.
func (app *Application) Terminate() {
	if !app.running {
		return
	}

	app.running = false
	app.cancel()
}

/* app.go ends here. */
