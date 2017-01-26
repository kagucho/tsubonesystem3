package context

import (
	"github.com/kagucho/tsubonesystem3/backend/db"
	"github.com/kagucho/tsubonesystem3/backend/handler/apiv0/token/backend"
	"github.com/kagucho/tsubonesystem3/backend/mail"
)

type Context struct {
	DB    db.DB
	Mail  mail.Mail
	Token backend.Backend
}
