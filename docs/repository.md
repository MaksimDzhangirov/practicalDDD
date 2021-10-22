# DDD на практике в Golang: Репозиторий

![intro](images/repository/intro.jpeg)
*Фото [George Kedenburg III](https://unsplash.com/@gk3) из [Unsplash](https://unsplash.com/)*

Сегодня трудно представить себе написание какого-либо приложения без доступа
к хранилищу данных во время выполнения. Наверное даже нельзя написать сценарии 
развертывания, поскольку им нужен доступ к файлам конфигурации, которые в некотором
смысле всё ещё являются типами хранилищ.

Всякий раз, когда вы пишете какое-то приложение, которое должно решить некую
проблему в реальном мире, вам необходимо подключиться к базе данных, внешнему 
API, какой-то системе кеширования, чему угодно. Это неизбежно.

С этой точки зрения неудивительно, иметь DDD шаблон, который решает такие 
задачи. Конечно, DDD не изобрело шаблон Репозиторий (`Repository`) и примеры 
его использования есть в другой литературе, но DDD добавило больше ясности.

> Другие статьи из DDD цикла:
> 1. [DDD на практике в Golang: Объект-значение](https://levelup.gitconnected.com/practical-ddd-in-golang-value-object-4fc97bcad70)
> 2. [DDD на практике в Golang: Сущности](https://levelup.gitconnected.com/practical-ddd-in-golang-entity-40d32bdad2a3)
> 3. [DDD на практике в Golang: Сервисы предметной области](https://levelup.gitconnected.com/practical-ddd-in-golang-domain-service-4418a1650274)
> 4. [Practical DDD in Golang: Domain Event](https://levelup.gitconnected.com/practical-ddd-in-golang-domain-event-de02ad492989)
> 5. [Practical DDD in Golang: Module](https://levelup.gitconnected.com/practical-ddd-in-golang-module-51edf4c319ec)
> 6. [Practical DDD in Golang: Aggregate](https://levelup.gitconnected.com/practical-ddd-in-golang-aggregate-de13f561e629)
> 7. [Practical DDD in Golang: Factory](https://levelup.gitconnected.com/practical-ddd-in-golang-factory-5ba135df6362)

## Предохранительный уровень

Предметно-ориентированное проектирование — это принцип, который мы можем 
применить ко многим аспектам разработки программного обеспечения и во многих 
местах.

Поскольку Репозиторий (`Repository`) всегда представляет собой структуру, в 
которой хранятся технические подробности о подключении к какому-либо внешнему 
источнику, он уже не относится к нашей бизнес-логике.

Но время от времени нам нужен доступ к Репозиторию (`Repository`) на уровне 
предметной области. Поскольку уровень предметной области находится в самом низу
и не взаимодействует с другими мы определяем репозиторий внутри него, но как
интерфейс.

```go
type Customer struct {
    id      uuid.UUID
    //
    // какие-то поля
    //
}

type Customers []model.Customer

type CustomerRepository interface {
    GetCustomer(ctx context.Context, ID uuid.UUID) (*model.Customer, error)
    Search(ctx context.Context, specification CustomerSpecification) (Customers, int, error)
    SaveCustomer(ctx context.Context, customer model.Customer) (*model.Customer, error)
    UpdateCustomer(ctx context.Context, customer model.Customer) (*model.Customer, error)
    DeleteCustomer(ctx context.Context, ID uuid.UUID) (*model.Customer, error)
}
```
*Пример контракта*

Этот интерфейс мы называем Контрактом, определяющий сигнатуры методом, которые
мы можем вызывать внутри нашего домена. Внутри вышеприведенного примера мы 
видим простой интерфейс, определяющий [CRUD](https://www.codecademy.com/articles/what-is-crud) 
методы.

Поскольку мы определили Репозиторий (`Repository`) как такой интерфейс, мы можем
использовать его везде внутри предметной области. Он всегда ожидает и возвращает
нам наши Сущности, в данном случае Клиента (`Customer`) и Клиентов (`Customers`)
(мне нравится определять такие конкретные коллекции в Go, связывая с ними различные
методы).

Сущность Customer не хранит никакой информации о типе хранилища: нет дескриптора
Go, определяющего структуру JSON, столбцы Gorm или что-либо подобное. Для этого
мы должны использовать инфраструктурный уровень.

```go
// уровень предметной области
type CustomerRepository interface {
    GetCustomer(ctx context.Context, ID uuid.UUID) (*model.Customer, error)
    Search(ctx context.Context, specification CustomerSpecification) (Customers, int, error)
    SaveCustomer(ctx context.Context, customer model.Customer) (*model.Customer, error)
    UpdateCustomer(ctx context.Context, customer model.Customer) (*model.Customer, error)
    DeleteCustomer(ctx context.Context, ID uuid.UUID) (*model.Customer, error)
}

// инфраструктурный уровень

type CustomerGorm struct {
    ID   int    `gorm:"primaryKey;column:id"`
    UUID string `gorm:"uniqueIndex;column:uuid"`
    //
    // какие-то поля
    //
}

func (c CustomerGorm) ToEntity() (model.Customer, error) {
    parsed, err := uuid.Parse(c.UUID)
    if err != nil {
        return model.Customer{}, err
    }

    return model.Customer{
        ID: parsed,
        //
        // какие-то поля
        //
    }, nil
}

type CustomerRepository struct {
    connection *gorm.DB
}

func (r *CustomerRepository) GetCustomer(ctx context.Context, ID uuid.UUID) (*model.Customer, error) {
    var row CustomerGorm
    err := r.connection.WithContext(ctx).Where("uuid = ?", ID).First(&row).Error
    if err != nil {
        return nil, err
    }

    customer, err := row.ToEntity()
    if err != nil {
        return nil, err
    }

    return &customer, nil
}

//
// другие методы
//
```
*Фактическая реализация на инфраструктурном уровне*

В вышеприведенном примере вы можете увидеть фрагмент реализации 
`CustomerRepository`. Внутри он использует `Gorm` для упрощения интеграции,
но вы также можете воспользоваться чистыми SQL запросами. В последнее время
я часто использую библиотеку [Ent](https://entgo.io/).

В этом примере вы видите две разные структуры: `Customer` и `CustomerGorm`. 
Первая — это Сущность, где мы хотим хранить нашу бизнес-логику, некоторые 
инварианты предметной области и правила. Она ничего не знает об используемой
базе данных.

Вторая структура — это Объект передачи данных (`Data Transfer Object`), который
определяет как наши данные передаются из и в хранилище. У этой структуры нет 
никаких других задач, кроме как сопоставить данные из базы данных с нашей 
Сущностью.

> Разделение этих двух структур является ключевым моментом при
> использовании репозитория в качестве предохранительного уровня в нашем 
> приложении. Это гарантирует, что технические детали, связанные со структурой 
> таблицы, никак не повлияют на нашу бизнес-логику.

К чему это приведёт? Во-первых, нам нужно поддерживать два типа структур: одну
для бизнес-логики, другую — для хранения. Кроме того я также добавлю третью 
структуру, которую будут использовать как Объект передачи данных в своём API.

Такой подход усложняет наше приложение и добавляет множество функций сопоставления,
подобных тем, что показаны в примере ниже. Также нужно протестировать такие 
методы соответствующим образом, чтобы избежать ошибок копирования и вставки.

```go
// уровень предметной области

type Customer struct {
    ID      uuid.UUID
    Person  *Person
    Company *Company
    Address Address
}

type Birthday time.Time

type Person struct {
    SSN       string
    FirstName string
    LastName  string
    Birthday  Birthday
}

type Company struct {
    Name               string
    RegistrationNumber string
    RegistrationDate   time.Time
}

type Address struct {
    Street   string
    Number   string
    Postcode string
    City     string
}

// инфраструктурный уровень

type CustomerGorm struct {
    ID        uint         `gorm:"primaryKey;column:id"`
    UUID      string       `gorm:"uniqueIndex;column:id"`
    PersonID  uint         `gorm:"column:person_id"`
    Person    *PersonGorm  `gorm:"foreignKey:PersonID"`
    CompanyID uint         `gorm:"column:company_id"`
    Company   *CompanyGorm `gorm:"foreignKey:CompanyID"`
    Street    string       `gorm:"column:street"`
    Number    string       `gorm:"column:number"`
    Postcode  string       `gorm:"column:postcode"`
    City      string       `gorm:"column:city"`
}

func (c CustomerGorm) ToEntity() (model.Customer, error) {
    parsed, err := uuid.Parse(c.UUID)
    if err != nil {
        return model.Customer{}, err
    }

    return model.Customer{
        ID:      parsed,
        Person:  c.Person.ToEntity(),
        Company: c.Company.ToEntity(),
        Address: model.Address{
            Street:   c.Street,
            Number:   c.Number,
            Postcode: c.Postcode,
            City:     c.City,
        },
    }, nil
}

type PersonGorm struct {
    ID        uint      `gorm:"primaryKey;column:id"`
    SSN       string    `gorm:"uniqueIndex;column:ssn"`
    FirstName string    `gorm:"column:first_name"`
    LastName  string    `gorm:"column:last_name"`
    Birthday  time.Time `gorm:"column:birthday"`
}

func (p *PersonGorm) ToEntity() *model.Person {
    if p == nil {
        return nil
    }

    return &model.Person{
        SSN:       p.SSN,
        FirstName: p.FirstName,
        LastName:  p.LastName,
        Birthday:  model.Birthday(p.Birthday),
    }
}

type CompanyGorm struct {
    ID                 uint      `gorm:"primaryKey;column:id"`
    Name               string    `gorm:"column:name"`
    RegistrationNumber string    `gorm:"column:registration_number"`
    RegistrationDate   time.Time `gorm:"column:registration_date"`
}

func (c *CompanyGorm) ToEntity() *model.Company {
    if c == nil {
        return nil
    }

    return &model.Company{
        Name:               c.Name,
        RegistrationNumber: c.RegistrationNumber,
        RegistrationDate:   c.RegistrationDate,
    }
}

func NewRow(customer model.Customer) dto.CustomerGorm {
    var person *dto.PersonGorm
    if customer.Person != nil {
        person = &dto.PersonGorm{
            SSN:       customer.Person.SSN,
            FirstName: customer.Person.FirstName,
            LastName:  customer.Person.LastName,
            Birthday:  time.Time(customer.Person.Birthday),
        }
    }

    var company *dto.CompanyGorm
    if customer.Company != nil {
        company = &dto.CompanyGorm{
            Name:               customer.Company.Name,
            RegistrationNumber: customer.Company.RegistrationNumber,
            RegistrationDate:   customer.Company.RegistrationDate,
        }
    }

    return dto.CustomerGorm{
        UUID:     uuid.NewString(),
        Person:   person,
        Company:  company,
        Street:   customer.Address.Street,
        Number:   customer.Address.Number,
        Postcode: customer.Address.Postcode,
        City:     customer.Address.City,
    }
}
```
*Предохранительный уровень и Объекты передачи данных*

Тем не менее, несмотря на все эти сложности, у нашего кода появляется новое 
преимущество. Мы можем представить наши Сущности внутри уровня предметной 
области так, чтобы они наилучшим образом описывали нашу бизнес-логику. Мы не 
ограничиваем их используемым хранилищем.

Мы можем использовать один тип идентификаторов внутри нашей бизнес-логики 
(например, UUID), а другой — для базы данных (беззнаковое целое). Это касается 
любых данных, которые мы хотим использовать для базы данных и бизнес-логики.

Всякий раз, когда мы вносим изменения в любой из этих слоёв, мы вероятно, 
будем вносить модифицировать также функции сопоставления, но ничего не меняем
в других слоях (или, хотя бы, не удаляем).

Мы можем захотеть перейти на MongoDB, Cassandra или любой другой тип 
хранилища. Мы можем перейти на внешний API, но это не повлияет на наш уровень 
предметной области.

## Сохранение

Мы используем Репозиторий в первую очередь для запросов. Он отлично работает с
другим DDD шаблоном, Спецификацией (`Specification`), которую вы возможно заметили
в примерах. Мы можем использовать его без спецификации, но это иногда облегчает 
нам жизнь.

Второй способ использования Репозитория — сохранение данных. Мы определяем какие
именно данные должны быть отправлены в хранилище для их сохранения, обновления
и даже удаления.

```go
func NewRow(customer model.Customer) dto.CustomerGorm {
    return dto.CustomerGorm{
        UUID:     uuid.NewString(),
        //
        // какие-то поля
        //
    }
}

type CustomerRepository struct {
    connection *gorm.DB
}

func (r *CustomerRepository) SaveCustomer(ctx context.Context, customer model.Customer) (*model.Customer, error) {
    row := NewRow(customer)
    err := r.connection.WithContext(ctx).Save(&row).Error
    if err != nil {
        return nil, err
    }

    customer, err = row.ToEntity()
    if err != nil {
        return nil, err
    }

    return &customer, nil
}

//
// другие методы
//
```
*Сохранение с генерацией UUID*

Иногда мы принимаем решение, что хотим иметь уникальные идентификаторы, которые
должны создаваться в приложении. В таких случаях Репозиторий подходящее место 
для этого. В вышеприведенном примере видно, что мы генерируем новый UUID перед 
сохранением записи в базе данных.

Мы можем использовать для этого целые числа, если хотим избежать автоматического
инкремента, используемого движком базы данных. В любом случае, если мы не хотим
полагаться на ключи в базе данных, необходимо создавать их внутри Репозитория.

```go
type CustomerRepository struct {
    connection *gorm.DB
}

func (r *CustomerRepository) CreateCustomer(ctx context.Context, customer model.Customer) (*model.Customer, error) {
    tx := r.connection.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    if err := tx.Error; err != nil {
        return nil, err
    }
    
    //
    // какой-то код
    //
    
    var total int64
    var err error
    if customer.Person != nil {
        err = tx.Model(dto.PersonGorm{}).Where("ssn = ?", customer.Person.SSN).Count(&total).Error
    } else if customer.Person != nil {
        err = tx.Model(dto.CompanyGorm{}).Where("registration_number = ?", customer.Company.RegistrationNumber).Count(&total).Error
    }
    if err != nil {
        tx.Rollback()
        return nil, err
    } else if total > 0 {
        tx.Rollback()
        return nil, errors.New("there is already such record in DB")
    }
    
    //
    // какой-то код
    //
    row := NewRow(customer)
    err = tx.Save(&row).Error
    if err != nil {
        tx.Rollback()
        return nil, err
    }
    
    err = tx.Commit().Error
    if err != nil {
        tx.Rollback()
        return nil, err
    }
    
    customer, err = row.ToEntity()
    if err != nil {
        tx.Rollback()
        return nil, err
    }
    
    return &customer, nil
}
```
*Транзакция в базе данных*

Также мы хотим использовать Репозиторий для транзакций. Каждый раз когда мы 
хотим сохранить какие-то данные и выполнить множество запросов, работающие 
с несколькими таблицами, лучше всего определить транзакцию, которая должна
находиться внутри Репозитория.

В вышеприведенном примере мы проверяем уникальность человека или компании.
Если они существуют в базе, мы возвращаем ошибку. Всё это мы можем определить
как часть одной транзакции и если что-то пошло не так, откатить её.

Здесь Репозиторий — идеальное место для такого кода. Преимущество заключается 
также в том, что если мы упростим операцию добавления в будущем и избавимся от 
транзакций, то нам не нужно будет менять Контракт Репозитория, а только код 
внутри.

## Типы Репозиториев

Будет ошибкой предполагать, что мы должны использовать Репозиторий только для
базы данных. Да, мы чаще всего используем его с базами данных, поскольку они 
в первую очередь используются для хранения данных, но сегодня становятся всё 
более популярны и другие типы.

Как было сказано ранее, мы можем использовать `MongoDB` или `Cassandra`. Репозиторий
можно применять для хранения нашего кэша и в этом случае примером может быть
`Redis`. Это может быть даже REST API или конфигурационный файл.

```go
// Репозиторий для Redis

type CustomerRedisRepository struct {
    client *redis.Client
}

    func (r *CustomerRedisRepository) GetCustomer(ctx context.Context, ID uuid.UUID) (*model.Customer, error) {
    data, err := r.client.Get(ctx, fmt.Sprintf("user-%s", ID.String())).Result()
    if err != nil {
        return nil, err
    }
    
    var row dto.CustomerJSON
    err = json.Unmarshal([]byte(data), &row)
    if err != nil {
        return nil, err
    }
    
    customer, err := row.ToEntity()
    if err != nil {
        return nil, err
    }
    
    return &customer, nil
}

// API

type CustomerRedisAPIRepository struct {
    client  *http.Client
    baseUrl string
}

func (r *CustomerRedisAPIRepository) GetCustomer(ctx context.Context, ID uuid.UUID) (*model.Customer, error) {
    resp, err := r.client.Get(path.Join(r.baseUrl, "users", ID.String()))
    if err != nil {
        return nil, err
    }
    
    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var row dto.CustomerJSON
    err = json.Unmarshal(data, &row)
    if err != nil {
        return nil, err
    }
    
    customer, err := row.ToEntity()
    if err != nil {
        return nil, err
    }
    
    return &customer, nil
}
```
*Различные типы хранилищ*

Теперь мы видим реальную пользу от разделения нашей бизнес-логики и технических
деталей. Интерфейс Репозитория не меняется, поэтому наш уровень предметной области
всегда может его использовать.

Но когда-нибудь наше приложение может вырасти до такой степени, что MySQL 
перестанет быть идеальным решением для него. Если использовались миграции, то
не нужно переживать, о том что будет затронута наша бизнес-логика, пока не 
меняются интерфейсы.

> Итак, Интерфейс Репозитория всегда должен работать с вашей бизнес-логикой, а
> ваша реализация Репозитория должна использовать внутренние структуры, которые
> вы можете позже сопоставить с Сущностями.

## Заключение

Репозиторий — это популярный шаблон, отвечающий за запросы и сохранение данных
внутри какого-либо хранилища данных. Это основной предохранительный механизм
внутри нашего приложения.

Мы определяем его как интерфейс внутри уровня предметной области и сохраняем
фактическую реализацию внутри инфраструктурного уровня. Здесь генерируются 
создаваемые приложением идентификаторы и выполняются транзакции.

> Другие статьи из DDD цикла:
> 1. [DDD на практике в Golang: Объект-значение](https://levelup.gitconnected.com/practical-ddd-in-golang-value-object-4fc97bcad70)
> 2. [DDD на практике в Golang: Сущности](https://levelup.gitconnected.com/practical-ddd-in-golang-entity-40d32bdad2a3)
> 3. [DDD на практике в Golang: Сервисы предметной области](https://levelup.gitconnected.com/practical-ddd-in-golang-domain-service-4418a1650274)
> 4. [Practical DDD in Golang: Domain Event](https://levelup.gitconnected.com/practical-ddd-in-golang-domain-event-de02ad492989)
> 5. [Practical DDD in Golang: Module](https://levelup.gitconnected.com/practical-ddd-in-golang-module-51edf4c319ec)
> 6. [Practical DDD in Golang: Aggregate](https://levelup.gitconnected.com/practical-ddd-in-golang-aggregate-de13f561e629)
> 7. [Practical DDD in Golang: Factory](https://levelup.gitconnected.com/practical-ddd-in-golang-factory-5ba135df6362)

## Полезные ссылки на источники:

* [https://martinfowler.com/](https://martinfowler.com/)
* [https://www.domainlanguage.com/](https://www.domainlanguage.com/)