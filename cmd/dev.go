package cmd

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Watch for changes, rebuild, and serve the site with live reload",
	Long:  `The dev command watches the source directory for changes, rebuilds the site, and serves it with live reload for local development.`,
	Run:   devFunc,
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin for local development
	},
}

var clients = make(map[*websocket.Conn]bool)
var clientsMu sync.Mutex

const liveReloadScript = `
<script>
  const ws = new WebSocket("ws://localhost:8080/ws");
  ws.onmessage = function(event) {
    if (event.data === "reload") {
      location.reload();
    }
  };
</script>
`

func devFunc(cmd *cobra.Command, args []string) {
	config := GetConfig()

	// check source exists
	if _, err := os.Stat(config.Src); os.IsNotExist(err) {
		log.Error("source directory doesn't exist", "directory", config.Src)
		os.Exit(1)
		return
	}

	opts := makeOpts(config)
	opts.Dev = true
	opts.DevScript = liveReloadScript

	// initial build
	if err := buildSite(config.Src, config.Dst, opts); err != nil {
		log.Error("initial build failed", "error", err)
		os.Exit(1)
		return
	}

	// run the http server
	go func() {
		log.Info("serving...", "port", config.Port)

		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				return
			}

			clientsMu.Lock()
			clients[conn] = true
			clientsMu.Unlock()
		})

		http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(config.Dst))))

		err := http.ListenAndServe(fmt.Sprintf(":%s", config.Port), nil)
		if err != nil {
			log.Error("listen error", "error", err)
			os.Exit(1)
			return
		}
	}()

	// run the watcher
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Error("watcher init error", "error", err)
			os.Exit(1)
		}
		defer watcher.Close()

		err = filepath.Walk(config.Src, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return watcher.Add(path)
			}
			return nil
		})
		if err != nil {
			log.Error("watcher init error", "error", err)
			os.Exit(1)
		}

		var lastBuild time.Time
		const debounceDuration = 500 * time.Millisecond

		for {
			select {
			case event := <-watcher.Events:
				if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0 {
					if time.Since(lastBuild) > debounceDuration {
						if err := buildSite(config.Src, config.Dst, opts); err != nil {
							log.Error("build failed", "error", err)
						} else {
							notifyClients()
						}
						lastBuild = time.Now()
					}
				}
			case err := <-watcher.Errors:
				log.Error("watcher error", "error", err)
			}
		}
	}()

	select {} // run forever
}

// Notify connected clients to reload
func notifyClients() {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte("reload"))
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}

func init() {
	rootCmd.AddCommand(devCmd)
}
