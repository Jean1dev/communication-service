package events

// import (
// 	"communication-service/application"
// 	"communication-service/infra/sockets"
// 	"encoding/json"
// )

// func MyNotificationsEventHandler(event sockets.EventMessage, c *sockets.ClientSocket) error {
// 	results := application.GetMyNotifications("jean")
// 	data, err := json.Marshal(results)

// 	if err != nil {
// 		return err
// 	}
// 	c.Egress <- data
// 	return nil
// }
