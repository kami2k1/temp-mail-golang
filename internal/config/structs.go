package config

type Dataconfig struct {
    HOST      string
    PORT      string
    APP_ENV   string
    IMAP_HOST string
    IMAP_PORT string
    STMP_USER string
    STMP_PASS string
    JWT_SECRET string
}

var (
	Config *Dataconfig
)
