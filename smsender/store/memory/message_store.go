package memory

import (
	"errors"
	"sync"
	"time"

	"github.com/minchao/smsender/smsender/model"
	"github.com/minchao/smsender/smsender/store"
)

type MessageStore struct {
	*Store
	messages []*model.Message
	sync.RWMutex
}

func NewMemoryMessageStore(store *Store) store.MessageStore {
	return &MessageStore{store, []*model.Message{}, sync.RWMutex{}}
}

func (s *MessageStore) Get(id string) store.Channel {
	storeChannel := make(store.Channel, 1)

	go func() {
		result := store.Result{}

		s.Lock()
		defer s.Unlock()

		if _, message, err := s.find(id); err != nil {
			result.Err = err
		} else {
			result.Data = message
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s *MessageStore) GetByIds(ids []string) store.Channel {
	storeChannel := make(store.Channel, 1)

	go func() {
		result := store.Result{}

		s.Lock()
		defer s.Unlock()

		var messages []*model.Message

		for _, id := range ids {
			if _, message, err := s.find(id); err == nil {
				messages = append(messages, message)
			}
		}

		result.Data = messages

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s *MessageStore) GetByProviderAndMessageId(provider, providerMessageId string) store.Channel {
	storeChannel := make(store.Channel, 1)

	go func() {
		result := store.Result{}

		s.Lock()
		defer s.Unlock()

		var message *model.Message
		for i, m := range s.messages {
			if provider == *m.Provider && providerMessageId == *m.ProviderMessageId {
				message = s.messages[i]
				break
			}
		}
		if message != nil {
			result.Data = message
		} else {
			result.Err = errors.New("message not found")
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s *MessageStore) Search(params map[string]interface{}) store.Channel {
	storeChannel := make(store.Channel, 1)

	go func() {
		result := store.Result{}

		s.Lock()
		defer s.Unlock()

		messages := []*model.Message{}
		length := len(s.messages)
		since, hasSince := params["since"] // ASC
		until, hasUntil := params["until"] // DESC, default

		for i := 0; i < length; i++ {
			var message *model.Message
			if hasSince {
				message = s.messages[i]

				if !message.CreatedTime.After(since.(time.Time)) {
					continue
				}
			} else {
				message = s.messages[length-i-1]

				if hasUntil && !message.CreatedTime.Before(until.(time.Time)) {
					continue
				}
			}

			if to, ok := params["to"]; ok {
				if message.To != to.(string) {
					continue
				}
			}
			if status, ok := params["status"]; ok {
				if message.Status.String() != status.(string) {
					continue
				}
			}

			messages = append(messages, message)

			if limit, ok := params["limit"]; ok {
				if len(messages) == limit.(int) {
					break
				}
			}
		}

		if hasSince {
			if num := len(messages); num > 1 {
				// Reverse the messages slice
				for i, j := 0, num-1; i < j; i, j = i+1, j-1 {
					messages[i], messages[j] = messages[j], messages[i]
				}
			}
		}

		result.Data = messages

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s *MessageStore) Save(message *model.Message) store.Channel {
	storeChannel := make(store.Channel, 1)

	go func() {
		result := store.Result{}

		s.Lock()
		defer s.Unlock()

		s.messages = append(s.messages, message)
		result.Data = &s.messages[len(s.messages)-1]

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s *MessageStore) Update(message *model.Message) store.Channel {
	storeChannel := make(store.Channel, 1)

	go func() {
		result := store.Result{}

		s.Lock()
		defer s.Unlock()

		if _, m, err := s.find(message.Id); err != nil {
			result.Err = err
		} else {
			*m = *message
			result.Data = m
		}

		storeChannel <- result
		close(storeChannel)
	}()

	return storeChannel
}

func (s *MessageStore) find(messageId string) (int64, *model.Message, error) {
	for index, message := range s.messages {
		if message.Id == messageId {
			return int64(index), message, nil
		}
	}

	return 0, nil, errors.New("message not found")
}
