package http

import (
	"log"
	"net/http"
	"strings"

	"github.com/open-falcon/mail-provider/config"
	"github.com/open-falcon/mail-provider/mail"
	"github.com/toolkits/smtp"
	"github.com/toolkits/web/param"
)

func configProcRoutes() {

	http.HandleFunc("/sender/mail", func(w http.ResponseWriter, r *http.Request) {
		cfg := config.Config()
		token := param.String(r, "token", "")
		if cfg.Http.Token != token {
			http.Error(w, "no privilege", http.StatusForbidden)
			return
		}

		tos := param.MustString(r, "tos")
		subject := param.MustString(r, "subject")
		content := param.MustString(r, "content")
		log.Println("prepare to send to: ", tos)

		var err error
		if cfg.Smtp.SSL {
			tosSlice := strings.Split(tos, ",")
			err = mail.Send(cfg.Smtp.Username, cfg.Smtp.Password, cfg.Smtp.From, subject, content, cfg.Smtp.Addr, tosSlice...)
		} else {
			tos = strings.Replace(tos, ",", ";", -1)
			s := smtp.New(cfg.Smtp.Addr, cfg.Smtp.Username, cfg.Smtp.Password)
			err = s.SendMail(cfg.Smtp.From, tos, subject, content)
		}

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			log.Println("send success")
			http.Error(w, "success", http.StatusOK)
		}

	})

}
