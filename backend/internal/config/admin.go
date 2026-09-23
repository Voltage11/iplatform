package config

type AdminConfig struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

// newAdminCinfig Конфигурация первой учетной записи админа TODO сделать валидацию
func newAdminCinfig() AdminConfig {
	return AdminConfig{
		Email:     getEnv("ADMIN_EMAIL", ""),
		Password:  getEnv("ADMIN_PASSWORD", ""),
		FirstName: getEnv("ADMIN_FIRST_NAME", ""),
		LastName:  getEnv("ADMIN_LAST_NAME", ""),
	}
}
