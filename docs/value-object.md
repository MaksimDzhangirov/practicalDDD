# DDD на практике в Golang: Объект-значение

Давайте начнём наш обзор практического использования предметно-ориентированного 
проектирования в Golang с наиболее важного шаблона — Объекта-значения.

![intro](images/value-object/intro.jpeg)
*Фото [Jason Leung](https://unsplash.com/@ninjason) из [Unsplash](https://unsplash.com/)*

Утверждение о том, что какой-то шаблон является наиболее важным, может показаться
преувеличенным, но я бы даже не стал спорить с ним. Впервые об [Объекте-значении](https://martinfowler.com/bliki/ValueObject.html)
я узнал из ["Большой Красной Книги"](https://www.amazon.de/-/en/Martin-Fowler/dp/0321127420) 
[Мартина Фаулера](https://twitter.com/martinfowler) (Martin Fowler). На тот момент это выглядело довольно просто и 
не очень интересно. В следующий раз я прочитал об этом в ["Большой синей книге"](https://www.amazon.de/-/en/Eric-J-Evans/dp/0321125215)
[Эрика Эванса](https://twitter.com/ericevans0?lang=en) (Eric Evans). С этого момента
шаблон начал приобретать все больший и больший смысл, и вскоре я уже не мог 
представить как писать свой код, не используя практически везде 
Объекты-значения.

## Просто, но элегантно

Объект-значение — на первый взгляд довольно простой шаблон. Он группирует несколько
атрибутов как единое целое, добавляя к ним определённое поведение. Это единое целое 
представляет собой определенную качественную или количественную величину,
которая существует в реальном мире, и его можно связать с другим более 
сложным объектом. Оно обладает определенным значением или характеристикой. Примером
может быть цвет или деньги (подтип Объекта-значения), номер телефона или любой
другой небольшой объект, представляющий собой какое-либо значение, как во 
фрагменте кода ниже.

```go
type Currency struct {
    ID uuid.UUID
    Code string
    HTML int
}

type Money struct {
    Value    float64
    Currency Currency
}

func (m Money) ToHTML() string {
    return fmt.Sprintf(`%.2f %d`, m.Value, m.Currency.HTML)
}

type Salutation string

func (s Salutation) IsPerson() bool {
    return s != "company"
}

type Color struct {
    Red   byte
    Green byte
    Blue  byte
}

func (c Color) ToCSS() string {
    return fmt.Sprintf(`rgb(%d, %d, %d`, c.Red, c.Green, c.Blue)
}

type Address struct {
    Street   string
    Number   int
    Suffix   string
    Postcode int
}

type Phone struct {
    CountryPrefix string
    AreaCode      string
    Number        string
}
```

В Golang Объекты-значения могут быть представлены в виде создаваемых пользователем
структур или путём расширения какого-либо примитивного типа. В обоих случаях идея 
состоит в обеспечении дополнительного поведения, уникального для этого отдельного
значения или группы значений. Во многих случаях Объект-значение может предоставлять
определенные методы для форматирования строк, описывающих как значения должны себя 
вести при JSON кодировании или декодировании. Тем не менее, основная цель этих методов
должна заключаться в поддержке бизнес-инвариантов, связанных с этой характеристикой 
или качеством в реальной жизни.

## Идентификация и равенство

Объект-значение не имеет никаких идентификационных данных и это его критическое
отличие от шаблона [Сущность](https://enterprisecraftsmanship.com/posts/entity-vs-value-object-the-ultimate-list-of-differences/) (`Entity`).
Шаблон Сущность имеет идентификатор, определяющий его уникальность. Если две
Сущности имеют одинаковый идентификатор, то мы можем говорить о них как об одном
и том же объекте. У объекта-значения нет такого идентификатора. У него есть только
несколько полей, которые позволяют лучше описать его значение. Чтобы проверить равны
ли два Объекта-значения, нужно проверить на равенство все его поля, как во 
фрагменте кода, показанном ниже.

```go
// проверяем на равенство Объекты-значения
func (c Color) EqualTo(other Color) bool {
    return c.Red == other.Red && c.Green == other.Green && c.Blue == other.Blue
}

// проверяем на равенство Объекты-значения
func (m Money) EqualTo(other Money) bool {
    return m.Value == other.Value && m.Currency.EqualTo(other.Currency)
}

// проверяем на равенство Сущности
func (c Currency) EqualTo(other Currency) bool {
    return c.ID.String() == other.ID.String()
}
```

В приведенном выше примере для структур `Money` и `Color` определены методы
`EqualTo`, которые проверяют на равенство все их поля. С другой стороны, 
Currency проверяет на равенство идентификаторы, которым в этом примере является
UUID.

Как вы возможно заметили, Объект-значение также может ссылаться на некоторую
Сущность, например, `Money` и `Currency` в этом примере. Он также может 
содержать другие Объекты-значения (например, структура `Coin` состоит из 
`Color` и `Money`) или задаваться в виде среза на коллекцию (`Colors`).

```go
type Coin struct {
    Value Money
    Color Color
}

type Colors []Color
```

В одном [Ограниченном Контексте](https://martinfowler.com/bliki/BoundedContext.html) у нас 
могут быть десятки объектов-значений. Тем не менее, некоторые из них могут 
быть Сущностями внутри других Ограниченных Контекстов. Примером может быть 
`Currency`. В простом веб-сервисе, где мы хотим отображать определённые суммы 
денег, мы можем рассматривать `Currency` как Объект-значение, связанное с 
`Money`, которые мы не планируем изменять. С другой стороны в сервисе `Payment`
мы хотим получать обновления в реальном времени с помощью некоторого API 
сервиса `Exchange`, где нам нужно будет использовать идентификаторы внутри 
модели предметной области. В этом случае мы будем использовать различные 
реализации `Currency` на разных сервисах.

```go
// Объект-значение в веб-сервисе
type Currency struct {
    Code string
    HTML int
}

// Сущность в сервисе Payment
type Currency struct {
    ID uuid.UUID
    Code string
    HTML int
}
```

Шаблон, который мы будем использовать, Объект-значение или Сущность, зависит от 
только от того, что этот объект из себя представляет в Ограниченном Контексте.
Если это многократно используемый объект, независимо хранящийся в базе данных,
может изменяться и задействован во многих других объектах или связан с некоторой
внешней Сущностью и его необходимо изменять при изменении внешней Сущности, то
мы говорим о Сущности. Но если объект описывает какое-то значение, принадлежит
определенной Сущности, является простой копией, получаемой из внешнего сервиса,
или не должен существовать независимо в базе данных, тогда это Объект-значение.

## Явное описание

Самая полезная особенность Объекта-значения — это его явное описание. Его проще
понять в случаях, когда исходные типы из Golang (или любого другого языка 
программирования) не поддерживают конкретное поведение или поддерживаемое 
поведение не является интуитивно понятным. Мы можем работать с клиентами во 
многих проектах, и они должны удовлетворять некоторым бизнес-инвариантам, 
например, быть совершеннолетними или представлять какое-либо юридическое лицо.
В таких случаях допустимо определять более ясные типы, например, `Birthday` и
`LegalForm`.

```go
type Birthday time.Time

func (b Birthday) IsYoungerThen(other time.Time) bool {
    return time.Time(b).After(other)
}

func (b Birthday) IsAdult() bool {
    return time.Time(b).AddDate(18, 0, 0).Before(time.Now())
}

const (
    Freelancer = iota
    Partnership
    LLC
    Corporation
)

type LegalForm int

func (s LegalForm) IsIndividual() bool {
    return s == Freelancer
}

func (s LegalForm) HasLimitedResponsibility() bool {
    return s == LLC || s == Corporation
}
```

Иногда Объект-значение не нужно явно определять как часть какой-либо другой 
Сущности или Объекта-значения. Тем не менее, мы можем определить Объект-значение
в виде вспомогательного объекта, чтобы упростить его дальнейшее использование
в коде. Например, Клиент (`Customer`) может быть физлицом (`Person`) или 
компанией (`Company`). В зависимости от типа Клиента меняется логика в 
приложении. Одним из лучших решений будет преобразование клиентов, используя 
вспомогательные объекты, чтобы с ними было проще работать.

```go
type Customer struct {
    ID        uuid.UUID
    Name      string
    LegalForm LegalForm
    Date      time.Time
}

func (c Customer) ToPerson() Person {
    return Person{
        FullName: c.Name,
        Birthday: Birthday(c.Date),
    }
}

func (c Customer) ToCompany() Company {
    return Company{
        Name:         c.Name,
        CreationDate: c.Date,
    }
}

type Person struct {
    FullName string
    Birthday Birthday
}

type Company struct {
    Name         string
    CreationDate time.Time
}
```

Хотя вариант с преобразованием можно использовать в некоторых проектах, в 
большинстве случаев это означает, что мы должны добавить эти Объекты-значения
в нашу модель предметной области. Фактически, каждый раз, когда мы замечаем, что
какая-то конкретная группа полей постоянно взаимодействует друг с другом, но она
находится внутри какой-то более крупной группы, то это знак. Мы должны 
сгруппировать их в Объект-значение и использовать его таким же образом внутри
нашей большой группы (которая после этого уменьшается).

## Неизменяемость

Объекты-значения неизменяемы. Нет ни единой повода, причины или другого 
аргумента для изменения состояния Объекта-значения в течение его жизненного 
цикла. Иногда несколько объектов могут содержать один и тот же Объект-значение
(хотя это не идеальное решение). В таких случаях мы определенно не хотим, чтобы 
Объекты-значения изменялись где-либо. Итак, всякий раз, когда мы хотим изменить
внутреннее состояние объекта-значения или объединить несколько из них, нам всегда
нужно возвращать новый экземпляр с новым состоянием, как во фрагменте кода ниже.

```go
// Неправильно. Состояние изменяется внутри объекта-значения
func (m Money) AddAmount(amount float64) {
    m.Value += amount
}

// Правильно. Возвращаем новый объект-значение с новым состоянием
func (m Money) WithAmount(amount float64) Money {
    return Money{
        Value:    m.Value + amount,
        Currency: m.Currency,
    }
}

// Неправильно. Состояние изменяется внутри объекта-значения
func (m *Money) Deduct(other Money) {
    m.Value -= other.Value
}

// Правильно. Возвращаем новый объект-значение с новым состоянием
func (m Money) DeductedWith(other Money) Money {
    return Money{
        Value:    m.Value - other.Value,
        Currency: m.Currency,
    }
}

// Неправильно. Состояние изменяется внутри объекта-значения
func (c *Color) KeepOnlyGreen() {
    c.Red = 0
    c.Blue = 0
}

// Правильно. Возвращаем новый объект-значение с новым состоянием
func (c Color) WithOnlyGreen() Color {
    return Color{
        Red:   0,
        Green: c.Green,
        Blue:  0,
    }
}
```

Во всех примерах единственный правильный способ — всегда возвращать новые 
экземпляры и оставлять старые нетронутыми. Хорошей практикой в Golang является
всегда передавать в методы значения, а не ссылки на Объекты-значения, чтобы 
случайно не изменить внутреннее состояние.

```go
func (m Money) Deduct(other Money) (Money, error) {
    if !m.Currency.EqualTo(other.Currency) {
        return Money{}, errors.New("currencies must be identical")
    }
    
    if other.Value > m.Value {
        return Money{}, errors.New("there is not enough amount to deduct")
    }
    return Money{
        Value:    m.Value - other.Value,
        Currency: m.Currency,
    }, nil
}
```

Неизменяемость означает, что не нужно постоянно проверять правильные ли 
значения хранятся в его полях в течение всего жизненного цикла, а только при
создании, как это показано в приведённом выше примере. Когда мы хотим создать
новый Объект-значение, мы всегда должны осуществить валидацию и вернуть ошибки,
если бизнес-инварианты не выполняются. Создавать Объект-значение нужно только 
в том случае, если проверка прошла успешна. С этого момента больше валидировать
его не нужно.

## Наличие поведения
Объекты-значения позволяют задавать различные варианты поведения. Его основная
цель — предоставить доступный интерфейс. Наличие объекта-значения без методов
заставляет задуматься о целесообразности его существования. Если объект-значение
используется в каком-то конкретном месте кода, то он предоставляет доступ к 
огромному числу дополнительных бизнес-инвариантов, намного лучше описывающих 
решаемую нами проблему.

```go
func (c Color) ToBrighter() Color {
    return Color{
        Red:   byte(math.Min(255, float64(c.Red+10))),
        Green: byte(math.Min(255, float64(c.Green+10))),
        Blue:  byte(math.Min(255, float64(c.Blue+10))),
    }
}

func (c Color) ToDarker() Color {
    return Color{
        Red:   byte(math.Max(255, float64(c.Red-10))),
        Green: byte(math.Max(255, float64(c.Green-10))),
        Blue:  byte(math.Max(255, float64(c.Blue-10))),
    }
}

func (c Color) Combine(other Color) Color {
    return Color{
        Red:   byte(math.Min(255, float64(c.Red+other.Red))),
        Green: byte(math.Min(255, float64(c.Green+other.Green))),
        Blue:  byte(math.Min(255, float64(c.Blue+other.Blue))),
    }
}

func (c Color) IsRed() bool {
    return c.Red == 255 && c.Green == 0 && c.Blue == 0
}

func (c Color) IsYellow() bool {
    return c.Red == 255 && c.Green == 255 && c.Blue == 0
}

func (c Color) IsMagenta() bool {
    return c.Red == 255 && c.Green == 0 && c.Blue == 255
}

func (c Color) ToCSS() string {
    return fmt.Sprintf(`rgb(%d, %d, %d`, c.Red, c.Green, c.Blue)
}
```

Декомпозиция всей модели предметной области на небольшие части, такие как 
Объекты-значения (и Сущности), делает код понятным и приближённым к бизнес-логике
в реальном мире. Каждый Объект-значение может описывать некоторые небольшие 
компоненты и поддерживать различные модели поведения подобно обычным 
бизнес-процессам. В конце концов, это значительно упрощает весь процесс unit 
тестирования и помогает охватить все пограничные случаи.

## Заключение

В реальном мире мы постоянно сталкиваемся с различными характеристиками, 
качественными, количественными величинами. Поскольку программное обеспечение 
пытается решить проблемы, существующие в реальном мире, использование таких 
показателей неизбежно. В нашей бизнес-логике для задания таких величин могут 
использоваться объекты-значения, представленные в этой статье.

## Полезные ссылки на источники:

* [https://martinfowler.com/](https://martinfowler.com/)
* [https://www.domainlanguage.com/](https://www.domainlanguage.com/)