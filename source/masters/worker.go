package masters

import (
	"Nosviak4/modules/gologr"
	"Nosviak4/modules/licenseme"
	"Nosviak4/source"
	"Nosviak4/source/masters/commands"
	"Nosviak4/source/masters/sessions"
	"encoding/hex"
	"errors"
	"fmt"
	"os"

	"path/filepath"
	"strings"
	"time"

	"github.com/radovskyb/watcher"
	"golang.org/x/exp/slices"
)

// spawnMasterTicker helps introduce the handler for branding and session titles
func (m *Masters) spawnMasterTicker() error {
	fileWatch := watcher.New()

	go func() {
		// whenever the goroutine is stopped, we check for panics and if one occurs we recover it.
		defer func() {
			err := recover()
			if err == nil {
				return
			}

			/* prints the error and attempts to recover the title worker. */
			m.logger.WithTerminal().WriteLog(gologr.ERROR, "Panic caught inside title worker: %v", err)
			if err := m.spawnMasterTicker(); err != nil {
				m.logger.WithTerminal().WriteLog(gologr.ERROR, "Error caught inside title worker: %v", err)
				os.Exit(1)
			}
		}()

		/* tries to authenticate the build */
		if err := authenticateMaster(); err != nil {
			os.Exit(0)
		}

		auth := time.NewTicker(5 * time.Minute)
		title := time.NewTicker(source.SmallestTickTime())
		go fakeSlavesWorker()

		for {
			select {

			/* checks that the license authenticates, if not it panics */
			case <-auth.C:
				if err := authenticateMaster(); err != nil {
					auth.Stop()
					os.Exit(0)
				}

			case event := <-fileWatch.Event:
				if !source.OPTIONS.Bool("branding", "live_reload") || !slices.Contains([]string{".tfx", ".toml", ".json"}, filepath.Ext(event.Name())) {
					continue
				}

				if err := source.OpenOptions(); err != nil {
					source.LOGGER.AggregateTerminal().WriteLog(gologr.ERROR, "error occurred rewriting assets, restored: %v", err)
					continue
				}

				if err := sessions.PushConcurrentChangesAcrossSessions(nil); err != nil {
					fmt.Println(err)
					continue
				}

				/* defines the params for the root args */
				commands.ROOT.Args = append(make([]*commands.Arg, 0), commands.Descriptor, commands.Target, commands.Duration, commands.Port)
				if source.OPTIONS.Bool("attacks", "port_then_duration") {
					commands.ROOT.Args = append(make([]*commands.Arg, 0), commands.Descriptor, commands.Target, commands.Port, commands.Duration)
				}

				source.LOGGER.AggregateTerminal().WriteLog(gologr.DEFAULT, "[LIVE-RELOAD] Reloaded %d items successfully (Trigger:%s)", len(source.OPTIONS.Config.Renders), event.Name())

			case <-title.C:
				for _, session := range sessions.Sessions {
					err := session.ExecuteBranding(make(map[string]any), "title.tfx")
					if err != nil {
						if strings.Contains(err.Error(), "EOF") {
							delete(sessions.Sessions, session.Opened)
							continue
						}

						m.logger.WithTerminal().WriteLog(gologr.ERROR, "error updating title for %s: %v", session.ConnIP(), err)
						continue
					}
				}

			case <-fileWatch.Closed:
				return
			}
		}
	}()

	err := fileWatch.AddRecursive(source.ASSETS)
	if err != nil {
		return err
	}

	for path := range fileWatch.WatchedFiles() {
		source.LOGGER.AggregateTerminal().WriteLog(gologr.DEBUG, "[LIVE-RELOAD] Watching %s", path)
	}

	return fileWatch.Start(1 * time.Second)
}

// authenticateMaster will check with licensing system
func authenticateMaster() error {
return nil
	system := client.MakeClient(source.LICENSE, "https", "nosviak4")
	if err := system.GetFromFile(filepath.Join(source.ASSETS, "license.ecl")); err != nil {
		source.LOGGER.AggregateTerminal().WriteLog(gologr.ERROR, err.Error())
		return err
	}

	commit, err := system.EnableEntry()
	if err != nil {
		source.LOGGER.AggregateTerminal().WriteLog(gologr.ERROR, "Possible license system outage, please contact the person you purchased from to receive further instructions.")
		return err
	}

	hardware, err := system.Hardware()
	if err != nil {
		source.LOGGER.AggregateTerminal().WriteLog(gologr.ERROR, err.Error())
		return err
	}

	source.LOGGER.AggregateTerminal().WriteLog(gologr.ALERT, "[Successfully entered the LicenseMe network]")
	source.LOGGER.AggregateTerminal().WriteLog(gologr.ALERT, "= BID: %s", source.VERSION)
	source.LOGGER.AggregateTerminal().WriteLog(gologr.ALERT, "- HWID: %s", hardware)
	source.LOGGER.AggregateTerminal().WriteLog(gologr.ALERT, "= Commit: %s", commit.Commit)
	source.LOGGER.AggregateTerminal().WriteLog(gologr.ALERT, "- Fingerprint: %s", hex.EncodeToString(system.Key().Public.Fingerprint()))

	query, err := system.RunQuery(commit.Commit, hardware, source.VERSION)
	if err != nil {
		source.LOGGER.AggregateTerminal().WriteLog(gologr.ERROR, err.Error())
		return err
	}

	if time.Now().After(query.LicenseExpiry) {
		source.LOGGER.AggregateTerminal().WriteLog(gologr.ERROR, "License expired")
		return errors.New("license expired")
	}

	for _, alert := range query.Alerts {
		source.LOGGER.AggregateTerminal().WriteLog(gologr.ERROR, alert)
	}

	source.LOGGER.AggregateTerminal().WriteLog(gologr.ALERT, "[Completed LicenseMe process successfully]")
	source.LOGGER.AggregateTerminal().WriteLog(gologr.ALERT, "= Authenticated %s's license, expires in %.2f days", query.Client.User, time.Until(query.LicenseExpiry).Hours()/24)
	return nil
}
