package questions

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mattn/go-isatty"
)

func PromptFormattedOptions(text string, def int, options ...string) (int, error) {
	var newOptions []string
	for i := range options {
		newOptions = append(newOptions, fmt.Sprintf("%d. %s\n", i+1, options[i]))
	}
	return PromptOptions(text+"\n", def, newOptions...)
}

func PromptOptions(text string, def int, options ...string) (int, error) {
	if len(options) == 1 {
		return 0, nil
	}

	PrintToTerm(text)
	for _, option := range options {
		PrintToTerm(option)
	}

	defString := ""
	if def >= 0 {
		defString = strconv.Itoa(def + 1)
	}

	for {
		ans, err := Prompt(fmt.Sprintf("Select Number [%s]: ", defString), defString)
		if err != nil {
			return 0, err
		}
		num, err := strconv.Atoi(ans)
		if err != nil {
			PrintfToTerm("Invalid number: %s\n", ans)
			continue
		}

		num--
		if num < 0 || num >= len(options) {
			PrintlnToTerm("Select a number between 1 and", +len(options))
			continue
		}

		return num, nil
	}
}

func PromptBool(text string, def bool) (bool, error) {
	msg := fmt.Sprintf("%s [y/N]: ", text)
	defStr := "n"
	if def {
		msg = fmt.Sprintf("%s [Y/n]: ", text)
		defStr = "y"
	}

	for {
		yn, err := Prompt(msg, defStr)
		if err != nil {
			return false, err
		}

		switch strings.ToLower(yn) {
		case "y":
			return true, nil
		case "n":
			return false, nil
		default:
			fmt.Println("Enter y or n")
		}
	}
}

func PrintToTerm(text ...interface{}) {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Print(text...)
	} else {
		fmt.Fprint(os.Stderr, text...)
	}
}

func PrintlnToTerm(text ...interface{}) {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Println(text...)
	} else {
		fmt.Fprintln(os.Stderr, text...)
	}
}

func PrintfToTerm(msg string, format ...interface{}) {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Printf(msg, format...)
	} else {
		fmt.Fprintf(os.Stderr, msg, format...)
	}
}

func Prompt(text, def string) (string, error) {
	for {
		PrintToTerm(text)
		answer, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			return "", err
		}

		answer = strings.TrimSpace(answer)
		if answer == "" {
			answer = def
		}

		if answer == "" {
			continue
		}

		return answer, nil
	}
}

func PromptOptional(text, def string) (string, error) {
	for {
		PrintToTerm(text)
		answer, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			return "", err
		}

		answer = strings.TrimSpace(answer)
		if answer == "" {
			answer = def
		}

		return answer, nil
	}
}
