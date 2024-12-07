package cmd

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"path/filepath"
	"shizuka/shizuka"
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
	var src, dst string
	port := portFlag

	if len(args) == 2 {
		src, dst = args[0], args[1]
	} else if len(args) == 0 {
		config := GetConfig()
		src, dst = config.Src, config.Dst
		if port == "" {
			port = config.Port
		}
	} else {
		cmd.PrintErrln("Usage: build [src] [dst] (provide both or neither)")
		return
	}

	// Ensure source and destination directories exist
	if _, err := os.Stat(src); os.IsNotExist(err) {
		cmd.PrintErr("Source directory (%s) does not exist\n", src)
		os.Exit(1)
		return
	}

	opts := shizuka.BuildOpts{
		Dev:       true,
		DevScript: liveReloadScript,
	}

	// Rebuild the site initially
	if err := buildSite(src, dst, &opts); err != nil {
		cmd.PrintErrln("Failed to build site:", err)
		os.Exit(1)
		return
	}

	// Serve the site
	go func() {
		cmd.Printf("Serving at http://localhost:%s\n", port)

		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				return
			}

			clientsMu.Lock()
			clients[conn] = true
			clientsMu.Unlock()
		})

		http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(dst))))

		err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
		if err != nil {
			cmd.PrintErrln("Failed to start server:", err)
			os.Exit(1)
			return
		}
	}()

	// Watch the source directory
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			cmd.PrintErrln("Failed to create watcher:", err)
			os.Exit(1)
		}
		defer watcher.Close()

		err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return watcher.Add(path)
			}
			return nil
		})
		if err != nil {
			cmd.PrintErrln("Failed to add directory to watcher:", err)
			os.Exit(1)
		}

		var lastBuild time.Time
		const debounceDuration = 500 * time.Millisecond

		for {
			select {
			case event := <-watcher.Events:
				if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0 {
					if time.Since(lastBuild) > debounceDuration {
						cmd.Println("Changes detected, rebuilding...")
						if err := buildSite(src, dst, &opts); err != nil {
							cmd.PrintErrln("Failed to build site:", err)
						} else {
							cmd.Println("Successfully built site")
							notifyClients()
						}
						lastBuild = time.Now()
					}
				}
			case err := <-watcher.Errors:
				cmd.PrintErrln("Watcher error:", err)
			}
		}
	}()

	select {} // Keep the function running
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
	devCmd.Flags().StringVarP(&portFlag, "port", "p", "", "Port to run the server on (overrides config)")
	rootCmd.AddCommand(devCmd)
}
