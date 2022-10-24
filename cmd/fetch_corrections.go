package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/achristie/save2db/internal/services/assessments"
	"github.com/achristie/save2db/internal/tui/progress"
	"github.com/achristie/save2db/pkg/platts"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"
)

var fcCmd = &cobra.Command{
	Use:   "corrections",
	Short: "Fetch deleted assessment data",
	Long:  `Fetch corrections (deletes) since t`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// initialize client
		main.client = platts.NewClient(config.Apikey, config.Username, config.Password)

		// setup TUI
		filters := make(map[string]string)
		filters["modifiedDate"] = ">= " + start
		main.p = tea.NewProgram(progress.New("FETCH CORRECTIONS", filters))

		// initialize assessments service
		as, err := assessments.New(ctx, db.GetDB(), config.Database.Name)
		if err != nil {
			fmt.Printf("assessments svc: %s", err)
			os.Exit(1)
		}

		// initialize Channel
		ch := make(chan platts.Result[platts.SymbolCorrection])

		go func() {
			main.getCorrections(ctx, startDate, ch)
			writeToSvc(ctx, &main, ch, as, true)
		}()
		main.p.Start()
	},
}

func init() {
	fetchCmd.AddCommand(fcCmd)
}

func (m *application) getCorrections(ctx context.Context, start time.Time, ch chan platts.Result[platts.SymbolCorrection]) {
	m.client.GetDeletes(start, 10000, ch)
}
