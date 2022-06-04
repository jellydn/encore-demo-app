package shorten_url

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"encore.dev/cron"
	"encore.dev/storage/sqldb"
)

type URL struct {
	ID    string // short-form URL id
	URL   string // complete URL, in long form
	Owner string // owner by
}

type ShortenParams struct {
	URL   string // the URL to shorten
	Owner string // Owner by, default is "anonymous"
}

// Shorten shortens a URL.
//encore:api public method=POST path=/url
func Shorten(ctx context.Context, p *ShortenParams) (*URL, error) {
	id, err := generateID()
	if err != nil {
		return nil, err
	} else if err := insert(ctx, id, p.URL, p.Owner); err != nil {
		return nil, err
	}
	return &URL{ID: id, URL: p.URL, Owner: p.Owner}, nil
}

//encore:api public method=GET path=/url/:id
func Get(ctx context.Context, id string) (*URL, error) {
	u := &URL{ID: id}
	err := sqldb.QueryRow(ctx, `
	    SELECT original_url, owner_by FROM url
	    WHERE id = $1
	`, id).Scan(&u.URL, &u.Owner)
	return u, err
}

// Delete guest links every day
var _ = cron.NewJob("delete-guest-link-every-day", cron.JobConfig{
	Title:    "Delete all anonymous links older every day",
	Every:    24 * cron.Hour,
	Endpoint: DeleteGuestLinks,
})

//encore:api private
func DeleteGuestLinks(ctx context.Context) error {
	_, err := sqldb.Exec(ctx, `
	    DELETE FROM url WHERE owner_by = $1
	`, "anonymous")
	return err
}

// generateID generates a random short ID.
func generateID() (string, error) {
	var data [6]byte // 6 bytes of entropy
	if _, err := rand.Read(data[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data[:]), nil
}

// insert inserts a URL into the database.
func insert(ctx context.Context, id, url string, owner string) error {
	_, err := sqldb.Exec(ctx, `
	    INSERT INTO url (id, original_url, owner_by)
	    VALUES ($1, $2, $3)
	`, id, url, owner)
	return err
}
