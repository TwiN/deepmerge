# deepmerge
![test](https://github.com/TwiN/deepmerge/workflows/test/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/TwiN/deepmerge)](https://goreportcard.com/report/github.com/TwiN/deepmerge)
[![Go version](https://img.shields.io/github/go-mod/go-version/TwiN/deepmerge.svg)](https://github.com/TwiN/deepmerge)
[![Go Reference](https://pkg.go.dev/badge/github.com/TwiN/deepmerge.svg)](https://pkg.go.dev/github.com/TwiN/deepmerge)

Go library for deep merging YAML or JSON files.


## Usage

### YAML
```go
package main

import (
	"github.com/TwiN/deepmerge"
)

func main() {
	dst := `
debug: true
client:
  insecure: true
users:
  - id: 1
    firstName: John
    lastName: Doe
  - id: 2
    firstName: Jane
    lastName: Doe`
	src := `
client:
  timeout: 5s
users:
  - id: 3
    firstName: Bob
    lastName: Smith`
	output, err := deepmerge.YAML([]byte(dst), []byte(src))
	if err != nil {
		panic(err)
	}
	println(string(output))
}
```

Output:
```yaml
client:
    insecure: true
    timeout: 5s
debug: true
users:
    - firstName: John
      id: 1
      lastName: Doe
    - firstName: Jane
      id: 2
      lastName: Doe
    - firstName: Bob
      id: 3
      lastName: Smith
```

### JSON
```go
package main

import (
	"github.com/TwiN/deepmerge"
)

func main() {
	dst := `{
  "debug": true,
  "client": {
    "insecure": true
  },
  "users": [
    {
      "id": 1,
      "firstName": "John",
      "lastName": "Doe"
    },
    {
      "id": 2,
      "firstName": "Jane",
      "lastName": "Doe"
    }
  ]
}`
	src := `{
  "client": {
    "timeout": "5s"
  },
  "users": [
    {
      "id": 3,
      "firstName": "Bob",
      "lastName": "Smith"
    }
  ]
}`
	output, err := deepmerge.JSON([]byte(dst), []byte(src))
	if err != nil {
		panic(err)
	}
	println(string(output))
}
```

Output:
```json
{
  "client": {
    "insecure": true,
    "timeout": "5s"
  },
  "debug": true,
  "users": [
    {
      "firstName": "John",
      "id": 1,
      "lastName": "Doe"
    },
    {
      "firstName": "Jane",
      "id": 2,
      "lastName": "Doe"
    },
    {
      "firstName": "Bob",
      "id": 3,
      "lastName": "Smith"
    }
  ]
}
```