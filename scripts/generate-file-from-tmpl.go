package main
import (
  "os"
  "text/template"
)

func checkError(err error) {
  if err != nil { panic(err) }
}

func main() {
  tmplFile := os.Getenv("TMPL_FILE_PATH")
  tmpl, err := template.ParseFiles(tmplFile)
  checkError(err)
  err = tmpl.Execute(os.Stdout, nil)
  checkError(err)
}
