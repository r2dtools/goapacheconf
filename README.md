# GoApacheConf

**GoApacheConf** is a Go library for parsing, modifying, and regenerating **Apache configuration files**.  
It makes it easy to work with Apache config blocks and directives programmatically in your Go applications.

Powered by [Participle v2](https://github.com/alecthomas/participle) under the hood.

---


## 📦 Installation

```bash
go get github.com/r2dtools/goapacheconf
```

---

## 🚀 Usage Example

Parse the entire Apache configuration and inspect blocks and directives:

```go
package main

import (
	"fmt"

	"github.com/r2dtools/goapacheconf"
)

func main() {
	// Load Apache config from /etc/apache2
	config, err := goapacheconf.GetConfig("/etc/apache2", "")
	if err != nil {
		panic(err)
	}

	// Find VirtualHost blocks
	blocks := config.FindVirtualHostBlocks()
	if len(blocks) == 0 {
		panic("virtual host block not found")
	}
	vBlock := blocks[0]

	// Get DocumentRoot
	fmt.Println(vBlock.GetDocumentRoot())

	// Find "ErrorLog" directives
	directives := vBlock.FindDirectives("ErrorLog")
	if len(directives) == 0 {
		panic("directive not found")
	}

	// Print the first ErrorLog directive value
	fmt.Println(directives[0].GetFirstValue())
}
```

👉 For more examples, check the tests.

---

## 🔧 Roadmap

- Support for additional Apache directive types  
- Advanced helpers for config manipulation  
- Built-in config validation  

---

## 📜 License

MIT License. See [LICENSE](./LICENSE) for details.
