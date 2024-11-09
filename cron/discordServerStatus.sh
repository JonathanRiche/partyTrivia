#!/bin/bash

# Function to display usage information
usage() {
    echo "Usage: $0 -u WEBSITE_URL -n WEBSITE_NAME -w WEBHOOK_URL [-h]"
    echo "  -u WEBSITE_URL   : The URL of the website to check"
    echo "  -n WEBSITE_NAME  : The name of the website"
    echo "  -w WEBHOOK_URL   : The Discord webhook URL"
    echo "  -h               : Display this help message"
    exit 1
}

# Parse command-line options
while getopts "u:n:w:h" opt; do
    case $opt in
        u) SITE_URL="$OPTARG" ;;
        n) WEBSITE_NAME="$OPTARG" ;;
        w) WEBHOOK_URL="$OPTARG" ;;
        h) usage ;;
        \?) echo "Invalid option -$OPTARG" >&2; usage ;;
    esac
done

# Check if required parameters are provided
if [ -z "$SITE_URL" ] || [ -z "$WEBSITE_NAME" ] || [ -z "$WEBHOOK_URL" ]; then
    echo "Error: Missing required parameters."
    usage
fi

response=$(curl -sL -w "%{http_code}" "$SITE_URL" -o /dev/null)
current_time=$(date "+%Y-%m-%d %H:%M:%S")

if [ $response == "200" ]; then
    status="Successful"
    emoji="✅"
    color="65280"  # Green in decimal
else
    status="Failed"
    emoji="⚠️"
    color="16744192"  # Orange in decimal
fi

json_payload=$(cat <<EOF
{
    "embeds": [{
        "title": "$emoji $WEBSITE_NAME Daily Check",
        "description": "Daily status check $status.",
        "color": $color,
        "fields": [
            {
                "name": "Site URL",
                "value": "$SITE_URL"
            },
            {
                "name": "Status Code",
                "value": "$response"
            },
            {
                "name": "Timestamp",
                "value": "$current_time"
            }
        ],
        "thumbnail": {
            "url": "$SITE_URL/favicon.ico"
        }
    }]
}
EOF
)

curl -H "Content-Type: application/json" -d "$json_payload" "$WEBHOOK_URL"

