package cmd

// func (m *Main) getDeletes(ctx context.Context, start time.Time) {
// 	// if len(symbols) > 0 {
// 	m.client.GetDeletes(mdc, start, 10000, m.ch)
// 	// } else {
// 	// 	m.client.GetHistoryByMDC(mdc, start, 10000, m.ch)
// 	// }
// 	m.p.Send(tui.StatusUpdater{Name: "Deletes", Status: tui.Status{Category: tui.INPROGRESS, Msg: "In Progress"}})
// }

// // Get Deleted Assessments and remove from `assessments` table
// func getDeletes() {
// 	data := make(chan platts.Result[platts.SymbolCorrection])
// 	client.GetDeletes(start, pageSize, data)
// 	a := []platts.Assessment{}
// 	p.Send(cli.StatusUpdater{Name: "Deletes", Status: cli.Status{Category: cli.INPROGRESS, Msg: "In Progress"}})

// 	for result := range data {
// 		if result.Err != nil {
// 			log.Printf("Error! %s", result.Err)
// 		} else {
// 			res := result.Message
// 			p.Send(cli.ProgressUpdater{Name: "Deletes", Percent: 1 / float64(res.Metadata.TotalPages)})
// 			log.Printf("Deletes: %d records received from page [%d] in [%s] (%d total records).",
// 				len(res.Results), res.Metadata.Page, res.Metadata.QueryTime, res.Metadata.Count)

// 			a = append(a, res.Flatten()...)
// 		}
// 	}
// 	if err := db.Remove(a); err != nil {
// 		log.Printf("Error removing records: %s", err)
// 		p.Send(cli.StatusUpdater{Name: "Deletes", Status: cli.Status{Category: cli.ERROR, Msg: "An error occured, please retry."}})
// 	}
// 	p.Send(cli.StatusUpdater{Name: "Deletes", Status: cli.Status{Category: cli.COMPLETED, Msg: fmt.Sprintf("Complete! Removed [%d records] from [assessments]", len(a))}})
// }
