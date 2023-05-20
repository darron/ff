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
- [ ] Had to turn off client side caching for miniredis tests - can we ONLY do that for tests?
- [ ] Redis Implementation of NewsStory
- [ ] Link NewsStory to Record
- [ ] Redis: Make sure to add Record to list - so we can grab groups of them.
- [ ] Redis: Make sure to add NewsStory to list - so we can groups of them.
- [ ] Display Record on the Web
- [ ] Import via line of CSV.
- [ ] Postgres Implementation of Record
- [ ] Postgres Implementation of NewsStory
- [ ] Get CRUD working for Record
- [ ] Get CRUD working for NewsStory
- [ ] Integrate News Stories with Web downloading.
- [ ] Integrate News Stories with AI Summaries.

And much, much more.