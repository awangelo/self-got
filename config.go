package main

import (
	"bufio"
	"os"
)

type config struct {
	Token  string
	Prefix string
}

func (c *config) loadConfig() error {
	file, err := os.Open("config.txt")
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
	file, err := os.Create("config.txt")
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
