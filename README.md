# truffle
Project to track queue of media that needs to be consumed.

## Examples
```bash
go install github.com/minkezhang/truffle/truffle@latest

# Add sample entry.
truffle add \
  --title=Sabikui \
  --corpus=anime \
  --score=6.3 \
  --provider=crunchyroll \
  --queued=true \
  --studio=OZ \
  --director="Atsushi Itagaki" \
  --writer="Sadayuki Murai" \
  --season=1 \
  --episode=4

# Start watching the next season of Sabikui. truffle supports partial title and
# corpus matching.
truffle bump \
  --title=Sabikui \
  --corpus=anime \
  --major

# Search the user database as well as the MAL API for similar entries.
truffle search \
  --title=Sabikui \
  --corpus=anime \
  --api=truffle\
  --api=mal

# Inspect the MAL entry directly.
truffle get --id=mal:anime/48414

# Re-rate the entry and mark the MAL entry as duplicate, which will be filtered
# out in searches.
truffle patch \
  --title=Sabikui \
  --corpus=anime \
  --score=6.4 \
  --link=mal:anime/48414

# Search for DB data merged with any linked duplicates.
truffle get --title=Sabikui

# Delete the entry.
truffle delete --title=Sabikui

```

## Uninstall

```bash
go clean -i github.com/minkezhang/truffle/truffle
```