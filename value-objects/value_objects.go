package value_objects

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"math"
	"time"
)

// Объект-значение в веб-сервисе
//type Currency struct {
//	Code string
//	HTML int
//}

// Сущность в сервисе Payment
type Currency struct {
	ID   uuid.UUID
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

type Coin struct {
	Value Money
	Color Color
}

type Colors []Color

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
//func (m Money) Deduct(other Money) {
//	m.Value -= other.Value
//}

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

func (m Money) Add(other Money) (Money, error) {
	if !m.Currency.EqualTo(other.Currency) {
		return Money{}, errors.New("currencies must be identical")
	}

	return Money{
		Value:    m.Value + other.Value,
		Currency: m.Currency,
	}, nil
}

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