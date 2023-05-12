# 🦜️🔗 Go LangChain

⚡ Building applications with LLMs through composability ⚡

> **Notice: The library is in a very early stage of development.**

## 📦 Installation

```bash
go get github.com/lab99x/go-langchain
```

## 🤔 What is this?

Large language models (LLMs) are emerging as a transformative technology, enabling
developers to build applications that they previously could not.
But using these LLMs in isolation is often not enough to
create a truly powerful app - the real power comes when you can combine them with other sources of computation or knowledge.

This library is aimed at assisting in the development of those types of applications.

## 📖 Full Documentation

> TODO: Add link to full documentation

## 🚀 Quick Start

To call a LLM, such as OpenAI GPT API, you can use the following code:

```go
package main

import (
    "context"
    "fmt"

    "github.com/lab99x/go-langchain/llms"
)

func main() {
    llm := &llms.OpenAI{} // OpenAI API key can be set in environment variable OPENAI_API_KEY
    resp, err := llm.Chat(context.Background(), "When the forbidden city was built?")
    if err != nil {
        panic(err)
    }
    fmt.Println(resp)
    // Output: The forbidden city in Beijing was built in 1406.
}
```

The Chains and Prompts are also implemented in this library. For example, you can use the following code to create a chain that calls OpenAI GPT API:

```go
package main

import (
    "context"
    "fmt"

    "github.com/lab99x/go-langchain/llms"
    "github.com/lab99x/go-langchain/prompts"
    "github.com/lab99x/go-langchain/chains"
)

func main() {
    llm := &llms.OpenAI{} // OpenAI API key can be set in environment variable OPENAI_API_KEY
    prompt_template := prompts.NewPromptTemplateByTemplate("Where is the capital city of {country}?")
    c := chains.NewLLMChain(llm, prompt_template)

    resp, err := c.RunText(context.Background(), "China")
    if err != nil {
        panic(err)
    }
    fmt.Println(resp)
    // Output: The capital city of China is Beijing.

    resp, err = c.Run(context.Background(), map[string]string{"country": "France"})
    if err != nil {
        panic(err)
    }
    fmt.Println(resp)
    // Output: The capital city of France is Paris.
}
```

## Relationship with Python LangChain

This library is the Go language implementation of [LangChain](https://github.com/hwchase17/langchain). 

The [LangChainHub](https://github.com/hwchase17/langchain-hub) is a central place for the serialized versions of these prompts, chains, and agents.

## 💁 Contributing

As an open source project in a rapidly developing field, we are extremely open to contributions, whether it be in the form of a new feature, improved infra, or better documentation.

> TODO: Add contributing guidelines
