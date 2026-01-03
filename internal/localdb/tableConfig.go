package localdb

var schemaStrList = []string{
	`CREATE TABLE IF NOT EXISTS user (
		id VARCHAR(36) NOT NULL,
		name VARCHAR(255) NOT NULL,
		PRIMARY KEY (id)
	)`,
	`CREATE TABLE IF NOT EXISTS schedulerJob (
		id VARCHAR(36) NOT NULL,
		name VARCHAR(255) NOT NULL,
		jobType INTEGER NOT NULL,
		year INTEGER,
		month INTEGER,
		day INTEGER,
		hour INTEGER,
		minute INTEGER,
		PRIMARY KEY (id)
	)`,
	`CREATE TABLE IF NOT EXISTS jobLog (
		id VARCHAR(36) NOT NULL,
		PRIMARY KEY (id)
	)`,
}
