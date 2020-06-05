package person

type Male interface {
	Get(id int64) error
	Put(id int64) error
}
