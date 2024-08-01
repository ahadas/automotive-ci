package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "net/http"
    "log"
    "strings"
)

type CommentEvent struct {
	Attributes ObjectAttributes `json:"object_attributes"`
	MR         MergeRequest     `json:"merge_request"`
}

type ObjectAttributes struct {
	Note string `json:"note"`
}

type MergeRequest struct {
        SourceBranch string `json:"source_branch"`
}

type StringSlice []string

func (s *StringSlice) String() string {
    return fmt.Sprint(*s)
}

func (s *StringSlice) Set(value string) error {
    *s = append(*s, value)
    return nil
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "only POST is allowed", http.StatusMethodNotAllowed)
        return
    }

    decoder := json.NewDecoder(r.Body)
    defer r.Body.Close()

    var event CommentEvent
    if err := decoder.Decode(&event); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    log.Printf("event: %+v\n", event)

    cmd := strings.Fields(event.Attributes.Note)
    log.Printf("command: %+v\n", cmd)
    if cmd[0] != "/build" {
        http.Error(w, "request does not start with /build", http.StatusBadRequest)
	return
    }

    var archs StringSlice
    fs := flag.NewFlagSet("", flag.ExitOnError)
    fs.Var(&archs, "arch", "The architectures to build on")
    fs.Parse(cmd[1:])

    allArchs := map[string]string{"amd64": "false", "arm64": "false"}
    for _, arch := range archs {
        allArchs[arch] = "true"
    }

    response := make(map[string]interface{})
    response["archs"] = allArchs
    response["git_revision"] = event.MR.SourceBranch
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", handleRequest)

    port := ":8080"
    log.Printf("Starting server on port %s", port)
    if err := http.ListenAndServe(port, mux); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

