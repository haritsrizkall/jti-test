package mysql

const (
	GET_ALL   = `SELECT id, number, provider FROM phones`
	GET_BY_ID = `SELECT id, number, provider FROM phones WHERE id = ?`
	UPDATE    = `UPDATE phones SET number = ?, provider = ? WHERE id = ?`
	STORE     = `INSERT INTO phones (number, provider) VALUES (?, ?)`
	DELETE    = `DELETE FROM phones WHERE id = ?`
)
