// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/anubhavg-icpl/agni/internal/gui"
)

// parseGUIFlags applies --port / --data-dir (and help) for GUI mode.
// GUI args are not handled by the CLI firecracker flags parser.
func parseGUIFlags(cfg *gui.Config) {
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--help" || arg == "-h":
			gui.PrintHelp()
			os.Exit(0)
		case arg == "--port" || arg == "-p":
			if i+1 >= len(args) {
				fmt.Fprintln(os.Stderr, "--port requires a value")
				os.Exit(2)
			}
			i++
			cfg.Port = args[i]
		case strings.HasPrefix(arg, "--port="):
			cfg.Port = strings.TrimPrefix(arg, "--port=")
		case arg == "--data-dir":
			if i+1 >= len(args) {
				fmt.Fprintln(os.Stderr, "--data-dir requires a value")
				os.Exit(2)
			}
			i++
			cfg.DataDir = args[i]
		case strings.HasPrefix(arg, "--data-dir="):
			cfg.DataDir = strings.TrimPrefix(arg, "--data-dir=")
		case arg == "--gui" || arg == "-g":
			// already handled in main
		default:
			// ignore unknown so --gui stays compatible with mixed invocations
		}
	}
}
