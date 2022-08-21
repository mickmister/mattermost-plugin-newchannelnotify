package main

import (
	"fmt"
	"sync"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
)

const defaultBotName = "newchannelbot"

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

	if channel.CreatorId == "" {
		p.API.LogDebug("Not creating post due to channel being created through automation.", "id", channel.Id)
		return
	}

	config := p.getConfiguration()
	ChannelPurpose := ""

	if config.BotUserName == "" {
		config.BotUserName = defaultBotName
	}

	if config.ChannelToPost == "" {
		return
	}

	if config.IncludeChannelPurpose && channel.Purpose != "" {
		ChannelPurpose = "\n **" + channel.Name + "'s Purpose:** " + channel.Purpose
	}

	newChannelName := channel.Name

	if channel.Type == model.ChannelTypeDirect || channel.Type == model.ChannelTypeGroup {
		return
	}

	if channel.Type == model.ChannelTypePrivate {
		if !config.IncludePrivateChannels {
			return
		}
		newChannelName += " [Private]"
	}

	p.ensureBotExists()
	bot, err := p.API.GetUserByUsername(config.BotUserName)
	if err != nil {
		p.API.LogError(err.Message)
		return
	}

	mainChannel, err := p.API.GetChannelByName(channel.TeamId, config.ChannelToPost, false)
	if err != nil {
		p.API.LogError(err.Message)
		return
	}

	creator, err := p.API.GetUser(channel.CreatorId)
	if err != nil {
		p.API.LogError(err.Message)
		return
	}

	post, err := p.API.CreatePost(&model.Post{
		ChannelId: mainChannel.Id,
		UserId:    bot.Id,
		Message:   fmt.Sprintf("%sHello there :wave:. You might want to check out the new channel ~%s created by @%s %s", config.Mention, newChannelName, creator.Username, ChannelPurpose),
	})

	if err != nil {
		p.API.LogError(err.Message)
		return
	}

	p.API.LogDebug(fmt.Sprintf("Created post %s", post.Id))
}
