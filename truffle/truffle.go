// Package main maintains a database of user media entries.
//
// To see how to use the tool, run
//
//   truffle help
//
// To see global flags, run
//
//   truffle flags
package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/subcommands"
	"github.com/minkezhang/truffle/api/go/config/utils"
	"github.com/minkezhang/truffle/client/truffle"
	"github.com/minkezhang/truffle/database"
	"github.com/minkezhang/truffle/truffle/commands/add"
	"github.com/minkezhang/truffle/truffle/commands/bump"
	"github.com/minkezhang/truffle/truffle/commands/common"
	"github.com/minkezhang/truffle/truffle/commands/get"
	"github.com/minkezhang/truffle/truffle/commands/patch"
	"github.com/minkezhang/truffle/truffle/commands/search"

	del "github.com/minkezhang/truffle/truffle/commands/delete"
)

var (
	home, _         = os.UserHomeDir()
	defaultFilename = filepath.Join(home, ".truffle/database.textproto")

	errCode = -1

	defaultConfig = utils.Config{
		MAL: utils.MAL{
			// clientID is the publically-known MAL API, per
			// https://github.com/SuperMarcus/myanimelist-api-specification#requests.
			//
			// This API key may also be generated by creating an
			// account and requesting developer access.
			ClientID: "6114d00ca681b7701d1e15fe11a4987e",

			// MaxResults sets a limit on the number of results
			// returned by the MAL API. It appears the MAL search
			// endpoints weights results by a matching heuristic,
			// meaning the first couple of matches are much more
			// likely to be useful for the user.
			SearchMaxResults: 50,

			PopularityCutoff: 2000,
		},
	}
)

var (
	fn   = flag.String("database", defaultFilename, "user database location")
	mock = flag.Bool("dry_run", false, "do not commit changes to database")
)

func read(fn string) ([]byte, error) {
	if err := os.MkdirAll(filepath.Dir(fn), 0777); err != nil {
		return nil, err
	}
	fp, err := os.OpenFile(fn, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	return ioutil.ReadAll(fp)
}

func write(fn string, data []byte) error {
	return ioutil.WriteFile(fn, data, 0666)
}

func main() {
	subcommands.ImportantFlag("database")
	subcommands.ImportantFlag("dry_run")

	flag.Parse()

	common := common.O{
		Output: subcommands.DefaultCommander.Output,
		Error:  subcommands.DefaultCommander.Error,
	}

	config, err := defaultConfig.PB()
	if err != nil {
		fmt.Fprintf(subcommands.DefaultCommander.Error, "could not import config: %v", err)
		os.Exit(errCode)
	}

	data, err := read(*fn)
	if err != nil {
		fmt.Fprintf(common.Error, "could not read file %v: %v", fn, err)
		os.Exit(errCode)
	}

	tdb, err := truffle.Unmarshal(data)
	if err != nil {
		fmt.Fprintf(common.Error, "could not read database: %v", err)
		os.Exit(errCode)
	}
	db := database.New(tdb, config)

	for _, c := range []subcommands.Command{
		subcommands.HelpCommand(),
		subcommands.FlagsCommand(),
		subcommands.CommandsCommand(),
		get.New(db, common),
		add.New(db, common),
		search.New(db, common),
		patch.New(db, common),
		bump.New(db, common),
		del.New(db, common),
	} {
		subcommands.Register(c, "")
	}

	status := subcommands.Execute(context.Background())

	if status == subcommands.ExitSuccess && !*mock {
		data, err := database.Marshal(db)
		if err != nil {
			fmt.Fprintf(common.Error, "could not marshal database: %v", err)
			os.Exit(errCode)
		}
		if err := write(*fn, data); err != nil {
			fmt.Fprintf(common.Error, "could not write to database: %v", err)
			os.Exit(errCode)
		}
	}

	os.Exit(int(status))
}
