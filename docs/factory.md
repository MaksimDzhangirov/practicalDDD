# DDD на практике в Golang: Фабрика

![intro](images/factory/intro.jpeg)
*Фото [Omar Flores](https://unsplash.com/@__itsflores) из [Unsplash](https://unsplash.com/)*

Когда я писал заголовок этой статьи, я пытался вспомнить первый шаблон проектирования,
который узнал из [«Банды четырех»](https://springframework.guru/gang-of-four-design-patterns/).
Я думаю это был один из следующих: [Фабричный метод](https://refactoring.guru/design-patterns/factory-method),
[Синглтон](https://refactoring.guru/design-patterns/singleton) или 
[Декоратор](https://refactoring.guru/design-patterns/decorator).

Я уверен у других разработчиков программного обеспечения существует похожая 
история. Когда они начали изучать шаблоны проектирования либо фабричный метод 
(`Factory Method`), либо [Абстрактная Фабрика](https://refactoring.guru/design-patterns/abstract-factory)
(`Abstract Factory`) были одними из первых трёх, о которых они узнали.

Сегодня любая производная шаблона Фабрики является неотъемлемой частью 
предметно-ориентированного проектирования. И его цель остается прежней даже 
спустя многие десятилетия.

> Другие статьи из DDD цикла:
> 1. [DDD на практике в Golang: Объект-значение](https://levelup.gitconnected.com/practical-ddd-in-golang-value-object-4fc97bcad70)
> 2. [DDD на практике в Golang: Сущности](https://levelup.gitconnected.com/practical-ddd-in-golang-entity-40d32bdad2a3)
> 3. [DDD на практике в Golang: Сервисы предметной области](https://levelup.gitconnected.com/practical-ddd-in-golang-domain-service-4418a1650274)
> 4. [DDD на практике в Golang: Событие предметной области](https://levelup.gitconnected.com/practical-ddd-in-golang-domain-event-de02ad492989)
> 5. [DDD на практике в Golang: Модуль](https://levelup.gitconnected.com/practical-ddd-in-golang-module-51edf4c319ec)
> 6. [DDD на практике в Golang: Агрегат](https://levelup.gitconnected.com/practical-ddd-in-golang-aggregate-de13f561e629)

## Сложная логика при создании

Мы используем шаблон Фабрика, если логика создания объекта сложна или для изоляции
процесса создания от другой бизнес-логики. В таких случаях гораздо лучше иметь 
выделенное место в коде, которое мы можем протестировать отдельно.

Когда я создаю Фабрику, в большинстве случаев, она является частью уровня 
предметной области. Таким образом, я могу использовать её везде в приложении.
Ниже приведён простой пример Фабрики.

```go
type Loan struct {
	ID uuid.UUID
	//
	// какие-то поля
	//
}

type LoanFactory interface {
	CreateShortTermLoan(specification LoanSpecification) Loan
	CreateLongTermLoan(specification LoanSpecification) Loan
}
```
*Пример шаблона Фабрика*

Шаблон Фабрика (`Factory`) тесно связан с шаблоном Спецификация `Specification`
(я расскажу о нём в следующих статьях). Здесь приведён небольшой пример с 
`LoanFactory`, `LoanSpecification` и `Loan`.

`LoanFactory` представляет собой шаблон `Factory` в DDD или более точно
фабричный метод (`Factory Method`). Он отвечает за создание и выдачу новых 
экземпляров Ссуды `Loan`, которая может меняться в зависимости от периода 
выплаты.

## Вариативность

Как уже говорилось, шаблон Фабрика можно реализовать по-разному. Форма, которая 
чаще всего используется, по крайней мере мной, - это Фабричный метод. В этом
случае мы предоставляем некоторые методы создающие нашу структуру.

```go
const (
    LongTerm = iota
    ShortTerm
)

type Loan struct {
    ID                    uuid.UUID
    Type                  int
    BankAccountID         uuid.UUID
    Amount                value_objects.Money
    RequiredLifeInsurance bool
}

type LoanFactory struct{}

func (f *LoanFactory) CreateShortTermLoan(bankAccountID uuid.UUID, amount value_objects.Money) Loan {
    return Loan{
        Type:          ShortTerm,
        BankAccountID: bankAccountID,
        Amount:        amount,
    }
}

func (f *LoanFactory) CreateLongTermLoan(bankAccountID uuid.UUID, amount value_objects.Money) Loan {
    return Loan{
        Type:                  LongTerm,
        BankAccountID:         bankAccountID,
        Amount:                amount,
        RequiredLifeInsurance: true,
    }
}
```
*Пример с фабричным методом*

В приведенном выше фрагменте кода `LoanFactory` теперь является конкретной 
реализацией фабричного метода. Он предоставляет два метода для создания 
экземпляров Сущности Ссуда (`Loan`).

В этом случае мы создаём один и тот же объект, но с различными значениями полей,
в зависимости от того является Ссуда (`Loan`) кратко- или долгосрочной. Разница
между этими двумя случаями может быть ещё более сложной и каждая дополнительная
особенность, которую нужно учесть при создании объекта, оправдывает 
существование этого шаблона.

```go
type Investment interface {
    Amount() value_objects.Money
}

type EtfInvestment struct {
    ID             uuid.UUID
    EtfID          uuid.UUID
    InvestedAmount value_objects.Money
    BankAccountID  uuid.UUID
}

func (e EtfInvestment) Amount() value_objects.Money {
    return e.InvestedAmount
}

type StockInvestment struct {
    ID               uuid.UUID
    CompanyID        uuid.UUID
    InvestedAmount value_objects.Money
    BankAccountID    uuid.UUID
}

func (s StockInvestment) Amount() value_objects.Money {
    return s.InvestedAmount
}

type InvestmentSpecification interface {
    Amount() value_objects.Money
    BankAccountID() uuid.UUID
    TargetID() uuid.UUID
}

type InvestmentFactory interface {
    Create(specification InvestmentSpecification) Investment
}

type EtfInvestmentFactory struct{}

func (f *EtfInvestmentFactory) Create(specification InvestmentSpecification) Investment {
    return EtfInvestment{
        EtfID:          specification.TargetID(),
        InvestedAmount: specification.Amount(),
        BankAccountID:  specification.BankAccountID(),
    }
}

type StockInvestmentFactory struct{}

func (f *StockInvestmentFactory) Create(specification InvestmentSpecification) Investment {
    return StockInvestment{
        CompanyID:      specification.TargetID(),
        InvestedAmount: specification.Amount(),
        BankAccountID:  specification.BankAccountID(),
    }
}
```
*Пример с Абстрактной Фабрикой*

В вышеприведенном примере дан фрагмент кода с шаблоном Абстрактная Фабрика. В
этом случае мы хотим создать несколько экземпляров интерфейса `Investment`.

Поскольку существует несколько реализаций этого интерфейса, сейчас идеальный 
момент для добавления шаблона Фабрика. И `EtfInvestmentFactory`, и 
`StockInvestmentFactory` создают экземпляры удовлетворяющие интерфейсу 
`Investment`.

В нашем коде мы можем сохранить их в некоторой карте интерфейсов 
`InvestmentFactory` и использовать их всякий раз, когда мы хотим создать
`Investment` из любого `BankAccount`.

Это идеальное место для использования абстрактной фабрики, поскольку мы должны
создавать некие объекты из определенного набора (на самом деле может существовать
ещё больше различных инвестиций).

## Преобразование

Мы можем использовать шаблон Фабрика на других уровнях. По крайней мере я его 
использую на инфраструктурном уровне и уровне представления. Там я преобразую
[Объекты для передачи данных](https://martinfowler.com/eaaCatalog/dataTransferObject.html)
(`Data Transfer Objects`) в Сущности и наоборот.

```go
// уровень предметной области
type CryptoInvestment struct {
    ID               uuid.UUID
    CryptoCurrencyID uuid.UUID
    InvestedMoney    value_objects.Money
    BankAccountID    uuid.UUID
}

// инфраструктурный уровень
type CryptoInvestmentGorm struct {
    ID                 int                 `gorm:"primaryKey;column:id"`
    UUID               string              `gorm:"column:uuid"`
    CryptoCurrencyID   int                 `gorm:"column:crypto_currency_id"`
    CryptoCurrency     CryptoCurrencyGorm  `gorm:"foreignKey:CryptoCurrencyID"`
    InvestedAmount     int                 `gorm:"column:amount"`
    InvestedCurrencyID int                 `gorm:"column:currency_id"`
    Currency           dto.CurrencyGorm    `gorm:"foreignKey:InvestedCurrencyID"`
    BankAccountID      int                 `gorm:"column:bank_account_id"`
    BankAccount        dto.BankAccountGorm `gorm:"foreignKey:BankAccountID"`
}

type CryptoInvestmentDBFactory struct {
}

func (f *CryptoInvestmentDBFactory) ToEntity(dto CryptoInvestmentGorm) (model.CryptoInvestment, error) {
    id, err := uuid.Parse(dto.UUID)
    if err != nil {
        return model.CryptoInvestment{}, err
    }
    
    cryptoId, err := uuid.Parse(dto.CryptoCurrency.UUID)
    if err != nil {
        return model.CryptoInvestment{}, err
    }
    
    currencyId, err := uuid.Parse(dto.Currency.UUID)
    if err != nil {
        return model.CryptoInvestment{}, err
    }
    
    accountId, err := uuid.Parse(dto.BankAccount.UUID)
    if err != nil {
        return model.CryptoInvestment{}, err
    }
    
    return model.CryptoInvestment{
        ID:               id,
        CryptoCurrencyID: cryptoId,
        InvestedMoney:    value_objects.NewMoney(dto.InvestedAmount, currencyId),
        BankAccountID:    accountId,
    }, nil
}
```
*Пример преобразования*

`CryptoInvestmentDBFactory` - это фабрика внутри инфраструктурного уровня, 
используемая для реконструкции объекта `CryptoInvestment`. Здесь показан 
только метод преобразования DTO в Сущность, но эта же Фабрика может иметь
метод преобразования Сущности (`Entity`) в DTO.

Поскольку `CryptoInvestmentDBFactory` использует структуру как для инфраструктуры
(`CryptoInvestmentGorm`), так и для предметной области (`CryptoInvestment`), 
она должна находиться внутри инфраструктурного уровня, поскольку у нас не может 
быть никаких зависимостей от других уровней внутри уровня предметной области.

Я всегда любил использовать UUID внутри бизнес-логики и выдавать только UUID в 
качестве ответа API. Но поскольку база данных плохо работает со строками или
двоичными данными в качестве первичных ключей, Фабрика кажется подходящим 
местом для выполнения этого преобразования.

## Заключение

Шаблон Фабрика (`Factory`) — это принцип, уходящий корнями в старые шаблоны из 
«Банды четырех». Мы можем реализовать его в виде Абстрактной Фабрики или 
фабричного метода.

Мы используем его в тех случаях, когда хотим отделить логику создания от другой
бизнес-логики. Мы также можем применять его для преобразования наших Сущностей 
в DTO и наоборот.

> Другие статьи из DDD цикла:
> 1. [DDD на практике в Golang: Объект-значение](https://levelup.gitconnected.com/practical-ddd-in-golang-value-object-4fc97bcad70)
> 2. [DDD на практике в Golang: Сущности](https://levelup.gitconnected.com/practical-ddd-in-golang-entity-40d32bdad2a3)
> 3. [DDD на практике в Golang: Сервисы предметной области](https://levelup.gitconnected.com/practical-ddd-in-golang-domain-service-4418a1650274)
> 4. [DDD на практике в Golang: Событие предметной области](https://levelup.gitconnected.com/practical-ddd-in-golang-domain-event-de02ad492989)
> 5. [DDD на практике в Golang: Модуль](https://levelup.gitconnected.com/practical-ddd-in-golang-module-51edf4c319ec)
> 6. [DDD на практике в Golang: Агрегат](https://levelup.gitconnected.com/practical-ddd-in-golang-aggregate-de13f561e629)

## Полезные ссылки на источники:

* [https://martinfowler.com/](https://martinfowler.com/)
* [https://www.domainlanguage.com/](https://www.domainlanguage.com/)