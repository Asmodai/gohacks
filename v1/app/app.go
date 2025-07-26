// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// app.go --- Application hacks.
//
// Copyright (c) 2021-2025 Paul Ward <paul@lisphacker.uk>
//
// Author:     Paul Ward <paul@lisphacker.uk>
// Maintainer: Paul Ward <paul@lisphacker.uk>
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation files
// (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge,
// publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
// BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
// ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
// mock:yes

package app

import (
	"github.com/Asmodai/gohacks/v1/config"
	"github.com/Asmodai/gohacks/v1/logger"
	"github.com/Asmodai/gohacks/v1/process"
	"github.com/Asmodai/gohacks/v1/semver"

	"context"
	"time"
)

const (
	// Time to sleep during main loop so we're a nice neighbour.
	EventLoopSleep time.Duration = 250 * time.Millisecond
)

// Application.
type Application interface {
	// Initialises the application object.
	//
	// This must be called, as it does several things to set up the
	// various facilities (such as logging) used by the application.
	Init()

	// Run the application.
	//
	// This enters the application's event loop, which will block until
	// the application is subsequently terminated.
	Run()

	// Terminate the application.
	//
	// Breaks out of the event loop, returning control back to the calling
	// function.
	Terminate()

	// Return the application's pretty name.
	Name() string

	// Return the application's version.
	Version() *semver.SemVer

	// Return the application's version control commit identifier.
	Commit() string

	// Return the application's context.
	Context() context.Context

	// Return the application's process manager instance.
	ProcessManager() process.Manager

	// Return the application's logger instance.
	Logger() logger.Logger

	// Return the application's configuration.
	Configuration() config.Config

	// Set the callback that will be invoked when the application starts.
	SetOnStart(OnSignalFn)

	// Set the callback that will be invoked when the application exits.
	//
	// If not set, then the default exit handler will be invoked.
	SetOnExit(OnSignalFn)

	// Set the callback that will be invoked when the application
	// receives a HUP signal.
	SetOnHUP(OnSignalFn)

	// Set the callback that will be invoked when the application
	// receives a USR1 signal.
	SetOnUSR1(OnSignalFn)

	// Set the callback that will be invoked when the application
	// receives a USR2 signal.
	SetOnUSR2(OnSignalFn)

	// Set the callback that will be invoked when the application
	// receives a WINCH signal.
	//
	// Be careful with this, as it will fire whenever the controlling
	// terminal is resized.
	SetOnWINCH(OnSignalFn)

	// Set the callback that will be invoked when the application
	// receives a CHLD signal.
	SetOnCHLD(OnSignalFn)

	// Set the callback that will be invoked whenever the event loop
	// fires.
	SetMainLoop(MainLoopFn)

	// Is the application running?
	IsRunning() bool

	// Is the application in 'debug' mode.
	IsDebug() bool
}

// Signal callback function type.
type OnSignalFn func(Application)

// Main loop callback function type.
type MainLoopFn func(Application)

// Application implementation.
type application struct {
	config    *Config       // App object configuration.
	appconfig config.Config // User's app configuration.

	onStart  OnSignalFn // Function called on app startup.
	onExit   OnSignalFn // Function called on app exit.
	onHUP    OnSignalFn // Function called when SIGHUP received.
	onUSR1   OnSignalFn // Function called when SIGUSR1 received.
	onUSR2   OnSignalFn // Function called when SIGUSR2 received.
	onWINCH  OnSignalFn // Function used when SIGWINCH received.
	onCHLD   OnSignalFn // Function used when SIGCHLD received.
	mainLoop MainLoopFn // Application main loop function.

	running bool               // Is the app running?
	ctx     context.Context    // Main context.
	cancel  context.CancelFunc // Context cancellation function.
}

// Create a new application.
func NewApplication(cnf *Config) Application {
	cnf.validate()

	// Set up a new parent context for the whole application.
	ctx, cancelFn := context.WithCancel(context.Background())

	obj := &application{
		config:   cnf,
		onStart:  defaultHandler,
		onExit:   defaultHandler,
		onHUP:    defaultOnHUP,
		onUSR1:   defaultHandler,
		onUSR2:   defaultHandler,
		onWINCH:  defaultHandler,
		onCHLD:   defaultHandler,
		mainLoop: defaultMainLoop,
		ctx:      ctx,
		cancel:   cancelFn,
	}

	if cnf.AppConfig != nil {
		obj.appconfig = config.Init(
			cnf.Name,
			cnf.Version,
			cnf.AppConfig,
			cnf.Validators,
			true,
		)
	}

	return obj
}

func (app *application) Init() {
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

// Return the application's pretty name.
func (app *application) Name() string {
	return app.config.Name
}

// Return the application's version.
func (app *application) Version() *semver.SemVer {
	return app.config.Version
}

// Return the application's version control commit identifier.
func (app *application) Commit() string {
	return app.config.Version.Commit
}

// Return the application's context.
func (app *application) Context() context.Context {
	return app.ctx
}

// Return the application's process manager instance.
func (app *application) ProcessManager() process.Manager {
	return app.config.ProcessManager
}

// Return the application's logger instance.
func (app *application) Logger() logger.Logger {
	return app.config.Logger
}

// Return the application's configuration.
func (app *application) Configuration() config.Config {
	return app.appconfig
}

// Set the callback that will be invoked when the application starts.
func (app *application) SetOnStart(fn OnSignalFn) {
	app.onStart = fn
}

// Set the callback that will be invoked when the application exits.
func (app *application) SetOnExit(fn OnSignalFn) {
	app.onExit = fn
}

// Set the callback that will be invoked when the application receives a HUP
// signal.
func (app *application) SetOnHUP(fn OnSignalFn) {
	app.onHUP = fn
}

// Set the callback that will be invoked when the application receives a USR1
// signal.
func (app *application) SetOnUSR1(fn OnSignalFn) {
	app.onUSR1 = fn
}

// Set the callback that will be invoked when the application receives a USR2
// signal.
func (app *application) SetOnUSR2(fn OnSignalFn) {
	app.onUSR2 = fn
}

// Set the callback that will be invoked when the application receives a WINCH
// signal.
func (app *application) SetOnWINCH(fn OnSignalFn) {
	app.onWINCH = fn
}

// Set the callback that will be invoked when the application receives a CHLD
// singal.
func (app *application) SetOnCHLD(fn OnSignalFn) {
	app.onCHLD = fn
}

// Set the callback that will be invoked whenever the event loop fires.
func (app *application) SetMainLoop(fn MainLoopFn) {
	app.mainLoop = fn
}

// Is the application running?
func (app *application) IsRunning() bool {
	return app.running
}

// Is the application using debug mode?
func (app *application) IsDebug() bool {
	return app.appconfig.IsDebug()
}

// Start the application.
func (app *application) Run() {
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
func (app *application) Terminate() {
	if !app.running {
		return
	}

	app.running = false
	app.cancel()
}

// app.go ends here.
