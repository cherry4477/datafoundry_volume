package message

import (
	"encoding/json"
	"errors"
	//"time"
	
	"github.com/asiainfoLDP/datahub_commons/mq"
	"github.com/asiainfoLDP/datahub_commons/log"
)

// todo: maybe it is better to package the message push and save time in the header
// the following time is the event happen time

type Message struct {
	Type     string      `json:"type,omitempty"`
	Receiver string      `json:"receiver,omitempty"`
	Sender   string      `json:"sender,omitempty"`
	Level    int         `json:"level,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	//Time     time.Time   `json:"time,omitempty"` 
		// this is the time to insert message into table.
		// if there are other times related to this messsage,
		// pls put them in Data field.
}

func PushMessageToQueue(queue mq.MessageQueue, mqTopic string, key []byte, message *Message) error {
	json_bytes, err := json.Marshal(message)
	if err != nil {
		log.DefaultLogger().Warningf("PushMessageToQueue Marshal error: %s.\nMessage=%s", err.Error(), string(json_bytes))
		return err
	}
	
	//err = queue.SendAsyncMessage(mqTopic, key, json_bytes)
	partition, offset, err := queue.SendSyncMessage(mqTopic, key, json_bytes)
	if err != nil {
		go func() {
			log.DefaultLogger().Warningf("PushMessageToQueue SendAsyncMessage error: %s.\nMessage=%s", err.Error(), string(json_bytes))
		}()
	} else {
		go func() {
			//log.DefaultLogger().Debugf("PushMessageToQueue succeeded. key=%s, Message=%s", key, string(json_bytes))
			log.DefaultLogger().Debugf("PushMessageToQueue succeeded. partition=%d, offset=%d, key=%s, Message=%s", partition, offset, key, string(json_bytes))
		}()
	}
	
	return err
}

func ParseJsonMessage(msgData []byte) (*Message, error) {
	if msgData == nil {
		return nil, errors.New("message data can't be nil")
	}

	msg := &Message{}
	err := json.Unmarshal(msgData, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil	
}

//=====================================
// 
//=====================================

type Email struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Content string `json:"content"`
	IsHTML  bool   `json:"ishtml"`
}

func ParseJsonEmail(msgData []byte) (*Email, error) {
	if msgData == nil {
		return nil, errors.New("message data can't be nil")
	}

	msg := &Email{}
	err := json.Unmarshal(msgData, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func PushMailToQueue(queue mq.MessageQueue, mqTopic string, key []byte, mail *Email) error {
	json_bytes, err := json.Marshal(mail)
	if err != nil {
		log.DefaultLogger().Warningf("PushMailToQueue Marshal error: %s.\nEmail=%s", err.Error(), string(json_bytes))
		return err
	}
	
	_, _, err = queue.SendSyncMessage(mqTopic, key, json_bytes)
	if err != nil {
		go func() {
			log.DefaultLogger().Warningf("PushMailToQueue SendAsyncMessage error: %s.\nEmail=%s", err.Error(), string(json_bytes))
		}()
	} else {
		go func() {
			//log.DefaultLogger().Debugf("PushMessageToQueue succeeded. key=%s, Message=%s", key, string(json_bytes))
			log.DefaultLogger().Debugf("PushMailToQueue succeeded. Email=%s", string(json_bytes))
		}()
	}
	
	return err
}


