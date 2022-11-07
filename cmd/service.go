package cmd

import (
	"fmt"
	"time"

	"github.com/achristie/save2db/internal/services"
	progress "github.com/achristie/save2db/internal/tui/progress"
	"github.com/achristie/save2db/pkg/platts"
)

func writeToSvc[T platts.Writeable](m *Application, ch chan platts.Result[T], svc services.Service, delete bool) {
	count := 0
	errCmd := progress.StatusCmd("ERROR: Please try again.")()
	m.p.Send(progress.StatusCmd("IN PROGRESS")())

	for result := range ch {
		if result.Err != nil {
			// config.errorLog.Printf("fetch: %s", result.Err)
			m.p.Send(errCmd)
			m.p.Quit()
			return
		}

		res := *result.Message
		m.p.Send(progress.ProgressMsg(1 / float64(res.GetTotalPages())))

		for _, r := range res.GetResults() {
			var err error
			if delete {
				_, err = svc.Remove(m.tx, r)
			} else {
				_, err = svc.Add(m.tx, r)
			}
			if err != nil {
				// config.errorLog.Printf("write: %s", err)
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
		// config.errorLog.Print(err)
		m.p.Quit()
	}

	m.p.Send(progress.StatusCmd(fmt.Sprintf("COMPLETE: %s", getCompletionMsg(delete, count)))())

	// make sure progress bar shows 100 before quitting :)
	time.Sleep(time.Second * 1)
	m.p.Quit()
}

func getCompletionMsg(delete bool, count int) string {
	switch delete {
	case true:
		return fmt.Sprintf("removed [%d] records from [%s]", count, config.Database.Name)
	default:
		return fmt.Sprintf("added [%d] records to [%s]", count, config.Database.Name)
	}
}
