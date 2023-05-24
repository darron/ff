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
- [x] Add GetAll() to interface.
- [x] Create GetAll() service endpoint.
- [x] Display GetAll() as simple table
- [x] Move JSON enpoints to `/api/v1`
- [x] From / click to display individual records
- [x] Display Record on the Web, with stories gotten from related keys.
- [x] Remove flags from config - just set default config.
- [ ] Add flags as neccessary to cmd - then run the required `WithOpt` config func.
- [ ] Better HTML templates - some sort of CSS framework.
- [ ] Includes for HTML headers and things?
- [ ] Need to deal with nullable booleans - there's lots of data we don't have.
- [ ] Make it actually be organized.
- [ ] Postgres/SQLite Implementation of Record
- [ ] Postgres/SQLite Implementation of NewsStory
- [ ] Get CRUD working for Record
- [ ] Get CRUD working for NewsStory
- [ ] Integrate News Stories with Web downloading.
- [ ] Integrate News Stories with AI Summaries.

And much, much more.