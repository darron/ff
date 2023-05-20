## firearms facts website

### GOALS

Show this data in a more complete format: https://bit.ly/canadian-multiple-murders
Be able to search and show summaries: by province, by City, same tabs as in the spreadsheet
Show graphs of data across time.
Grab web page data and store in DB.
Summaries of web pages by OpenAI.
Summarize ALL of the news articles for a particular record.
Store in Postgres.
Served by Cloudflare.
Import spreadsheet via CSV.

### TODO
[ ] Had to turn off client side caching for miniredis tests - can we ONLY do that for tests?
[ ] Redis Implementation of NewsStory
[ ] Link NewsStory to Record
[ ] Redis: Make sure to add Record to list - so we can grab groups of them.
[ ] Redis: Make sure to add NewsStory to list - so we can groups of them.
[ ] Display Record on the Web
[ ] Import via line of CSV.
[ ] Postgres Implementation of Record
[ ] Postgres Implementation of NewsStory
[ ] Get CRUD working for Record
[ ] Get CRUD working for NewsStory
[ ] Integrate News Stories with Web downloading.
[ ] Integrate News Stories with AI Summaries.

And much, much more.