package handwritingio

import (
	"fmt"
	"net/url"
	"os"
)

func ExampleClient_List() {

	u, err := url.Parse(os.Getenv("HANDWRITINGIO_API_URL"))
	if err != nil {
		fmt.Println("Make sure you have your environment variable (HANDWRITINGIO_API_URL) set correctly")
		fmt.Println(err)
		return
	}

	c := NewClient(u)

	offset := 0
	limit := 5
	order_by := "title"
	order_direction := "asc"
	handwritings, err := c.ListHandwritings(offset, limit, order_by, order_direction)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("ListHandwritings returned %d Handwritings", len(handwritings))
	fmt.Printf("%#v", handwritings)

	// Output:
	// ListHandwritings returned 5 Handwritings
}
