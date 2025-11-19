package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gaspoll/internal/config"
	"gaspoll/internal/db"
)

func main() {
	var dir string
	flag.StringVar(&dir, "dir", "migrations", "migrations directory")
	flag.Parse()

	action := "up"
	if flag.NArg() > 0 {
		action = strings.ToLower(flag.Arg(0))
	}

	if action != "up" && action != "down" {
		log.Fatalf("unsupported action %q, want up or down", action)
	}

	cfg := config.Load()
	dbConn, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer dbConn.Close()

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("failed to read migrations dir %s: %v", dir, err)
	}

	var files []string
	for _, e := range entries {
		name := e.Name()
		if action == "up" && strings.HasSuffix(name, ".up.sql") {
			files = append(files, filepath.Join(dir, name))
		}
		if action == "down" && strings.HasSuffix(name, ".down.sql") {
			files = append(files, filepath.Join(dir, name))
		}
	}

	sort.Strings(files)
	if len(files) == 0 {
		fmt.Println("no migrations to apply")
		return
	}

	ctx := context.Background()
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			log.Fatalf("failed to read migration %s: %v", f, err)
		}
		sqlText := string(b)
		fmt.Printf("applying %s...\n", f)
		if _, err := dbConn.ExecContext(ctx, sqlText); err != nil {
			log.Fatalf("failed applying %s: %v", f, err)
		}
	}

	fmt.Println("migrations applied")
}
