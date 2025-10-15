#!/bin/bash

PWD=$(pwd)
LOGS_DIR="$PWD/../logs"
ARCHIVES_DIR="$LOGS_DIR/archives"
TIMESTAMP=$(date +"%Y%m%s%d-%H%M")
MAX_SIZE=$((5 * 1024 * 1024)) # 5MB in bytes

# Detect OS type for stat compatibility
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    get_file_size() {
        stat -f%z "$1"
    }
else
    # Linux
    get_file_size() {
        stat -c%s "$1"
    }
fi

# Function to archive and compress log files
archive_log() {
    local log_file=$1
    local service_name=$2
    local log_size=$(get_file_size "$log_file")

    if (( log_size > MAX_SIZE )); then
        local archive_subdir="$ARCHIVES_DIR/$service_name"
        mkdir -p "$archive_subdir"
        local archive_file="$archive_subdir/${service_name}-${TIMESTAMP}.log.gz"
        
        # Compress and move the log file to the archives directory
        gzip -c "$log_file" > "$archive_file"
        
        # Clear the original log file
        : > "$log_file"
        
        echo "✅ Archived and compressed $log_file → $archive_file"
    else
        echo "ℹ️  $log_file is under the size limit; no action taken."
    fi
}

# Loop through all log files (excluding archives)
for log_file in "$LOGS_DIR"/*.log; do
    [ -e "$log_file" ] || continue
    service_name=$(basename "$log_file" .log)
    archive_log "$log_file" "$service_name"
done