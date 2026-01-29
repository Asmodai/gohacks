// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// app.go --- Application hacks.
//
// Copyright (c) 2021-2026 Paul Ward <paul@lisphacker.uk>
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

// * Comments:

// * Package:

package app

// * Imports:

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Asmodai/gohacks/config"
	"github.com/Asmodai/gohacks/events"
	"github.com/Asmodai/gohacks/logger"
	"github.com/Asmodai/gohacks/process"
	"github.com/Asmodai/gohacks/responder"
	"github.com/Asmodai/gohacks/semver"
	"gitlab.com/tozd/go/errors"
)

// * Constants:

const (
	// Time to sleep during main loop so we're a nice neighbour.
	eventLoopSleep time.Duration = 250 * time.Millisecond

	// Responder type for the application object.
	responderType string = "app.Application"
)

// * Code:

// ** Interface:

// Application.
type Application interface {
	// Cannot be invoked once the application has been initialised.
	ParseConfig()

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

	// Set the parent context for the application.
	//
	// Danger.  Setting this while the application is running can cause
	// unintended side effects due to the old context's cancel function
	// being executed.
	//
	// It is advisable to run this prior to initialisation.
	SetContext(context.Context)

	// Return the application's process manager instance.
	ProcessManager() process.Manager

	// Return the application's logger instance.
	Logger() logger.Logger

	// Return the application's configuration.
	Configuration() config.Config

	// Set the callback that will be invoked when the application starts.
	//
	// If not set, then the default startup handler will be invoked.
	//
	// This cannot be set once the application has been initialised.
	SetOnStart(CallbackFn)

	// Set the callback that will be invoked when the application exits.
	//
	// If not set, then the default exit handler will be invoked.
	//
	/// This cannot be set once the application has been initialised.
	SetOnExit(CallbackFn)

	// Set the callback that will be invoked whenever the event loop
	// fires.
	//
	// If not set, then the default main loop callback will be invoked.
	//
	// This cannot be set once the application has been initialised.
	SetMainLoop(MainLoopFn)

	// Is the application running?
	IsRunning() bool

	// Is the application in 'debug' mode.
	IsDebug() bool

	// Add a responder to the application's responder chain.
	AddResponder(responder.Respondable) (responder.Respondable, error)

	// Remove a responder from the application's responder chain.
	RemoveResponder(responder.Respondable) bool

	// Send an event to the application's responder.
	//
	// Event will be consumed by the first responder that handles it.
	SendFirstResponder(events.Event) (events.Event, bool)

	// Send an event to all the application's responders.
	SendAllResponders(events.Event) []events.Event

	// Return the name of the application's responder chain.
	//
	// Implements `Respondable`.
	Type() string

	// Ascertain if any of the application's responders will respond to
	// an event.
	//
	// The first responder found that responds to the event will result
	// in `true` being returned.
	//
	// Implements `Respondable`.
	RespondsTo(events.Event) bool

	// Send an event to the application's responders.
	//
	// The first object that can respond to the event will consume it.
	//
	// Implements `Respondable`.
	Invoke(events.Event) events.Event
}

// ** Types:

// Signal callback function type.
type CallbackFn func(Application)

// Main loop callback function type.
type MainLoopFn func(Application)

// Application implementation.
type application struct {
	mu sync.RWMutex

	config    *Config       // App object configuration.
	appconfig config.Config // User's app configuration.

	lgr  logger.Logger   // Logger instance.
	pmgr process.Manager // Process manager instance.

	onStart    CallbackFn      // Function called on app startup.
	onExit     CallbackFn      // Function called on app exit.
	mainLoop   MainLoopFn      // Application main loop function.
	responders responder.Chain // Responder chain.

	running     atomic.Bool // Is the app running?
	initialised atomic.Bool // Has the app been initialised?

	ctx    context.Context    // Main context.
	cancel context.CancelFunc // Context cancellation function.
}

// ** Methods:

// Parse the application's config (if available).
//
// Cannot be invoked once the application has been initialised.
func (app *application) ParseConfig() {
	if app.initialised.Load() {
		// Already initialised.
		return
	}

	app.mu.Lock()
	defer app.mu.Unlock()

	if app.config != nil {
		app.appconfig.Parse()
	}
}

// Initialises the application object.
//
// This must be called, as it does several things to set up the
// various facilities (such as logging) used by the application.
func (app *application) Init() {
	if !app.initialised.CompareAndSwap(false, true) {
		// Already initialised.
		return
	}

	app.mu.Lock()
	defer app.mu.Unlock()

	if app.ctx == nil {
		panic("Attempt made to initialise application with nil context")
	}

	// Get components from DI.
	app.lgr = logger.MustGetLogger(app.ctx)
	app.pmgr = process.MustGetManager(app.ctx)

	if app.config != nil {
		app.lgr.SetLogFile(app.appconfig.LogFile())
		app.lgr.SetDebug(app.appconfig.IsDebug())
	}

	// Propagate context and logger.
	app.pmgr.SetContext(app.ctx)
	app.pmgr.SetLogger(app.lgr)

	// Install signals.
	app.installSignals()

	// Let the user know we're ready.
	// Safe to use this directly, object is read/write locked.
	app.lgr.Info(
		"Application initialised.",
		"type", "init",
		"name", app.config.Name,
		"version", app.config.Version,
		"commit", app.config.Version.Commit,
	)
}

// Return the application's pretty name.
func (app *application) Name() string {
	app.mu.RLock()
	defer app.mu.RUnlock()

	return app.config.Name
}

// Return the application's version.
func (app *application) Version() *semver.SemVer {
	app.mu.RLock()
	defer app.mu.RUnlock()

	return app.config.Version
}

// Return the application's version control commit identifier.
func (app *application) Commit() string {
	app.mu.RLock()
	defer app.mu.RUnlock()

	return app.config.Version.Commit
}

// Return the application's context.
func (app *application) Context() context.Context {
	app.mu.RLock()
	defer app.mu.RUnlock()

	return app.ctx
}

// Set the parent context for the application.
//
// Danger.  Setting this while the application is running can cause
// unintended side effects due to the old context's cancel function being
// executed.
//
// It is advisable to run this prior to initialisation.
func (app *application) SetContext(ctx context.Context) {
	if app.initialised.Load() {
		app.Logger().Warn(
			"Attempt made to set context after app initialisation",
		)

		return
	}

	app.mu.Lock()
	defer app.mu.Unlock()

	if app.cancel != nil {
		app.cancel()
	}

	app.ctx, app.cancel = context.WithCancel(ctx)
}

// Return the application's process manager instance.
func (app *application) ProcessManager() process.Manager {
	app.mu.RLock()
	defer app.mu.RUnlock()

	return app.pmgr
}

// Return the application's logger instance.
func (app *application) Logger() logger.Logger {
	app.mu.RLock()
	defer app.mu.RUnlock()

	return app.lgr
}

// Return the application's configuration.
func (app *application) Configuration() config.Config {
	app.mu.RLock()
	defer app.mu.RUnlock()

	return app.appconfig
}

// Set the callback that will be invoked when the application starts.
//
// This cannot be set once the application has been initialised.
func (app *application) SetOnStart(callback CallbackFn) {
	if app.initialised.Load() {
		app.Logger().Warn(
			"Attempt made to set 'on start' callback after app initialisation",
		)

		return
	}

	app.mu.Lock()
	defer app.mu.Unlock()

	app.onStart = callback
}

// Set the callback that will be invoked when the application exits.
//
// This cannot be set once the application has been initialised.
func (app *application) SetOnExit(callback CallbackFn) {
	if app.initialised.Load() {
		app.Logger().Warn(
			"Attempt made to set `on exit' callback after app initialisation",
		)

		return
	}

	app.mu.Lock()
	defer app.mu.Unlock()

	app.onExit = callback
}

// Set the callback that will be invoked whenever the event loop fires.
//
// This cannot be set once the application has been initialised.
func (app *application) SetMainLoop(callback MainLoopFn) {
	if app.initialised.Load() {
		app.Logger().Warn(
			"Attempt made to set `on exit' callback after app initialisation",
		)

		return
	}

	app.mu.Lock()
	defer app.mu.Unlock()

	app.mainLoop = callback
}

// Is the application running?
func (app *application) IsRunning() bool {
	return app.running.Load()
}

// Is the application using debug mode?
func (app *application) IsDebug() bool {
	app.mu.RLock()
	defer app.mu.RUnlock()

	return app.appconfig.IsDebug()
}

// Start the application.
func (app *application) Run() {
	if !app.running.CompareAndSwap(false, true) {
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
	if app.running.CompareAndSwap(true, false) {
		app.cancel()
	}
}

// Add a responder to the application's responder chain.
func (app *application) AddResponder(sel responder.Respondable) (responder.Respondable, error) {
	app.mu.Lock()
	defer app.mu.Unlock()

	result, err := app.responders.Add(sel)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

// Remove a responder from the application's responder chain.
func (app *application) RemoveResponder(sel responder.Respondable) bool {
	app.mu.Lock()
	defer app.mu.Unlock()

	return app.responders.Remove(sel)
}

// Send an event to the application's responder.
//
// Event will be consumed by the first responder that handles it.
func (app *application) SendFirstResponder(evt events.Event) (events.Event, bool) {
	app.mu.RLock()
	defer app.mu.RUnlock()

	return app.responders.SendFirst(evt)
}

// Send an event to all the application's responders.
func (app *application) SendAllResponders(evt events.Event) []events.Event {
	app.mu.RLock()
	defer app.mu.RUnlock()

	return app.responders.SendAll(evt)
}

// Return the name of the application's responder chain.
//
// Implements `Respondable`.
func (app *application) Type() string {
	app.mu.RLock()
	defer app.mu.RUnlock()

	return responderType
}

// Ascertain if any of the application's responders will respond to an event.
//
// The first responder found that responds to the event will result in `true`
// being returned.
//
// Implements `Respondable`.
func (app *application) RespondsTo(event events.Event) bool {
	app.mu.RLock()
	defer app.mu.RUnlock()

	return app.responders.RespondsTo(event)
}

// Send an event to the application's responders.
//
// The first object that can respond to the event will consume it.
//
// Implements `Respondable`.
func (app *application) Invoke(event events.Event) events.Event {
	app.mu.RLock()
	defer app.mu.RUnlock()
	result, _ := app.responders.SendFirst(event)

	return result
}

// ** Functions:

// Create a new application.
func NewApplication(cnf *Config) Application {
	cnf.validate()

	// Temptation here is to set up a context and cancel func... don't.
	// If the user fails to provide one, then there will be a panic
	// when Init is called.  This is what we want to have happen.
	//
	// In short, this just sets up the config and default handlers,
	// the rest needs to be manually done by the user.  This is to
	// enforce the use of dependency injection.

	obj := &application{
		responders: *responder.NewChain("Application"),
		config:     cnf,
	}

	if cnf.AppConfig != nil {
		obj.appconfig = config.Init(
			cnf.Name,
			cnf.Version,
			cnf.AppConfig,
			cnf.Validators,
			cnf.RequireCLI,
		)
	}

	obj.onStart = defaultHandler
	obj.onExit = defaultHandler
	obj.mainLoop = defaultMainLoop

	return obj
}

// * app.go ends here.
