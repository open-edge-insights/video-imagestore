/*
Copyright (c) 2019 Intel Corporation.

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package submanager

import (
	eismsgbus "EISMessageBus/eismsgbus"
	common "IEdgeInsights/ImageStore/common"
	"errors"
	"strings"
	"github.com/golang/glog"
)

// SubManager - SubManager of type struct
type SubManager struct {
	subscribers map[string]*eismsgbus.Subscriber
	clientMap   map[string]*eismsgbus.MsgbusClient
	subConfig   map[string]interface{}
	writers     map[string]common.Writer
}
// NewSubManager - function to initialize a new SubManager
func NewSubManager()(*SubManager){
	var subMgr SubManager
	subMgr.init()
	return &subMgr
}

func (subMgr *SubManager) init() {
	subMgr.writers = make(map[string]common.Writer)
	subMgr.subscribers = make(map[string]*eismsgbus.Subscriber)
	subMgr.clientMap = make(map[string]*eismsgbus.MsgbusClient)
}

func (subMgr *SubManager) close() {
	glog.Infof("-- Closing message bus context-- \n")
	for _, client := range subMgr.clientMap {
		client.Close()
	}
}

// RegWriterInterface - RegWriterInterface function
func (subMgr *SubManager) RegWriterInterface(name string, writer common.Writer) {
	subMgr.writers[name] = writer
}

// RegSubscriberList - RegSubscriberList function
func (subMgr *SubManager) RegSubscriberList(subConfig map[string]interface{}) {
	subMgr.subConfig = subConfig
}

// StartAllSubscribers - function to create subscription object for all the topics
// in topics array
func (subMgr *SubManager) StartAllSubscribers(topics []string) error {

	glog.Infof("-- subscribe to topics : %v\n", topics)
	for _, topic := range topics {
		infoMap, ok := subMgr.subConfig[topic].(map[string]interface{})
		if ok == false {
			errorMessage := "Error in getting the topic info of " + topic
			return errors.New(errorMessage)
		}

		glog.Infof("-- Info map for topic %v : %v -- \n", topic, infoMap)
		client, err := eismsgbus.NewMsgbusClient(infoMap)
		if err != nil {
			glog.Infof("-- Error initializing message bus context: %v\n", err)
			errorMessage := "-- Error initializing message bus context: " + err.Error()
			return errors.New(errorMessage)
		}

		subMgr.clientMap[topic] = client
		subTopics := strings.Split(topic, "/")
		subscriber, err := client.NewSubscriber(subTopics[1])
		if err != nil {
			glog.Infof("-- Error subscribing to topic: %v\n", err)
			return err
		}

		subMgr.subscribers[topic] = subscriber
	}

	return nil
}

// ReceiveFromAll - function to start new go routine which receives a frame from the given subscription
// topic and writes it to a storage
func (subMgr *SubManager) ReceiveFromAll() {
	for topicName, subscriber := range subMgr.subscribers {
		go Receive(topicName, subMgr.writers[topicName], subscriber)
	}
}

// Receive - function to receive image for given topic name and put it into storage
func Receive(topicName string, writer common.Writer, subscriber *eismsgbus.Subscriber) {

	for {
		select {
		case msg := <-subscriber.MessageChannel:
			glog.Infof("\n-- Received Message: %v\n", msg.Data)
			imgHandle, ok := msg.Data[common.ImageHandle].(string)
			if ok == false {
				errMessage := "Missing image handle for topic " + topicName
				glog.Infof(errMessage)
				continue
			}

			if msg.Blob != nil {
				_, err := writer.Store(msg.Blob, imgHandle)

				if err != nil {
					errMessage := "Error In storing the image %s from topic %s & Error %s"
					glog.Errorf(errMessage, msg.Data[common.ImageHandle], topicName, err)
				} else {
					glog.Infof("Image with handle %s stored successfully", imgHandle)
				}
			} else {
				errMessage := "Empty image for handle %s from topic %s"
				glog.Errorf(errMessage, msg.Data[common.ImageHandle], topicName)
			}

		case err := <-subscriber.ErrorChannel:
			glog.Infof("-- Error receiving message: %v\n", err)
			if err != nil {
				errMessage := "Error while receiving from topic: %s, Error: %s"
				glog.Infof(errMessage, topicName, err)
			}
		}
	}
}

// StopAllSubscribers - function to close all subscriber objects
func (subMgr *SubManager) StopAllSubscribers() {
	for _, sub := range subMgr.subscribers {
		sub.Close()
	}
}
