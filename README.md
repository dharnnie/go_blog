# go_blog
Please help see coding practices. Thank you!!!

How code is organised

File : main.go
	Of course this is the main file... It contains everything a golang main file should have. Actually Just the routes.


Folder/Package : database

Contains two files
   - database.go
   - database2.go

   Basically I simply seprated the two files so I don't have 5 million lines of code in one file. Thankfully golang handles that gracefully...

   * They contain all the functions that speak directly to the database.

Folder/Package : handlers
	- handler.go

	Contains functions that binds to my routes. Functions here use funcs from [package > database]...

Folder/Package : templates
	Contains everything that may end with .html , .css , and .js . And Oh!! There are no JavaScript codes here.

	WARNING - The interfaces suck - Be prepared!!!
