# Lomb 

### UNIX Philosophy

Out of the projects that I have worked on in my spare time, `me` has been quite successful and Lomb has been largely a disaster. Why? This is not easy to answer.

Maybe if I could even run Lomb it might not be so bad. So the distribution story of Python is one thing.

I think the problem is largely that everything is tangled all together. It's also super annoying that when I want to do anything I have to make changes to the codebase. I also have this polished UI that I obviously need to maintain.

I think a unix philosophy `lomb` would make more sense. It would be a collection of scripts that have a dependency on a folder where logs, books, etc. are kept.

In fact, the part that works is the `lomb-preprocessor` script.

It would be nice to have commands and just work on files. That would simplify things a lot.

Books could just live in `/var/lomb/books`. Logs could just live in `/var/lomb/logs`. Subtitles could live in `/var/lomb/subtitles`. And so on.

This would work great. Lomb preprocessor could be instructed to send the book there when it's done.

I could then have
```
lomb revise
lomb revise -s <source language> -t <target language> -w <list of words as text file>`
lomb revise -s <source language> -t <target language> -c <list of words as csv>

lomb read // continues reading last text
lomb read --list // lists texts to fzf

lomb watch <video file> // matches subtitles automatically and opens browser
```

And so on. This would make it very easy to add new features to lomb.

The only thing these scripts really have in common is the `/var/lomb` directory and perhaps some types. Other than that they are largely independent.

Perhaps I should start with lomb revise, since I kind of need it.

The most annoying thing to implement will be some sort of service that returns chunks for a lemma. It will likely have to be a `docker-compose` thing with `mongodb` and a `go` frontend.

### `lomb revise`

For now, I just need some easy way to revise the words I come across in Danish books. They don't have to be ranked or anything. 

In the future, there should be a log of the events.

After that, it should support reading a csv file with two columns: the lemma and the PoR.

### `lomb-score`

This will be a python script responsible for doing ETL, training models, scoring words, etc.

```
lomb-score etl
lomb-score train
lomb-score score
```

### `lomb read`

Basically load the reader with some html.

By convention, all texts will end in `<source lang>.<target lang>.html`.

```
lomb read --list
lomb read <file>
```

### lomb watch

Launches the video page with some video. The video must end in `<source lang>.mp4`. Then there must be subtitles in `/var/lomb/subtitles`.

`lomb watch <video file>`

### misc

Obviously `/var/lomb` and other folders would be backed up via syncthing.

It would be trivial to have lomb preprocessor dump all translated subtitle files, books, etc. in there.

I could also make tarballs of the thing and back these up (e.g. using `borg`).

### `go` vs `python`

`python` has a horrible distribution story. I think it makes sense to ONLY use python for training models and scoring. I can also see having a service eventually where a csv file gets sent and it comes back scored.

This should hopefully make the whole thing more resilient. Also, I can easily see the use of becoming very good at `go` given my current career. I can see the use of having good skeletons and documentation of the tools I use at work. Then I can piggyback off that for this project.

### Ideas for reverse indexes

So the idea is to avoid the use of mongodb.

To be honest, if I only have, say ~1000 files, scanning those files is not super traumatic (maybe ~10 minutes?). In practice those ~1000 files _are_ the database, and anything else is ETL.

The use case of having a reverse index is to look for examples of words I am studying. So if we are going to be efficient about this, keeping a reverse index of all the words is unnecessary. 

Anyway, so all I really need to do is have a command that reads all the files and creates a reverse index somewhere. It can / should be stored in `sqlite`.

`lomb index`

As a naive approach:
- Keep a table with all the files that have been indexed and when they were last indexed.
- Keep a table with the chunks (both langs and reference to source file)
- Need to index again? Do it in memory, then nuke db tables and insert.

If that takes too long, can look into alternatives like keeping a list of which books have been read since last index and just look at those. Frankly, this can just be inferred from the log.

### Cool features of Lomb

1. Drill books before reading them
2. Revision lemmas have frequency based on corpus
3. Probability of forgetting
4. Reverse index (lemma to chunk) of revision lemmas
5. Ignore lemmas automatically (if they are scrolled)
6. Video player
7. Proxy server that cleans up dictionary sites
8. Text reader

(1) is kinda trivial. Just build a reverse index of the file and serve it using the revision view. Cache it so it's fast (use `gob` so it's in binary).
(2) this is also trivial, it's the count of the chunks.
(3) this is NOT trivial at all. I guess I could use Python for this one. However I would need to calculate the PoR and share it with go somehow. I think it would need to be a `systemctl` service
(4) how to compute it is easy (e.g. use `goquery` and look at the `data` attributes). Storing it... I guess sqlite with an index.
(5) keep a list of ignored lemmas which gets persisted on blacklisting new ones.
(6) i already have this so I just need to serve it from `go` and log.
(7) also trivial. i could add a `blacklisted_css` field to the `config.json` `dictionaries` entries.
(8) also trivial.

I think with `go` and `gob` I sort of get to implement my own database.
