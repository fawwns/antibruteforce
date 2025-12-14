package postgres

func (p *Postgres) Migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS whitelist (
    id SERIAL PRIMARY KEY,
    cidr CIDR NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS blacklist (
    id SERIAL PRIMARY KEY,
    cidr CIDR NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);`,
	}

	for _, q := range queries {
		_, err := p.DB.Exec(q)
		if err != nil {
			return err
		}
	}

	return nil
}
