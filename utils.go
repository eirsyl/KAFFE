package main

import (
	"net"

	"fmt"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
)

func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

func PostToSlack(token, channel string, ip net.IP) error {
	api := slack.New(token)
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Text: fmt.Sprintf("Moccamaster listening on %v", ip.String()),
	}
	params.Attachments = []slack.Attachment{attachment}
	channelID, _, err := api.PostMessage(channel, ":coffee: Moccamaster has access to the world!", params)
	if err != nil {
		return err
	}
	log.Infof("Sent ip update message to channel: %v", channelID)
	return nil
}
