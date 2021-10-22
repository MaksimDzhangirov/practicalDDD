# DDD на практике в Golang: Событие предметной области

![intro](images/domain-event/intro.jpeg)
*Фото [Anthony DELANOIX](https://unsplash.com/@anthonydelanoix) из [Unsplash](https://unsplash.com/)*

Во многих случаях Сущности — наилучший способ описать что-либо при 
предметно-ориентированном проектировании. Вместе с объектами-значениями они 
предоставляют наиболее полную картину рассматриваемой предметной области.

Иногда отличный способ описать рассматриваемую предметную область — использовать
события, происходящие в ней. На самом деле я всё чаще пытаюсь определить события,
а затем Сущности, связанные с ними.

Хотя Эрик Эванс не рассмотрел шаблон Событие предметной области (`Domain Event`)
в первом издании своей книги, сегодня сложно представить уровень предметной 
области без использования событий.

Шаблон Событие предметной области описывает возникающие события в нашем коде. 
Мы можем использовать его для представления любого явления из реального мира,
которое имеет отношение к нашей бизнес-логике. Сегодня всё в деловом мире связано
с какими-то событиями.

> Другие статьи из DDD цикла:
> 1. [DDD на практике в Golang: Объект-значение](https://levelup.gitconnected.com/practical-ddd-in-golang-value-object-4fc97bcad70)
> 2. [DDD на практике в Golang: Сущности](https://levelup.gitconnected.com/practical-ddd-in-golang-entity-40d32bdad2a3)
> 3. [DDD на практике в Golang: Сервисы предметной области](https://levelup.gitconnected.com/practical-ddd-in-golang-domain-service-4418a1650274)

## Им может быть что угодно

Событиями предметной области может быть что угодно, но они должны удовлетворять
некоторым правилам. Во-первых, они неизменяемы. Чтобы следовать этому правилу
я всегда использую приватные поля внутри структуры `Event`, даже если я не большой
поклонник приватных полей и геттеров в Go. По крайней мере, у событий не так много
геттеров.

Одно конкретное событие может произойти только один раз. Это означает, что мы
можем только один раз создать Сущность Заказ (`Order`) с неким идентификатором, 
поэтому только один раз наш код может инициировать событие (`Event`), описывающее 
создание этого Заказа.

Любое другое Событие для этого Заказа будет иметь другой тип. Любое другое, описывающее
создание, Событие будет относиться к другому Заказу.

Каждое Событие практически описывает то, что уже произошло. Оно представляет 
прошлое. Это означает, что мы запускаем событие `OrderCreated`, когда уже создали
`Order`, а не до этого.

```go
// Интерфейс Event для описания События предметной области
type Event interface {
    Name() string
}

// Событие GeneralError
type GeneralError string

func NewGeneralError(err error) Event {
    return GeneralError(err.Error())
}

func (e GeneralError) Name() string {
    return "event.general.error"
}

// Интерфейс OrderEvent для описания События предметной области, связанных с Заказом
type OrderEvent interface {
    Event
    OrderID() uuid.UUID
}

// Событие OrderDispatched
type OrderDispatched struct {
    orderID uuid.UUID
}

func (e OrderDispatched) Name() string {
    return "event.order.dispatched"
}

func (e OrderDispatched) OrderID() uuid.UUID {
    return e.orderID
}

// Событие OrderDelivered
type OrderDelivered struct {
    orderID uuid.UUID
}

func (e OrderDelivered) Name() string {
    return "event.order.delivery.success"
}

func (e OrderDelivered) OrderID() uuid.UUID {
    return e.orderID
}

// Событие OrderDeliveryFailed
type OrderDeliveryFailed struct {
    orderID uuid.UUID
}

func (e OrderDeliveryFailed) Name() string {
    return "event.order.delivery.failed"
}

func (e OrderDeliveryFailed) OrderID() uuid.UUID {
    return e.orderID
}
```
*Простые События*

В приведённом выше примере кода показаны простые События предметной области. 
Этот код — один из миллиардов реализующих их на Go. В некоторых случаях, как здесь
для `GeneralError`, я использовал простые строки.

Но иногда я создавал сложные объекты. Или мне приходилось расширить основной 
интерфейс `Event` каким-то более специфическим, чтобы добавить дополнительные 
методы, например, как в случае с `OrderEvent`.

Для События предметной области, поскольку оно является интерфейсом, не нужно 
реализовывать какие-либо методы. Им может быть что угодно. Как я уже говорил, 
иногда я использую строки, но достаточно и чего-либо другого. Для обобщения 
время от времени я все же объявляю интерфейс `Event`.

## Старый друг

Событие предметной области как шаблон, не является чем-то новым, а
является всего лишь другим представлением шаблона Наблюдатель (`Observer`). 
Шаблон Наблюдатель включает Издателя (`Publisher`), Подписчика (`Subscriber`)
как основных исполнителей и конечно же Событие (`Event`).

Событие предметной области использует ту же логику. Подписчик (`Subscriber`) 
или Обработчик Cобытий (`Event Handler`) — это структура, реагирующая на 
конкретное событие домена, на которое подписана. Издатель (`Publisher`) - это
структура, уведомляющая все Обработчики Событий (`Event Handlers`) о том, что 
какое-то событие произошло.

Издатель — это точка входа для запуска любого События. Он содержит все 
Обработчики Событий и предоставляет простой интерфейс для любого сервиса 
предметной области (`Domain Service`), фабрики (`Factory`) или других объектов,
которые хотят опубликовать какое-либо событие.

```go
// Интерфейс EventHandler, описывающий любой объект, который должен быть
// уведомлен о каком-либо Event
type EventHandler interface {
	Notify(event Event)
}

// EventPublisher - основная структура, уведомляющая все EventHandler
type EventPublisher struct {
    handlers map[string][]EventHandler
}

// Метод Subscribe подписывает EventHandler на определённое событие (Event)
func (e *EventPublisher) Subscribe(handler EventHandler, events ...Event) {
    for _, event := range events {
        handlers := e.handlers[event.Name()]
        handlers = append(handlers, handler)
        e.handlers[event.Name()] = handlers
    }
}

// Метод Notify уведомляет подписанный EventHandler о том, что произошло определенное событие (Event)
func (e *EventPublisher) Notify(event Event) {
    for _, handler := range e.handlers[event.Name()] {
        handler.Notify(event)
    }
}
```
*Пример с `EventHandler` и `EventPublisher`*

Приведенный выше фрагмент кода показывает остальную часть шаблона Событие 
предметной области. Интерфейс `EventHandler` представляет собой любую структуру, 
которая должна реагировать на какое-либо событие. У него есть только один 
метод `Notify`, который ожидает событие в качестве аргумента.

Структура `EventPublisher` более сложная. Она предоставляет общий метод `Notify`,
который отвечает за уведомление всех Обработчиков Событий, подписанных на него.
Ещё один метод `Subscribe` позволяет любому `EventHandler` подписаться на любое 
событие.

Структура `EventPublisher` может быть менее сложной. Вместо того, чтобы давать
возможность `EventHandler` подписаться на конкретное событие `Event`, используя 
map, у него может быть простой массив с `EventHandler`. Он будет уведомлять все
обработчики о любом событии.

в общем случае мы должны публиковать События предметной области синхронно на 
уровне предметной области. Но иногда по какой-то причине я могу запускать их 
асинхронно. Для этой цели я использую [Goroutine](https://tour.golang.org/concurrency/1).

```go
type Event interface {
    Name() string
    IsAsynchronous() bool
}

type EventPublisher struct {
    handlers map[string][]EventHandler
}

func (e *EventPublisher) Notify(event Event) {
    if event.IsAsynchronous() {
        go e.notify(event) // запускаем код в отдельной горутине
    }

    e.notify(event) // синхронный вызов
}

func (e *EventPublisher) notify(event Event) {
    for _, handler := range e.handlers[event.Name()] {
        handler.Notify(event)
    }
}
```
*Запускаем События асинхронно*

В приведенном выше примере показан один из вариантов асинхронной публикации 
событий. Чтобы реализовать оба подхода, я часто определяю метод внутри 
интерфейса `Event`, который позже предоставляет мне информацию должен ли я 
запускать событие синхронно или нет.

## Создание

Моя самая большая дилемма заключалась в том, где правильное место для 
создания события. И, честно говоря, я задавал их везде. Единственное правило, 
которое у меня было, заключалось в том, что объекты, отслеживающие состояния, 
не могли уведомлять `EventPublisher`.

Сущности (`Entity`), Объекты-значения (`Value Objects`) и Агрегаты 
(`Aggregates`) (которые мы рассмотрим в следующей статье) являются объектами, 
отслеживающими состояния. С этой точки зрения они не должны содержать внутри 
себя `EventPublisher`, и передача его в качестве аргумента их методам
я всегда считал безобразным кодом.

Кроме того, я не использую объекты, отслеживающие состояния, в качестве 
обработчиков событий (`EventHandlers`). Если мне нужно было бы что-то сделать с
какой-либо Сущностью (`Entity`), когда происходит конкретное Событие 
(`Event`), я бы создал `EventHandler`, содержащий репозиторий (`Repository`).
Из репозитория можно получить Сущность, которую следует модифицировать.

Тем не менее создание объектов `Event` внутри какого-либо метода Агрегата 
(`Aggregate`) - это нормально. Иногда я инициирую их внутри метода Сущности 
(`Entity`) и возвращаю как результат. Затем я используя структуры, не 
хранящие состояния, например, Сервисы предметной области (`Domain Service`) 
или Фабрики (`Factory`) для уведомления `EventPublisher`.

```go
type Order struct {
    id uuid.UUID
    //
    // какие-то поля
    //
    isDispatched bool
    deliverAddress value_objects.Address
}

func (o Order) ID() uuid.UUID {
    return o.id
}

func (o Order) ChangeAddress(address value_objects.Address) events.Event {
    if o.isDispatched {
        return events.NewDeliveryAddressChangeFailed(o.ID())
    }
    //
    // какой-то код
    //
    return events.NewDeliveryAddressChanged(o.ID())
}

type OrderService struct {
    repository repository.OrderRepository
    publisher events.EventPublisher
}

func (s *OrderService) Create(order entity.Order) (*entity.Order, error) {
    result, err := s.repository.Create(order)
    if err != nil {
        return nil, err
    }
    //
    // обновляем адрес в базе данных
    //
    s.publisher.Notify(events.NewOrderCreated(result.ID()))
    
    return result, nil
}

func (s *OrderService) ChangeAddress(order entity.Order, address value_objects.Address) {
    evt := order.ChangeAddress(address)
    
    s.publisher.Notify(evt) // публикуем события только внутри объект, не хранящих состояние
}
```
*Создание Событий*

В приведенном выше примере Агрегат `Order` содержит метод для обновления
адресов доставки. Результатом работы этого метода может быть Событие (`Event`).
Это означает, что `Order` может создавать некоторые события, но не более.

С другой стороны, `OrderService` может как создавать События, так и публиковать
их. Он также может инициировать события, которые получает от `Order`, при обновлении
адреса доставки. Это возможно, поскольку он содержит `EventPublisher`.

## События на других уровнях

Мы можем прослушивать события на других уровнях, например, прикладных операций,
представления или инфраструктуры. Мы также можем определить отдельные События,
которые будут относиться только к этим уровням. В таких случаях мы не говорим о
Событиях предметной области.

Простым примером являются события на уровне прикладных операций. После создания
Заказа (`Order`) в большинстве случаев мы должны отправить клиенту электронное
письмо (`Email`). Хотя это может выглядеть как бизнес-правило, отправка 
электронных писем всегда зависит от приложения.

В приведенном ниже примере показан простой код с `EmailEvent`. Как вы наверное 
догадались электронное письмо (`Email`) может иметь различные состояния 
и переход от одного к другому всегда выполняется во время некоторых событий 
(`Events`).

```go
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
            s.publisher.Nofity(event)
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
```
*Пример с Событиями прикладных операций*

Иногда мы хотим инициировать событие предметной области вне нашего [Ограниченного
контекста](https://martinfowler.com/bliki/BoundedContext.html). Эти события 
предметной области являются внутренними событиями для нашего Ограниченного 
контекста, но они являются внешними для других.

Хотя эта тема относится больше к стратегическому предметно-ориентированному 
проектированию, я коснусь её здесь. Чтобы создать Событие вне нашего Микросервиса,
мы можем использовать какой-то сервис обмена сообщениями, например, 
[SQS](https://aws.amazon.com/sqs/).

```go
// инфраструктурный уровень
import (
    "encoding/json"
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
```
*Передаём внутренние События во внешний мир*

В приведенном выше фрагменте кода есть `EventSQSHandler`, простая структура 
в инфраструктурном уровне, которая отправляет сообщение в очередь SQS всякий раз,
когда происходит какое-либо событие. Она публикует только названия событий 
без каких-либо конкретных деталей.

Публикуя внутренние События во внешний мир, мы также можем прослушивать внешние
События и сопоставлять их с внутренними. Для этого я всегда создаю какой-то
Сервис в инфраструктурном уровне, который прослушивает события извне.

```go
// инфраструктурный уровень
type SQSService struct {
    svs       *sqs.SQS
    publisher *EventPublisher
    stopChannel chan bool
}

// Run запускает прослушивание SQS сообщений
func (s *SQSService) Run(event Event) {
    eventChan := make(chan Event)
    
MessageLoop:
    for {
        s.listen(eventChan)
    
        select {
        case event := <- eventChan:
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
```
*Прослушиваем внешние События*

В приведенном выше примере показан SQSService внутри инфраструктурного уровня.
Этот Сервис прослушивает SQS сообщения и сопоставляет их с внутренними 
событиями, если это возможно.

Я нечасто использовал этот подход, но в некоторых случаях он того стоил. Например, 
если несколько микросервисов должны отреагировать на создание Заказа (`Order`)
или когда регистрируется Клиент (`Customer`).

## Заключение

События предметной области — это неотъемлемая часть нашей логики предметной 
области. Сегодня все в деловом мире привязано к определенным событиям, поэтому 
описание нашей модели предметной области с помощью событий является хорошей
практикой.

Шаблон Событие предметной области — это просто реализация шаблона Наблюдатель.
Он может быть создан внутри многих объектов, но должен инициироваться только из
объектов, не хранящих состояние. Другие уровни также могут использовать 
события предметной области или свои собственные.

> Другие статьи из DDD цикла:
> 1. [DDD на практике в Golang: Объект-значение](https://levelup.gitconnected.com/practical-ddd-in-golang-value-object-4fc97bcad70)
> 2. [DDD на практике в Golang: Сущности](https://levelup.gitconnected.com/practical-ddd-in-golang-entity-40d32bdad2a3)
> 3. [DDD на практике в Golang: Сервисы предметной области](https://levelup.gitconnected.com/practical-ddd-in-golang-domain-service-4418a1650274)

## Полезные ссылки на источники:

* [https://martinfowler.com/](https://martinfowler.com/)
* [https://www.domainlanguage.com/](https://www.domainlanguage.com/)