In order to use this program, you need Postgres and Go
Goose would be helpful for migration tracking
Gator is for aggregating specific RSS feeds and creating posts with that data

Setup includes either creating a personal server via something like psql or something comparable
then creating a config file, I used the home directory with json parsing, with two needed fields 
example:
{"db_url":"postgres://username:passwword@localhost:5432/gator?sslmode=disable", "current_user_name":"username"}

dburl is used to identify where the server and database are located, and the current user name is used to determine who is CURRENTLY logged in