package openai

import (
	"strings"
	"time"

	"github.com/golly-go/golly"
)

type AIContext interface {
	String() string
	Message() Message
	MessageRole() MessageRole
}

type DefaultAIContext struct {
	Role MessageRole `json:"role"`

	Who     string `json:"who"`
	Action  string `json:"action"`
	Content string `json:"content"`

	CreatedAt time.Time `json:"createdAt"`
}

func NewDefaultContext(who, action, content string) DefaultAIContext {
	return NewDefaultContextWithRole(RoleSystem, who, action, content)
}

func NewContextForContent(content string) DefaultAIContext {
	return NewDefaultContext("", "", content)
}

func NewDefaultContentRole(role MessageRole, content string) DefaultAIContext {
	return DefaultAIContext{Role: role, Who: "", Action: "", Content: content, CreatedAt: time.Now()}
}

func NewDefaultContextWithRole(role MessageRole, who, action, content string) DefaultAIContext {
	return DefaultAIContext{Role: role, Who: who, Action: action, Content: content, CreatedAt: time.Now()}
}

func NewDefaultSystemContext(content string) DefaultAIContext {
	return DefaultAIContext{Role: RoleSystem, Content: content, CreatedAt: time.Now()}
}

func (c DefaultAIContext) MessageRole() MessageRole {
	return c.Role
}

func (c DefaultAIContext) String() string {
	var pieces []string

	if c.Content == "" {
		return ""
	}

	if c.Who != "" {
		pieces = append(pieces, c.Who)
	}

	if c.Action != "" {
		pieces = append(pieces, c.Action)
	}

	pieces = append(pieces, c.Content)

	return strings.Join(pieces, " ")
}

func (c DefaultAIContext) Message() Message {
	if c.Role == "" {
		c.Role = RoleSystem
	}

	return Message{Role: c.Role, Content: c.String()}
}

type AIContexts []AIContext

func (c AIContexts) Messages() Messages {
	return golly.Map(c, func(c AIContext) Message { return c.Message() })
}

func (c *AIContexts) Append(aiContexts ...AIContext) *AIContexts {
	*c = append(*c, aiContexts...)
	return c
}
