package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mcdafydd/omw/backend"
	"github.com/mcdafydd/omw/cmd"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pelletier/go-toml"
)

func main() {
	layoutEvent := "2006-1-2 15:4"
	items := backend.SavedItems{}

	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	oldPath := fmt.Sprintf("%s/%s/omw.log", home, cmd.DefaultDir)
	newPath := fmt.Sprintf("%s/%s/omw.toml", home, cmd.DefaultDir)

	r, err := os.Open(oldPath)
	defer r.Close()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		item := backend.SavedEntry{}
		line := strings.Fields(scanner.Text())
		if len(line) < 2 {
			continue
		}
		ts, err := time.Parse(layoutEvent, strings.Join(line[:2], " "))
		if err != nil {
			continue
		}
		entry := strings.Join(line[2:], " ")
		item.ID = uuid.New().String()
		item.Start = ts
		item.Task = entry
		items.Entries = append(items.Entries, item)
	}
	b, err := toml.Marshal(items)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", string(b))
	fmt.Fprintf(os.Stderr, "\nTo finish migrating your data, save the converted output as %s", newPath)
	return
}
