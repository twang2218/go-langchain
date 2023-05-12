package callbacks

type Callback struct {
	OnLLMStart    func(prompts string)
	OnLLMNewToken func(token string)
	OnLLMEnd      func(resp string)
	OnLLMError    func(err error)
	OnChainStart  func(inputs map[string]string)
	OnChainEnd    func(resp string)
	OnChainError  func(err error)
	OnToolStart   func(input string)
	OnToolEnd     func(resp string)
	OnToolError   func(err error)
	OnAgentAction func(action string)
	OnAgentFinish func(resp string)
}
