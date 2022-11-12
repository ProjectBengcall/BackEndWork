package services

import (
	"bengcall/features/transaction/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
	"golang.org/x/oauth2"
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

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "./token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func (ss *transactionService) Transaction(newTrx domain.TransactionCore, newDtl []domain.DetailCore) (domain.TransactionDetail, error) {
	var invo int
	rand.Seed(time.Now().UnixNano())
	invo = rand.Intn(100000)
	newTrx.Invoice = invo
	str := strconv.Itoa(invo)

	event := &calendar.Event{
		Summary:     str,
		Location:    newTrx.Address,
		Description: newTrx.Phone,
		Start: &calendar.EventDateTime{
			DateTime: "2022-06-24T11:00:00+05:30",
			TimeZone: "Asia/Kolkata",
		},
		End: &calendar.EventDateTime{
			DateTime: "2022-06-24T11:15:00+05:30",
			TimeZone: "Asia/Kolkata",
		},
		//Recurrence: []string{"RRULE:FREQ=DAILY;COUNT=2"},
		Attendees: []*calendar.EventAttendee{
			&calendar.EventAttendee{Email: "gerdo.tewel@gmail.com"},
		},
	}

	ctx := context.Background()

	b, err := ioutil.ReadFile("./credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	calendarId := "primary"
	event, err = srv.Events.Insert(calendarId, event).Do()
	if err != nil {
		log.Fatalf("Unable to create event. %v\n", err)
	}
	fmt.Printf("Event created: %s\n", event.HtmlLink)

	res, err := ss.qry.Post(newTrx, newDtl)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return domain.TransactionDetail{}, errors.New("rejected from database")
		}
		return domain.TransactionDetail{}, errors.New("some problem on database")
	}

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
