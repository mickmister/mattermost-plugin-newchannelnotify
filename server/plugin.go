package main

import (
	"fmt"
	"sync"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

const defaultBotName = "New-Channel-Bot"

type NewChannelNotifyPlugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

func (p *NewChannelNotifyPlugin) OnActivate() error {
	p.API.LogInfo("Plugin loaded")
	return nil
}

func (p *NewChannelNotifyPlugin) ChannelHasBeenCreated(c *plugin.Context, channel *model.Channel) {
	log := fmt.Sprintf("ChannelHasBeenCreated for channel with id [%s], type [%s] triggerd", channel.Id, channel.Type)
	p.API.LogDebug(log)

	config := p.getConfiguration()

	if config.BotUserName == "" {
		config.BotUserName = defaultBotName
	}

	if config.ChannelToPost == "" {
		config.ChannelToPost = model.DEFAULT_CHANNEL
	}

	newChannelName := channel.Name

	if channel.Type == model.CHANNEL_DIRECT || channel.Type == model.CHANNEL_GROUP {
		return
	}

	if channel.Type == model.CHANNEL_PRIVATE {
		if config.IncludePrivateChannels == false {
			return
		}
		newChannelName += " [Private]"
	}

	p.ensureBotExists()
	bot, err := p.API.GetUserByUsername(config.BotUserName)

	mainChannel, err := p.API.GetChannelByName(channel.TeamId, config.ChannelToPost, false)
	if err != nil {
		p.API.LogError(err.Message)
	}

	creator, err := p.API.GetUser(channel.CreatorId)
	if err != nil {
		p.API.LogError(err.Message)
	}

	post, err := p.API.CreatePost(&model.Post{
		ChannelId: mainChannel.Id,
		UserId:    bot.Id,
		Message:   fmt.Sprintf("@channel Hello there :wave:. You might want to check out the new channel ~%s created by @%s :).", newChannelName, creator.Username),
	})

	p.API.LogDebug(fmt.Sprintf("Created post %s", post.Id))

	if err != nil {
		p.API.LogError(err.Message)
	}
}
