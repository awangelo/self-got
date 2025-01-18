package main

import (
	"bufio"
	"os"
	"os/user"
	"path/filepath"
)

type config struct {
	Token  string
	Prefix string
}

func (c *config) getConfigPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(usr.HomeDir, ".config", "self-got")
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.txt"), nil
}

func (c *config) loadConfig() error {
	configPath, err := c.getConfigPath()
	if err != nil {
		return err
	}

	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		c.Token = scanner.Text()
	}
	if scanner.Scan() {
		c.Prefix = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (c *config) createConfig() error {
	configPath, err := c.getConfigPath()
	if err != nil {
		return err
	}

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(c.Token + "\n" + c.Prefix)
	if err != nil {
		return err
	}

	return nil
}

func (c *config) isValid() bool {
	return c.Token != "" && c.Prefix != ""
}
