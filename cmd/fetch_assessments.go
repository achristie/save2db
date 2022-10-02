package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/achristie/save2db/pkg/platts"
	"github.com/achristie/save2db/pkg/tui/progress"
	"github.com/achristie/save2db/services"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"
)

var faCmd = &cobra.Command{
	Use:   "assessments",
	Short: "Fetch assessment data",
	Long:  `Fetch assessments either by MDC (Market Data category) or Symbol(s) since t`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// initialize assessments service
		as, err := services.NewAssessmentsService(ctx, db.GetDB(), config.DBSelection)
		if err != nil {
			fmt.Printf("assessments svc: %s", err)
			os.Exit(1)
		}

		// initialize Channel
		ch := make(chan platts.Result[platts.SymbolHistory])

		// setup TUI
		filters := make(map[string]string)
		if mdc != "" {
			filters["MDC"] = "= " + mdc
		} else {
			filters["Symbol"] = "in " + strings.Join(symbols, ", ")
		}
		filters["modifiedDate"] = ">= " + start
		main.p = tea.NewProgram(progress.New("FETCH ASSESSMENTS", filters))

		// fetch and store
		go func() {
			main.getAssessments(ctx, mdc, symbols, startDate, ch)
			writeToSvc(ctx, &main, ch, as)
		}()

		// start TUI
		main.p.Start()
	},
}

func init() {
	fetchCmd.AddCommand(faCmd)
}

func (m *Main) getAssessments(ctx context.Context, mdc string, symbols []string, start time.Time, ch chan platts.Result[platts.SymbolHistory]) {
	if len(symbols) > 0 {
		m.client.GetHistoryBySymbol(symbols, start, 10000, ch)
	} else {
		m.client.GetHistoryByMDC(mdc, start, 10000, ch)
	}
}
