package handwritingio

import (
	"fmt"
	"net/url"
	"os"
)

func ExampleClient_ListHandwritings() {

	u, err := url.Parse(os.Getenv("HANDWRITINGIO_API_URL"))
	if err != nil {
		fmt.Println("Make sure you have your environment variable HANDWRITINGIO_API_URL set correctly")
		fmt.Println(err)
		return
	}

	var params = DefaultHandwritingListParams
	params.Limit = 5
	c := NewClient(u)

	handwritings, err := c.ListHandwritings(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("ListHandwritings returned %d Handwritings", len(handwritings))

	// Output:
	// ListHandwritings returned 5 Handwritings
}

func ExampleClient_GetHandwriting() {
	id := "2D5S46A80003"

	u, err := url.Parse(os.Getenv("HANDWRITINGIO_API_URL"))
	if err != nil {
		fmt.Println("Make sure you have your environment variable HANDWRITINGIO_API_URL set correctly")
		fmt.Println(err)
		return
	}

	c := NewClient(u)

	handwriting, err := c.GetHandwriting(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(handwriting.Title)

	// Output:
	// Perry
}
