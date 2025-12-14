package postgres

import (
	"context"
)

func (p *Postgres) AddToWhitelist(ctx context.Context, cidr string) error {
	_, err := p.DB.ExecContext(ctx, `INSERT INTO whitelist (cidr) VALUES ($1) ON CONFLICT DO NOTHING`, cidr)
	return err
}

func (p *Postgres) AddToBlacklist(ctx context.Context, cidr string) error {
	_, err := p.DB.ExecContext(ctx, `INSERT INTO blacklist (cidr) VALUES ($1) ON CONFLICT DO NOTHING`, cidr)
	return err
}

func (p *Postgres) RemoveFromWhitelist(ctx context.Context, cidr string) error {
	_, err := p.DB.ExecContext(ctx, `DELETE FROM whitelist WHERE cidr = $1`, cidr)
	return err
}

func (p *Postgres) RemoveFromBlacklist(ctx context.Context, cidr string) error {
	_, err := p.DB.ExecContext(ctx, `DELETE FROM blacklist WHERE cidr = $1`, cidr)
	return err
}

func (p *Postgres) GetWhitelist(ctx context.Context) ([]string, error) {
	rows, err := p.DB.QueryContext(ctx, `SELECT cidr FROM whitelist`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []string
	for rows.Next() {
		var cidr string
		rows.Scan(&cidr)
		res = append(res, cidr)
	}
	return res, nil
}

func (p *Postgres) GetBlacklist(ctx context.Context) ([]string, error) {
	rows, err := p.DB.QueryContext(ctx, `SELECT cidr FROM blacklist`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []string
	for rows.Next() {
		var cidr string
		rows.Scan(&cidr)
		res = append(res, cidr)
	}
	return res, nil
}
