package history

import "github.com/parakeet-nest/parakeet/llm"

type Messages interface {
	Get(id string) (llm.MessageRecord, error)
	GetMessage(id string) (llm.Message, error)
	GetAll() ([]llm.MessageRecord, error)

	GetLastNMessages(n int) ([]llm.Message, error)

	GetAllMessages(patterns ... string) ([]llm.Message, error)
	
	GetAllMessagesOfSession(sessionId string, patterns ... string) ([]llm.Message, error)


	Save(messageRecord llm.MessageRecord) (llm.MessageRecord, error)
	SaveMessage(id string, message llm.Message) (llm.MessageRecord, error)

	SaveMessageWithSession(sessionId, messageId string, message llm.Message) (llm.MessageRecord, error)


	RemoveMessage(id string) error
	RemoveAllMessages() error
	RemoveTopMessage() error
	
	RemoveAllMessagesOfSession(sessionId string) error
	RemoveTopMessageOfSession(sessionId string) error
	
	KeepLastN(n int) error
	KeepLastNOfSession(sessionId string, n int) error
}



