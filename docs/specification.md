# DDD на практике в Golang: Спецификация

![intro](images/specification/intro.jpeg)
*Фото [Esteban Castle](https://unsplash.com/@estebancastle) из [Unsplash](https://unsplash.com/)*

Существует не так много алгоритмических структур, которые я реализую с 
удовольствием. Первой такой стало упрощенное ORM в Go, когда у нас его не было.

С другой стороны я много лет использовал ORM. В какой-то момент, когда вы 
зависите от ORM, возникает неизбежная необходимость применения 
[QueryBuilder](https://www.doctrine-project.org/projects/doctrine-dbal/en/latest/reference/query-builder.html).
Вот где используется шаблон Спецификация (`Specification`).

Сложно найти какой-либо шаблон, который мы использовали бы так же часто, как и 
Спецификацию, но не произнося его название. Я думаю, что сложнее только написать
приложение без использования этого шаблона.

Шаблон Спецификация имеет множество различных применений. Мы можем использовать
его для запросов, создания или валидации. Мы можем написать универсальный код,
выполняющий всю эту работу, или написать свою реализацию для каждого случая.

> Другие статьи из DDD цикла:
> 1. [DDD на практике в Golang: Объект-значение](https://levelup.gitconnected.com/practical-ddd-in-golang-value-object-4fc97bcad70)
> 2. [DDD на практике в Golang: Сущности](https://levelup.gitconnected.com/practical-ddd-in-golang-entity-40d32bdad2a3)
> 3. [DDD на практике в Golang: Сервисы предметной области](https://levelup.gitconnected.com/practical-ddd-in-golang-domain-service-4418a1650274)
> 4. [DDD на практике в Golang: Событие предметной области](https://levelup.gitconnected.com/practical-ddd-in-golang-domain-event-de02ad492989)
> 5. [DDD на практике в Golang: Модуль](https://levelup.gitconnected.com/practical-ddd-in-golang-module-51edf4c319ec)
> 6. [DDD на практике в Golang: Агрегат](https://levelup.gitconnected.com/practical-ddd-in-golang-aggregate-de13f561e629)
> 7. [DDD на практике в Golang: Фабрика](https://levelup.gitconnected.com/practical-ddd-in-golang-factory-5ba135df6362)
> 8. [DDD на практике в Golang: Репозиторий](https://levelup.gitconnected.com/practical-ddd-in-golang-repository-d308c9d79ba7)

## Для валидации

Первый вариант использования шаблона Спецификация — это валидация. В первую
очередь мы проверяем данные в формах, но это происходит на уровне представления.
Иногда мы выполняем её во время создания, например, для Объектов-значений.

На уровне предметной области мы можем использовать Спецификации для проверки 
состояний Сущности и фильтрации их из коллекции. Итак, валидация на уровне 
предметной области уже имеет более широкое применение, чем только проверка 
пользовательского ввода.

```go
type MaterialType = string

const Plastic = "plastic"

type Product struct {
    ID            uuid.UUID
    Material      MaterialType
    IsDeliverable bool
    Quantity      int
}

type ProductSpecification interface {
    IsValid(product Product) bool
}

type AndSpecification struct {
    specifications []ProductSpecification
}

func NewAndSpecification(specifications ...ProductSpecification) ProductSpecification {
    return AndSpecification{
        specifications: specifications,
    }
}

func (s AndSpecification) IsValid(product Product) bool {
    for _, specification := range s.specifications {
        if !specification.IsValid(product) {
            return false
        }
    }

    return true
}

type HasAtLeast struct {
    pieces int
}

func NewHasAtLeast(pieces int) ProductSpecification {
    return HasAtLeast{
        pieces: pieces,
    }
}

func (h HasAtLeast) IsValid(product Product) bool {
    return product.Quantity >= h.pieces
}

func IsPlastic(product Product) bool {
    return product.Material == Plastic
}

func IsDeliverable(product Product) bool {
    return product.IsDeliverable
}

type FunctionSpecification func(product Product) bool

func (fs FunctionSpecification) IsValid(product Product) bool {
    return fs(product)
}

func main() {
    spec := model.NewAndSpecification(
        model.NewHasAtLeast(10),
        model.FunctionSpecification(model.IsPlastic),
        model.FunctionSpecification(model.IsDeliverable),
    )

    fmt.Println(spec.IsValid(model.Product{}))
    // выводит: false

    fmt.Println(spec.IsValid(model.Product{
        Material:      model.Plastic,
        IsDeliverable: true,
        Quantity:      50,
    }))
    // выводит: true
}
```
*Использование Спецификации для валидации данных*

В вышеприведенном примере задан интерфейс `ProductSpecification`. Он определяет
только один метод IsValid, который ожидает экземпляры `Product` и в результате
возвращает логическое значение, если `Product` соответствует правилам проверки.

Простая реализация этого интерфейса - `HasAtLeast`, который проверяет минимальное
количество продукта. Более интересными валидаторами являются две функции: 
`IsPlastic` и `IsDeliverable`.

Мы можем обернуть эти функции особым типом `FunctionSpecification`. Этот тип
использует функцию с такой же сигнатурой как у двух упомянутых выше. Кроме того,
он предоставляет методы, соответствующие интерфейсу `ProductSpecification`.

Этот пример показывает особенность Go, где мы можем определить функцию как
тип и добавить к нему метод, чтобы он мог неявно реализовать некоторый 
интерфейс. В нашем случае создаётся метод `IsValid`, который выполняет 
встроенную функцию.

Кроме того, существует также одна не похожая на другие спецификация 
`AndSpecification`. Такая структура позволяет нам использовать объект, который 
реализует интерфейс `ProductSpecification` и объединяет все входящие в неё 
Спецификации, используя логическое "И".

```go
type OrSpecification struct {
    specifications []ProductSpecification
}

func NewOrSpecification(specifications ...ProductSpecification) ProductSpecification {
    return OrSpecification{
        specifications: specifications,
    }
}

func (s OrSpecification) IsValid(product Product) bool {
    for _, specification := range s.specifications {
        if specification.IsValid(product) {
            return true
        }
    }

    return false
}

type NotSpecification struct {
    specification ProductSpecification
}

func NewNotSpecification(specification ProductSpecification) ProductSpecification {
    return NotSpecification{
        specification: specification,
    }
}

func (s NotSpecification) IsValid(product Product) bool {
    return !s.specification.IsValid(product)
}
```
*Дополнительные Спецификации*

В вышеприведенном фрагменте кода описаны две дополнительные Спецификации.
Одна из них `OrSpecification`. Она, как и `AndSpecification`, объединяет все 
входящие в неё Спецификации. Просто в данном случае используется логическое 
"ИЛИ" вместо "И".

Последняя - `NotSpecification`, логически инвертирует результат переданной
Спецификации. `NotSpecification` можно было бы задать с помощью 
`FunctionSpecification`, но я не хотел её слишком усложнять.

## Для запросов

Я уже упоминал в этой статье о применении шаблона Спецификация как части ORM. 
Во многих случаях вам не нужно будет реализовывать Спецификации для этого 
варианта использования, по крайней мере, если вы применяете какую-либо ORM.

Отличную реализацию Спецификации в виде предикатов я нашёл в библиотеке [Ent](https://entgo.io/)
от Facebook. С того момента я не писал спецификации для запросов.

Тем не менее, когда вы обнаружите, что ваш запрос для Репозитория на уровне
предметной области может быть слишком сложным, вам понадобятся дополнительные 
способы фильтрации желаемых объектов. Реализация может выглядеть как показано в 
примере ниже.

```go
type ProductSpecification interface {
    Query() string
    Value() []interface{}
}

type AndSpecification struct {
    specifications []ProductSpecification
}

func NewAndSpecification(specifications ...ProductSpecification) ProductSpecification {
    return AndSpecification{
        specifications: specifications,
    }
}

func (s AndSpecification) Query() string {
    var queries []string
    for _, specification := range s.specifications {
        queries = append(queries, specification.Query())
    }

    query := strings.Join(queries, " AND ")

    return fmt.Sprintf("(%s)", query)
}

func (s AndSpecification) Value() []interface{} {
    var values []interface{}
    for _, specification := range s.specifications {
        values = append(values, specification.Value()...)
    }
    return values
}

type OrSpecification struct {
    specifications []ProductSpecification
}

func NewOrSpecification(specifications ...ProductSpecification) ProductSpecification {
    return OrSpecification{
        specifications: specifications,
    }
}

func (s OrSpecification) Query() string {
    var queries []string
    for _, specification := range s.specifications {
        queries = append(queries, specification.Query())
    }

    query := strings.Join(queries, " OR ")

    return fmt.Sprintf("(%s)", query)
}

func (s OrSpecification) Value() []interface{} {
    var values []interface{}
    for _, specification := range s.specifications {
        values = append(values, specification.Value()...)
    }
    return values
}

type HasAtLeast struct {
    pieces int
}

func NewHasAtLeast(pieces int) ProductSpecification {
    return HasAtLeast{
        pieces: pieces,
    }
}

func (h HasAtLeast) Query() string  {
    return "quantity >= ?"
}

func (h HasAtLeast) Value() []interface{} {
    return []interface{}{h.pieces}
}

func IsPlastic() string {
    return "material = 'plastic'"
}

func IsDeliverable() string {
    return "deliverable = 1"
}

type FunctionSpecification func() string

func (fs FunctionSpecification) Query() string {
    return fs()
}

func (fs FunctionSpecification) Value() []interface{} {
    return nil
}

func main() {

    spec := infrastructure.NewOrSpecification(
        infrastructure.NewAndSpecification(
            infrastructure.NewHasAtLeast(10),
            infrastructure.FunctionSpecification(infrastructure.IsPlastic),
            infrastructure.FunctionSpecification(infrastructure.IsDeliverable),
        ),
        infrastructure.NewAndSpecification(
            infrastructure.NewHasAtLeast(100),
            infrastructure.FunctionSpecification(infrastructure.IsPlastic),
        ),
    )
    
    fmt.Println(spec.Query())
    // выводит: ((quantity >= ? AND material = 'plastic' AND deliverable = 1) OR (quantity >= ? AND material = 'plastic'))
    
    fmt.Println(spec.Value())
    // выводит: [10 100]
}
```
*Пример выполнения запроса*

В новой реализации интерфейс `ProductSpecification` предоставляет два метода: 
`Query` и `Values`. Мы используем из для получения строки запроса для конкретной
Спецификации и значений, которые она может содержать.

Опять же мы видим дополнительные спецификации, `AndSpecification` и 
`OrSpecification`. В этом случае они объединяют все входящие запросы в 
зависимости от оператора ("AND" или "OR") и все значения.

Наличие такой Спецификации на уровне домена вызывает вопросы. Как видите из 
выходных данных Спецификации предоставляют синтаксис, схожий с SQL, в котором
слишком много технических деталей.

В этом случае решением, вероятно, будет определить интерфейсы на уровне 
предметной области и фактические реализации на инфраструктурном уровне.

Или модифицировать код так, чтобы спецификация содержала информацию об имени поля,
операции и значении. Затем мы создадим некий сопоставитель на инфраструктурном 
уровне, который сможет преобразовать такую Спецификацию в SQL запрос.

## Для создания

Один из простейших вариантов использования Спецификации — создание сложного 
объекта, значения которого каждый раз сильно отличаются. В таких случаях мы
можем комбинировать его с шаблоном Фабрика (`Factory`) или использовать внутри Сервиса
предметной области (`Domain Service`).

```go
type ProductSpecification interface {
    Create(product model.Product) model.Product
}

type AndSpecification struct {
    specifications []ProductSpecification
}

func NewAndSpecification(specifications ...ProductSpecification) ProductSpecification {
    return AndSpecification{
        specifications: specifications,
    }
}

func (s AndSpecification) Create(product model.Product) model.Product {
    for _, specification := range s.specifications {
        product = specification.Create(product)
    }
    return product
}

type HasAtLeast struct {
    pieces int
}

func NewHasAtLeast(pieces int) ProductSpecification {
    return HasAtLeast{
        pieces: pieces,
    }
}

func (h HasAtLeast) Create(product model.Product) model.Product {
    product.Quantity = h.pieces
    return product
}

func IsPlastic(product model.Product) model.Product {
    product.Material = model.Plastic
    return product
}

func IsDeliverable(product model.Product) model.Product {
    product.IsDeliverable = true
    return product
}

type FunctionSpecification func(product model.Product) model.Product

func (fs FunctionSpecification) Create(product model.Product) model.Product {
    return fs(product)
}

func main() {
    spec := create.NewAndSpecification(
        create.NewHasAtLeast(10),
        create.FunctionSpecification(create.IsPlastic),
        create.FunctionSpecification(create.IsDeliverable),
    )
    
    fmt.Printf("%+v", spec.Create(model.Product{
        ID: uuid.New(),
    }))
    // выводит: {ID:befaf2b9-73cd-44cf-95f1-5fba087e46d9 Material:plastic IsDeliverable:true Quantity:10}
}
```
*Пример использования для создания объектов*

В этом примере показан третий вариант использования Спецификации. В этом случае
`ProductSpecification` поддерживает единственный метод, `Create`, который ожидает
`Product`, модифицирует его и возвращает обратно.

Опять, `AndSpecification` позволяет применить изменения, определённые в 
нескольких спецификациях, но здесь нет `OrSpecification`. Мне не удалось найти
реальный вариант его использования или алгоритм, когда бы он понадобился, при 
создании объекта.

Даже хотя здесь он не приводится мы можем создать `NotSpecification`, который
будет работать с определенными типами данных, например, логическими. Тем не 
менее, для данного примера, я не придумал как его можно было бы использовать.

## Заключение

Спецификация — это шаблон, который мы используем везде, во многих разных 
случаях. В настоящее время не просто обеспечить валидацию на уровне предметной
области без использования спецификации.

Спецификацию мы также можем использовать для запросов объектов из 
соответствующего хранилища. Сегодня они являются частью ORM. Третий вариант 
использования — создание сложных экземпляров, где мы можем комбинировать его
с шаблоном Фабрика (`Factory`).

> Другие статьи из DDD цикла:
> 1. [DDD на практике в Golang: Объект-значение](https://levelup.gitconnected.com/practical-ddd-in-golang-value-object-4fc97bcad70)
> 2. [DDD на практике в Golang: Сущности](https://levelup.gitconnected.com/practical-ddd-in-golang-entity-40d32bdad2a3)
> 3. [DDD на практике в Golang: Сервисы предметной области](https://levelup.gitconnected.com/practical-ddd-in-golang-domain-service-4418a1650274)
> 4. [DDD на практике в Golang: Событие предметной области](https://levelup.gitconnected.com/practical-ddd-in-golang-domain-event-de02ad492989)
> 5. [DDD на практике в Golang: Модуль](https://levelup.gitconnected.com/practical-ddd-in-golang-module-51edf4c319ec)
> 6. [DDD на практике в Golang: Агрегат](https://levelup.gitconnected.com/practical-ddd-in-golang-aggregate-de13f561e629)
> 7. [DDD на практике в Golang: Фабрика](https://levelup.gitconnected.com/practical-ddd-in-golang-factory-5ba135df6362)
> 8. [DDD на практике в Golang: Репозиторий](https://levelup.gitconnected.com/practical-ddd-in-golang-repository-d308c9d79ba7)

## Полезные ссылки на источники:

* [https://martinfowler.com/](https://martinfowler.com/)
* [https://www.domainlanguage.com/](https://www.domainlanguage.com/)