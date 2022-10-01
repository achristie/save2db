package cmd

import (
	"context"
	"log"

	"github.com/achristie/save2db/pkg/platts"
	"github.com/achristie/save2db/pkg/tui/fetch"
	"github.com/achristie/save2db/services"
)

func writeToSvc[T platts.Writeable](ctx context.Context, m *Main, ch chan platts.Result[T], svc services.Service) {
	count := 0

	for result := range ch {
		if result.Err != nil {
			log.Printf("fetch: %s", result.Err)
			// m.p.Send(tui.StatusUpdater{Name: "Symbols", Status: tui.Status{Category: tui.ERROR, Msg: fmt.Sprint(result.Err)}})
			m.p.Quit()
		}

		res := result.Message
		m.p.Send(fetch.ProgressMsg(1 / float64(result.Message.GetTotalPages())))

		for _, r := range res.GetResults() {
			_, err := svc.Add(ctx, m.tx, r)
			if err != nil {
				log.Printf("write: %s", err)
				m.p.Quit()
			}

			count += 1
		}
	}

	// m.p.Send(tui.StatusUpdater{Name: "Symbols", Status: tui.Status{Category: tui.COMPLETED, Msg: fmt.Sprintf("Complete! Added [%d records] to [assessments]", count)}})
	m.tx.Commit()
	// m.p.Quit()
}
