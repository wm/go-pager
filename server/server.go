package main

import (
	"github.com/wm/go_pager"
	"github.com/wm/go_pager/oncall"
	"github.com/codegangsta/martini"
)

func main() {
	config := &go_pager.Config{}
	config.LoadFromFile("server/config.json")

	m := martini.New()

	// Setup middleware
	m.Use(martini.Recovery())
	m.Use(martini.Logger())

	m.Map(config)

	// Setup routes
	r := martini.NewRouter()
	r.Get(`/`, Welcome)
	r.Get(`/oncall/`, oncall.Welcome)
	r.Post(`/oncall/attempt_call`, oncall.AttemptCall)
	r.Post(`/oncall/screen_for_machine`, oncall.ScreenForMachine)
	r.Post(`/oncall/complete_call`, oncall.CompleteCall)

	// Add the router action
	m.Action(r.Handle)
	m.Run()
}

func Welcome() string {
    return `Welcome to go_pager`
}

