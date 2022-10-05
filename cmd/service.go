package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/achristie/save2db/pkg/platts"
	progress "github.com/achristie/save2db/pkg/tui/progress"
	"github.com/achristie/save2db/services"
)

func writeToSvc[T platts.Writeable](ctx context.Context, m *Main, ch chan platts.Result[T], svc services.Service, delete bool) {
	count := 0
	errCmd := progress.StatusCmd("ERROR: check the log file for details.")()
	m.p.Send(progress.StatusCmd("IN PROGRESS")())

	for result := range ch {
		if result.Err != nil {
			log.Printf("fetch: %s", result.Err)
			m.p.Send(errCmd)
			m.p.Quit()
		}

		res := *result.Message
		m.p.Send(progress.ProgressMsg(1 / float64(res.GetTotalPages())))

		for _, r := range res.GetResults() {
			var err error
			if delete {
				_, err = svc.Remove(ctx, m.tx, r)
			} else {
				_, err = svc.Add(ctx, m.tx, r)
			}
			if err != nil {
				log.Printf("write: %s", err)
				m.p.Send(errCmd)
				m.p.Quit()
			}

			count += 1
		}
	}

	// all good so proceed with commit
	err := m.tx.Commit()
	if err != nil {
		m.p.Send(errCmd)
		log.Print(err)
		m.p.Quit()
	}

	m.p.Send(progress.StatusCmd(fmt.Sprintf("COMPLETE: %s", getCompletionMsg(delete, count)))())

	// make sure progress bar shows 100 before quitting :)
	time.Sleep(time.Second * 1)
	m.p.Quit()
}

func getCompletionMsg(delete bool, count int) string {
	var s string
	if delete && config.DBSelection == "SQLite" {
		if config.DBSelection == "SQLITE" {
			s = fmt.Sprintf("removed [%d] records from [%s/%s]", count, config.DBSelection, config.Path)
		} else {

			s = fmt.Sprintf("removed [%d] records from [%s/%s]", count, config.DBSelection, config.DBName)
		}
	} else {
		if config.DBSelection == "SQLite" {
			s = fmt.Sprintf("added [%d] records to [%s/%s]", count, config.DBSelection, config.Path)
		} else {

			s = fmt.Sprintf("added [%d] records to [%s/%s]", count, config.DBSelection, config.DBName)
		}
	}
	return s
}
