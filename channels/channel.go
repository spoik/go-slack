package channels

type Channel struct {
	Id int64 `bun:",pk,autoincrement"`
	Name string `bun:",notnull,unique"`
}
