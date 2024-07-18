package entity

type DataSource struct {
	Id       int     `sql:"dso_id"`
	Name     *string `sql:"dso_name"`
	Username *string `sql:"dso_username"`
	Password *string `sql:"dso_password"`
	Token    *string `sql:"dso_token"`
	Struct   *string `sql:"dso_struct"`
	Url      *string `sql:"dso_url"`
	Port     *string `sql:"dso_port"`
	Type     int     `sql:"dso_type"`
}
