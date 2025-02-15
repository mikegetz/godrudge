DISCLAIMER:

This project is an independent, unofficial implementation and is not affiliated with or endorsed by Drudge. It is provided for educational and experimental purposes only.

Please note that this implementation relies on parsing the website’s HTML DOM, which, unlike a formal API, does not adhere to a fixed contract. As a result, any modifications to the website’s structure may cause the implementation to break. Users should be aware of these limitations when using this package.

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


