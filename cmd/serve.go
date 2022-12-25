package cmd

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/akerl/otgkey/keyboard"
	"github.com/spf13/cobra"
)

var inputMatch = regexp.MustCompile("^[123]$")

func newServeHandler(key string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "missing bearer token", http.StatusForbidden)
			return
		}
		authToken := strings.TrimPrefix(authHeader, "Bearer ")
		if authToken != key {
			http.Error(w, "invalid bearer token", http.StatusUnauthorized)
			return
		}
		newInput := req.URL.Query().Get("id")
		if !inputMatch.MatchString(newInput) {
			http.Error(w, "invalid input id", http.StatusBadRequest)
			return
		}
		d := keyboard.NewDevice("/dev/hidg0")
		args := []string{"scrolllock", "scrolllock", newInput}
		for _, x := range args {
			err := d.SendString(x)
			if err != nil {
				http.Error(w, "failed to send keypress", http.StatusInternalServerError)
				fmt.Printf("failed to send keypress: %s\n", err)
				return
			}
		}
		w.Write([]byte("done"))
		return
	}
}

func serveRunner(_ *cobra.Command, _ []string) error {
	key := os.Getenv("KVMTOGGLE_KEY")
	if key == "" {
		return fmt.Errorf("no KVMTOGGLE_KEY provided")
	}

	h := newServeHandler(key)

	mux := http.NewServeMux()
	mux.HandleFunc("/toggle", h)
	http.ListenAndServe(":8080", mux)
	return nil
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run web server to listen for toggle commands",
	RunE:  serveRunner,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
