# card-cli

A command-line tool for fetching, decoding, and modifying V2/V3 character cards (chara cards). Character cards are PNG
images that contain embedded JSON metadata defining AI character personalities, backgrounds, and behaviors.

## Features

- **Fetch cards** from multiple supported platforms
- **Decode cards** to extract and view JSON metadata
- **Inject JSON** to modify card metadata
- **Batch operations** for processing multiple cards at once
- **Template-based naming** for organized card storage
- **Pretty printing** with progress bars and styled output

## Installation

```bash
go install github.com/r3dpixel/card-cli@latest
```

Or clone the repository and build locally:

```bash
git clone https://github.com/r3dpixel/card-cli.git
cd card-cli
go build -o card-cli main.go
```

## Supported Sources

The tool currently supports fetching from the following platforms:

- **CharacterTavern** - Character sharing platform
- **ChubAI** - AI character hub
- **NyaiMe** - Character card repository
- **PepHop** - Character chat platform
- **WyvernChat** - AI chat service
- **Pygmalion** - Character AI service
- **JannyAI** - AI character platform

Check available sources and their status:

```bash
card-cli sources
```

## Commands

### fetch - Download character cards

Fetch cards from one or more URLs and save them locally.

```bash
# Basic usage - fetch a single card
card-cli fetch https://characterhub.org/characters/example/card-name

# Fetch multiple cards
card-cli fetch https://url1.com/card1 https://url2.com/card2 https://url3.com/card3

# Specify output directory
card-cli fetch -o ./my-cards https://example.com/card

# Custom file naming format
card-cli fetch -f "{{NAME}}_{{DATE}}" https://example.com/card
```

**File naming tokens:**

- `{{SOURCE}}` - Platform source (e.g., ChubAI)
- `{{PLATFORM_ID}}` - Unique platform identifier
- `{{NAME}}` - Character name
- `{{DATE}}` - Download date
- `{{VERSION}}` - Card version

### decode - Extract JSON metadata

Decode a character card to view or save its JSON metadata.

```bash
# Output JSON to stdout
card-cli decode character.png

# Save to file
card-cli decode -o metadata.json character.png

# Pretty print with indentation
card-cli decode -p character.png

# Pretty print with stable key sorting
card-cli decode -p -s character.png
```

### inject - Modify card metadata

Replace the JSON metadata in an existing character card.

```bash
# Inject new metadata into a card
card-cli inject character.png new_metadata.json
```

**Note:** This modifies the original card file in-place.

## Examples

### Example 1: Bulk download from multiple sources

```bash
# Create a list of URLs
URLS=(
  "https://characterhub.org/characters/user1/alice"
  "https://characterhub.org/characters/user2/bob"
  "https://chub.ai/characters/creator/charlie"
)

# Fetch all with organized naming
card-cli fetch -o ./collection -f "{{SOURCE}}/{{NAME}}_{{PLATFORM_ID}}" "${URLS[@]}"
```

### Example 2: Card metadata workflow

```bash
# 1. Download a card
card-cli fetch https://example.com/character/alice

# 2. Extract metadata for editing
card-cli decode -o alice_meta.json -p alice.png

# 3. Edit alice_meta.json with your preferred editor
vim alice_meta.json

# 4. Inject modified metadata back
card-cli inject alice.png alice_meta.json
```

### Example 3: Batch processing with shell scripting

```bash
#!/bin/bash
# Process all PNG files in current directory

for card in *.png; do
  echo "Processing: $card"

  # Decode to JSON
  card-cli decode -o "${card%.png}.json" -p "$card"

  # Optional: modify JSON files here

  # Re-inject if needed
  # card-cli inject "$card" "${card%.png}.json"
done
```

### Example 4: Validate card integrity

```bash
# Decode a card to check if it contains valid metadata
if card-cli decode suspicious_card.png > /dev/null 2>&1; then
  echo "Card is valid"
else
  echo "Card is corrupted or invalid"
fi
```

## Output Format

The tool provides colored output with progress bars for batch operations:

- **Green**: Successful operations
- **Red**: Failed operations
- **Yellow**: Warnings or invalid URLs
- **Blue**: Information and statistics

## Error Handling

The tool handles various error conditions:

- Invalid URLs are reported separately
- Failed downloads are tracked and reported
- Corrupted cards are identified during decode operations
- File permission errors are handled gracefully