package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"text/template"

	"github.com/olekukonko/tablewriter"
)

type tagFlags []string

func (i *tagFlags) String() string {
	return "my string representation"
}

func (i *tagFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var tags tagFlags

// stolen from https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func openBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}

func createAlias(service *Service, storage *Storage) {
	log.Printf("creating service alias %s -> %s (%+v)\n", service.Alias, service.Link, service.Tags)

	if err := storage.AliasService(service); err != nil {
		log.Fatalf("failed to create service alias: %v\n", err)
	}
}

// TODO: add entries into log table
func launch(alias string, storage *Storage) {
	log.Printf("searching for specific alias: %s\n", alias)
	service := storage.byAlias(alias)
	if service != nil {
		log.Printf("opening: %s\n", service.Link)
		openBrowser(service.Link)
		return
	}

	// not found - we assume the last part was a dynamic argument and the rest a
	// prefix we can match in DB
	partials := strings.Split(alias, "/")
	prefix := strings.Join(partials[:len(partials)-1], "/")
	argument := partials[len(partials)-1]
	log.Printf("searching for dynamic alias: %s\n", prefix)
	var url bytes.Buffer
	service = storage.byAlias(prefix)
	if service != nil {
		log.Printf("found it - will use templating\n")
		t, _ := template.New("alias").Parse(service.Link)
		if err := t.Execute(&url, argument); err != nil {
			log.Printf("failed to execute template: %v\n")
		}

		log.Printf("opening: %s\n", url.String())
		openBrowser(url.String())
		return
	}

	log.Printf("could not find any matching alias")
}

func lsCmd(storage *Storage) {
	defer storage.db.Close()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Alias", "Link", "Description", "Tags", "Created"})
	table.SetBorder(false)
	table.SetRowLine(false)
	for i, service := range storage.List(100) {
		table.Append([]string{strconv.Itoa(i), service.Alias, service.Link, service.Description, strings.Join(service.Tags, ", "), service.CreatedAt})
	}

	table.Render()
}

func main() {
	aliasCmd := flag.NewFlagSet("alias", flag.ExitOnError)
	aliasCmd.Var(&tags, "tag", "link tags")
	description := aliasCmd.String("description", "new alias", "description")

	if len(os.Args) < 2 {
		log.Fatalln("No command provided")
	}

	// TODO: make it a CLI argument
	dbPath := "/tmp/goto.1.db"
	log.Printf("initializing data backend [driver=%s path=%v]\n", DB_DRIVER, dbPath)
	storage, err := NewStorage(dbPath, false)
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	// create tables
	if err := storage.Init(); err != nil {
		log.Fatalf("failed to init DB: %v\n", err)
	}

	switch os.Args[1] {
	// tTODO: edit command
	case "alias":
		alias := os.Args[2]
		link := os.Args[3]
		aliasCmd.Parse(os.Args[4:])

		service := &Service{
			Link:        link,
			Alias:       alias,
			Description: *description,
			Tags:        tags,
		}
		createAlias(service, storage)
	case "ls":
		lsCmd(storage)
	default:
		launch(os.Args[1], storage)
	}
}
