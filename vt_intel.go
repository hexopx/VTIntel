package main

import (
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/fatih/color"
    "github.com/spf13/viper"
)

// VTResponse (partial, for Intelligence search)
type VTResponse struct {
    Data []struct {
        ID         string `json:"id"`
        Attributes struct {
            MeaningfulName string `json:"meaningful_name"`
        } `json:"attributes"`
    } `json:"data"`
}

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
    color.Green(`


██╗   ██╗████████╗████████╗
██║   ██║╚══██╔══╝╚══██╔══╝
██║   ██║   ██║      ██║       
╚██╗ ██╔╝   ██║      ██║   
 ╚████╔╝    ██║      ██║   
  ╚═══╝     ╚═╝      ╚═╝   

→ VirusTotal Intelligence Tool (v1.0)

→ https://github.com/hexopx
 
    `)

    // Get Search Query
    var query string
    fmt.Print("Enter your search query: ") 
    fmt.Scanln(&query)
    intelligenceSearch(apiKey, query)
}

func intelligenceSearch(apiKey, query string) {
    endpoint := fmt.Sprintf("https://www.virustotal.com/api/v3/intelligence/search?query=%s&limit=25", query)

    req, err := http.NewRequest("GET", endpoint, nil)
    if err != nil {
        color.Red("Error creating request: %v", err)
        return
    }
    req.Header.Set("x-apikey", apiKey)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        color.Red("Error sending request: %v", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        var data VTResponse
        err = json.NewDecoder(resp.Body).Decode(&data)
        if err != nil {
            color.Red("Error decoding response: %v", err)
            return
        }

        if len(data.Data) > 0 {
            for _, item := range data.Data {
                color.Green("%s | %s", item.ID, item.Attributes.MeaningfulName)
            }
        } else {
            color.Yellow("No results found.")
        }
    } else {
        color.Red("Error: %s", resp.Status)
    }
}
