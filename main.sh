#!/usr/bin/env bash

export GITHUB_BRANCH=${GITHUB_REF##*heads/}
export SLACK_ICON=${SLACK_ICON:-"https://avatars1.githubusercontent.com/u/63812035?s=200&v=4"}
export SLACK_USERNAME=${SLACK_USERNAME:-"PM's whisperer"}
export SLACK_TITLE=${SLACK_TITLE:-"Notification"}
export COMMIT_MESSAGE=$(cat "/github/workflow/event.json" | jq .commits | jq '.[0].message' -r)

if [[ -n "$user_slack_channel" ]]; then
	export SLACK_CHANNEL="$user_slack_channel"
fi


if [[ -n "$SITE_NAME" ]]; then
    export SITE_TITLE="Site"
fi


if [[ -z "$SLACK_MESSAGE" ]]; then
	export SLACK_MESSAGE="$COMMIT_MESSAGE"
fi

slack-notify "$@"