package modules

type User struct {
	ID    int     `db:"id"`
	Name  string  `db:"name"`
	Email *string `db:"email"`
	Age   *int    `db:"age"`
	City  *string `db:"city"`
}
