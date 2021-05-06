# Dwitter Go API

I'm trying to make a Twitter clone using Go, GraphQL, and Prisma with a Postgresql backend.

I use `make` to make my workflow easier (no pun intended), the commands are:

`make run` : Go through the entire workflow (no pun intended). This will delete the database, make a fresh database with the schema, and run the code on this database.

`make migrate` : Run through only the deletion and migration of the database. This will delete the database and make a fresh database with the schema.

`make clean` : Run through only the deletion of the database. This will delete the database.

The prisma schema is [here](./prisma/schema.prisma)

The GraphQL part is still to be implemented, and this will probably have a frontend built with Vue and Tailwind.

The abstraction to the database operations with Prisma are not set in stone. This is my first working attempt at them, and further optimizations might be possible.

This (hopefully) will be integrated into a bigger Dwitter project later.