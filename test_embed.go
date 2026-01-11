//go:build ignore

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
