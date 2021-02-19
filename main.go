package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	match "github.com/alexpantyukhin/go-pattern-match"

	"k0s.io/pkg/cli/agent"
	"k0s.io/pkg/cli/bcrypt"
	"k0s.io/pkg/cli/buildkite"
	"k0s.io/pkg/cli/caddy"
	"k0s.io/pkg/cli/chassis"
	"k0s.io/pkg/cli/client"
	"k0s.io/pkg/cli/dohserver"
	"k0s.io/pkg/cli/gitd"
	"k0s.io/pkg/cli/gos"
	"k0s.io/pkg/cli/gost"
	"k0s.io/pkg/cli/hub"
	"k0s.io/pkg/cli/k16s"
	"k0s.io/pkg/cli/mnt"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	exe := strings.TrimSuffix(filepath.Base(os.Args[0]), ".exe")

	osargs := append([]string{exe}, os.Args[1:]...)

	// arg parse using rust-style match
	// https://github.com/ylxdzsw/v2socks/blob/master/src/main.rs
	// https://github.com/alexpantyukhin/go-pattern-match
	match.Match(osargs).

		// hub -> hub
		// agent -> agent
		// client -> client
		When([]interface{}{"dohserver", match.ANY}, func() {
			run(dohserver.Run(osargs[1:]))
		}).
		When([]interface{}{"bcrypt", match.ANY}, func() {
			run(bcrypt.Run(osargs[1:]))
		}).
		When([]interface{}{"k16s", match.ANY}, func() {
			log.Fatalln(k16s.Run(osargs[1:]))
		}).
		When([]interface{}{"mnt", match.ANY}, func() {
			log.Fatalln(mnt.Run(osargs[1:]))
		}).
		When([]interface{}{"gos", match.ANY}, func() {
			run(gos.Run(osargs[1:]))
		}).
		When([]interface{}{"buildkite-agent", match.ANY}, func() {
			run(buildkite.Run(osargs[1:]))
		}).
		When([]interface{}{"caddy", match.ANY}, func() {
			run(caddy.Run(osargs[1:]))
		}).
		When([]interface{}{"chassis", match.ANY}, func() {
			log.Fatalln(chassis.Run(osargs[1:]))
		}).
		When([]interface{}{"client", match.ANY}, func() {
			log.Fatalln(client.Run(osargs[1:]))
		}).
		When([]interface{}{"hub", match.ANY}, func() {
			log.Fatalln(hub.Run(osargs[1:]))
		}).
		When([]interface{}{"hub2", match.ANY}, func() {
			log.Fatalln(hub.Run2(osargs[1:]))
		}).
		When([]interface{}{"agent", match.ANY}, func() {
			log.Fatalln(agent.Run(osargs[1:]))
		}).
		When([]interface{}{"gitd", match.ANY}, func() {
			log.Fatalln(gitd.Run(osargs[1:]))
		}).
		When([]interface{}{"gost", match.ANY}, func() {
			gost.Main(osargs[1:])
		}).

		// conntroll hub -> hub
		// conntroll agent -> agent
		// conntroll client -> client
		// k0s hub -> hub
		// k0s agent -> agent
		// k0s client -> client
		// * hub -> hub
		// * agent -> agent
		// * client -> client
		When([]interface{}{match.ANY, "dohserver", match.ANY}, func() {
			run(dohserver.Run(osargs[2:]))
		}).
		When([]interface{}{match.ANY, "bcrypt", match.ANY}, func() {
			run(bcrypt.Run(osargs[2:]))
		}).
		When([]interface{}{match.ANY, "k16s", match.ANY}, func() {
			log.Fatalln(k16s.Run(osargs[2:]))
		}).
		When([]interface{}{match.ANY, "mnt", match.ANY}, func() {
			log.Fatalln(mnt.Run(osargs[2:]))
		}).
		When([]interface{}{match.ANY, "gos", match.ANY}, func() {
			run(gos.Run(osargs[2:]))
		}).
		When([]interface{}{match.ANY, "buildkite-agent", match.ANY}, func() {
			run(buildkite.Run(osargs[2:]))
		}).
		When([]interface{}{match.ANY, "caddy", match.ANY}, func() {
			run(caddy.Run(osargs[2:]))
		}).
		When([]interface{}{match.ANY, "chassis", match.ANY}, func() {
			log.Fatalln(chassis.Run(osargs[2:]))
		}).
		When([]interface{}{match.ANY, "client", match.ANY}, func() {
			log.Fatalln(client.Run(osargs[2:]))
		}).
		When([]interface{}{match.ANY, "hub", match.ANY}, func() {
			log.Fatalln(hub.Run(osargs[2:]))
		}).
		When([]interface{}{match.ANY, "hub2", match.ANY}, func() {
			log.Fatalln(hub.Run2(osargs[2:]))
		}).
		When([]interface{}{match.ANY, "agent", match.ANY}, func() {
			log.Fatalln(agent.Run(osargs[2:]))
		}).
		When([]interface{}{match.ANY, "gitd", match.ANY}, func() {
			log.Fatalln(gitd.Run(osargs[2:]))
		}).
		When([]interface{}{match.ANY, "gost", match.ANY}, func() {
			gost.Main(osargs[2:])
		}).

		// k0s -> client
		// k0s hub -> hub
		// k0s agent -> agent
		When([]interface{}{"k0s", match.ANY}, func() {
			log.Fatalln(client.Run(osargs[1:]))
		}).

		// conntroll -> usage
		When(match.ANY, usage).
		Result()
}

func run(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func usage() {
	fmt.Println(`please specify one of the subcommands: 
- agent
- hub
- client`)
}
