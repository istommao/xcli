package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

type Alias struct {
    Name    string `json:"name"`
    Command string `json:"command"`
}

var aliases []Alias


func getAliasFilePath() (string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    aliasDir := filepath.Join(homeDir, ".xcliup")

    if err := os.MkdirAll(aliasDir, os.ModePerm); err != nil {
        return "", err
    }
    return filepath.Join(aliasDir, "aliases.json"), nil
}

// Load aliases from file
func loadAliases() error {
    aliasFile, err := getAliasFilePath()
    if err != nil {
        return err
    }

    file, err := ioutil.ReadFile(aliasFile)
    if err != nil {
        if os.IsNotExist(err) {
            return nil
        }
        return err
    }
    return json.Unmarshal(file, &aliases)
}

// Save aliases to file
func saveAliases() error {
    aliasFile, err := getAliasFilePath()
    if err != nil {
        return err
    }

    data, err := json.MarshalIndent(aliases, "", "  ")
    if err != nil {
        return err
    }
    return ioutil.WriteFile(aliasFile, data, 0644)
}

// Show all aliases
func showAliases() {
    if len(aliases) == 0 {
        fmt.Println("No aliases found.")
        return
    }
    for _, alias := range aliases {
        fmt.Printf("%s: %s\n", alias.Name, alias.Command)
    }
}

// Create a new alias
func createAlias(name string, command string) {
    for _, alias := range aliases {
        if alias.Name == name {
            fmt.Printf("Alias '%s' already exists.\n", name)
            return
        }
    }
    aliases = append(aliases, Alias{Name: name, Command: command})
    saveAliases()
    fmt.Printf("Created alias '%s' for command '%s'.\n", name, command)
}

// Show a specific alias
func showAlias(name string) {
    for _, alias := range aliases {
        if alias.Name == name {
            fmt.Printf("%s: %s\n", alias.Name, alias.Command)
            return
        }
    }
    fmt.Printf("Alias '%s' not found.\n", name)
}

// Rename an alias
func renameAlias(oldName string, newName string) {
    for i, alias := range aliases {
        if alias.Name == oldName {
            aliases[i].Name = newName
            saveAliases()
            fmt.Printf("Renamed alias '%s' to '%s'.\n", oldName, newName)
            return
        }
    }
    fmt.Printf("Alias '%s' not found.\n", oldName)
}

// Update an alias
func updateAlias(name string, newCommand string) {
    for i, alias := range aliases {
        if alias.Name == name {
            aliases[i].Command = newCommand
            saveAliases()
            fmt.Printf("Updated alias '%s' to command '%s'.\n", name, newCommand)
            return
        }
    }
    fmt.Printf("Alias '%s' not found.\n", name)
}

// Run a command or script associated with an alias
func runAlias(name string, args []string) {
    for _, alias := range aliases {
        if alias.Name == name {
            command := alias.Command
            if len(args) > 0 {
                command = strings.Replace(command, "$1", args[0], 1)
            }

            cmd := exec.Command("sh", "-c", command)
            cmd.Stdout = os.Stdout
            cmd.Stderr = os.Stderr
            if err := cmd.Run(); err != nil {
                fmt.Printf("Error running command '%s': %s\n", alias.Command, err)
            }
            return
        }
    }
    fmt.Printf("Alias '%s' not found.\n", name)
}

func main() {
    if err := loadAliases(); err != nil {
        fmt.Printf("Error loading aliases: %s\n", err)
        return
    }

    if len(os.Args) < 2 {
        fmt.Println("Usage: xcli [command] [args]")
        return
    }

    command := os.Args[1]

    switch command {
    case "show":
        if len(os.Args) == 3 && os.Args[2] != "list" {
            showAlias(os.Args[2])
        } else {
            showAliases()
        }
    case "new":
        if len(os.Args) < 4 {
            fmt.Println("Usage: xcli new [alias name] [Your Command or Script Path]")
            return
        }
        createAlias(os.Args[2], os.Args[3])
    case "rename":
        if len(os.Args) != 4 {
            fmt.Println("Usage: xcli rename [alias old name] [new name]")
            return
        }
        renameAlias(os.Args[2], os.Args[3])
    case "update":
        if len(os.Args) < 4 {
            fmt.Println("Usage: xcli update [alias name] [New Command or Script Path]")
            return
        }
        updateAlias(os.Args[2], os.Args[3])
    case "run":
        if len(os.Args) < 3 {
            fmt.Println("Usage: xcli run [alias name] [args...]")
            return
        }
        runAlias(os.Args[2], os.Args[3:])
    default:
        fmt.Println("Unknown command:", command)
    }
}
