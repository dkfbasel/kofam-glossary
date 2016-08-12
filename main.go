// The find service will scrape all supplied urls for hyperlinks, filter out all document
// links (as specified by the allowed file extension), download those documents
// and calculate the size and md5 checksum for later file comparison
package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strings"
	"time"

	charmaps "golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"

	"github.com/spf13/viper"
)

// define a configuration struct
var (
	OutputFile           string // the file to write the results to
	KofamGlossaryAddress string // the url to the kofam glossary
)

// check if there is some configuration
func init() {

	// look for a file called config
	viper.SetConfigName("config")

	// set some default values for our service
	viper.SetDefault("kofam", "http://kofam.de/de/glossar")
	viper.SetDefault("output", "./kofam-glossar-[TIMESTAMP].csv")

	// read the config file
	viper.ReadInConfig()

	// get the downloadDirectory and allowedFileExtensions into a global variable
	OutputFile = viper.GetString("output")
	KofamGlossaryAddress = viper.GetString("kofam")
}

func main() {

	// fetch all glossary items from kofam
	items, err := FetchGlossaryItemsFrom(KofamGlossaryAddress)
	if err != nil {
		panic("could not fetch glossary items")
	}

	// replace the [TIMESTAMP] variable in output file name
	OutputFile = strings.Replace(OutputFile, "[TIMESTAMP]", time.Now().Format("20060102-150405"), -1)

	// create and open the file for writing
	file, err := os.Create(OutputFile)
	if err != nil {
		panic("could not create the output file: " + OutputFile)
	}

	defer file.Close()

	// use windows 1252 encoding for file
	fileInWindows1252Encoding := encodeToWindows1252(file)

	// create a csv writer
	csvWriter := csv.NewWriter(fileInWindows1252Encoding)

	// set settings for usage with Microsoft Excel
	csvWriter.Comma = ';'
	csvWriter.UseCRLF = true

	// write title string
	csvWriter.Write([]string{"Titel", "Englisch", "Definition", "Quelle", "Gesamter Inhalt", "Url"})

	// write a row for each glossary item
	for _, item := range items {
		csvWriter.Write([]string{
			item.Name, item.English, item.Description,
			item.Source, item.ContentFull, item.Url,
		})
	}

	csvWriter.Flush()
}

func encodeToWindows1252(file *os.File) io.Writer {
	// create a new io.writer for the file
	writeToFile := bufio.NewWriter(file)

	// create an encoder to change the file-encoding to windows format
	writer := transform.NewWriter(writeToFile, charmaps.Windows1252.NewEncoder())
	return writer
}
