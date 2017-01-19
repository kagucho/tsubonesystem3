package context

import (
	"github.com/kagucho/tsubonesystem3/db"
	"github.com/kagucho/tsubonesystem3/handler/apiv0/token/backend"
	"github.com/kagucho/tsubonesystem3/mail"
)

type Context struct {
	DB    db.DB
	Mail  mail.Mail
	Token backend.Backend
}
