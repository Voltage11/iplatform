package filterbool

type Filter string

const (
	FilterTrue  Filter = "true"
	FilterFalse Filter = "false"
	FilterAll   Filter = "all"
)

func (f Filter) IsTrue() bool  { return f == FilterTrue }
func (f Filter) IsFalse() bool { return f == FilterFalse }

// IsAll возвращает true для явного "all" и для пустого значения
func (f Filter) IsAll() bool { return f == FilterAll || f == "" }

func (f Filter) String() string { return string(f) }

// Bool возвращает указатель на значение фильтра
// Для FilterAll и пустого значения возвращает nil — значит, фильтр не задан.
func (f Filter) Bool() *bool {
	switch f {
	case FilterTrue:
		b := true
		return &b
	case FilterFalse:
		b := false
		return &b
	default:
		return nil
	}
}

// NewFilterBool парсит строку в Filter
// Неизвестные значения = FilterAll
func NewFilterBool(s string) Filter {
	switch Filter(s) {
	case FilterTrue:
		return FilterTrue
	case FilterFalse:
		return FilterFalse
	default:
		return FilterAll
	}
}
