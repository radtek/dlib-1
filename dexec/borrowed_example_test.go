// MODIFIED: This file is copied verbatim from Go 1.15.5 os/exec/example_test.go,
// MODIFIED: except for lines marked "MODIFIED".
//
// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dexec_test // MODIFIED

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	exec "github.com/datawire/dlib/dexec" // MODIFIED
)

func ExampleLookPath() {
	path, err := exec.LookPath("fortune")
	if err != nil {
		log.Fatal("installing fortune is in your future")
	}
	fmt.Printf("fortune is available at %s\n", path)
}

func ExampleCommandContext() { // MODIFIED
	cmd := exec.CommandContext(context.Background(), "tr", "a-z", "A-Z") // MODIFIED
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("in all caps: %q\n", out.String())
}

func ExampleCommandContext_environment() { // MODIFIED
	cmd := exec.CommandContext(context.Background(), "prog") // MODIFIED
	cmd.Env = append(os.Environ(),
		"FOO=duplicate_value", // ignored
		"FOO=actual_value",    // this value is used
	)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func ExampleCmd_Output() {
	out, err := exec.CommandContext(context.Background(), "date").Output() // MODIFIED
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The date is %s\n", out)
}

func ExampleCmd_Run() {
	cmd := exec.CommandContext(context.Background(), "sleep", "1") // MODIFIED
	log.Printf("Running command and waiting for it to finish...")
	err := cmd.Run()
	log.Printf("Command finished with error: %v", err)
}

func ExampleCmd_Start() {
	cmd := exec.CommandContext(context.Background(), "sleep", "5") // MODIFIED
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	log.Printf("Command finished with error: %v", err)
}

func ExampleCmd_StdoutPipe() {
	cmd := exec.CommandContext(context.Background(), "echo", "-n", `{"Name": "Bob", "Age": 32}`) // MODIFIED
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	var person struct {
		Name string
		Age  int
	}
	if err := json.NewDecoder(stdout).Decode(&person); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s is %d years old\n", person.Name, person.Age)
}

func ExampleCmd_StdinPipe() {
	cmd := exec.CommandContext(context.Background(), "cat") // MODIFIED"
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "values written to stdin are passed to cmd's standard input")
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", out)
}

func ExampleCmd_StderrPipe() {
	cmd := exec.CommandContext(context.Background(), "sh", "-c", "echo stdout; echo 1>&2 stderr") // MODIFIED
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	slurp, _ := ioutil.ReadAll(stderr)
	fmt.Printf("%s\n", slurp)

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}

func ExampleCmd_CombinedOutput() {
	cmd := exec.CommandContext(context.Background(), "sh", "-c", "echo stdout; echo 1>&2 stderr") // MODIFIED
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", stdoutStderr)
}

func ExampleCommandContext_timeout() { // MODIFIED
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	if err := exec.CommandContext(ctx, "sleep", "5").Run(); err != nil {
		// This will fail after 100 milliseconds. The 5 second sleep
		// will be interrupted.
	}
}
