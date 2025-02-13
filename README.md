DISCLAIMER:

This project is just a fun project and is not an official client to Drudge. The implementation is fragile since it relies on parsing the HTML DOM of the website which has no strict contract unlike an API. A small change to the website can easily break this.

Example Usage:

```
package main

import (
	"fmt"

	"github.com/MGuitar24/godrudge"
)

func main() {
	client := godrudge.NewClient()
	err := client.Parse()
	if err != nil {
		fmt.Println("Error parsing", err)
	}
	client.PrintDrudge()
}


```

![image](https://github.com/user-attachments/assets/58a0f545-3f1a-480d-8106-ebf3425b502d)


