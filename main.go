package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/izquiratops/bundler-tools/tools/commands"
)

type BundlerFlags struct {
	Entry  string
	Dist   string
	Minify bool
}

type ClearFlags struct {
	Dist string
}

func defineBundlerFlags(fs *flag.FlagSet) *BundlerFlags {
	cf := &BundlerFlags{}

	fs.StringVar(&cf.Entry, "entry", "src/main.js", "Path to the entry point")
	fs.StringVar(&cf.Dist, "dist", "dist/", "Path to the distribution directory")
	fs.BoolVar(&cf.Minify, "minify", false, "Enable JavaScript minification")

	return cf
}

func defineCleanFlags(fs *flag.FlagSet) *BundlerFlags {
	cf := &BundlerFlags{}

	fs.StringVar(&cf.Dist, "dist", "dist/", "Path to the distribution directory")

	return cf
}

func printBundlerHelp() {
	fmt.Println("Usage: build [options] or serve [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  -entry\tPath to the entry point (default: src/main.js)")
	fmt.Println("  -dist\tPath to the distribution directory (default: dist/)")
	fmt.Println("  -minify\tEnable JavaScript minification (default: false)")
}

func printCleanHelp() {
	fmt.Println("Usage: clean [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  -dist\tPath to the distribution directory (default: dist/)")
}

func printHelp() {
	fmt.Println("Usage: <command> [options]")
	fmt.Println("<command> help\tShows the command help message")
	fmt.Println("\nCommands:")
	fmt.Println("  build\tBuild the project")
	fmt.Println("  serve\tServe the project")
	fmt.Println("  clean\tClean the distribution directory")
}

func main() {
	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	buildFlags := defineBundlerFlags(buildCmd)
	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
	serveFlags := defineBundlerFlags(serveCmd)
	cleanCmd := flag.NewFlagSet("clean", flag.ExitOnError)
	cleanFlags := defineCleanFlags(cleanCmd)

	switch os.Args[1] {
	case "build":
		if len(os.Args) > 2 && os.Args[2] == "help" {
			printBundlerHelp()
			return
		}
		buildCmd.Parse(os.Args[2:])
		commands.Build(buildFlags.Entry, buildFlags.Dist, buildFlags.Minify)

	case "serve":
		if len(os.Args) > 2 && os.Args[2] == "help" {
			printBundlerHelp()
			return
		}
		serveCmd.Parse(os.Args[2:])
		commands.Serve(serveFlags.Entry, serveFlags.Dist, serveFlags.Minify)

	case "clean":
		if len(os.Args) > 2 && os.Args[2] == "help" {
			printCleanHelp()
			return
		}
		cleanCmd.Parse(os.Args[2:])
		commands.Clean(cleanFlags.Dist)

	case "help":
		printHelp()

	default:
		log.Fatal("expected an action subcommand (build, serve or clean)")
	}
}
