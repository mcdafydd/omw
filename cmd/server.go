// Copyright © 2019 David McPike
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

// Body holds data sent by a web client
type Body struct {
	Args []string `json:"args"`
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Server starts the hotkey-listener and GUI pop-up",
	Long: `Server starts the primary functionality of omw.  It:
	- Creates a server on port 31337 that provides a  with the
	headless Chrome window and provides our HTML/JS GUI

	- Listens for the global hotkey <LEFT_SHIFT>+<RIGHT_SHIFT> that
	triggers the GUI`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("server called")
		if len(args) > 0 {
			fmt.Fprintf(os.Stderr, "Unused arguments provided after server command\n")
			os.Exit(1)
		}
		return Run(args)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// Run starts the HTTP server
func Run(args []string) error {
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)

	r := mux.NewRouter()
	r.HandleFunc("/omw/{command:report}", OmwHandler).Methods("GET").Queries("start", "{start}", "end", "{end}", "format", "{format}")
	r.HandleFunc("/omw/{command:report}", OmwHandler).Methods("GET").Queries("start", "{start}", "end", "{end}")
	r.HandleFunc("/omw/{command}", OmwHandler).Methods("GET")
	r.HandleFunc("/omw/{command}", OmwHandler).Methods("OPTIONS")
	r.HandleFunc("/omw/{command}", OmwHandler).Methods("POST")
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(setCorbHeaderMiddleware)
	r.Use(setCorsOriginMiddleware)

	port := os.Getenv("OMW_PORT")
	if port == "" {
		port = "31337"
	}
	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("127.0.0.1:%s", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go func() {
		log.Println(fmt.Sprintf("Listening on http://%s", srv.Addr))
		log.Fatal(srv.ListenAndServe())
	}()

	<-sigc

	log.Println("\nShutting down the server...")

	srv.Shutdown(context.Background())

	return nil
}

func setCorbHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func setCorsOriginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:31337")
		w.Header().Set("Vary", "Origin")
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// OmwHandler executes the request and decodes the body
func OmwHandler(w http.ResponseWriter, r *http.Request) {
	body := Body{}
	vars := mux.Vars(r)
	decErr := json.NewDecoder(r.Body).Decode(&body)
	if decErr == io.EOF {
		decErr = nil // ignore EOF errors caused by empty response body
	}
	if r.Method == http.MethodOptions {
		log.Println("Handling preflight request", vars["command"])
		return
	}

	log.Println("Vars", vars)
	log.Println("Body", r.Body)
	switch vars["command"] {
	case "add", "a":
		if err := server.Add(body.Args); err != nil {
			w.WriteHeader(http.StatusConflict)
		}
	case "break", "b":
		if err := server.Add([]string{"break", "**"}); err != nil {
			w.WriteHeader(http.StatusConflict)
		}
	case "edit", "e":
		if err := server.Edit(); err != nil {
			w.WriteHeader(http.StatusConflict)
		}
	case "ignore", "i":
		if err := server.Add([]string{"ignore", "***"}); err != nil {
			w.WriteHeader(http.StatusConflict)
		}
	case "hello", "h":
		if err := server.Hello(); err != nil {
			w.WriteHeader(http.StatusConflict)
		}
	case "report", "r":
		re := regexp.MustCompile(`(?P<date>20[12][0-9]-[0-9][1-9]-[0123][1-9])T(?P<time>[0-9][0-9]:[0-9][0-9]:[0-9][0-9])[-+](?P<tz>[012][0-9]:[034][05])`)
		matchStart := re.FindStringSubmatch(vars["start"])
		matchEnd := re.FindStringSubmatch(vars["end"])
		w.Header().Set("Content-Type", "application/json")
		if matchStart == nil || matchEnd == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		format := "json"
		if vars["format"] == "fc" {
			format = "fc"
		}
		output, err := server.Report(vars["start"], vars["end"], format)
		if err != nil {
			io.WriteString(w, err.Error())
		} else {
			io.WriteString(w, output)
		}
	case "stretch", "s":
		if err := server.Stretch(); err != nil {
			w.WriteHeader(http.StatusConflict)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	return
}
