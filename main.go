package main

import (
    "fmt"
    "io"
    "math/rand"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "time"
)

var musicDir = "/path/to/mounted/cloud/disk"

func main() {
    http.HandleFunc("/", dirHandler)
    http.HandleFunc("/browse", browseHandler)
    http.HandleFunc("/stream", streamHandler)
    fmt.Println("Starting server on :8080")
    http.ListenAndServe(":8080", nil)
}

func dirHandler(w http.ResponseWriter, r *http.Request) {
    dirs, err := os.ReadDir(musicDir)
    if err != nil {
        http.Error(w, "Unable to read directories", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "text/html")
    fmt.Fprintln(w, "<html><body><h1>Select Directory</h1><ul>")
    for _, dir := range dirs {
        if dir.IsDir() {
            fmt.Fprintf(w, "<li><a href=\"/browse?dir=%s\">%s</a></li>", dir.Name(), dir.Name())
        }
    }
    fmt.Fprintln(w, "</ul></body></html>")
}

func browseHandler(w http.ResponseWriter, r *http.Request) {
    dir := r.URL.Query().Get("dir")
    if dir == "" {
        http.Error(w, "No directory selected", http.StatusBadRequest)
        return
    }

    fullPath := filepath.Join(musicDir, dir)
    items, err := os.ReadDir(fullPath)
    if err != nil {
        http.Error(w, "Unable to read directory", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "text/html")
    fmt.Fprintln(w, "<html><body><h1>Browse Directory</h1><ul>")
    for _, item := range items {
        if item.IsDir() {
            fmt.Fprintf(w, "<li><a href=\"/browse?dir=%s\">%s</a></li>", filepath.Join(dir, item.Name()), item.Name())
        } else {
            fmt.Fprintf(w, "<li>%s</li>", item.Name())
        }
    }
    fmt.Fprintf(w, "</ul><a href=\"/stream?dir=%s\">Stream from this directory</a></body></html>", dir)
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
    dir := r.URL.Query().Get("dir")
    if dir == "" {
        http.Error(w, "No directory selected", http.StatusBadRequest)
        return
    }

    files, err := filepath.Glob(filepath.Join(musicDir, dir, "*.mp3"))
    if err != nil || len(files) == 0 {
        http.Error(w, "No music files found", http.StatusInternalServerError)
        return
    }

    rand.Seed(time.Now().UnixNano())
    file := files[rand.Intn(len(files))]

    cmd := exec.Command("sox", file, "-t", "mp3", "-", "fade", "t", "5")
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        http.Error(w, "Unable to process file", http.StatusInternalServerError)
        return
    }

    if err := cmd.Start(); err != nil {
        http.Error(w, "Unable to start processing", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "audio/mpeg")
    io.Copy(w, stdout)
    cmd.Wait()
}
