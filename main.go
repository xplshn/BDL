// main.go // This is the main entrypoint, which calls all the different functions //>
package main

import (
	"fmt"
	"os"
	"runtime"
)

// Repositories contains the URLs for fetching metadata.
var Repositories []string

// MetadataURLs contains the URLs for fetching metadata.
var MetadataURLs []string

// Array for storing a variable that fsearch and info use.
var validatedArch = [3]string{}

func init() {
	switch runtime.GOARCH {
	case "amd64":
		Repositories = append(Repositories, "https://bin.ajam.dev/x86_64_Linux/")
		Repositories = append(Repositories, "https://raw.githubusercontent.com/xplshn/Handyscripts/master/")
		Repositories = append(Repositories, "https://bin.ajam.dev/x86_64_Linux/Baseutils/")
		MetadataURLs = append(MetadataURLs, "https://bin.ajam.dev/x86_64_Linux/METADATA.json")
		MetadataURLs = append(MetadataURLs, "https://bin.ajam.dev/x86_64_Linux/Baseutils/METADATA.json")
		MetadataURLs = append(MetadataURLs, "https://api.github.com/repos/xplshn/Handyscripts/contents")
		validatedArch = [3]string{"x86_64_Linux", "x86_64", "x86_64-Linux"}
	case "arm64":
		Repositories = append(Repositories, "https://bin.ajam.dev/aarch64_arm64_Linux/")
		Repositories = append(Repositories, "https://raw.githubusercontent.com/xplshn/Handyscripts/master/")
		Repositories = append(Repositories, "https://bin.ajam.dev/aarch64_arm64_Linux/Baseutils/")
		MetadataURLs = append(MetadataURLs, "https://bin.ajam.dev/aarch64_arm64_Linux/METADATA.json")
		MetadataURLs = append(MetadataURLs, "https://bin.ajam.dev/aarch64_arm64_Linux/Baseutils/METADATA.json")
		MetadataURLs = append(MetadataURLs, "https://api.github.com/repos/xplshn/Handyscripts/contents")
		validatedArch = [3]string{"aarch64_arm64_Linux", "aarch64_arm64", "aarch64-Linux"}
	default:
		fmt.Println("Unsupported architecture:", runtime.GOARCH)
		os.Exit(1)
	}
}

var installUseCache = true

const RMetadataURL = "https://raw.githubusercontent.com/Azathothas/Toolpacks/main/metadata.json"
const RNMetadataURL = "https://bin.ajam.dev/METADATA.json"
const VERSION = "1.3"

///// YOU MAY CHANGE THESE TO POINT TO ANOTHER PLACE.

const (
	// Cache size limit & handling.
	MaxCacheSize     = 10
	BinariesToDelete = 5
	// TMPDIR is the directory for storing temporary files.
	TEMP_DIR = "/tmp/bigdl_cached" // TODO: Don't hardcode this value.
)

func printHelp() {
	fmt.Println("Usage: bigdl [-vh] {list|install|remove|update|run|info|search|tldr} [args...]")
	fmt.Println("\nOptions:")
	fmt.Println("  -h, --help     Show this help message")
	fmt.Println("  -v, --version  Show the version number")
	fmt.Println("\nCommands:")
	fmt.Println("  list           List all available binaries")
	fmt.Println("  install, add   Install a binary")
	fmt.Println("  remove, del    Remove a binary")
	fmt.Println("  update         Update binaries, by checking their SHA against the repo's SHA.")
	fmt.Println("  run            Run a binary")
	fmt.Println("  info           Show information about a specific binary")
	fmt.Println("  search         Search for a binary - (not all binaries have metadata. Use list to see all binaries)")
	fmt.Println("  tldr           Show a brief description & usage examples for a given program/command")
	fmt.Println("\nExamples:")
	fmt.Println("  bigdl install micro")
	fmt.Println("  bigdl remove bed")
	fmt.Println("  bigdl info jq")
	fmt.Println("  bigdl search editor")
	fmt.Println("  bigdl tldr gum")
	fmt.Println("  bigdl run --verbose neofetch")
	fmt.Println("  bigdl run --silent micro")
	fmt.Println("  bigdl run btop")
	fmt.Println("\nVersion:", VERSION)
}

func main() {
	// Check for flags directly in the main function
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			fmt.Println("bigdl", VERSION)
			os.Exit(0)
		case "--help", "-h":
			printHelp()
			os.Exit(0)
		}
	}

	// If no arguments are received, show the usage text
	if len(os.Args) < 2 {
		fmt.Println("Usage: bigdl [-vh] {list|install|remove|update|run|info|search|tldr} [args...]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "find_url":
		if len(os.Args) < 3 {
			fmt.Println("Usage: bigdl find_url <binary>")
			os.Exit(1)
		}
		findURLCommand(os.Args[2])
	case "list":
		binaries, err := listBinaries()
		if err != nil {
			fmt.Println("Error listing binaries:", err)
			os.Exit(1)
		}
		for _, binary := range binaries {
			fmt.Println(binary)
		}
	case "install", "add":
		if len(os.Args) < 3 {
			fmt.Printf("Usage: bigdl %s <binary> [install_dir] [install_message]\n", os.Args[1])
			os.Exit(1)
		}
		binaryName := os.Args[2]
		var installDir, installMessage string
		if len(os.Args) > 3 {
			installDir = os.Args[3]
		}
		if len(os.Args) > 4 {
			installMessage = os.Args[4]
		}
		err := installCommand(binaryName, []string{installDir}, installMessage)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}
	case "remove", "del":
		if len(os.Args) != 3 {
			fmt.Printf("Usage: bigdl %s <binary>\n", os.Args[1])
			os.Exit(1)
		}
		remove(os.Args[2])
	case "run":
		if len(os.Args) < 3 {
			fmt.Println("Usage: bigdl run [--verbose, --silent] <binary> [args...]")
			os.Exit(1)
		}
		RunFromCache(os.Args[2], os.Args[3:])
	case "tldr":
		if len(os.Args) < 3 {
			fmt.Println("Usage: bigdl tldr <page> [args...]")
			os.Exit(1)
		}
		RunFromCache("tlrc", os.Args[2:])
	case "info":
		if len(os.Args) != 3 {
			fmt.Println("Usage: bigdl info <binary>")
			os.Exit(1)
		}
		binaryName := os.Args[2]
		binaryInfo, err := getBinaryInfo(binaryName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Print the fields
		if binaryInfo.Name != "" {
			fmt.Printf("Name: %s\n", binaryInfo.Name)
		}
		if binaryInfo.Description != "" {
			fmt.Printf("Description: %s\n", binaryInfo.Description)
		}
		if binaryInfo.Repo != "" {
			fmt.Printf("Repo: %s\n", binaryInfo.Repo)
		}
		if binaryInfo.Size != "" {
			fmt.Printf("Size: %s\n", binaryInfo.Size)
		}
		if binaryInfo.SHA256 != "" {
			fmt.Printf("SHA256: %s\n", binaryInfo.SHA256)
		}
		if binaryInfo.B3SUM != "" {
			fmt.Printf("B3SUM: %s\n", binaryInfo.B3SUM)
		}
		if binaryInfo.Source != "" {
			fmt.Printf("Source: %s\n", binaryInfo.Source)
		}
	case "search":
		if len(os.Args) != 3 {
			fmt.Println("Usage: bigdl search <search-term>")
			os.Exit(1)
		}
		searchTerm := os.Args[2]
		fSearch(searchTerm)
	case "update":
		if len(os.Args) > 2 {
			// Bulk update with list of programs to update
			programsToUpdate := os.Args[2:]
			if err := update(programsToUpdate); err != nil {
				fmt.Printf("Error updating programs: %v\n", err)
				os.Exit(1)
			}
		} else {
			// Update by listing files from InstallDir
			if err := update(nil); err != nil {
				fmt.Printf("Error updating programs: %v\n", err)
				os.Exit(1)
			}
		}
	default:
		fmt.Printf("bigdl: Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
