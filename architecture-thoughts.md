# Why build this?

I think there is a need for a language learning tool that is realia-oriented.


# Architecture thoughts

- It should be possible to get a lot of functionality out of `lomb` as command-line tool. Lemmatization can be done with lemmatization lists; instead of training models I can use a simple SRS model, logging can be done to sqlite. All of that can be done in `go`. And really if it fulfils most of my needs, I might want to stop there.
- Ok, this is the stuff I'd need Python for:
    - Lemmatization of non-Western languages (jieba for Chinese). Since this is only needed for the processing step, I could write a simple wrapper and call it from `go` directly.
    - Lemmatization using spacy. Same as above.
    - Model training. Actually, this could also be done from `go` because it's really just linear regression under the hood.
- Stuff that can be done in `go`.
    - Event logging.
    - Look at revision events and apply SRS.
    - Lemmatization using github.com/aaaton/golem/
    - Serving simple versions of revise for a text file or many
    - Reader, video viewer

So, come to think of it, the whole thing can really be done in `go` except for some lemmatization.

Why `go` and not `python`?
- Become better at one thing vs mediocre at two.
- Can copy paste code from `raffle`, e.g. gorm idioms.
- `go` is built for good software development practices, whereas `python` has tooling on top to achieve this. So it is simpler to develop something maintainable with `go`.
