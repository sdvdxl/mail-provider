package http

import (
	"log"
	"net/http"
	"strings"

	"github.com/open-falcon/mail-provider/config"
	"github.com/toolkits/web/param"
	"github.com/open-falcon/mail-provider/mail"
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
		tosSlice := strings.Split(tos, ",")
		log.Println("prepare to send to: ", tos)
		log.Println(cfg.Smtp)

		err := mail.Send(cfg.Smtp.From, cfg.Smtp.Password, subject, content, cfg.Smtp.Addr, tosSlice...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, "success", http.StatusOK)
		}
	})

}
