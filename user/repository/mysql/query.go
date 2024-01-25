package mysql

const (
	GET_BY_EMAIL = `SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = ?`
	STORE        = `INSERT INTO users (name, email, password) VALUES (?, ?, ?)`
)
