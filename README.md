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
- [x] Add flags as neccessary to cmd - then run the required `WithOpt` config func.
- [x] Allow actual files to be downloaded if there's no route.
- [x] Better HTML templates - some sort of CSS framework.
- [x] Prometheus endpoint for /metrics
- [x] Includes for HTML headers and things?
- [x] Need to deal with nullable booleans - there's lots of data we don't have.
- [x] e.Use(middleware.RequestID())
- [x] Super hacky: Be able to see different types of Records: firearms only, mass shooting only, OIC firearms only, licensed mass shootings
- [x] Add totals to group pages?
- [x] Add tests for Group pages - get test coverage back to +80%
- [x] ACLs for /api/ endpoints - protected by JWT requirement.
- [x] Dockerfile so we can run it in Docker.
- [x] Docker Compose file
- [x] Add autoTLS - https://echo.labstack.com/cookbook/auto-tls/
- [x] Do we need a page cache? https://github.com/victorspringer/http-cache - added https://github.com/SporkHubr/echo-http-cache
- [x] Too many Redis connections - need to close connections as much as possible.
- [x] Import with TLS - needs to handle new ports and domain name and HTTPS
- [x] For Records and News Stories - more efficient way to download them all at once when using Redis. Right now we're opening and closing a new connection for each one.
- [x] Be able to enable / disable traces and profiling using flags.
- [ ] Some sort of Analytics
- [ ] Better Groups/Tags interface.
- [ ] Add groups by Province and City?
- [ ] Make it actually be organized.
- [ ] Postgres/SQLite Implementation of Record
- [ ] Postgres/SQLite Implementation of NewsStory
- [ ] Will need to do some more performance tuning once we have more traffic.
- [ ] Monitoring for uptime.
- [ ] skaffold + k8s files?
- [ ] Need some additional JWT work: actually check the claims - how to get them generated and asssigned to me?
- [ ] Do I need to add some contexts to track requests?
- [ ] Add /healthz which tests for health of DB?
- [ ] Deploy to domain name?
- [ ] Put Behind Cloudflare?
- [ ] Integrate StatsCan homicide records.
- [ ] Get CRUD working for Record
- [ ] Get CRUD working for NewsStory
- [ ] Add Search?
- [ ] Integrate News Stories with Web downloading.
- [ ] Integrate News Stories with AI Summaries.

And much, much more.