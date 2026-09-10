# Description

This project lets you browse different artists and get information about
their different members, concerts, albums...

# Authors

LEMATELOT Amélie
FEDAOUI Yohann

# Usage

## Server setup

`go run .`

## Webapp interface

For now, artists are displayed in the central frame and you click on each
to get more detailed info
Later, functionality will be added to search artists by name, album name,
concert location, etc.
Just type text in the input text area and it will suggest possible artists
that match.

# Implementation details

The server's API endpoints replicate the API endpoints at heroku
(and add/translate some data where needed).
Some JS handles API requests but it will be rolled back as templates can do
what the JS is doing.

# Potential improvements

- Search bar
- Caching info from the API
