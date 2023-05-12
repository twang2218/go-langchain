package prompts

import (
	"fmt"
	"strings"
)

type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

const (
	MessageTypeHuman  = "human"
	MessageTypeAI     = "ai"
	MessageTypeSystem = "system"
	MessageTypeChat   = "chat"
)

func NewHumanMessage(content string) Message {
	return Message{
		Type:    MessageTypeHuman,
		Content: content,
	}
}

func NewAIMessage(content string) Message {
	return Message{
		Type:    MessageTypeAI,
		Content: content,
	}
}

func NewSystemMessage(content string) Message {
	return Message{
		Type:    MessageTypeSystem,
		Content: content,
	}
}

func NewChatMessage(content string) Message {
	return Message{
		Type:    MessageTypeChat,
		Content: content,
	}
}

type MessagePromptTemplate struct {
	PromptTemplate
	Type string `json:"type"`
}

func (t MessagePromptTemplate) FormatMessage(values map[string]string) Message {
	prompt := t.PromptTemplate.Format(values)
	return Message{
		Type:    t.Type,
		Content: prompt,
	}
}

func (t MessagePromptTemplate) Format(values map[string]string) string {
	return fmt.Sprintf("%s: %s", t.Type, t.PromptTemplate.Format(values))
}

func NewMessagePromptTemplate(template, t string) *MessagePromptTemplate {
	t = strings.TrimSpace(t)
	switch t {
	case MessageTypeHuman:
	case MessageTypeAI:
	case MessageTypeSystem:
		//	合法类型不改变
		break
	default:
		//	不合法类型默认为 human
		t = MessageTypeHuman
	}
	//  返回 MessagePromptTemplate
	return &MessagePromptTemplate{
		PromptTemplate: NewPromptTemplateByTemplate(template),
		Type:           t,
	}
}

func NewHumanMessagePromptTemplate(template string) *MessagePromptTemplate {
	return NewMessagePromptTemplate(template, MessageTypeHuman)
}

func NewAIMessagePromptTemplate(template string) *MessagePromptTemplate {
	return NewMessagePromptTemplate(template, MessageTypeAI)
}

func NewSystemMessagePromptTemplate(template string) *MessagePromptTemplate {
	return NewMessagePromptTemplate(template, MessageTypeSystem)
}

type ChatMessagePromptTemplate struct {
	Templates []MessagePromptTemplate
	Variables []string
	Type      string
}

func NewChatMessagePromptTemplate(templates []MessagePromptTemplate) *ChatMessagePromptTemplate {
	var variablesSet = map[string]bool{}
	var variables []string
	for _, template := range templates {
		for _, variable := range template.Variables {
			if _, ok := variablesSet[variable]; !ok {
				variablesSet[variable] = true
				variables = append(variables, variable)
			}
		}
	}

	return &ChatMessagePromptTemplate{
		Templates: templates,
		Variables: variables,
		Type:      MessageTypeChat,
	}
}

func (t ChatMessagePromptTemplate) FormatMessages(values map[string]string) []Message {
	var messages []Message
	for _, template := range t.Templates {
		messages = append(messages, template.FormatMessage(values))
	}
	return messages
}

func (t ChatMessagePromptTemplate) Format(values map[string]string) string {
	var messages []string
	for _, template := range t.Templates {
		messages = append(messages, template.Format(values))
	}
	return strings.Join(messages, "\n")
}

func (t ChatMessagePromptTemplate) GetVariables() []string {
	return t.Variables
}
