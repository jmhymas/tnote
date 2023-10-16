package main

import "fmt"
import "os"
import "flag"
import "time"

func main() {

    numOfNotesExpected, numOfFlags := getFlags()
    notes := os.Args[numOfFlags + 1:]
    numOfNotesActual := len(notes)
    
    //This check ensures all notes were wrapped in quotation marks (can't check directly since they get stripped before they are handed to GO)
    if numOfNotesActual != numOfNotesExpected{
        errorMessage := fmt.Sprintf("tnote expected %v note(s) but found %v.", numOfNotesExpected, numOfNotesActual)
        panic(errorMessage)
    }

    writeNotesToFile(notes)
}

func getFlags() (numOfNotesExpected int, numOfFlags int){
    flag.IntVar(&numOfNotesExpected, "n", 1, "this flag sets the number of notes being passed per call")
    flag.Parse()
    return numOfNotesExpected, flag.NFlag()
}

func createNewFile(filepath string, date string) (file *os.File){
    file, err := os.Create(filepath)
    if err != nil {
        panic("Cannot create file. Directory path not found")
    }
    header := "# " + date + "\n"
    file.WriteString(header)
    return file
}

func getNotesLocation() (notes_loc string){
    notes_loc, exists := os.LookupEnv("NOTES_LOC")
    if exists == false{
        fmt.Println("NOTES_LOC environment variable not found. Using default notes folder in the home directory")
        notes_loc = os.Getenv("HOME") + "/notes"
    }

    _ , err := os.Stat(notes_loc)
    if err != nil {
        panic("No notes directory found.")
    }

    return notes_loc
}

func writeNotesToFile(notes []string) {
    notes_loc := getNotesLocation()
    date := string(time.Now().Format("01-02-06"))
    filename := date + ".md"
    filepath := notes_loc + string(os.PathSeparator) + filename
    timestamp := time.Now().Format("[2006-01-02T15:04:05]")

    file, err := os.OpenFile(filepath,os.O_APPEND|os.O_WRONLY, 0644 )
    if err != nil {
        fmt.Printf("Creating new file %v\n", filename )
        file = createNewFile(filepath, date)
    }

    defer file.Close()

    for _, note := range notes{
        stringToAppend := note
        stringToAppend = "\n### " + string(timestamp) + " " + stringToAppend + "\n"

        _, err = file.WriteString(stringToAppend)
        if err != nil {
        panic("couldn't write note to file")
        }
    }
}