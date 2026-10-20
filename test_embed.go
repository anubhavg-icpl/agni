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

//go:build ignore

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
	"embed"
	"fmt"
	"io/fs"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Check root of embed
	fmt.Println("=== Root of embed.FS ===")
	entries, _ := fs.ReadDir(assets, ".")
	for _, e := range entries {
		fmt.Printf("  %s (dir: %v)\n", e.Name(), e.IsDir())
	}

	// Check using fs.Sub
	fmt.Println("\n=== After fs.Sub 'frontend/dist' ===")
	distFS, err := fs.Sub(assets, "frontend/dist")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	entries2, _ := fs.ReadDir(distFS, ".")
	for _, e := range entries2 {
		fmt.Printf("  %s (dir: %v)\n", e.Name(), e.IsDir())
	}

	// List _app directory
	fmt.Println("\n=== List _app directory ===")
	entries3, err := fs.ReadDir(distFS, "_app")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		for _, e := range entries3 {
			fmt.Printf("  %s (dir: %v)\n", e.Name(), e.IsDir())
		}
	}

	// Try to read a JS file
	fmt.Println("\n=== Try reading _app/immutable/entry/start.By4TUTd6.js ===")
	content, err := fs.ReadFile(distFS, "_app/immutable/entry/start.By4TUTd6.js")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Success! Size: %d bytes\n", len(content))
	}
}
