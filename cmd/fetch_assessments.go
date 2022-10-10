package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/achristie/save2db/internal/services/assessments"
	"github.com/achristie/save2db/internal/tui/progress"
	"github.com/achristie/save2db/pkg/platts"
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

		// initialize client
		main.client = platts.NewClient(config.Apikey, config.Username, config.Password, config.errorLog, config.infoLog)

		// initialize assessments service
		as, err := assessments.New(ctx, db.GetDB(), config.DBSelection)
		if err != nil {
			fmt.Printf("assessments svc: %s", err)
			os.Exit(1)
		}

		// initialize Channel
		ch := make(chan platts.Result[platts.SymbolHistory])

		// use watchlist if set
		if watchlist != "" {
			wl, err := getWatchlist(ctx, watchlist)
			if err != nil {
				fmt.Printf("get watchlist: %v\n", err)
				os.Exit(1)
			}
			if wl.Metadata.Count == 0 {
				fmt.Printf("could not find watchlist named %q\n", watchlist)
				os.Exit(1)
			}
			symbols = wl.Results[0].Payload
		}

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
			writeToSvc(ctx, &main, ch, as, false)
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

func getWatchlist(ctx context.Context, name string) (*platts.Watchlist, error) {
	c := platts.NewClient(config.Fake, config.Username, config.Password, config.errorLog, config.infoLog)
	wl, err := c.GetWatchlistByName(name)
	if err != nil {
		return nil, err
	}
	return wl, nil
}
