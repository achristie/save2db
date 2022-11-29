package platts_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/achristie/save2db/internal/assert"
	"github.com/achristie/save2db/pkg/platts"
)

func TestGetSymbols(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		symRes        platts.SymbolResponse
		pagingDetails platts.PagingDetails
		pageOptions   platts.PageOptions
		wantErr       error
		statusCode    int
	}{
		{name: "Valid Response",
			symRes: platts.SymbolResponse{
				Metadata: platts.SymbolResponseMetadata{TotalPages: 100},
				Results:  []platts.SymbolResponseResults{{Symbol: "awc"}}},
			statusCode: http.StatusOK,
			body: `{
				"metadata": {
					"total_pages": 100,
				},
				"results": [{
					"symbol": "awc",
				}]
			}`,
			// {name: "Invalid Token", },
			// {name: "", },
			// {name: "Page Details", },
			// {name: "Page Options", },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			ts := httptest.NewServer(mux)
			// http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// w.WriteHeader(tt.statusCode)
			// w.Header().Set("Content-Type", "application/json")
			// w.Write([]byte(tt.body))
			// }))
			defer ts.Close()

			mux.HandleFunc("/market-data/reference-data/v3/search", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(tt.body))
			}))

			c := platts.NewClient("test", "test", "test")
			c.BaseURL = ts.URL + "/"

			sr, _, err := c.GetSymbols("brent", time.Now(), platts.PageOptions{})
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, &tt.symRes.Results[0].Symbol, &sr.Results[0].Symbol)
			assert.Equal(t, &tt.symRes.Metadata.TotalPages, &sr.Metadata.TotalPages)
		})
	}

	// fmt.Printf("%+v", sr)
	// fmt.Printf("%+v", pd.Url)

	// var target platts.SymbolResponse

	// ch := platts.FetchAll[platts.SymbolResponse](c, pd, 2)
	// for res := range ch {
	// 	if res.Error != nil {
	// 		t.Error(res.Error)
	// 		continue
	// 	}

	// 	fmt.Printf("%+v\n", len(res.Message.Results))
}

// close(ch)
