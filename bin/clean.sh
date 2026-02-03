#!/bin/bash
# Mole - Clean command.
# Runs cleanup modules with optional sudo.
# Supports dry-run and whitelist.

set -euo pipefail

export LC_ALL=C
export LANG=C

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
<<<<<<< HEAD
source "$SCRIPT_DIR/../lib/common.sh"
source "$SCRIPT_DIR/../lib/clean_brew.sh"
source "$SCRIPT_DIR/../lib/clean_caches.sh"
source "$SCRIPT_DIR/../lib/clean_apps.sh"
source "$SCRIPT_DIR/../lib/clean_dev.sh"
source "$SCRIPT_DIR/../lib/clean_user_apps.sh"
source "$SCRIPT_DIR/../lib/clean_system.sh"
source "$SCRIPT_DIR/../lib/clean_user_data.sh"
=======
source "$SCRIPT_DIR/../lib/core/common.sh"

source "$SCRIPT_DIR/../lib/core/sudo.sh"
source "$SCRIPT_DIR/../lib/clean/brew.sh"
source "$SCRIPT_DIR/../lib/clean/caches.sh"
source "$SCRIPT_DIR/../lib/clean/apps.sh"
source "$SCRIPT_DIR/../lib/clean/dev.sh"
source "$SCRIPT_DIR/../lib/clean/app_caches.sh"
source "$SCRIPT_DIR/../lib/clean/system.sh"
source "$SCRIPT_DIR/../lib/clean/user.sh"
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b

SYSTEM_CLEAN=false
DRY_RUN=false
PROTECT_FINDER_METADATA=false
IS_M_SERIES=$([[ "$(uname -m)" == "arm64" ]] && echo "true" || echo "false")

<<<<<<< HEAD
# Protected Service Worker domains (web-based editing tools)
=======
EXPORT_LIST_FILE="$HOME/.config/mole/clean-list.txt"
CURRENT_SECTION=""
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
readonly PROTECTED_SW_DOMAINS=(
    "capcut.com"
    "photopea.com"
    "pixlr.com"
)
<<<<<<< HEAD
readonly FINDER_METADATA_SENTINEL="FINDER_METADATA"
# Default whitelist patterns (preselected, user can disable)
declare -a DEFAULT_WHITELIST_PATTERNS=(
    "$HOME/Library/Caches/ms-playwright*"
    "$HOME/.cache/huggingface*"
    "$HOME/.m2/repository/*"
    "$HOME/.ollama/models/*"
    "$HOME/Library/Caches/com.nssurge.surge-mac/*"
    "$HOME/Library/Application Support/com.nssurge.surge-mac/*"
    "$HOME/Library/Caches/org.R-project.R/R/renv/*"
    "$FINDER_METADATA_SENTINEL"
)
=======

>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
declare -a WHITELIST_PATTERNS=()
WHITELIST_WARNINGS=()
if [[ -f "$HOME/.config/mole/whitelist" ]]; then
    while IFS= read -r line; do
        # shellcheck disable=SC2295
        line="${line#"${line%%[![:space:]]*}"}"
        # shellcheck disable=SC2295
        line="${line%"${line##*[![:space:]]}"}"
        [[ -z "$line" || "$line" =~ ^# ]] && continue

        [[ "$line" == ~* ]] && line="${line/#~/$HOME}"
<<<<<<< HEAD

        # Security: reject path traversal attempts
        if [[ "$line" =~ \.\. ]]; then
            WHITELIST_WARNINGS+=("Path traversal not allowed: $line")
            continue
        fi

        # Path validation with support for spaces and wildcards
        # Allow: letters, numbers, /, _, ., -, @, spaces, and * anywhere in path
        if [[ ! "$line" =~ ^[a-zA-Z0-9/_.@\ *-]+$ ]]; then
            WHITELIST_WARNINGS+=("Invalid path format: $line")
            continue
        fi

        # Require absolute paths (must start with /)
        if [[ "$line" != /* ]]; then
            WHITELIST_WARNINGS+=("Must be absolute path: $line")
            continue
        fi

        # Reject paths with consecutive slashes (e.g., //)
        if [[ "$line" =~ // ]]; then
            WHITELIST_WARNINGS+=("Consecutive slashes: $line")
            continue
        fi

        # Prevent critical system directories
=======
        line="${line//\$HOME/$HOME}"
        line="${line//\$\{HOME\}/$HOME}"
        if [[ "$line" =~ \.\. ]]; then
            WHITELIST_WARNINGS+=("Path traversal not allowed: $line")
            continue
        fi

        if [[ "$line" != "$FINDER_METADATA_SENTINEL" ]]; then
            if [[ ! "$line" =~ ^[a-zA-Z0-9/_.@\ *-]+$ ]]; then
                WHITELIST_WARNINGS+=("Invalid path format: $line")
                continue
            fi

            if [[ "$line" != /* ]]; then
                WHITELIST_WARNINGS+=("Must be absolute path: $line")
                continue
            fi
        fi

        if [[ "$line" =~ // ]]; then
            WHITELIST_WARNINGS+=("Consecutive slashes: $line")
            continue
        fi

>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
        case "$line" in
            / | /System | /System/* | /bin | /bin/* | /sbin | /sbin/* | /usr/bin | /usr/bin/* | /usr/sbin | /usr/sbin/* | /etc | /etc/* | /var/db | /var/db/*)
                WHITELIST_WARNINGS+=("Protected system path: $line")
                continue
                ;;
        esac

        duplicate="false"
        if [[ ${#WHITELIST_PATTERNS[@]} -gt 0 ]]; then
            for existing in "${WHITELIST_PATTERNS[@]}"; do
                if [[ "$line" == "$existing" ]]; then
                    duplicate="true"
                    break
                fi
            done
        fi
        [[ "$duplicate" == "true" ]] && continue
        WHITELIST_PATTERNS+=("$line")
    done < "$HOME/.config/mole/whitelist"
else
    WHITELIST_PATTERNS=("${DEFAULT_WHITELIST_PATTERNS[@]}")
fi
<<<<<<< HEAD

if [[ ${#WHITELIST_PATTERNS[@]} -gt 0 ]]; then
    for entry in "${WHITELIST_PATTERNS[@]}"; do
        if [[ "$entry" == "$FINDER_METADATA_SENTINEL" ]]; then
            PROTECT_FINDER_METADATA=true
            break
        fi
    done
fi
total_items=0
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b

# Expand whitelist patterns once to avoid repeated tilde expansion in hot loops.
expand_whitelist_patterns() {
    if [[ ${#WHITELIST_PATTERNS[@]} -gt 0 ]]; then
        local -a EXPANDED_PATTERNS
        EXPANDED_PATTERNS=()
        for pattern in "${WHITELIST_PATTERNS[@]}"; do
            local expanded="${pattern/#\~/$HOME}"
            EXPANDED_PATTERNS+=("$expanded")
        done
        WHITELIST_PATTERNS=("${EXPANDED_PATTERNS[@]}")
    fi
}
expand_whitelist_patterns

if [[ ${#WHITELIST_PATTERNS[@]} -gt 0 ]]; then
    for entry in "${WHITELIST_PATTERNS[@]}"; do
        if [[ "$entry" == "$FINDER_METADATA_SENTINEL" ]]; then
            PROTECT_FINDER_METADATA=true
            break
        fi
    done
fi

# Section tracking and summary counters.
total_items=0
TRACK_SECTION=0
SECTION_ACTIVITY=0
files_cleaned=0
total_size_cleaned=0
whitelist_skipped_count=0

# shellcheck disable=SC2329
note_activity() {
    if [[ "${TRACK_SECTION:-0}" == "1" ]]; then
        SECTION_ACTIVITY=1
    fi
}

CLEANUP_DONE=false
# shellcheck disable=SC2329
cleanup() {
    local signal="${1:-EXIT}"
    local exit_code="${2:-$?}"

    if [[ "$CLEANUP_DONE" == "true" ]]; then
        return 0
    fi
    CLEANUP_DONE=true

    stop_inline_spinner 2> /dev/null || true

    cleanup_temp_files

    stop_sudo_session

    show_cursor
}

trap 'cleanup EXIT $?' EXIT
trap 'cleanup INT 130; exit 130' INT
trap 'cleanup TERM 143; exit 143' TERM

start_section() {
    TRACK_SECTION=1
    SECTION_ACTIVITY=0
    CURRENT_SECTION="$1"
    echo ""
    echo -e "${PURPLE_BOLD}${ICON_ARROW} $1${NC}"

    if [[ "$DRY_RUN" == "true" ]]; then
        ensure_user_file "$EXPORT_LIST_FILE"
        echo "" >> "$EXPORT_LIST_FILE"
        echo "=== $1 ===" >> "$EXPORT_LIST_FILE"
    fi
}

end_section() {
<<<<<<< HEAD
    if [[ $TRACK_SECTION -eq 1 && $SECTION_ACTIVITY -eq 0 ]]; then
=======
    stop_section_spinner

    if [[ "${TRACK_SECTION:-0}" == "1" && "${SECTION_ACTIVITY:-0}" == "0" ]]; then
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
        echo -e "  ${GREEN}${ICON_SUCCESS}${NC} Nothing to clean"
    fi
    TRACK_SECTION=0
}

# shellcheck disable=SC2329
normalize_paths_for_cleanup() {
    local -a input_paths=("$@")
    local -a unique_paths=()

    for path in "${input_paths[@]}"; do
        local normalized="${path%/}"
        [[ -z "$normalized" ]] && normalized="$path"
        local found=false
        if [[ ${#unique_paths[@]} -gt 0 ]]; then
            for existing in "${unique_paths[@]}"; do
                if [[ "$existing" == "$normalized" ]]; then
                    found=true
                    break
                fi
            done
        fi
        [[ "$found" == "true" ]] || unique_paths+=("$normalized")
    done

    local sorted_paths
    if [[ ${#unique_paths[@]} -gt 0 ]]; then
        sorted_paths=$(printf '%s\n' "${unique_paths[@]}" | awk '{print length "|" $0}' | LC_ALL=C sort -n | cut -d'|' -f2-)
    else
        sorted_paths=""
    fi

    local -a result_paths=()
    while IFS= read -r path; do
        [[ -z "$path" ]] && continue
        local is_child=false
        if [[ ${#result_paths[@]} -gt 0 ]]; then
            for kept in "${result_paths[@]}"; do
                if [[ "$path" == "$kept" || "$path" == "$kept"/* ]]; then
                    is_child=true
                    break
                fi
            done
        fi
        [[ "$is_child" == "true" ]] || result_paths+=("$path")
    done <<< "$sorted_paths"

    if [[ ${#result_paths[@]} -gt 0 ]]; then
        printf '%s\n' "${result_paths[@]}"
    fi
}

# shellcheck disable=SC2329
get_cleanup_path_size_kb() {
    local path="$1"

    if [[ -f "$path" && ! -L "$path" ]]; then
        if command -v stat > /dev/null 2>&1; then
            local bytes
            bytes=$(stat -f%z "$path" 2> /dev/null || echo "0")
            if [[ "$bytes" =~ ^[0-9]+$ && "$bytes" -gt 0 ]]; then
                echo $(((bytes + 1023) / 1024))
                return 0
            fi
        fi
    fi

    if [[ -L "$path" ]]; then
        if command -v stat > /dev/null 2>&1; then
            local bytes
            bytes=$(stat -f%z "$path" 2> /dev/null || echo "0")
            if [[ "$bytes" =~ ^[0-9]+$ && "$bytes" -gt 0 ]]; then
                echo $(((bytes + 1023) / 1024))
            else
                echo 0
            fi
            return 0
        fi
    fi

    get_path_size_kb "$path"
}

# Classification helper for cleanup risk levels
# shellcheck disable=SC2329
classify_cleanup_risk() {
    local description="$1"
    local path="${2:-}"

    # HIGH RISK: System files, preference files, require sudo
    if [[ "$description" =~ [Ss]ystem || "$description" =~ [Ss]udo || "$path" =~ ^/System || "$path" =~ ^/Library ]]; then
        echo "HIGH|System files or requires admin access"
        return
    fi

    # HIGH RISK: Preference files that might affect app functionality
    if [[ "$description" =~ [Pp]reference || "$path" =~ /Preferences/ ]]; then
        echo "HIGH|Preference files may affect app settings"
        return
    fi

    # MEDIUM RISK: Installers, large files, app bundles
    if [[ "$description" =~ [Ii]nstaller || "$description" =~ [Aa]pp.*[Bb]undle || "$description" =~ [Ll]arge ]]; then
        echo "MEDIUM|Installer packages or app data"
        return
    fi

    # MEDIUM RISK: Old backups, downloads
    if [[ "$description" =~ [Bb]ackup || "$description" =~ [Dd]ownload || "$description" =~ [Oo]rphan ]]; then
        echo "MEDIUM|Backup or downloaded files"
        return
    fi

    # LOW RISK: Caches, logs, temporary files (automatically regenerated)
    if [[ "$description" =~ [Cc]ache || "$description" =~ [Ll]og || "$description" =~ [Tt]emp || "$description" =~ [Tt]humbnail ]]; then
        echo "LOW|Cache/log files, automatically regenerated"
        return
    fi

    # DEFAULT: MEDIUM
    echo "MEDIUM|User data files"
}

# shellcheck disable=SC2329
safe_clean() {
    if [[ $# -eq 0 ]]; then
        return 0
    fi

    local description
    local -a targets

    if [[ $# -eq 1 ]]; then
        description="$1"
        targets=("$1")
    else
        description="${*: -1}"
        targets=("${@:1:$#-1}")
    fi

    local -a valid_targets=()
    for target in "${targets[@]}"; do
        # Optimization: If target is a glob literal and parent dir missing, skip it.
        if [[ "$target" == *"*"* && ! -e "$target" ]]; then
            local base_path="${target%%\**}"
            local parent_dir
            if [[ "$base_path" == */ ]]; then
                parent_dir="${base_path%/}"
            else
                parent_dir=$(dirname "$base_path")
            fi

            if [[ ! -d "$parent_dir" ]]; then
                # debug_log "Skipping nonexistent parent: $parent_dir for $target"
                continue
            fi
        fi
        valid_targets+=("$target")
    done

    if [[ ${#valid_targets[@]} -gt 0 ]]; then
        targets=("${valid_targets[@]}")
    else
        targets=()
    fi
    if [[ ${#targets[@]} -eq 0 ]]; then
        return 0
    fi

    local removed_any=0
    local total_size_kb=0
    local total_count=0
    local skipped_count=0
    local removal_failed_count=0
    local permission_start=${MOLE_PERMISSION_DENIED_COUNT:-0}

    local show_scan_feedback=false
    if [[ ${#targets[@]} -gt 20 && -t 1 ]]; then
        show_scan_feedback=true
        stop_section_spinner
        MOLE_SPINNER_PREFIX="  " start_inline_spinner "Scanning ${#targets[@]} items..."
    fi

    local -a existing_paths=()
    for path in "${targets[@]}"; do
        local skip=false

<<<<<<< HEAD
        # Hard-coded protection for critical apps (cannot be disabled by user)
        case "$path" in
            *clash* | *Clash* | *surge* | *Surge* | *mihomo* | *openvpn* | *OpenVPN*)
                skip=true
                ((skipped_count++))
                ;;
        esac

        [[ "$skip" == "true" ]] && continue

        # Check user-defined whitelist
        if [[ ${#WHITELIST_PATTERNS[@]} -gt 0 ]]; then
            for w in "${WHITELIST_PATTERNS[@]}"; do
                # Match both exact path and glob pattern
                # shellcheck disable=SC2053
                if [[ "$path" == "$w" ]] || [[ $path == $w ]]; then
                    skip=true
                    ((skipped_count++))
                    break
                fi
            done
=======
        if should_protect_path "$path"; then
            skip=true
            ((skipped_count++))
            log_operation "clean" "SKIPPED" "$path" "protected"
        fi

        [[ "$skip" == "true" ]] && continue

        if is_path_whitelisted "$path"; then
            skip=true
            ((skipped_count++))
            log_operation "clean" "SKIPPED" "$path" "whitelist"
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
        fi
        [[ "$skip" == "true" ]] && continue
        [[ -e "$path" ]] && existing_paths+=("$path")
    done

    if [[ "$show_scan_feedback" == "true" ]]; then
        stop_section_spinner
    fi

    debug_log "Cleaning: $description, ${#existing_paths[@]} items"

    # Enhanced debug output with risk level and details
    if [[ "${MO_DEBUG:-}" == "1" && ${#existing_paths[@]} -gt 0 ]]; then
        # Determine risk level for this cleanup operation
        local risk_info
        risk_info=$(classify_cleanup_risk "$description" "${existing_paths[0]}")
        local risk_level="${risk_info%%|*}"
        local risk_reason="${risk_info#*|}"

        debug_operation_start "$description"
        debug_risk_level "$risk_level" "$risk_reason"
        debug_operation_detail "Item count" "${#existing_paths[@]}"

        # Log sample of files (first 10) with details
        if [[ ${#existing_paths[@]} -le 10 ]]; then
            debug_operation_detail "Files to be removed" "All files listed below"
        else
            debug_operation_detail "Files to be removed" "Showing first 10 of ${#existing_paths[@]} files"
        fi
    fi

    if [[ $skipped_count -gt 0 ]]; then
        ((whitelist_skipped_count += skipped_count))
    fi

    if [[ ${#existing_paths[@]} -eq 0 ]]; then
        return 0
    fi

    if [[ ${#existing_paths[@]} -gt 1 ]]; then
        local -a normalized_paths=()
        while IFS= read -r path; do
            [[ -n "$path" ]] && normalized_paths+=("$path")
        done < <(normalize_paths_for_cleanup "${existing_paths[@]}")

        if [[ ${#normalized_paths[@]} -gt 0 ]]; then
            existing_paths=("${normalized_paths[@]}")
        else
            existing_paths=()
        fi
    fi

    local show_spinner=false
    if [[ ${#existing_paths[@]} -gt 10 ]]; then
        show_spinner=true
        local total_paths=${#existing_paths[@]}
        if [[ -t 1 ]]; then MOLE_SPINNER_PREFIX="  " start_inline_spinner "Scanning items..."; fi
    fi

    local cleaning_spinner_started=false

    # For larger batches, precompute sizes in parallel for better UX/stat accuracy.
    if [[ ${#existing_paths[@]} -gt 3 ]]; then
        local temp_dir
        temp_dir=$(create_temp_dir)

        local dir_count=0
        local sample_size=$((${#existing_paths[@]} > 20 ? 20 : ${#existing_paths[@]}))
        local max_sample=$((${#existing_paths[@]} * 20 / 100))
        [[ $max_sample -gt $sample_size ]] && sample_size=$max_sample

        for ((i = 0; i < sample_size && i < ${#existing_paths[@]}; i++)); do
            [[ -d "${existing_paths[i]}" ]] && ((dir_count++))
        done

        # Heuristic: mostly files -> sequential stat is faster than subshells.
        if [[ $dir_count -lt 5 && ${#existing_paths[@]} -gt 20 ]]; then
            if [[ -t 1 && "$show_spinner" == "false" ]]; then
                MOLE_SPINNER_PREFIX="  " start_inline_spinner "Scanning items..."
                show_spinner=true
            fi

            local idx=0
            local last_progress_update
            last_progress_update=$(get_epoch_seconds)
            for path in "${existing_paths[@]}"; do
                local size
                size=$(get_cleanup_path_size_kb "$path")
                [[ ! "$size" =~ ^[0-9]+$ ]] && size=0

<<<<<<< HEAD
            if ((${#pids[@]} >= MOLE_MAX_PARALLEL_JOBS)); then
                wait "${pids[0]}" 2> /dev/null || true
                pids=("${pids[@]:1}")
                ((completed++))
                # Update progress every 10 items for smoother display
                if [[ -t 1 ]] && ((completed % 10 == 0)); then
                    stop_inline_spinner
                    MOLE_SPINNER_PREFIX="  " start_inline_spinner "Scanning items ($completed/$total_paths)..."
=======
                if [[ "$size" -gt 0 ]]; then
                    echo "$size 1" > "$temp_dir/result_${idx}"
                else
                    echo "0 0" > "$temp_dir/result_${idx}"
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
                fi

                ((idx++))
                if [[ $((idx % 20)) -eq 0 && "$show_spinner" == "true" && -t 1 ]]; then
                    update_progress_if_needed "$idx" "${#existing_paths[@]}" last_progress_update 1 || true
                    last_progress_update=$(get_epoch_seconds)
                fi
            done
        else
            local -a pids=()
            local idx=0
            local completed=0
            local last_progress_update
            last_progress_update=$(get_epoch_seconds)
            local total_paths=${#existing_paths[@]}

<<<<<<< HEAD
        # Read results using same index
        idx=0
        for path in "${existing_paths[@]}"; do
            local result_file="$temp_dir/result_${idx}"
            if [[ -f "$result_file" ]]; then
                read -r size count < "$result_file" 2> /dev/null || true
                if [[ "$count" -gt 0 && "$size" -gt 0 ]]; then
                    if [[ "$DRY_RUN" != "true" ]]; then
                        # Handle symbolic links separately (only remove the link, not the target)
                        if [[ -L "$path" ]]; then
                            rm "$path" 2> /dev/null || true
                        else
                            safe_remove "$path" true || true
=======
            if [[ ${#existing_paths[@]} -gt 0 ]]; then
                for path in "${existing_paths[@]}"; do
                    (
                        local size
                        size=$(get_cleanup_path_size_kb "$path")
                        [[ ! "$size" =~ ^[0-9]+$ ]] && size=0
                        local tmp_file="$temp_dir/result_${idx}.$$"
                        if [[ "$size" -gt 0 ]]; then
                            echo "$size 1" > "$tmp_file"
                        else
                            echo "0 0" > "$tmp_file"
                        fi
                        mv "$tmp_file" "$temp_dir/result_${idx}" 2> /dev/null || true
                    ) &
                    pids+=($!)
                    ((idx++))

                    if ((${#pids[@]} >= MOLE_MAX_PARALLEL_JOBS)); then
                        wait "${pids[0]}" 2> /dev/null || true
                        pids=("${pids[@]:1}")
                        ((completed++))

                        if [[ "$show_spinner" == "true" && -t 1 ]]; then
                            update_progress_if_needed "$completed" "$total_paths" last_progress_update 2 || true
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
                        fi
                    fi
                done
            fi

            if [[ ${#pids[@]} -gt 0 ]]; then
                for pid in "${pids[@]}"; do
                    wait "$pid" 2> /dev/null || true
                    ((completed++))

                    if [[ "$show_spinner" == "true" && -t 1 ]]; then
                        update_progress_if_needed "$completed" "$total_paths" last_progress_update 2 || true
                    fi
                done
            fi
        fi

        # Read results back in original order.
        # Start spinner for cleaning phase
        if [[ "$DRY_RUN" != "true" && ${#existing_paths[@]} -gt 0 && -t 1 ]]; then
            MOLE_SPINNER_PREFIX="  " start_inline_spinner "Cleaning..."
            cleaning_spinner_started=true
        fi
        idx=0
        if [[ ${#existing_paths[@]} -gt 0 ]]; then
            for path in "${existing_paths[@]}"; do
                local result_file="$temp_dir/result_${idx}"
                if [[ -f "$result_file" ]]; then
                    read -r size count < "$result_file" 2> /dev/null || true
                    local removed=0
                    if [[ "$DRY_RUN" != "true" ]]; then
                        if safe_remove "$path" true; then
                            removed=1
                        fi
                    else
                        removed=1
                    fi

                    if [[ $removed -eq 1 ]]; then
                        if [[ "$size" -gt 0 ]]; then
                            ((total_size_kb += size))
                        fi
                        ((total_count += 1))
                        removed_any=1
                    else
                        if [[ -e "$path" && "$DRY_RUN" != "true" ]]; then
                            ((removal_failed_count++))
                        fi
                    fi
                fi
                ((idx++))
            done
        fi

    else
        # Start spinner for cleaning phase (small batch)
        if [[ "$DRY_RUN" != "true" && ${#existing_paths[@]} -gt 0 && -t 1 ]]; then
            MOLE_SPINNER_PREFIX="  " start_inline_spinner "Cleaning..."
            cleaning_spinner_started=true
        fi
        local idx=0
        if [[ ${#existing_paths[@]} -gt 0 ]]; then
            for path in "${existing_paths[@]}"; do
                local size_kb
                size_kb=$(get_cleanup_path_size_kb "$path")
                [[ ! "$size_kb" =~ ^[0-9]+$ ]] && size_kb=0

                local removed=0
                if [[ "$DRY_RUN" != "true" ]]; then
<<<<<<< HEAD
                    # Handle symbolic links separately (only remove the link, not the target)
                    if [[ -L "$path" ]]; then
                        rm "$path" 2> /dev/null || true
                    else
                        safe_remove "$path" true || true
                    fi
=======
                    if safe_remove "$path" true; then
                        removed=1
                    fi
                else
                    removed=1
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
                fi

                if [[ $removed -eq 1 ]]; then
                    if [[ "$size_kb" -gt 0 ]]; then
                        ((total_size_kb += size_kb))
                    fi
                    ((total_count += 1))
                    removed_any=1
                else
                    if [[ -e "$path" && "$DRY_RUN" != "true" ]]; then
                        ((removal_failed_count++))
                    fi
                fi
                ((idx++))
            done
        fi
    fi

    if [[ "$show_spinner" == "true" || "$cleaning_spinner_started" == "true" ]]; then
        stop_inline_spinner
    fi

    local permission_end=${MOLE_PERMISSION_DENIED_COUNT:-0}
    # Track permission failures in debug output (avoid noisy user warnings).
    if [[ $permission_end -gt $permission_start && $removed_any -eq 0 ]]; then
        debug_log "Permission denied while cleaning: $description"
    fi
    if [[ $removal_failed_count -gt 0 && "$DRY_RUN" != "true" ]]; then
        debug_log "Skipped $removal_failed_count items, permission denied or in use, for: $description"
    fi

    if [[ $removed_any -eq 1 ]]; then
        # Stop spinner before output
        stop_section_spinner

        local size_human=$(bytes_to_human "$((total_size_kb * 1024))")

        local label="$description"
        if [[ ${#targets[@]} -gt 1 ]]; then
            label+=" ${#targets[@]} items"
        fi

        if [[ "$DRY_RUN" == "true" ]]; then
            echo -e "  ${YELLOW}${ICON_DRY_RUN}${NC} $label${NC}, ${YELLOW}$size_human dry${NC}"

            local paths_temp=$(create_temp_file)

            idx=0
            if [[ ${#existing_paths[@]} -gt 0 ]]; then
                for path in "${existing_paths[@]}"; do
                    local size=0

                    if [[ -n "${temp_dir:-}" && -f "$temp_dir/result_${idx}" ]]; then
                        read -r size count < "$temp_dir/result_${idx}" 2> /dev/null || true
                    else
                        size=$(get_cleanup_path_size_kb "$path" 2> /dev/null || echo "0")
                    fi

                    [[ "$size" == "0" || -z "$size" ]] && {
                        ((idx++))
                        continue
                    }

                    echo "$(dirname "$path")|$size|$path" >> "$paths_temp"
                    ((idx++))
                done
            fi

            # Group dry-run paths by parent for a compact export list.
            if [[ -f "$paths_temp" && -s "$paths_temp" ]]; then
                sort -t'|' -k1,1 "$paths_temp" | awk -F'|' '
                {
                    parent = $1
                    size = $2
                    path = $3

                    parent_size[parent] += size
                    if (parent_count[parent] == 0) {
                        parent_first[parent] = path
                    }
                    parent_count[parent]++
                }
                END {
                    for (parent in parent_size) {
                        if (parent_count[parent] > 1) {
                            printf "%s|%d|%d\n", parent, parent_size[parent], parent_count[parent]
                        } else {
                            printf "%s|%d|1\n", parent_first[parent], parent_size[parent]
                        }
                    }
                }
                ' | while IFS='|' read -r display_path total_size child_count; do
                    local size_human=$(bytes_to_human "$((total_size * 1024))")
                    if [[ $child_count -gt 1 ]]; then
                        echo "$display_path  # $size_human, $child_count items" >> "$EXPORT_LIST_FILE"
                    else
                        echo "$display_path  # $size_human" >> "$EXPORT_LIST_FILE"
                    fi
                done

                rm -f "$paths_temp"
            fi
        else
            echo -e "  ${GREEN}${ICON_SUCCESS}${NC} $label${NC}, ${GREEN}$size_human${NC}"
        fi
        ((files_cleaned += total_count))
        ((total_size_cleaned += total_size_kb))
        ((total_items++))
        note_activity
    fi

    return 0
}

start_cleanup() {
    # Set current command for operation logging
    export MOLE_CURRENT_COMMAND="clean"
    log_operation_session_start "clean"

    if [[ -t 1 ]]; then
        printf '\033[2J\033[H'
    fi
    printf '\n'
<<<<<<< HEAD
    echo -e "${PURPLE}Clean Your Mac${NC}"
    echo ""

    if [[ "$DRY_RUN" != "true" && -t 0 ]]; then
        echo -e "${YELLOW}☻${NC} First time? Run ${GRAY}mo clean --dry-run${NC} first to preview changes"
    fi

    if [[ "$DRY_RUN" == "true" ]]; then
        echo -e "${YELLOW}Dry Run Mode${NC} - Preview only, no deletions"
=======
    echo -e "${PURPLE_BOLD}Clean Your Mac${NC}"
    echo ""

    if [[ "$DRY_RUN" != "true" && -t 0 ]]; then
        echo -e "${GRAY}${ICON_WARNING} Use --dry-run to preview, --whitelist to manage protected paths${NC}"
    fi

    if [[ "$DRY_RUN" == "true" ]]; then
        echo -e "${YELLOW}Dry Run Mode${NC}, Preview only, no deletions"
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
        echo ""
        SYSTEM_CLEAN=false

        ensure_user_file "$EXPORT_LIST_FILE"
        cat > "$EXPORT_LIST_FILE" << EOF
# Mole Cleanup Preview - $(date '+%Y-%m-%d %H:%M:%S')
#
# How to protect files:
# 1. Copy any path below to ~/.config/mole/whitelist
# 2. Run: mo clean --whitelist
#
# Example:
#   /Users/*/Library/Caches/com.example.app
#

EOF
        return
    fi

    if [[ -t 0 ]]; then
<<<<<<< HEAD
        echo -ne "${PURPLE}${ICON_ARROW}${NC} System caches need sudo — ${GREEN}Enter${NC} continue, ${GRAY}Space${NC} skip: "

        # Use read_key to properly handle all key inputs
        local choice
        choice=$(read_key)

        # Check for cancel (ESC or Q)
        if [[ "$choice" == "QUIT" ]]; then
            echo -e " ${GRAY}Cancelled${NC}"
            exit 0
        fi

        # Space = skip
        if [[ "$choice" == "SPACE" ]]; then
            echo -e " ${GRAY}Skipped${NC}"
            echo ""
            SYSTEM_CLEAN=false
        # Enter = yes, do system cleanup
        elif [[ "$choice" == "ENTER" ]]; then
            printf "\r\033[K" # Clear the prompt line
            if request_sudo_access "System cleanup requires admin access"; then
                SYSTEM_CLEAN=true
                echo -e "${GREEN}${ICON_SUCCESS}${NC} Admin access granted"
                echo ""
                # Start sudo keepalive with robust parent checking
                # Store parent PID to ensure keepalive exits if parent dies
                parent_pid=$$
                (
                    # Initial delay to let sudo cache stabilize after password entry
                    # This prevents immediately triggering Touch ID again
                    sleep 2

                    local retry_count=0
                    while true; do
                        # Check if parent process still exists first
                        if ! kill -0 "$parent_pid" 2> /dev/null; then
                            exit 0
                        fi

                        if ! sudo -n true 2> /dev/null; then
                            ((retry_count++))
                            if [[ $retry_count -ge 3 ]]; then
                                exit 1
                            fi
                            sleep 5
                            continue
                        fi
                        retry_count=0
                        sleep 30
                    done
                ) 2> /dev/null &
                SUDO_KEEPALIVE_PID=$!
=======
        if sudo -n true 2> /dev/null; then
            SYSTEM_CLEAN=true
            echo -e "${GREEN}${ICON_SUCCESS}${NC} Admin access already available"
            echo ""
        else
            echo -ne "${PURPLE}${ICON_ARROW}${NC} System caches need sudo. ${GREEN}Enter${NC} continue, ${GRAY}Space${NC} skip: "

            local choice
            choice=$(read_key)

            # ESC/Q aborts, Space skips, Enter enables system cleanup.
            if [[ "$choice" == "QUIT" ]]; then
                echo -e " ${GRAY}Canceled${NC}"
                exit 0
            fi

            if [[ "$choice" == "SPACE" ]]; then
                echo -e " ${GRAY}Skipped${NC}"
                echo ""
                SYSTEM_CLEAN=false
            elif [[ "$choice" == "ENTER" ]]; then
                printf "\r\033[K" # Clear the prompt line
                if ensure_sudo_session "System cleanup requires admin access"; then
                    SYSTEM_CLEAN=true
                    echo -e "${GREEN}${ICON_SUCCESS}${NC} Admin access granted"
                    echo ""
                else
                    SYSTEM_CLEAN=false
                    echo ""
                    echo -e "${YELLOW}Authentication failed${NC}, continuing with user-level cleanup"
                fi
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
            else
                SYSTEM_CLEAN=false
                echo -e " ${GRAY}Skipped${NC}"
                echo ""
            fi
<<<<<<< HEAD
        else
            # Other keys (including arrow keys) = skip, no message needed
            SYSTEM_CLEAN=false
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
        fi
    else
        echo ""
        echo "Running in non-interactive mode"
<<<<<<< HEAD
        echo "  • System-level cleanup skipped (requires interaction)"
        echo "  • User-level cleanup will proceed automatically"
        echo ""
    fi
}

# Clean Service Worker CacheStorage with domain protection

perform_cleanup() {
    echo -e "${BLUE}${ICON_ADMIN}${NC} $(detect_architecture) | Free space: $(get_free_space)"

    # Pre-check TCC permissions upfront (delegated to clean_caches module)
    check_tcc_permissions

    # Show whitelist info if patterns are active
    local active_count=${#WHITELIST_PATTERNS[@]}
    if [[ $active_count -gt 2 ]]; then
        local custom_count=$((active_count - 2))
        echo -e "${BLUE}${ICON_SUCCESS}${NC} Whitelist: $custom_count custom + 2 core patterns active"
    elif [[ $active_count -eq 2 ]]; then
        echo -e "${BLUE}${ICON_SUCCESS}${NC} Whitelist: 2 core patterns active"
=======
        if sudo -n true 2> /dev/null; then
            SYSTEM_CLEAN=true
            echo "  ${ICON_LIST} System-level cleanup enabled, sudo session active"
        else
            SYSTEM_CLEAN=false
            echo "  ${ICON_LIST} System-level cleanup skipped, requires sudo"
        fi
        echo "  ${ICON_LIST} User-level cleanup will proceed automatically"
        echo ""
    fi
}

perform_cleanup() {
    # Test mode skips expensive scans and returns minimal output.
    local test_mode_enabled=false
    if [[ "${MOLE_TEST_MODE:-0}" == "1" ]]; then
        test_mode_enabled=true
        if [[ "$DRY_RUN" == "true" ]]; then
            echo -e "${YELLOW}Dry Run Mode${NC}, Preview only, no deletions"
            echo ""
        fi
        echo -e "${GREEN}${ICON_LIST}${NC} User app cache"
        if [[ ${#WHITELIST_PATTERNS[@]} -gt 0 ]]; then
            local -a expanded_defaults
            expanded_defaults=()
            for default in "${DEFAULT_WHITELIST_PATTERNS[@]}"; do
                expanded_defaults+=("${default/#\~/$HOME}")
            done
            local has_custom=false
            for pattern in "${WHITELIST_PATTERNS[@]}"; do
                local is_default=false
                local normalized_pattern="${pattern%/}"
                for default in "${expanded_defaults[@]}"; do
                    local normalized_default="${default%/}"
                    [[ "$normalized_pattern" == "$normalized_default" ]] && is_default=true && break
                done
                [[ "$is_default" == "false" ]] && has_custom=true && break
            done
            [[ "$has_custom" == "true" ]] && echo -e "${GREEN}${ICON_SUCCESS}${NC} Protected items found"
        fi
        if [[ "$DRY_RUN" == "true" ]]; then
            echo ""
            echo "Potential space: 0.00GB"
        fi
        total_items=1
        files_cleaned=0
        total_size_cleaned=0
    fi

    if [[ "$test_mode_enabled" == "false" ]]; then
        echo -e "${BLUE}${ICON_ADMIN}${NC} $(detect_architecture) | Free space: $(get_free_space)"
    fi

    if [[ "$test_mode_enabled" == "true" ]]; then
        local summary_heading="Test mode complete"
        local -a summary_details
        summary_details=()
        summary_details+=("Test mode - no actual cleanup performed")
        print_summary_block "$summary_heading" "${summary_details[@]}"
        printf '\n'
        return 0
    fi

    # Pre-check TCC permissions to avoid mid-run prompts.
    check_tcc_permissions

    if [[ ${#WHITELIST_PATTERNS[@]} -gt 0 ]]; then
        local predefined_count=0
        local custom_count=0

        for pattern in "${WHITELIST_PATTERNS[@]}"; do
            local is_predefined=false
            for default in "${DEFAULT_WHITELIST_PATTERNS[@]}"; do
                local expanded_default="${default/#\~/$HOME}"
                if [[ "$pattern" == "$expanded_default" ]]; then
                    is_predefined=true
                    break
                fi
            done

            if [[ "$is_predefined" == "true" ]]; then
                ((predefined_count++))
            else
                ((custom_count++))
            fi
        done

        if [[ $custom_count -gt 0 || $predefined_count -gt 0 ]]; then
            local summary=""
            [[ $predefined_count -gt 0 ]] && summary+="$predefined_count core"
            [[ $custom_count -gt 0 && $predefined_count -gt 0 ]] && summary+=" + "
            [[ $custom_count -gt 0 ]] && summary+="$custom_count custom"
            summary+=" patterns active"

            echo -e "${BLUE}${ICON_SUCCESS}${NC} Whitelist: $summary"

            if [[ "$DRY_RUN" == "true" ]]; then
                for pattern in "${WHITELIST_PATTERNS[@]}"; do
                    [[ "$pattern" == "$FINDER_METADATA_SENTINEL" ]] && continue
                    echo -e "  ${GRAY}→ $pattern${NC}"
                done
            fi
        fi
    fi

    if [[ -t 1 && "$DRY_RUN" != "true" ]]; then
        local fda_status=0
        has_full_disk_access
        fda_status=$?
        if [[ $fda_status -eq 1 ]]; then
            echo ""
            echo -e "${GRAY}${ICON_WARNING}${NC} ${GRAY}Tip: Grant Full Disk Access to your terminal in System Settings for best results${NC}"
        fi
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
    fi

    total_items=0
    files_cleaned=0
    total_size_cleaned=0

    local had_errexit=0
    [[ $- == *e* ]] && had_errexit=1

    # Allow per-section failures without aborting the full run.
    set +e

    # ===== 1. Deep system cleanup (if admin) =====
    if [[ "$SYSTEM_CLEAN" == "true" ]]; then
        start_section "Deep system"
<<<<<<< HEAD
        # Deep system cleanup (delegated to clean_system module)
        clean_deep_system
=======
        clean_deep_system
        clean_local_snapshots
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
        end_section
    fi

    if [[ ${#WHITELIST_WARNINGS[@]} -gt 0 ]]; then
        echo ""
        for warning in "${WHITELIST_WARNINGS[@]}"; do
            echo -e "  ${GRAY}${ICON_WARNING}${NC} Whitelist: $warning"
        done
    fi

<<<<<<< HEAD
    # ===== 2. User essentials =====
    start_section "User essentials"
    # User essentials cleanup (delegated to clean_user_data module)
    clean_user_essentials
    end_section

    start_section "Finder metadata"
    # Finder metadata cleanup (delegated to clean_user_data module)
=======
    start_section "User essentials"
    clean_user_essentials
    scan_external_volumes
    end_section

    start_section "Finder metadata"
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
    clean_finder_metadata
    end_section

    # ===== 3. macOS system caches =====
    start_section "macOS system caches"
<<<<<<< HEAD
    # macOS system caches cleanup (delegated to clean_user_data module)
    clean_macos_system_caches
=======
    clean_macos_system_caches
    clean_recent_items
    clean_mail_downloads
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
    end_section

    # ===== 4. Sandboxed app caches =====
    start_section "Sandboxed app caches"
<<<<<<< HEAD
    # Sandboxed app caches cleanup (delegated to clean_user_data module)
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
    clean_sandboxed_app_caches
    end_section

    # ===== 5. Browsers =====
    start_section "Browsers"
<<<<<<< HEAD
    # Browser caches cleanup (delegated to clean_user_data module)
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
    clean_browsers
    end_section

    # ===== 6. Cloud storage =====
    start_section "Cloud storage"
<<<<<<< HEAD
    # Cloud storage caches cleanup (delegated to clean_user_data module)
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
    clean_cloud_storage
    end_section

    # ===== 7. Office applications =====
    start_section "Office applications"
<<<<<<< HEAD
    # Office applications cleanup (delegated to clean_user_data module)
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
    clean_office_applications
    end_section

    # ===== 8. Developer tools =====
    start_section "Developer tools"
<<<<<<< HEAD
    # Developer tools cleanup (delegated to clean_dev module)
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
    clean_developer_tools
    end_section

    # ===== 9. Development applications =====
    start_section "Development applications"
<<<<<<< HEAD
    # User GUI applications cleanup (delegated to clean_user_apps module)
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
    clean_user_gui_applications
    end_section

    # ===== 10. Virtualization tools =====
    start_section "Virtual machine tools"
<<<<<<< HEAD
    # Virtualization tools cleanup (delegated to clean_user_data module)
=======
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
    clean_virtualization_tools
    end_section

    # ===== 11. Application Support logs and caches cleanup =====
    start_section "Application Support"
<<<<<<< HEAD
    # Clean logs, Service Worker caches, Code Cache, Crashpad, stale updates, Group Containers
    clean_application_support_logs
    end_section

    # ===== 12. Orphaned app data cleanup =====
    # Only touch apps missing from scan + 60+ days inactive
    # Skip protected vendors, keep Preferences/Application Support
    start_section "Uninstalled app data"

    # Check if we have permission to access Library folders
    # Use simple ls test instead of find to avoid hanging
    local has_library_access=true
    if ! ls "$HOME/Library/Caches" > /dev/null 2>&1; then
        has_library_access=false
    fi

    if [[ "$has_library_access" == "false" ]]; then
        note_activity
        echo -e "  ${YELLOW}${ICON_WARNING}${NC} Skipped: No permission to access Library folders"
        echo -e "  ${GRAY}Tip: Grant 'Full Disk Access' to iTerm2/Terminal in System Settings${NC}"
    else

        local -r ORPHAN_AGE_THRESHOLD=60 # 60 days - good balance between safety and cleanup

        # Build list of installed application bundle identifiers
        MOLE_SPINNER_PREFIX="  " start_inline_spinner "Scanning installed applications..."
        local installed_bundles=$(create_temp_file)

        # Simplified: only scan primary locations (reduces scan time by ~70%)
        local -a search_paths=(
            "/Applications"
            "$HOME/Applications"
        )

        # Scan for .app bundles with timeout protection
        for search_path in "${search_paths[@]}"; do
            [[ -d "$search_path" ]] || continue
            while IFS= read -r app; do
                [[ -f "$app/Contents/Info.plist" ]] || continue
                bundle_id=$(defaults read "$app/Contents/Info.plist" CFBundleIdentifier 2> /dev/null || echo "")
                [[ -n "$bundle_id" ]] && echo "$bundle_id" >> "$installed_bundles"
            done < <(run_with_timeout 10 find "$search_path" -maxdepth 2 -type d -name "*.app" 2> /dev/null || true)
        done

        # Get running applications and LaunchAgents with timeout protection
        local running_apps=$(run_with_timeout 5 osascript -e 'tell application "System Events" to get bundle identifier of every application process' 2> /dev/null || echo "")
        echo "$running_apps" | tr ',' '\n' | sed -e 's/^ *//;s/ *$//' -e '/^$/d' >> "$installed_bundles"

        run_with_timeout 5 find ~/Library/LaunchAgents /Library/LaunchAgents \
            -name "*.plist" -type f 2> /dev/null | while IFS= read -r plist; do
            basename "$plist" .plist
        done >> "$installed_bundles" 2> /dev/null || true

        # Deduplicate
        sort -u "$installed_bundles" -o "$installed_bundles"

        local app_count=$(wc -l < "$installed_bundles" 2> /dev/null | tr -d ' ')
        stop_inline_spinner
        echo -e "  ${GREEN}${ICON_SUCCESS}${NC} Found $app_count active/installed apps"

        # Track statistics
        local orphaned_count=0
        local total_orphaned_kb=0

        # Check if bundle is orphaned - conservative approach
        is_orphaned() {
            local bundle_id="$1"
            local directory_path="$2"

            # Skip system-critical and protected apps
            if should_protect_data "$bundle_id"; then
                return 1
            fi

            # Check if app exists in our scan
            if grep -q "^$bundle_id$" "$installed_bundles" 2> /dev/null; then
                return 1
            fi

            # Extra check for system bundles
            case "$bundle_id" in
                com.apple.* | loginwindow | dock | systempreferences | finder | safari)
                    return 1
                    ;;
            esac

            # Skip major vendors
            case "$bundle_id" in
                com.adobe.* | com.microsoft.* | com.google.* | org.mozilla.* | com.jetbrains.* | com.docker.*)
                    return 1
                    ;;
            esac

            # Check file age - only clean if 60+ days inactive
            # Use modification time (mtime) instead of access time (atime)
            # because macOS disables atime updates by default for performance
            if [[ -e "$directory_path" ]]; then
                local last_modified_epoch=$(get_file_mtime "$directory_path")
                local current_epoch=$(date +%s)
                local days_since_modified=$(((current_epoch - last_modified_epoch) / 86400))

                if [[ $days_since_modified -lt $ORPHAN_AGE_THRESHOLD ]]; then
                    return 1
                fi
            fi

            return 0
        }

        # Unified orphaned resource scanner (caches, logs, states, webkit, HTTP, cookies)
        MOLE_SPINNER_PREFIX="  " start_inline_spinner "Scanning orphaned app resources..."

        # Define resource types to scan
        # CRITICAL: NEVER add LaunchAgents or LaunchDaemons (breaks login items/startup apps)
        local -a resource_types=(
            "$HOME/Library/Caches|Caches|com.*:org.*:net.*:io.*"
            "$HOME/Library/Logs|Logs|com.*:org.*:net.*:io.*"
            "$HOME/Library/Saved Application State|States|*.savedState"
            "$HOME/Library/WebKit|WebKit|com.*:org.*:net.*:io.*"
            "$HOME/Library/HTTPStorages|HTTP|com.*:org.*:net.*:io.*"
            "$HOME/Library/Cookies|Cookies|*.binarycookies"
        )

        orphaned_count=0

        for resource_type in "${resource_types[@]}"; do
            IFS='|' read -r base_path label patterns <<< "$resource_type"

            # Check both existence and permission to avoid hanging
            if [[ ! -d "$base_path" ]]; then
                continue
            fi

            # Quick permission check - if we can't ls the directory, skip it
            if ! ls "$base_path" > /dev/null 2>&1; then
                continue
            fi

            # Build file pattern array
            local -a file_patterns=()
            IFS=':' read -ra pattern_arr <<< "$patterns"
            for pat in "${pattern_arr[@]}"; do
                file_patterns+=("$base_path/$pat")
            done

            # Scan and clean orphaned items
            for item_path in "${file_patterns[@]}"; do
                # Use shell glob (no ls needed)
                # Limit iterations to prevent hanging on directories with too many files
                local iteration_count=0
                local max_iterations=100

                for match in $item_path; do
                    [[ -e "$match" ]] || continue

                    # Safety: limit iterations to prevent infinite loops on massive directories
                    ((iteration_count++))
                    if [[ $iteration_count -gt $max_iterations ]]; then
                        break
                    fi

                    # Extract bundle ID from filename
                    local bundle_id=$(basename "$match")
                    bundle_id="${bundle_id%.savedState}"
                    bundle_id="${bundle_id%.binarycookies}"

                    if is_orphaned "$bundle_id" "$match"; then
                        # Use timeout to prevent du from hanging on large/problematic directories
                        local size_kb
                        size_kb=$(run_with_timeout 2 du -sk "$match" 2> /dev/null | awk '{print $1}' || echo "0")
                        if [[ -z "$size_kb" || "$size_kb" == "0" ]]; then
                            continue
                        fi
                        safe_clean "$match" "Orphaned $label: $bundle_id"
                        ((orphaned_count++))
                        ((total_orphaned_kb += size_kb))
                    fi
                done
            done
        done

        stop_inline_spinner

        if [[ $orphaned_count -gt 0 ]]; then
            local orphaned_mb=$(echo "$total_orphaned_kb" | awk '{printf "%.1f", $1/1024}')
            echo "  ${GREEN}${ICON_SUCCESS}${NC} Cleaned $orphaned_count items (~${orphaned_mb}MB)"
            note_activity
        fi

        rm -f "$installed_bundles"

    fi # end of has_library_access check

    end_section

    # ===== 13. Apple Silicon optimizations =====
    if [[ "$IS_M_SERIES" == "true" ]]; then
        start_section "Apple Silicon updates"
        safe_clean /Library/Apple/usr/share/rosetta/rosetta_update_bundle "Rosetta 2 cache"
        safe_clean ~/Library/Caches/com.apple.rosetta.update "Rosetta 2 user cache"
        safe_clean ~/Library/Caches/com.apple.amp.mediasevicesd "Apple Silicon media service cache"
        # Skip: iCloud sync cache, may affect device pairing
        # safe_clean ~/Library/Caches/com.apple.bird.lsuseractivity "User activity cache"
        end_section
    fi

    # ===== 14. iOS device backups =====
    start_section "iOS device backups"
    # iOS device backups check (delegated to clean_user_data module)
    check_ios_device_backups
    end_section

    # ===== 15. Time Machine failed backups =====
    start_section "Time Machine failed backups"
    # Time Machine failed backups cleanup (delegated to clean_system module)
    clean_time_machine_failed_backups
=======
    clean_application_support_logs
    end_section

    # ===== 12. Orphaned app data cleanup (60+ days inactive, skip protected vendors) =====
    start_section "Uninstalled app data"
    clean_orphaned_app_data
    clean_orphaned_system_services
    end_section

    # ===== 13. Apple Silicon optimizations =====
    clean_apple_silicon_caches

    # ===== 14. iOS device backups =====
    start_section "iOS device backups"
    check_ios_device_backups
    end_section

    # ===== 15. Time Machine incomplete backups =====
    start_section "Time Machine incomplete backups"
    clean_time_machine_failed_backups
    end_section

    # ===== 16. Large files to review (report only) =====
    start_section "Large files to review"
    check_large_file_candidates
>>>>>>> a5c7abd2276eb9bd376e877b2068a3e4064cdc9b
    end_section

    # ===== Final summary =====
    echo ""

    local summary_heading=""
    local summary_status="success"
    if [[ "$DRY_RUN" == "true" ]]; then
        summary_heading="Dry run complete - no changes made"
    else
        summary_heading="Cleanup complete"
    fi

    local -a summary_details=()

    if [[ $total_size_cleaned -gt 0 ]]; then
        local freed_gb
        freed_gb=$(echo "$total_size_cleaned" | awk '{printf "%.2f", $1/1024/1024}')

        if [[ "$DRY_RUN" == "true" ]]; then
            local stats="Potential space: ${GREEN}${freed_gb}GB${NC}"
            [[ $files_cleaned -gt 0 ]] && stats+=" | Items: $files_cleaned"
            [[ $total_items -gt 0 ]] && stats+=" | Categories: $total_items"
            summary_details+=("$stats")

            {
                echo ""
                echo "# ============================================"
                echo "# Summary"
                echo "# ============================================"
                echo "# Potential cleanup: ${freed_gb}GB"
                echo "# Items: $files_cleaned"
                echo "# Categories: $total_items"
            } >> "$EXPORT_LIST_FILE"

            summary_details+=("Detailed file list: ${GRAY}$EXPORT_LIST_FILE${NC}")
            summary_details+=("Use ${GRAY}mo clean --whitelist${NC} to add protection rules")
        else
            local summary_line="Space freed: ${GREEN}${freed_gb}GB${NC}"

            if [[ $files_cleaned -gt 0 && $total_items -gt 0 ]]; then
                summary_line+=" | Items cleaned: $files_cleaned | Categories: $total_items"
            elif [[ $files_cleaned -gt 0 ]]; then
                summary_line+=" | Items cleaned: $files_cleaned"
            elif [[ $total_items -gt 0 ]]; then
                summary_line+=" | Categories: $total_items"
            fi

            summary_details+=("$summary_line")

            if [[ $(echo "$freed_gb" | awk '{print ($1 >= 1) ? 1 : 0}') -eq 1 ]]; then
                local movies
                movies=$(echo "$freed_gb" | awk '{printf "%.0f", $1/4.5}')
                if [[ $movies -gt 0 ]]; then
                    if [[ $movies -eq 1 ]]; then
                        summary_details+=("Equivalent to ~$movies 4K movie of storage.")
                    else
                        summary_details+=("Equivalent to ~$movies 4K movies of storage.")
                    fi
                fi
            fi

            local final_free_space=$(get_free_space)
            summary_details+=("Free space now: $final_free_space")
        fi
    else
        summary_status="info"
        if [[ "$DRY_RUN" == "true" ]]; then
            summary_details+=("No significant reclaimable space detected, system already clean.")
        else
            summary_details+=("System was already clean; no additional space freed.")
        fi
        summary_details+=("Free space now: $(get_free_space)")
    fi

    if [[ $had_errexit -eq 1 ]]; then
        set -e
    fi

    # Log session end with summary
    log_operation_session_end "clean" "$files_cleaned" "$total_size_cleaned"

    print_summary_block "$summary_heading" "${summary_details[@]}"
    printf '\n'
}

main() {
    for arg in "$@"; do
        case "$arg" in
            "--debug")
                export MO_DEBUG=1
                ;;
            "--dry-run" | "-n")
                DRY_RUN=true
                export MOLE_DRY_RUN=1
                ;;
            "--whitelist")
                source "$SCRIPT_DIR/../lib/manage/whitelist.sh"
                manage_whitelist "clean"
                exit 0
                ;;
        esac
    done

    start_cleanup
    hide_cursor
    perform_cleanup
    show_cursor
    exit 0
}

main "$@"
