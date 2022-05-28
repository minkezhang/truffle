# tracker
Project to track queue of media that needs to be consumed.

## Examples
```bash
go install github.com/minkezhang/truffle/truffle

# Add sample entry.
truffle add \
  --titles=Sabikui \
  --corpus=anime \
  --score=6.3 \
  --providers=crunchyroll \
  --queued=true \
  --studios=OZ \
  --directors="Atsushi Itagaki" \
  --writers="Sadayuki Murai" \
  --season=1 \
  --episode=4

# Start watching the next season of Sabikui. rt supports partial title and
# corpus matching.
truffle bump \
  --title=Sabikui \
  --corpus=anime \
  --major

# Re-rate the entry.
truffle patch \
  --title=Sabikui \
  --corpus=anime \
  --score=6.4

truffle get --title=Sabikui

# Search the user database as well as the MAL API for similar entries.
truffle search \
  --title=Sabikui \
  --corpus=anime \
  --trackers=truffle\
  --trackers=mal

# Delete the entry.
truffle delete --title=Sabikui

```

## Uninstall

```bash
go clean -i github.com/minkezhang/truffle/truffle
```

## Proposed API

```
* Filter(title_asc, corpus_asc, ...): return subset of collection
* Recommend()
```
