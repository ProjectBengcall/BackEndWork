package services

import (
	"bengcall/features/transaction/domain"
	"bengcall/utils/helper"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type transactionService struct {
	qry domain.Repository
}

func New(repo domain.Repository) domain.Service {
	return &transactionService{
		qry: repo,
	}
}

func (ss *transactionService) Transaction(newTrx domain.TransactionCore, newDtl []domain.DetailCore) (domain.TransactionDetail, error) {
	var invo int
	rand.Seed(time.Now().UnixNano())
	invo = rand.Intn(100000)
	newTrx.Invoice = invo
	start := newTrx.Schedule + "T08:00:00+07:00"
	end := newTrx.Schedule + "T09:00:00+07:00"

	res, err := ss.qry.Post(newTrx, newDtl)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return domain.TransactionDetail{}, errors.New("rejected from database")
		}
		return domain.TransactionDetail{}, errors.New("some problem on database")
	}

	str := strconv.Itoa(invo)
	//t := newTrx.Schedule.String()
	// p := newTrx.Schedule.Add(time.Hour * 1)
	// e := p.String()
	//fmt.Println(t)

	event := &calendar.Event{
		Summary:     "Bengcall Invoice Number: " + str,
		Location:    newTrx.Address,
		Description: "Contact Admin: 081234567890",
		Start: &calendar.EventDateTime{
			DateTime: start,
			TimeZone: "Asia/Jakarta",
		},
		End: &calendar.EventDateTime{
			DateTime: end,
			TimeZone: "Asia/Jakarta",
		},
		//Recurrence: []string{"RRULE:FREQ=DAILY;COUNT=2"},
		Attendees: []*calendar.EventAttendee{
			&calendar.EventAttendee{Email: res.Email},
		},
	}

	ctx := context.Background()

	// b, err := ioutil.ReadFile("./credentials.json")
	// if err != nil {
	// 	log.Fatalf("Unable to read client secret file: %v", err)
	// }

	client_id := os.Getenv("GOOGLE_CLIENT_ID")
	project := os.Getenv("GOOGLE_PROJECT_ID")
	secret := os.Getenv("GOOGLE_CLIENT_SECRET")
	b := `{"installed":{"client_id":"` + client_id + `","project_id":"` + project + `","auth_uri":"https://accounts.google.com/o/oauth2/auth","token_uri":"https://oauth2.googleapis.com/token","auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs","client_secret":"` + secret + `","redirect_uris":["http://localhost"]}}`
	bt := []byte(b)

	//fmt.Println(b)

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(bt, calendar.CalendarScope)
	if err != nil {
		log.Error(errors.New("Unable to parse client secret file to config"))
	}
	client := helper.GetClient(config)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Error(errors.New("Unable to retrieve Calendar client"))
	}

	calendarId := "primary"
	event, err = srv.Events.Insert(calendarId, event).Do()
	if err != nil {
		log.Error(errors.New("Unable to create event"))
	}
	fmt.Printf("Event created: %s\n", event.HtmlLink)

	return res, nil
}

func (ss *transactionService) Success(ID uint) error {
	err := ss.qry.PutScss(ID)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return errors.New("Rejected from Database")
		}
		return errors.New("Some Problem on Database")
	}

	return nil
}

func (ss *transactionService) Status(updateStts domain.TransactionCore, ID uint) (domain.TransactionCore, error) {
	res, err := ss.qry.PutStts(updateStts, ID)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return domain.TransactionCore{}, errors.New("Rejected from Database")
		}
		return domain.TransactionCore{}, errors.New("Some Problem on Database")
	}

	return res, nil
}

func (ss *transactionService) Comment(updateCmmt domain.TransactionCore, ID uint) (domain.TransactionCore, error) {
	res, err := ss.qry.PutCmmt(updateCmmt, ID)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return domain.TransactionCore{}, errors.New("Rejected from Database")
		} else if strings.Contains(err.Error(), "id") {
			return domain.TransactionCore{}, errors.New("There's no ID")
		}
		return domain.TransactionCore{}, errors.New("Some Problem on Database")
	}

	return res, nil
}

func (ss *transactionService) All() ([]domain.TransactionAll, error) {
	res, err := ss.qry.GetAll()
	if err != nil {
		log.Error(err.Error())
		if strings.Contains(err.Error(), "table") {
			return nil, errors.New("Database Error")
		} else if strings.Contains(err.Error(), "found") {
			return nil, errors.New("No Data")
		}
	}

	return res, nil
}

func (ss *transactionService) History(userID uint) ([]domain.TransactionHistory, error) {
	res, err := ss.qry.GetHistory(userID)
	if err != nil {
		log.Error(err.Error())
		if strings.Contains(err.Error(), "table") {
			return nil, errors.New("Database Error")
		} else if strings.Contains(err.Error(), "found") {
			return nil, errors.New("No Data")
		}
	}
	return res, nil
}

func (ss *transactionService) My(userID uint) (domain.TransactionHistory, error) {
	res, err := ss.qry.GetMy(userID)
	if err != nil {
		log.Error(err.Error())
		if strings.Contains(err.Error(), "table") {
			return domain.TransactionHistory{}, errors.New("Database Error")
		} else if strings.Contains(err.Error(), "found") {
			return domain.TransactionHistory{}, errors.New("No Data")
		}
	}
	return res, nil
}

func (ss *transactionService) Detail(ID uint) (domain.TransactionDetail, []domain.DetailCores, error) {
	res, dtl, err := ss.qry.GetDetail(ID)
	if err != nil {
		log.Error(err.Error())
		if strings.Contains(err.Error(), "table") {
			return domain.TransactionDetail{}, nil, errors.New("Database Error")
		} else if strings.Contains(err.Error(), "found") {
			return domain.TransactionDetail{}, nil, errors.New("No Data")
		}
	}

	return res, dtl, nil
}

func (ss *transactionService) Cancel(ID uint) error {
	err := ss.qry.Delete(ID)
	if err != nil {
		if strings.Contains(err.Error(), "table") {
			return errors.New("database error")
		} else if strings.Contains(err.Error(), "found") {
			return errors.New("no data")
		}
	}
	return nil
}
