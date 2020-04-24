# Slack Notify - GitHub Action
[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#active)


A [GitHub Action](https://github.com/features/actions) to send a message to a Slack channel.

## Usage

You can use this action after any other action. Here is an example setup of this action:

1. Create a `.github/workflows/slack-notify.yml` file in your GitHub repo.
2. Add the following code to the `slack-notify.yml` file.

```yml
on: push
name: Notify Slack Demo
jobs:
  slackNotification:
    name: Notify Slack
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - name: Notify Slack
      uses: drilonrecica/action-slack-notify@0.0.3
      env:
        SLACK_WEBHOOK: ${{ secrets.SLACK_WEBHOOK }}
```

3. Create `SLACK_WEBHOOK` secret [GitHub Action's Secret](https://help.github.com/en/actions/configuring-and-managing-workflows/creating-and-storing-encrypted-secrets). You can [generate a Slack incoming webhook token from here](https://slack.com/apps/A0F7XDUAZ-incoming-webhooks).


## Environment Variables

By default, action is designed to run with minimal configuration but you can alter Slack notification using following environment variables:

Variable       | Default                                               | Purpose
---------------|-------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------
SLACK_CHANNEL  | Set during Slack webhook creation                     | Specify Slack channel in which message needs to be sent
PRIMARY_SLACK_CHANNEL  | Set during Slack webhook creation                     | You might want to send some of the notifications to primary channel as well instead of "log" channel
SLACK_USERNAME | `PM's Bot`                                               | The name of the sender of the message. Does not need to be a "real" username
SLACK_ICON     | ![default Avatar](https://avatars1.githubusercontent.com/u/63812035?s=32&v=4) | User/Bot icon shown with Slack message
SLACK_COLOR    | `good` (green)                                        | You can pass an RGB value like `#efefef` which would change color on left side vertical line of Slack message.
SLACK_MESSAGE  | Generated from git commit message.                    | The main Slack message in attachment. It is advised not to override this.
SLACK_TITLE    | Message                                               | Title to use before main Slack message

You can see the action block with all variables as below:

```yml
    - name: Slack notifier
      uses: pmagentur/pm_slack_notifier@v0.1
      env:
        SLACK_CHANNEL: deploy-logs
        SLACK_COLOR: '#FF0000'
        SLACK_ICON: https://avatars1.githubusercontent.com/u/63812035?s=200&v=4
        SLACK_MESSAGE: 'Post Content :rocket:'
        SLACK_TITLE: 'Post Title'
        SLACK_USERNAME: 'PM's whisperer'
        SLACK_WEBHOOK: ${{ secrets.SLACK_WEBHOOK }}
```

## License

[MIT](LICENSE) © 2019 PM Agenture

<sub><sup>The project is forked from the following [project] [https://github.com/rtCamp/action-slack-notify)</sup></sub>
