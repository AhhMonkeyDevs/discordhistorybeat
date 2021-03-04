package beater

import (
	"encoding/json"
	"fmt"
	"github.com/AhhMonkeyDevs/discordgo-lite"
)

func getChannels(token string, guildID string) []discordgo.Channel {

	data := make(chan []byte)

	discordgo.NewRestRequest().
		Method("GET").
		Token(token).
		Route("guilds").
		Guild(guildID).
		Route("channels").
		Callback(data).
		Enqueue()

	bytes := <-data

	var channels []discordgo.Channel
	_ = json.Unmarshal(bytes, &channels)

	return channels

}

type ChannelIterator struct {
	token     string
	channelID string
	lastID    string
}

func GetChannelIterator(token string, channelID string, startID string) *ChannelIterator {
	iterator := ChannelIterator{
		token:     token,
		channelID: channelID,
		lastID:    startID,
	}
	return &iterator
}

func (i *ChannelIterator) Next() []discordgo.Message {
	data := make(chan []byte)

	discordgo.NewRestRequest().
		Token(i.token).
		Route("channels").
		Channel(i.channelID).
		Route("messages").
		Query(fmt.Sprintf("after=%s&limit=%d", i.lastID, 100)).
		Callback(data).
		Enqueue()

	response := <-data

	var messages []discordgo.Message
	_ = json.Unmarshal(response, &messages)

	i.lastID = messages[0].Id

	return messages
}
