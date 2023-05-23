## firearms facts website

### GOALS

1. Show this data in a more complete format: https://bit.ly/canadian-multiple-murders
2. Be able to search and show summaries: by province, by City, same tabs as in the spreadsheet
3. Show graphs of data across time.
4. Grab web page data and store in DB.
5. Summaries of web pages by OpenAI.
6. Summarize ALL of the news articles for a particular record.
7. Store in Postgres.
8. Served by Cloudflare.
9. Import spreadsheet via CSV.
10. Use some of the new Golang patterns I've been learning.

### TODO
- [x] Add logger to config
- [x] Add some debug logging
- [x] Be able to insert fake data via curl.
- [x] Redis: Make sure to add Record to list - so we can grab groups of them.
- [x] Switch to spf13/cobra
- [x] Split out the CLI.
- [x] Be able to insert Faked data into Repository via CLI
- [ ] Had to turn off client side caching for miniredis tests - can we ONLY do that for tests? https://github.com/redis/rueidis#disable-client-side-caching
- [x] Redis Implementation of NewsStory
- [x] Link NewsStory to Record
- [x] Redis: Make sure to add NewsStory to list - so we can see groups of them.
- [x] Import via line of CSV.
- [ ] Convert config to cobra flags - pass that in?
- [ ] Need to deal with nullable booleans - there's lots of data we don't have.
- [ ] Add GetAll() to interface.
- [ ] Create GetAll() service endpoint.
- [ ] Display GetAll() as simple table - click to display individual records.
- [ ] Display Record on the Web, with stories gotten from related keys.
- [ ] Postgres Implementation of Record
- [ ] Postgres Implementation of NewsStory
- [ ] Get CRUD working for Record
- [ ] Get CRUD working for NewsStory
- [ ] Integrate News Stories with Web downloading.
- [ ] Integrate News Stories with AI Summaries.

And much, much more.