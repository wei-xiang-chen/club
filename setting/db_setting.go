package setting

type DBSetting struct {
	DBtype       string
	Username     string
	Password     string
	Host         string
	Part         int
	DBName       string
	TablePrefix  string
	Charset      string
	ParseTime    bool
	MaxIdleConns int
	MaxOpenConns int
}

var (
	DatabaseSetting *DBSetting
)

func (s *Setting) ReadDBSetting() error {
	err := s.vp.UnmarshalKey("Database", &DatabaseSetting)
	if err != nil {
		return err
	}

	return nil
}
