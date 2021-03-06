package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/vrgl117-games/roms-manager/gamelist"
)

func resetVisibility(gamelistFiles []*gamelist.File) {
	for _, gamelistFile := range gamelistFiles {
		for j := range gamelistFile.Games {
			log.WithFields(log.Fields{"rom": gamelistFile.Games[j].RomName}).Debug("resetting visibility")
			gamelistFile.Games[j].Hidden = false
			gamelistFile.Games[j].Reason = ""
		}
		log.WithFields(log.Fields{"path": gamelistFile.ShortPath}).Info("all games were marked as visible")
	}
}

func NewResetVisibilityCmd() *cli.Command {
	return &cli.Command{
		Name:  "reset-visibility",
		Usage: "set all games from a list of gamelist.xml files to visible",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "gamelist",
				Usage:    "path to the <gamelist.xml> file(s)",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("debug") {
				log.SetLevel(log.DebugLevel)
			}

			gamelistFiles := []*gamelist.File{}
			for _, gamelistPath := range c.StringSlice("gamelist") {
				gamelistFile, err := gamelist.New(gamelistPath)
				if err != nil {
					return fmt.Errorf("unable to open: %s %v", gamelistPath, err)
				}

				log.WithFields(log.Fields{"games": len(gamelistFile.Games), "path": gamelistFile.ShortPath}).Info("gamelist loaded")
				gamelistFiles = append(gamelistFiles, gamelistFile)
			}

			resetVisibility(gamelistFiles)

			for _, gamelistFile := range gamelistFiles {
				if err := gamelistFile.Save(); err != nil {
					return fmt.Errorf("unable to save: %s %v", gamelistFile.Path, err)
				}
			}

			return nil
		},
	}
}
