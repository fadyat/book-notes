package main

import "fmt"

type Customer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (c *Customer) String() string {
	return fmt.Sprintf("%s (%d)", c.Name, c.ID)
}
