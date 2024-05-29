package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "path"
    "strings"

    "github.com/fatih/color"
    "github.com/spf13/viper"
)

func main() {
    // Load Configuration
    viper.SetConfigName("vtconfig") 
    viper.SetConfigType("yaml")    
    viper.AddConfigPath(".")       
    if err := viper.ReadInConfig(); err != nil {
        color.Red("Error reading configuration file: %v", err)
        return
    }

    apiKey := viper.GetString("api_key")

    // ASCII Art
    color.Yellow(`

 _  _  ____    ____   __   _  _  __ _  __     __    __   ____  ____  ____ 
/ )( \(_  _)  (    \ /  \ / )( \(  ( \(  )   /  \  / _\ (    \(  __)(  _ \
\ \/ /  )(     ) D ((  O )\ /\ //    // (_/\(  O )/    \ ) D ( ) _)  )   /
 \__/  (__)   (____/ \__/ (_/\_)\_)__)\____/ \__/ \_/\_/(____/(____)(__\_)

 Version 1.0 | https://github.com/hexopx
 
    `)

    // Get Sample SHA256
    var sampleSHA256 string
    fmt.Print("Enter Sample SHA256: ")
    fmt.Scanln(&sampleSHA256)

    // VirusTotal API Endpoint
    endpoint := fmt.Sprintf("https://www.virustotal.com/api/v3/files/%s/download", sampleSHA256)

    // Create HTTP Request
    req, err := http.NewRequest("GET", endpoint, nil)
    if err != nil {
        color.Red("Error creating request: %v", err)
        return
    }
    req.Header.Set("x-apikey", apiKey)

    // Send Request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        color.Red("Error sending request: %v", err)
        return
    }
    defer resp.Body.Close()

    // Handle Response
    if resp.StatusCode == http.StatusOK {
        contentDisposition := resp.Header.Get("Content-Disposition")
        filename := sampleSHA256
        if contentDisposition != "" {
            if strings.Contains(contentDisposition, "filename=") {
                filename = contentDisposition[strings.Index(contentDisposition, "filename=")+len("filename="):]
                filename = strings.Trim(filename, "\"")
            }
        }

        outFilePath := path.Join(".", filename)
        outFile, err := os.Create(outFilePath)
        if err != nil {
            color.Red("Error creating file: %v", err)
            return
        }
        defer outFile.Close()

        _, err = io.Copy(outFile, resp.Body)
        if err != nil {
            color.Red("Error writing file: %v", err)
            return
        }
        color.Green("Sample downloaded successfully!")
    } else {
        color.Red("Error: %s", resp.Status)
    }
}
