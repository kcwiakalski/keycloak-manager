package model

import "github.com/alecthomas/kong"

var CLI struct {
	Port   int    `short:"p" help:"Port on which Keycloak Admin Api is available"`
	Host   string `short:"h" help:"Server hosting Keycloak instalation"`
	User   string `short:"u" help:"Username with administrative rights"`
	Pass   string `help:"Password for user with administrative rights. It is highly discuraged to use this flag directly"`
	Realm  string `short:"r" help:"Realm holding user with administrative rights, usually the same as realm that is target for operation"`
	Client struct {
		File   string `short:"f" help:"Path to file with client configuration" default:"client-config.json"`
		Mode   string `short:"m" help:"Indicates what should be done with config file" enum:"diff,apply" default:"diff"`
		Output string `short:"o" help:"For diff flag indicates file name what will hold operations to apply" default:"client-config-change.json"`
	} `cmd help:"Operates on client configuration" default:"1"`
}

var Ctx *kong.Context
