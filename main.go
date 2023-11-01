package main

import (
	"bufio"
	"flag"
	"fmt"
	"idaru/url"
	"log"
	"os"
	"strings"
	"time"
)

var _AUTHOR_ = "siriil"
var _VERSION_ = "v1.0"

func main() {
	// Script parameters
	flag_version := flag.Bool("version", false, "Display version")
	flag_v := flag.Bool("v", false, "Alias for -version")

	flag_verbose := flag.Bool("verbose", false, "Display log messages on stderr")
	flag_vb := flag.Bool("vb", false, "Alias for -verbose")

	flag_merge := flag.Bool("merge", false, "Merge all keys param on one query")
	flag_m := flag.Bool("m", false, "Alias for -merge")

	flag_filterParam := flag.Bool("filterParam", false, "Filter queries if it has param")
	flag_fP := flag.Bool("fP", false, "Alias for -filterParam")

	flag_show := flag.Bool("show", false, "Show sitemap")
	flag_sh := flag.Bool("sh", false, "Alias for -show")

	flag_outputJson := flag.Bool("outputJson", false, "Save sitemap to JSON file")
	flag_oJ := flag.Bool("oJ", false, "Alias for -outputJson")

	flag_inputFiles := []string{}
	flag.Var((*ArrayFlagString)(&flag_inputFiles), "inputFile", "Specify an input file")
	flag.Var((*ArrayFlagString)(&flag_inputFiles), "iF", "Alias for -inputFile")

	flag_add := []string{}
	flag.Var((*ArrayFlagString)(&flag_add), "add", "Specify 'key=value' to add (all key is '*')")
	flag.Var((*ArrayFlagString)(&flag_add), "a", "Alias for -add")

	flag_set := []string{}
	flag.Var((*ArrayFlagString)(&flag_set), "set", "Specify 'key=value' to set (all key is '*')")
	flag.Var((*ArrayFlagString)(&flag_set), "s", "Alias for -set")

	flag.Parse()
	// End of script parameters

	*flag_version = (*flag_version || *flag_v)
	if *flag_version {
		Version()
		os.Exit(0)
	}

	// Create a new log.Logger with custom formatting
	var logger *log.Logger
	Nullout := NullWriter{}
	*flag_verbose = (*flag_verbose || *flag_vb)
	if *flag_verbose {
		logger = log.New(os.Stdout, "", 0)
	} else {
		logger = log.New(Nullout, "", 0)
	}

	// Check if URLs have been provided through any input method, including stdin and files
	if len(flag_inputFiles) == 0 && !isStdinAvailable() {
		customError(logger, "You must provide URLs via stdin or files.")
	}

	var scanners []*bufio.Scanner
	// URLs from stdin
	if isStdinAvailable() {
		scanners = append(scanners, bufio.NewScanner(os.Stdin))
	}
	// Input files
	for _, inputFile := range flag_inputFiles {
		file, err := os.Open(inputFile)
		if err != nil {
			customWarning(logger, "Error opening the input file %s: %v", inputFile, err)
			//continue
		} else {
			scanners = append(scanners, bufio.NewScanner(file))
		}
		defer file.Close()
	}

	urlsTmpFilename := ".urls.tmp.idaru." + fmt.Sprintf("%d", time.Now().Unix())
	// Open (or create if it doesn't exist) the file in write mode
	file, err := os.Create(urlsTmpFilename)
	if err != nil {
		customError(logger, "Error creating the file:", err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)

	// Use the scanners
	*flag_filterParam = (*flag_filterParam || *flag_fP)
	for _, scanner := range scanners {
		for scanner.Scan() {
			urlStr := scanner.Text()
			if url.ValidateURL(urlStr, *flag_filterParam) {
				customInfo(logger, "Valid URL: %s", urlStr)
				_, err := writer.WriteString(urlStr + "\n")
				if err != nil {
					customWarning(logger, "Error writing to the file:", err)
				}
			} else {
				customWarning(logger, "Invalid URL: %s", urlStr)
			}
		}
		if err := scanner.Err(); err != nil {
			customWarning(logger, "Error reading from the input: %v", err)
			continue
		}
	}
	writer.Flush()

	// Save the URLs to the sitemap
	sitemap := url.Init()
	urls, _ := url.GetFromFile(urlsTmpFilename)
	sitemap.Add(urls)

	// Remove the temporary file
	err = os.Remove(urlsTmpFilename)
	if err != nil {
		customWarning(logger, "Error removing the file:", err)
	}

	*flag_merge = (*flag_merge || *flag_m)
	if *flag_merge {
		sitemap.MergeKeysParam()
	}

	for _, element := range flag_set {
		parts := strings.Split(element, "=")
		if len(parts) == 2 {
			key, value := parts[0], parts[1]
			sitemap.SetValueParam(key, value)
		}
	}

	for _, element := range flag_add {
		parts := strings.Split(element, "=")
		if len(parts) == 2 {
			key, value := parts[0], parts[1]
			sitemap.AddValueParam(key, value)
		}
	}

	*flag_outputJson = (*flag_outputJson || *flag_oJ)
	if *flag_outputJson {
		sitemap.SaveToJson("sitemap.json")
	}

	*flag_show = (*flag_show || *flag_sh)
	if *flag_show {
		sitemap.ShowTree()
	} else {
		sitemap.Show()
	}

}

// Function to check if stdin is available
func isStdinAvailable() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

// Custom function for informational messages
func customInfo(logger *log.Logger, format string, v ...interface{}) {
	message := fmt.Sprintf("[+] "+format, v...)
	logger.Print(message)
}

// Custom function for warning messages
func customWarning(logger *log.Logger, format string, v ...interface{}) {
	message := fmt.Sprintf("[!] "+format, v...)
	logger.Print(message)
}

// Custom function for error messages
func customError(logger *log.Logger, format string, v ...interface{}) {
	message := fmt.Sprintf("[x] "+format, v...)
	logger.Fatal(message)
}

// Define a value type that accepts multiple values for an option
type ArrayFlagString []string

func (s *ArrayFlagString) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *ArrayFlagString) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// Define a Dev Null using the io package
type NullWriter struct{}

func (nw NullWriter) Write(p []byte) (int, error) {
	// Simply discards the data without doing anything
	return len(p), nil
}

// Version displays the ASCII art and terminates the program
func Version() {
	asciiArt := `
.___    .___                   
|   | __| _/____ _______ __ __ 
|   |/ __ |\__  \\_  __ \  |  \
|   / /_/ | / __ \|  | \/  |  /
|___\____ |(____  /__|  |____/ 
         \/     \/             
                      by ` + _AUTHOR_ + `

[+] Version is ` + _VERSION_
	fmt.Println(asciiArt)
	fmt.Println()
}
