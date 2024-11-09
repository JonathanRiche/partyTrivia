#!/bin/bash

# Function to display usage information
usage() {
    echo "Usage: $0 -u WEBSITE_URL -n WEBSITE_NAME -w WEBHOOK_URL [-h]"
    echo "  -u WEBSITE_URL   : The URL of the website to check"
    echo "  -n WEBSITE_NAME  : The name of the website"
    echo "  -w WEBHOOK_URL   : The Microsoft Teams webhook URL"
    echo "  -h               : Display this help message"
    exit 1
}

# Parse command-line options
while getopts "u:n:w:h" opt; do
    case $opt in
        u) SITE_URL="$OPTARG" ;;
        n) WEBSITE_NAME="$OPTARG" ;;
        w) TEAMS_WEBHOOK_URL="$OPTARG" ;;
        h) usage ;;
        \?) echo "Invalid option -$OPTARG" >&2; usage ;;
    esac
done

# Check if required parameters are provided
if [ -z "$SITE_URL" ] || [ -z "$WEBSITE_NAME" ] || [ -z "$TEAMS_WEBHOOK_URL" ]; then
    echo "Error: Missing required parameters."
    usage
fi

response=$(curl -sL -w "%{http_code}" "$SITE_URL" -o /dev/null)
current_time=$(date "+%Y-%m-%d %H:%M:%S")

if [ $response == "200" ]; then
    json_payload=$(cat <<EOF
{
    "@type": "MessageCard",
    "@context": "http://schema.org/extensions",
    "themeColor": "00FF00",
    "summary": "$WEBSITE_NAME Daily Check - Successful",
    "sections": [{
        "activityTitle": "✅ $WEBSITE_NAME Daily Check",
        "activitySubtitle": "Daily status check completed successfully.",
        "activityImage": "$SITE_URL/favicon.ico",
        "facts": [{
            "name": "Site URL:",
            "value": "$SITE_URL"
        }, {
            "name": "Status Code:",
            "value": "$response"
        }, {
            "name": "Timestamp:",
            "value": "$current_time"
        }],
        "markdown": true
    }],
    "potentialAction": [{
        "@type": "OpenUri",
        "name": "Visit Site",
        "targets": [{
            "os": "default",
            "uri": "$SITE_URL"
        }]
    }]
}
EOF
)
else
    json_payload=$(cat <<EOF
{
    "@type": "MessageCard",
    "@context": "http://schema.org/extensions",
    "themeColor": "FFA500",
    "summary": "$WEBSITE_NAME Daily Check - Failed",
    "sections": [{
        "activityTitle": "⚠️ $WEBSITE_NAME Daily Check",
        "activitySubtitle": "The daily status check has detected an issue.",
        "activityImage": "$SITE_URL/favicon.ico",
        "facts": [{
            "name": "Site URL:",
            "value": "$SITE_URL"
        }, {
            "name": "Status Code:",
            "value": "$response"
        }, {
            "name": "Timestamp:",
            "value": "$current_time"
        }],
        "markdown": true
    }],
    "potentialAction": [{
        "@type": "OpenUri",
        "name": "Check Site",
        "targets": [{
            "os": "default",
            "uri": "$SITE_URL"
        }]
    }]
}
EOF
)
fi

curl -H "Content-Type: application/json" -d "$json_payload" "$TEAMS_WEBHOOK_URL"

