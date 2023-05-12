package prompts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessagePromptTemplate(t *testing.T) {
	// Human
	var template Template = NewMessagePromptTemplate("Hello {name}!", MessageTypeHuman)
	vars := template.GetVariables()
	assert.Equal(t, 1, len(vars))
	assert.Equal(t, "name", vars[0])
	values := map[string]string{
		"name": "World",
	}
	assert.Equal(t, "human: Hello World!", template.Format(values))
	assert.Equal(t, Message{
		Type:    MessageTypeHuman,
		Content: "Hello World!",
	}, template.(*MessagePromptTemplate).FormatMessage(values))

	// AI
	template = NewMessagePromptTemplate("Hello {name}!", MessageTypeAI)
	vars = template.GetVariables()
	assert.Equal(t, 1, len(vars))
	assert.Equal(t, "name", vars[0])
	values = map[string]string{
		"name": "World",
	}
	assert.Equal(t, "ai: Hello World!", template.Format(values))
	assert.Equal(t, Message{
		Type:    MessageTypeAI,
		Content: "Hello World!",
	}, template.(*MessagePromptTemplate).FormatMessage(values))

	// System
	template = NewMessagePromptTemplate("Hello {name}!", MessageTypeSystem)
	vars = template.GetVariables()
	assert.Equal(t, 1, len(vars))
	assert.Equal(t, "name", vars[0])
	values = map[string]string{
		"name": "World",
	}
	assert.Equal(t, "system: Hello World!", template.Format(values))
	assert.Equal(t, Message{
		Type:    MessageTypeSystem,
		Content: "Hello World!",
	}, template.(*MessagePromptTemplate).FormatMessage(values))
}

func TestChatMessagePromptTemplate(t *testing.T) {
	// Chat
	templates := []MessagePromptTemplate{
		*NewMessagePromptTemplate("You're an AI bot, name: {name}!", MessageTypeSystem),
		*NewMessagePromptTemplate("Hello, my name is {name}!", MessageTypeAI),
		*NewMessagePromptTemplate("Hello, {name}, {question}!", MessageTypeHuman),
	}
	var template Template = NewChatMessagePromptTemplate(templates)
	vars := template.GetVariables()
	assert.Equal(t, 2, len(vars))
	assert.Equal(t, "name", vars[0])
	assert.Equal(t, "question", vars[1])
	values := map[string]string{
		"name":     "Stone",
		"question": "how are you",
	}
	assert.Equal(t, `system: You're an AI bot, name: Stone!
ai: Hello, my name is Stone!
human: Hello, Stone, how are you!`, template.Format(values))
	assert.Equal(t, []Message{
		{
			Type:    MessageTypeSystem,
			Content: `You're an AI bot, name: Stone!`,
		},
		{
			Type:    MessageTypeAI,
			Content: `Hello, my name is Stone!`,
		},
		{
			Type:    MessageTypeHuman,
			Content: `Hello, Stone, how are you!`,
		},
	}, template.(*ChatMessagePromptTemplate).FormatMessages(values))
}
