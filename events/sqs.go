package events

// инфраструктурный уровень
import (
	"encoding/json"
	"fmt"
	"log"

	//
	// какой-то импорт
	//
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// EventSQSHandler передаёт внутренние события во внешний мир
type EventSQSHandler struct {
	svc *sqs.SQS
}

// Notify передаёт события через SQS
func (e *EventSQSHandler) Notify(event Event) {
	data := map[string]string{
		"event": event.Name(),
	}

	body, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	_, err = e.svc.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(body)),
		QueueUrl:    &e.svc.Endpoint,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// инфраструктурный уровень
type SQSService struct {
	svc         *sqs.SQS
	publisher   *EventPublisher
	stopChannel chan bool
}

// Run запускает прослушивание SQS сообщений
func (s *SQSService) Run(event Event) {
	eventChan := make(chan Event)

MessageLoop:
	for {
		s.listen(eventChan)

		select {
		case event := <-eventChan:
			s.publisher.Notify(event)
		case <-s.stopChannel:
			break MessageLoop
		}
	}

	close(eventChan)
	close(s.stopChannel)
}

// Stop останавливает прослушивание SQS сообщений
func (s *SQSService) Stop() {
	s.stopChannel <- true
}

func (s *SQSService) listen(eventChan chan Event) {
	go func() {
		message, err := s.svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			//
			// какой-то код
			//
		})
		fmt.Println(message)

		var event Event
		if err != nil {
			log.Print(err)
			event = NewGeneralError(err)
			return
		} else {
			//
			// извлечь сообщение
			//
		}

		eventChan <- event
	}()
}
