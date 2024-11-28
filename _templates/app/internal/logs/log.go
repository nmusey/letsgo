package logs

type Log struct {
	Message   string `db:"message"`
	Timestamp string `db:"time"`
    Severity  string `db:"severity"`
}
