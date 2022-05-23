package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

var TimeoutError = errors.New("timed out running program")
var TimeoutDuration time.Duration

func execute(source []byte) (string, error) {
	// create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), TimeoutDuration)
	defer cancel()

	// create a command within the context. If the context is done
	// before the command completes, the process will be killed
	cmd := exec.CommandContext(ctx, aspenPath, "--stdin")

	// hook up stdout, stdin and stderr
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	stdin := bytes.NewBuffer(source)
	cmd.Stdin = stdin

	// run the command
	if err := cmd.Run(); err != nil {
		var e *exec.ExitError
		if errors.As(err, &e) {
			switch e.ProcessState.ExitCode() {
			case -1:
				// process was terminated by a kill signal because it timed out
				return "", TimeoutError
			case 1:
				// program exited with error code 1, return stderr
				return stderr.String(), nil
			}
		}
		// an error actually occurred, forward it
		return "", err
	}

	return stdout.String(), nil
}

func run(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		// must be a post request
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log.Print(err)
		http.Error(w, "could not read request", http.StatusInternalServerError)
		return
	}

	output, err := execute(body)

	if err != nil {
		if errors.Is(err, TimeoutError) {
			http.Error(w, TimeoutError.Error(), http.StatusUnprocessableEntity)
			return
		} else {
			log.Print(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	// Allow CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	io.WriteString(w, output)
	w.WriteHeader(http.StatusCreated)
}

var aspenPath string

func main() {
	http.HandleFunc("/run", run)

	port, ok := os.LookupEnv("PLAYGROUND_PORT")
	if !ok {
		port = "8080"
	}

	timeout, ok := os.LookupEnv("PLAYGROUND_TIMEOUT_DURATION")

	if duration, err := time.ParseDuration(timeout); !ok || err != nil {
		TimeoutDuration = 4000 * time.Millisecond
	} else {
		TimeoutDuration = duration
	}

	aspenPath, ok = os.LookupEnv("PLAYGROUND_ASPEN_PATH")
	if !ok {
		log.Fatal("environment variable PLAYGROUND_ASPEN_PATH is not set")
	}

	log.Printf("listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
