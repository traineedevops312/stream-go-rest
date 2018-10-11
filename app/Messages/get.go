package Messages

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Stream-golang/app/Cassandra"
	"github.com/Stream-golang/app/Stream"
	"github.com/Stream-golang/app/Users"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

// Get -- handles GET request to /messages/ to fetch all messages
// params:
// w - response writer for building JSON payload response
// r - request reader to fetch form data or url params (unused here)
func Get(w http.ResponseWriter, r *http.Request) {
	var messageList []Message
	var enrichedMessages []Message
	var userList []gocql.UUID
	var err error
	m := map[string]interface{}{}

	globalMessages, err := Stream.Client.FlatFeed("messages", "global")
	// fetch from Stream
	if err == nil {
		activities, err := globalMessages.Activities(nil)
		if err == nil {
			fmt.Println("Fetching activities from Stream")
			for _, activity := range activities.Activities {
				fmt.Println(activity)
				userID, _ := gocql.ParseUUID(activity.Actor)
				messageID, _ := gocql.ParseUUID(activity.Object)
				messageList = append(messageList, Message{
					ID:      messageID,
					UserID:  userID,
					Message: activity.MetaData["message"],
				})
				userList = append(userList, userID)
			}
		}
	}
	// if Stream fails, pull from database instead
	if err != nil {
		fmt.Println("Fetching activities from Database")
		query := "SELECT id,userID,message FROM messages"
		iterable := Cassandra.Session.Query(query).Iter()
		for iterable.MapScan(m) {
			userID := m["userid"].(gocql.UUID)
			messageList = append(messageList, Message{
				ID: m["id"].(gocql.UUID),
				//UserID:  userID,
				Message: m["message"].(string),
			})
			userList = append(userList, userID)
			m = map[string]interface{}{}
		}
	}

	names := Users.Enrich(userList)
	for _, message := range messageList {
		message.UserFullName = names[message.UserID.String()]
		enrichedMessages = append(enrichedMessages, message)
	}
	fmt.Println("message list after enrichment", enrichedMessages)

	json.NewEncoder(w).Encode(AllMessagesResponse{Messages: enrichedMessages})
}

// GetOne -- handles GET request to /messages/{message_uuid} to fetch one message
// params:
// w - response writer for building JSON payload response
// r - request reader to fetch form data or url params
func GetOne(w http.ResponseWriter, r *http.Request) {
	var message Message
	var errs []string
	var found bool = false

	vars := mux.Vars(r)
	id := vars["message_uuid"]

	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		m := map[string]interface{}{}
		query := "SELECT id,userID,message FROM messages WHERE id=? LIMIT 1 ALLOW FILTERING"
		iterable := Cassandra.Session.Query(query, uuid).Consistency(gocql.One).Iter()

		for iterable.MapScan(m) {
			found = true
			userID := m["userid"].(gocql.UUID)
			names := Users.Enrich([]gocql.UUID{userID})
			fmt.Println("names", names)
			message = Message{
				ID:           userID,
				UserID:       m["userid"].(gocql.UUID),
				UserFullName: names[userID.String()],
				Message:      m["message"].(string),
			}
		}
		if !found {
			errs = append(errs, "Message not found")
		}
	}

	if found {
		json.NewEncoder(w).Encode(GetMessageResponse{Message: message})
	} else {
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
	}
}
