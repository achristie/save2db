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

func writeToSvc[T platts.Writeable](ctx context.Context, m *Main, ch chan platts.Result[T], svc services.Service) {
	count := 0
	errCmd := progress.StatusCmd("ERROR: check the log file for details.")()
	m.p.Send(progress.StatusCmd("IN PROGRESS")())

	for result := range ch {
		if result.Err != nil {
			log.Printf("fetch: %s", result.Err)
			m.p.Send(errCmd)
			m.p.Quit()
		}

		res := result.Message
		m.p.Send(progress.ProgressMsg(1 / float64(result.Message.GetTotalPages())))

		for _, r := range res.GetResults() {
			_, err := svc.Add(ctx, m.tx, r)
			if err != nil {
				log.Printf("write: %s", err)
				m.p.Send(errCmd)
				m.p.Quit()
			}

			count += 1
		}
	}

	m.tx.Commit()
	s := fmt.Sprintf("added [%d] records to [%s]", count, config.DBSelection)
	m.p.Send(progress.StatusCmd(fmt.Sprintf("COMPLETE: %s", s))())

	// make sure progress bar shows 100 before quitting :)
	time.Sleep(time.Second * 2)
	m.p.Quit()
}
