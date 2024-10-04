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
    http.HandleFunc("/stream", streamHandler)
    fmt.Println("Starting server on :8080")
    http.ListenAndServe(":8080", nil)
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
    files, err := filepath.Glob(filepath.Join(musicDir, "*.mp3"))
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
