package whatsapp

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

var (
	client *whatsmeow.Client
	qrChan chan string
	mu     sync.Mutex
)

func Init() error {
	mu.Lock()
	defer mu.Unlock()

	if client != nil {
		client.Disconnect()
	}

	dbLog := waLog.Stdout("Database", "WARN", true)
	container, err := sqlstore.New(context.Background(), "sqlite", "file:wa.db?_pragma=foreign_keys(1)", dbLog)
	if err != nil {
		return err
	}
	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		return err
	}
	clientLog := waLog.Stdout("Client", "WARN", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)

	qrChan = make(chan string, 1)

	client.AddEventHandler(eventHandler)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChanInt, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			return err
		}
		go func() {
			for evt := range qrChanInt {
				if evt.Event == "code" {
					select {
					case qrChan <- evt.Code:
					default:
						// If channel is full, discard old and push new
						<-qrChan
						qrChan <- evt.Code
					}
				} else {
					fmt.Println("WhatsApp Login event:", evt.Event)
				}
			}
		}()
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			return err
		}
	}
	return nil
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		_ = v
		// fmt.Println("Received a message!", v.Message.GetConversation())
	}
}

func GetQR() string {
	select {
	case qr := <-qrChan:
		// push it back so it can be fetched again if needed, but it changes every 20s
		// wait, actually if we push it back, we just keep returning the same one until replaced.
		// Let's just return it and the channel will get the new one when it updates.
		return qr
	case <-time.After(50 * time.Millisecond):
		return ""
	}
}

func IsConnected() bool {
	if client == nil {
		return false
	}
	return client.IsConnected() && client.IsLoggedIn()
}

func Logout() error {
	mu.Lock()
	defer mu.Unlock()

	if client == nil {
		return nil
	}
	client.Logout(context.Background())
	client.Disconnect()
	client = nil

	// Re-init so we get a new QR
	mu.Unlock() // unlock early so Init can lock
	err := Init()
	mu.Lock()
	return err
}

func SendMessage(phone, message string) error {
	if !IsConnected() {
		return fmt.Errorf("whatsapp not connected")
	}

	// Normalize phone number
	phone = strings.ReplaceAll(phone, "+", "")
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	if strings.HasPrefix(phone, "0") {
		phone = "62" + phone[1:]
	}

	jid := types.JID{
		User:   phone,
		Server: types.DefaultUserServer,
	}

	msg := &waE2E.Message{
		Conversation: proto.String(message),
	}
	
	_, err := client.SendMessage(context.Background(), jid, msg)
	return err
}
