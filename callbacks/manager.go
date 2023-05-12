package callbacks

const (
	EventLLMStart    = "llm_start"
	EventLLMNewToken = "llm_new_token"
	EventLLMEnd      = "llm_end"
	EventLLMError    = "llm_error"
	EventChainStart  = "chain_start"
	EventChainEnd    = "chain_end"
	EventChainError  = "chain_error"
	EventToolStart   = "tool_start"
	EventToolEnd     = "tool_end"
	EventToolError   = "tool_error"
	EventAgentAction = "agent_action"
	EventAgentFinish = "agent_finish"
)

type CallbackManager []Callback

func NewCallbackManager(cbs ...Callback) *CallbackManager {
	return (*CallbackManager)(&cbs)
}

func (c *CallbackManager) On(event string, args ...any) {
	if c == nil {
		return
	}
	for _, cb := range *c {
		switch event {
		case EventLLMStart:
			if cb.OnLLMStart != nil {
				cb.OnLLMStart(args[0].(string))
			}
		case EventLLMNewToken:
			if cb.OnLLMNewToken != nil {
				cb.OnLLMNewToken(args[0].(string))
			}
		case EventLLMEnd:
			if cb.OnLLMEnd != nil {
				cb.OnLLMEnd(args[0].(string))
			}
		case EventLLMError:
			if cb.OnLLMError != nil {
				cb.OnLLMError(args[0].(error))
			}
		case EventChainStart:
			if cb.OnChainStart != nil {
				cb.OnChainStart(args[0].(map[string]string))
			}
		case EventChainEnd:
			if cb.OnChainEnd != nil {
				cb.OnChainEnd(args[0].(string))
			}
		case EventChainError:
			if cb.OnChainError != nil {
				cb.OnChainError(args[0].(error))
			}
		case EventToolStart:
			if cb.OnToolStart != nil {
				cb.OnToolStart(args[0].(string))
			}
		case EventToolEnd:
			if cb.OnToolEnd != nil {
				cb.OnToolEnd(args[0].(string))
			}
		case EventToolError:
			if cb.OnToolError != nil {
				cb.OnToolError(args[0].(error))
			}
		case EventAgentAction:
			if cb.OnAgentAction != nil {
				cb.OnAgentAction(args[0].(string))
			}
		case EventAgentFinish:
			if cb.OnAgentFinish != nil {
				cb.OnAgentFinish(args[0].(string))
			}
		}
	}
}

func (c *CallbackManager) OnLLMStart(prompt string) {
	c.On(EventLLMStart, prompt)
}

func (c *CallbackManager) OnLLMNewToken(token string) {
	c.On(EventLLMNewToken, token)
}

func (c *CallbackManager) OnLLMEnd(resp string) {
	c.On(EventLLMEnd, resp)
}

func (c *CallbackManager) OnLLMError(err error) {
	c.On(EventLLMError, err)
}

func (c *CallbackManager) OnChainStart(inputs map[string]string) {
	c.On(EventChainStart, inputs)
}

func (c *CallbackManager) OnChainEnd(resp string) {
	c.On(EventChainEnd, resp)
}

func (c *CallbackManager) OnChainError(err error) {
	c.On(EventChainError, err)
}

func (c *CallbackManager) OnToolStart(input string) {
	c.On(EventToolStart, input)
}

func (c *CallbackManager) OnToolEnd(resp string) {
	c.On(EventToolEnd, resp)
}

func (c *CallbackManager) OnToolError(err error) {
	c.On(EventToolError, err)
}

func (c *CallbackManager) OnAgentAction(action string) {
	c.On(EventAgentAction, action)
}

func (c *CallbackManager) OnAgentFinish(resp string) {
	c.On(EventAgentFinish, resp)
}
