# Tatoeba Parser / Indexer

This package can parse and index sentences from [Tatoeba](https://tatoeba.org).

## Supported Search Engines

* MeiliSearch
* Elasticsearch

## How to use

You need to install and run an instance of the desired search engines.

### Working with MeiliSearch

Run the following command to index in MeiliSearch:

```bash
go run . meilisearch
```

MeiliSearch accepts the following arguments:

<pre>
   --api-key          will ask you to enter the API key
   --host             host url (default: 127.0.0.1:7700)
-i --index            index name (default: tatoeba)
-d --download-files   download files needed to index Tatoeba's sentences
</pre>

### Working with Elasticsearch

Run the following command to index in Elasticsearch:

```bash
go run . elasticsearch
```

Elasticsearch accepts the following arguments:

<pre>
   --host             host url (default: 127.0.0.1:9200)
-w --workers          the number of workers. Maximum [your maximum workers available will be printed here] (default: 2)
-b --flush-bytes      the flush threshold in bytes (default: 1000000)
-i --index            index name (default: tatoeba)
-d --download-files   download files needed to index Tatoeba's sentences
</pre>

## Roadmap

- [ ] Add tests
- [ ] Add tags
- [ ] Some sentences are not indexed on MeiliSearch, find what's happening
- [x] Adding the search engine Elasticsearch

## Why this project

I planned to build a desktop application with Flutter and use the index built with this project.  
Obviously, the application will be open source on Github as well as this one. Feel free to contribute, it will be a pleasure to work with everyone :)

As soon as I have a beginning of application, I will link the project repository here! Stay tuned!

## Buy me a coffee

If you like this project, it is much appreciated :)

<a href="https://www.buymeacoffee.com/cronos87" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-red.png" alt="Buy Me A Coffee" width="217"></a>
