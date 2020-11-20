# Tatoeba Parser / Indexer

This package can parse and index sentences from [Tatoeba](https://tatoeba.org).

It only supports [MeiliSearch](https://www.meilisearch.com) as search-engine at the moment.

## How to use

You need to install and run an instance of MeiliSearch.

Then you need to download the archives from Tatoeba. You can use the following commands:

```bash
wget https://downloads.tatoeba.org/exports/sentences_detailed.tar.bz2
```

```bash
wget https://downloads.tatoeba.org/exports/links.tar.bz2
```

```bash
wget https://downloads.tatoeba.org/exports/sentences_with_audio.tar.bz2
```

Extract the archives like this:

```bash
tar xjf sentences_detailed.tar.bz2 && tar xjf links.tar.bz2 && tar xjf sentences_with_audio.tar.bz2
```

If you have an instance of MeiliSearch running on your local, simple run:

```bash
go run .
```

In case it's on a remote server, you can specify the host passing the parameter `-host`.

## Roadmap

* Add tests
* Add tags
