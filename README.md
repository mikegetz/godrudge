DISCLAIMER:

This project is just a fun project and is not an official client to Drudge. The implementation is fragile since it relies on parsing the HTML DOM of the website which has no strict contract unlike an API. A small change to the website can easily break this.

Example Usage:

```
package main

import (
	"github.com/MGuitar24/go-drudge"
)

func main() {
	client := drudge.NewClient()
	client.Parse()
	client.PrintHeadlines()
}

```

![image](https://github.com/user-attachments/assets/f1fd1cab-d4ba-449c-b4aa-fb9526645d0b)
