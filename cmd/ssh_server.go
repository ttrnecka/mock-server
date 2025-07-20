package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/spf13/viper"
)

const (
	DUMMY_PASS = "secret"
	PROMPT     = "root> "
	EXIT_CMD   = "exit"
)

var prompt string

func sshServer(port int) error {
	ssh.Handle(defaultHandler)
	setPrompt()
	fmt.Printf("Running mock ssh server for module %s on port %d\nPress Ctrl+C to exit.\n", viper.GetString("module"), port)
	err := ssh.ListenAndServe(fmt.Sprintf(":%d", port), nil,
		ssh.PasswordAuth(defaultPasswordAuth),
	)
	return err
}

// func defaultHandler(s ssh.Session) {
// 	printPrompt(s)
// 	scanner := bufio.NewScanner(s)
// 	scanner.Split(bufio.ScanBytes)
// 	command := []byte{}

// 	for scanner.Scan() {
// 		if scanner.Text() == string('\n') || scanner.Text() == string('\r') {
// 			if string(command) == EXIT_CMD {
// 				s.Exit(0)
// 			}
// 			fmt.Fprintf(s, "\n%s", output(command))
// 			printPrompt(s)
// 			command = []byte{}
// 		} else {
// 			fmt.Fprint(s, scanner.Text())
// 			command = append(command, scanner.Bytes()...)
// 		}
// 	}
// 	if err := scanner.Err(); err != nil {
// 		fmt.Fprintln(os.Stderr, "reading standard input:", err)
// 	}
// }

func defaultHandler(s ssh.Session) {
	// printPrompt(s)
	command := s.RawCommand()

	fmt.Fprint(s, "\n")
	if command == EXIT_CMD {
		s.Exit(0)
	}

	if command != "" {
		// One-shot command mode
		result, code := output([]byte(command))
		fmt.Fprintf(s, "%s\n", result)
		s.Exit(code)
	}

	s.Write([]byte("Welcome to SSH shell. Type 'exit' to quit.\n\n"))
	// reader := bufio.NewReader(s)
	// for {
	// 	printPrompt(s)
	// 	line, err := reader.ReadString('\n')
	// 	if err == io.EOF {
	// 		break
	// 	} else if err != nil {
	// 		s.Stderr().Write([]byte("Error reading input.\n"))
	// 		break
	// 	}

	// 	input := strings.TrimSpace(line)
	// 	if input == "" {
	// 		continue
	// 	}

	// 	if input == EXIT_CMD {
	// 		s.Write([]byte("Goodbye!\n"))
	// 		break
	// 	}

	// 	result, _ := output([]byte(input))
	// 	fmt.Fprintf(s, "%s\n", result)
	// }
	// s.Exit(0)
	scanner := bufio.NewScanner(s)
	scanner.Split(bufio.ScanBytes)
	cmd := []byte{}
	printPrompt(s)
	for scanner.Scan() {
		if scanner.Text() == string('\n') || scanner.Text() == string('\r') {
			fmt.Fprint(s, "\n")
			if string(cmd) == EXIT_CMD {
				s.Exit(0)
			}
			result, _ := output(cmd)
			fmt.Fprintf(s, "%s\n", result)
			// s.Exit(code)
			printPrompt(s)
			cmd = []byte{}
		} else {
			fmt.Fprint(s, scanner.Text())
			cmd = append(cmd, scanner.Bytes()...)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		s.Exit(0)
	}
}

func defaultPasswordAuth(ctx ssh.Context, pass string) bool {
	return pass == DUMMY_PASS
}

func output(cmd []byte) (string, int) {
	command := strings.ReplaceAll(string(cmd), " ", "_")
	cmd_file := fmt.Sprintf("%s/%s/commands/%s", MODULES_FOLDER, viper.GetString("module"), command)
	if _, err := os.Stat(cmd_file); os.IsNotExist(err) {
		return fmt.Sprintf("%s: command does not exist", cmd), 127
	}

	b, err := os.ReadFile(cmd_file) // just pass the file name
	if err != nil {
		return err.Error(), 1
	}
	return string(b), 0
}

func printPrompt(s ssh.Session) {
	fmt.Fprintf(s, "%s ", prompt)
}

func setPrompt() {
	prompt = viper.GetString("prompt")
	if prompt == "" {
		prompt = PROMPT
	}
}
