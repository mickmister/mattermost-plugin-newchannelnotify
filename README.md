# Mattermost New Channel Notify Plugin

A plugin for Mattermost to notify all users about newly created channels.

![screenshot](https://i.imgur.com/SII7ZEi.png)

## Installation

__Requires Mattermost 5.10 or higher.__

Download the [latest release here](https://gitlab.com/thepill/mattermost-plugin-newchannelnotify/uploads/84e1c6feb323a0ccae1a19921466ed73/matttermost-plugin-newchannelnotify-0.9.2.tar.gz)

In production, deploy and upload your plugin via the [System Console](https://about.mattermost.com/default-plugin-uploads).

Optionally, change `settings` under the plugins settings menu in System Console:
- Bot Name
- Channel to Post
- Include private channels
- Mention
- IncludeChannelPurpose

## Developing 

Use `make dist` to build distributions of the plugin that you can upload to a Mattermost server.

Use `make deploy` to deploy the plugin to your local server. Before running `make deploy` you need to set a few environment variables:

```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_USERNAME=admin
export MM_ADMIN_PASSWORD=password
```
