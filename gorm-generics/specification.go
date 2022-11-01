package main

import (
	"fmt"
	"strings"
)

type Specification interface {
	GetQuery() string
	GetValues() []any
}

// joinSpecification - это действующая реализация интерфейса Specification
// Она используется для операторов AND и OR
type joinSpecification struct {
	specifications []Specification
	separator      string
}

// GetQuery объединяет все подзапросы
func (s joinSpecification) GetQuery() string {
	queries := make([]string, 0, len(s.specifications))

	for _, spec := range s.specifications {
		queries = append(queries, spec.GetQuery())
	}

	return strings.Join(queries, fmt.Sprintf(" %s ", s.separator))
}

// GetValues объединяет все подзначения
func (s joinSpecification) GetValues() []any {
	values := make([]any, 0)

	for _, spec := range s.specifications {
		values = append(values, spec.GetValues()...)
	}

	return values
}

// And передаёт AND оператор в виде Specification
func And(specifications ...Specification) Specification {
	return joinSpecification{
		specifications: specifications,
		separator:      "AND",
	}
}

// notSpecification отрицает под Specification
type notSpecification struct {
	Specification
}

// GetQuery отрицает подзапрос
func (s notSpecification) GetQuery() string {
	return fmt.Sprintf(" NOT (%s)", s.Specification.GetQuery())
}

// Not передаёт NOT оператор в виде Specification
func Not(specification Specification) Specification {
	return notSpecification{
		specification,
	}
}

// binaryOperatorSpecification определяет бинарный оператор как Specification
// Он используется для операторов =, >, <, >=, <=.
type binaryOperatorSpecification[T any] struct {
	field    string
	operator string
	value    T
}

// GetQuery создаёт запрос для бинарного оператора
func (s binaryOperatorSpecification[T]) GetQuery() string {
	return fmt.Sprintf("%s %s ?", s.field, s.operator)
}

// GetValues возвращает значение для бинарного оператора
func (s binaryOperatorSpecification[T]) GetValues() []any {
	return []any{s.value}
}

// Equal передаёт оператор равенства в виде Specification
func Equal[T any](field string, value T) Specification {
	return binaryOperatorSpecification[T]{
		field:    field,
		operator: "=",
		value:    value,
	}
}
