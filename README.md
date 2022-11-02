# go-scraper
![build](https://github.com/CalebTracey/go-scraper/actions/workflows/build.yml/badge.svg)
![test](https://github.com/CalebTracey/go-scraper/actions/workflows/test.yml/badge.svg)


Web scraper written in Go. Intended to collect business data based on location.

[go-colly](https://github.com/gocolly/colly) used to scrape yellowpages for basic data.

Google geocode finds the longitude and latitude with scraped address. 

### Current available search parameters:
  * <i>Search terms - "terms"*</i>
  * <i>"city"*</i>
  * <i>"state"*</i>
  
 <i>*Required</i>

### Example request body:
<img src=./docs/req-body.png  alt="" width="400"/>

### Run configuration:
<br/>
<img src=./docs/run-conf.png  alt="" width="400"/>
