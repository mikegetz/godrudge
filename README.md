# Table of Contents

- [DISCLAIMER](#disclaimer)
- [Example Usage](#example-usage)
  - [How to access the parsed headlines from the client](#how-to-access-the-parsed-headlines-from-the-client)
  - [How to print Drudge to stdout with ANSI links](#prints-drudge-to-stdout-with-ansi-links)

# DISCLAIMER:

This project is an independent, unofficial implementation and is not affiliated with or endorsed by Drudge. It is provided for educational and experimental purposes only.

Please note that this implementation relies on parsing the website’s HTML DOM, which, unlike a formal API, does not adhere to a fixed contract. As a result, any modifications to the website’s structure may cause the implementation to break. Users should be aware of these limitations when using this package.

![image](https://github.com/user-attachments/assets/58a0f545-3f1a-480d-8106-ebf3425b502d)

# Example Usage:

## How to access the parsed headlines from the client

```
package main

import (
    "fmt"

    "github.com/MGuitar24/godrudge"
)

func main() {
    client := godrudge.NewClient()
    err := client.ParseRSS()
    if err != nil {
        fmt.Println("Error parsing", err)
    }

    client.PrintDrudge(true)

    //Access the first headline title from Column 1
    fmt.Println(client.Page.HeadlineColumns[0][0].Title)
    fmt.Println(client.Page.HeadlineColumns[0][0].Href)

    //Access the first headline title from Column 2
    fmt.Println(client.Page.HeadlineColumns[1][0].Title)
    fmt.Println(client.Page.HeadlineColumns[1][0].Href)

    //Access the first headline title from Column 3
    fmt.Println(client.Page.HeadlineColumns[2][0].Title)
    fmt.Println(client.Page.HeadlineColumns[2][0].Href)
}
```

## Prints Drudge to stdout with ANSI links

```
package main

import (
    "fmt"

    "github.com/MGuitar24/godrudge"
)

func main() {
    client := godrudge.NewClient()
    err := client.ParseRSS()
    if err != nil {
        fmt.Println("Error parsing", err)
    }

    // set to true to print without ANSI links
    client.PrintDrudge(false)
}
```