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
