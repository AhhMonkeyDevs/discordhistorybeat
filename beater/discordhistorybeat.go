package beater

import (
	"fmt"
	"github.com/AhhMonkeyDevs/discordgo-lite"
	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"
	"time"

	"github.com/AhhMonkeyDevs/discordhistorybeat/config"
)

// discordhistorybeat configuration.
type discordhistorybeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// New creates an instance of discordhistorybeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &discordhistorybeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

// Run starts discordhistorybeat.
func (bt *discordhistorybeat) Run(b *beat.Beat) error {
	logp.Info("discordhistorybeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	fmt.Println(bt.config.Token)

	channels := getChannels(bt.config.Token, bt.config.GuildID)
	fmt.Println(channels)

	running := true
	go func() {

		for i, channel := range channels {
			if channel.Type != 0 {
				continue
			}
			fmt.Printf("Starting channel '%s' (%d of %d)\n", channel.Name, i, len(channels))

			i := GetChannelIterator(bt.config.Token, channel.Id, bt.config.StartID)
			count := 0
			for running {
				fmt.Printf("\rProcessed %d messages", count)

				messages := i.Next()
				if len(messages) == 0 {
					break
				}

				count += len(messages)

				events := make([]beat.Event, 100)
				for k, message := range messages {
					eventFields := messageToFields(&message)
					eventFields.Put("event", "MESSAGE_HISTORY")
					event := beat.Event{
						Timestamp: time.Now(),
						Fields:    eventFields,
					}
					events[k] = event
				}
				bt.client.PublishAll(events)
			}
			fmt.Printf("\nCompleed channel %s\n", channel.Name)
		}

		fmt.Println("\nAll done")

	}()

	<-bt.done

	running = false

	return nil

}

// Stop stops discordhistorybeat.
func (bt *discordhistorybeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

func messageToFields(msg *discordgo.Message) common.MapStr {

	mf := GetMessageFormatter(*msg)

	return common.MapStr{
		"id":                    msg.Id,
		"channel_id":            msg.ChannelID,
		"guild_id":              msg.GuildID,
		"author_id":             mf.getAuthorID(),
		"author_type":           mf.getAuthorType(),
		"type":                  msg.Type,
		"user_mentions":         mf.getUserMentions(),
		"role_mentions":         mf.getRoleMentions(),
		"channel_mentions":      mf.getChannelMentions(),
		"content":               mf.getContent(),
		"tts":                   msg.Tts,
		"mention_everyone":      msg.MentionsEveryone,
		"link_hostnames":        mf.extractLinks(),
		"has":                   mf.getHasArray(),
		"attachment_filenames":  mf.getAttachmentFilenames(),
		"attachment_extensions": mf.getAttachmentExtensions(),
		"referenced_message":    mf.getMessageReference(),
		"created_timestamp":     GetTimestamp(msg.Timestamp),
		"edited_timestamp":      GetTimestamp(msg.EditedTimestamp),
	}
}
