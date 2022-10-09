package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/achristie/save2db/internal/services"
	progress "github.com/achristie/save2db/internal/tui/progress"
	"github.com/achristie/save2db/pkg/platts"
)

func writeToSvc[T platts.Writeable](ctx context.Context, m *Main, ch chan platts.Result[T], svc services.Service, delete bool) {
	count := 0
	errCmd := progress.StatusCmd("ERROR: Please try again.")()
	m.p.Send(progress.StatusCmd("IN PROGRESS")())

	for result := range ch {
		if result.Err != nil {
			config.errorLog.Printf("fetch: %s", result.Err)
			m.p.Send(errCmd)
			m.p.Quit()
			return
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
				config.errorLog.Printf("write: %s", err)
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
		config.errorLog.Print(err)
		m.p.Quit()
	}

	m.p.Send(progress.StatusCmd(fmt.Sprintf("COMPLETE: %s", getCompletionMsg(delete, count)))())

	// make sure progress bar shows 100 before quitting :)
	time.Sleep(time.Second * 1)
	m.p.Quit()
}

func getCompletionMsg(delete bool, count int) string {
	var s string
	if delete && config.dbSelection == "SQLite" {
		if config.dbSelection == "SQLITE" {
			s = fmt.Sprintf("removed [%d] records from [%s/%s]", count, config.dbSelection, config.path)
		} else {

			s = fmt.Sprintf("removed [%d] records from [%s/%s]", count, config.dbSelection, config.dbName)
		}
	} else {
		if config.dbSelection == "SQLite" {
			s = fmt.Sprintf("added [%d] records to [%s/%s]", count, config.dbSelection, config.path)
		} else {

			s = fmt.Sprintf("added [%d] records to [%s/%s]", count, config.dbSelection, config.dbName)
		}
	}
	return s
}
