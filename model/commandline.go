package model

type CLI struct {
	Port    int      `short:"p" help:"Port on which Keycloak Admin Api is available"`
	Server  string   `short:"s" help:"Server hosting Keycloak instalation"`
	User    string   `short:"u" help:"Username with administrative rights"`
	Pass    string   `help:"Password for user with administrative rights.If password is not provided via command-line user will be prompted for it. It is highly discuraged to use this flag directly."`
	Realm   string   `short:"r" help:"Realm holding user with administrative rights, usually the same as realm that is target for operation"`
	Version struct{} `cmd help:"Print tool version"`

	Client struct {
		File   string `short:"f" help:"Path to file with client configuration" default:"client-config.json"`
		Mode   string `short:"m" help:"Indicates what should be done with config file" enum:"diff,apply" default:"diff"`
		Output string `short:"o" help:"For diff flag indicates file name what will hold operations to apply" default:"client-config-change.json"`
	} `cmd help:"Operates on client configuration" default:"1"`
}
