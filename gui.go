// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// Copyright 2024 Anubhav Gain. All Rights Reserved.
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

//go:build gui

package main

import (
	"embed"

	"github.com/firecracker-microvm/firectl/internal/gui"
)

//go:embed frontend/build
var assets embed.FS

// runGUI starts the GUI mode with embedded frontend
func runGUI() {
	cfg := gui.DefaultConfig()
	launcher := gui.NewLauncher(cfg)

	if err := launcher.Start(); err != nil {
		panic(err)
	}
	defer launcher.Stop()

	launcher.PrintBanner()

	// In full GUI mode, we would launch Wails with embedded assets here
	// For now, just run the API server
	// TODO: Integrate Wails for native desktop window

	launcher.WaitForShutdown()
}

// printGUIHelp prints help for GUI mode
func printGUIHelp() {
	gui.PrintHelp()
}
