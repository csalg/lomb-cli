# Lomb 

### UNIX Philosophy

Out of the projects that I have worked on in my spare time, `me` has been quite successful and Lomb has been largely a disaster. Why? This is not easy to answer.

Maybe I could even run Lomb it might not be so bad.

I think the problem is largely that everything is tangled all together. It's also super annoying that when I want to do anything I have to make changes to the codebase.

I think a unix philosophy lomb would make more sense.

In fact, the part that works is the `lomb-preprocessor` script.

It would be nice to have commands and just work on files. That would simplify things a lot.

Books could just live in `/var/lomb/books`. Logs could just live in `/var/lomb/logs`. Subtitles could live in `/var/lomb/subtitles`. And so on.

This would work great. Lomb preprocessor could be instructed to send the book there when it's done.

I could then have
```
lomb revise
lomb revise -s <source language> -t <target language> -w <list of words as text file>`
lomb revise -s <source language> -t <target language> -m <model> -c <list of words as csv>

lomb read // continues reading last text
lomb read --list // lists texts to fzf

lomb watch <video file> // matches subtitles automatically and opens browser
```

And so on. This would make it very easy to add new features to lomb.

The only thing these scripts really have in common is the `/var/lomb` directory and perhaps some types. Other than that they are largely independent.

Perhaps I should start with lomb revise, since I kind of need it.

### `lomb revise`

For now, I just need some easy way to revise the words I come across in Danish books.
