package main

type Customer struct {
	ID   int
	Name string
}

func (c *Customer) toMap() map[string]interface{} {
	return map[string]interface{}{
		"id":   c.ID,
		"name": c.Name,
	}
}
