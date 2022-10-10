package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"text/template"
)

var (
	recv  *string = flag.String("addr", "addr.csv", "Adressdatei")
	brief *string = flag.String("brief", "serienbrief.tmpl", "Template-Datei")
)

type Recipient struct {
	LastName string
	Fach     string
	Mf       bool
}

func main() {
	flag.Parse()
	recvFile, err := os.Open(*recv)
	if err != nil {
		log.Fatalf("could not open address file %w", err)
	}
	defer recvFile.Close()

	csvReader := csv.NewReader(recvFile)
	tokens, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalf("could not read from address file %w", err)
	}

	addr := []Recipient{}
	for _, line := range tokens {
		if len(line) != 3 {
			continue
		}
		b, err := strconv.ParseBool(line[2])
		if err != nil {
			log.Printf("invalid bool value %w", err)
			b = true
		}
		r := Recipient{
			LastName: line[0],
			Fach:     line[1],
			Mf:       b,
		}
		addr = append(addr, r)
	}

	content, err := os.ReadFile(*brief)
	if err != nil {
		log.Fatalf("could not read template letter %w", err)
	}
	t := template.Must(template.New("Brief").Parse(string(content)))
	for _, a := range addr {
		out, err := os.Create(fmt.Sprintf("%s.txt", a.Fach))
		if err != nil {
			log.Printf("could not create letter %w", err)
		}
		defer out.Close()

		err = t.Execute(out, a)
		if err != nil {
			log.Printf("error in template %w", err)
		}
	}
}
