package main

import (
	"errors"
	"flag"
	"github.com/emersion/go-smtp"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

// The Backend implements SMTP server methods.
type Backend struct{}

func (bkd *Backend) NewSession(_ smtp.ConnectionState, _ string) (smtp.Session, error) {
	return &Session{}, nil
}

// A Session is returned after EHLO.
type Session struct {
	to string
}

func (s *Session) AuthPlain(username, password string) error {
	if username != "username" || password != "password" {
		return errors.New("Invalid username or password")
	}
	return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	log.Println("Mail from:", from)
	return nil
}

func (s *Session) Rcpt(to string) error {
	s.to = to
	log.Println("Rcpt to:", to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	if b, err := ioutil.ReadAll(r); err != nil {
		return err
	} else {
		_ = b
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		log.Println("save to path:", dir)
		if err != nil {
			return err
		}
		f, err := os.Create(path.Join(dir, s.to+".eml"))
		if err != nil {
			return err
		}
		defer func() {
			err := f.Close()
			if err != nil {
				return
			}
		}()

		_, err = f.WriteString(string(b))
		if err != nil {
			panic(err)
			return err
		}
	}
	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}

func main() {
	be := &Backend{}

	s := smtp.NewServer(be)

	domain := ""
	addr := ""

	flag.StringVar(&domain, "domain", "localhost", "Почтовый домен")
	flag.StringVar(&addr, "addr", ":1025", "Адрес")

	flag.Parse()

	s.Addr = addr
	s.Domain = domain
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	log.Println("Starting server at", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
