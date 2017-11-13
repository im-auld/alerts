package main

import (
	"fmt"

	"github.com/im-auld/alerts/client"
)

func main() {
	c := client.NewAlertClient()
	fmt.Println(c.GetAlertsForUser(1))
	fmt.Println(c.GetAlertsForUser(2))
}
