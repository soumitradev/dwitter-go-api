# Dwitter Go API

I'm trying to make a Twitter clone using Go, GraphQL, and Prisma with a Postgresql backend.

I use `make` to make my workflow easier (no pun intended), the commands are:

`make run` : Run only the API. Don't migrate databases.

`make migrate` : Run through only the deletion and migration of the database. This will delete the database and make a fresh database with the schema given.

`make clean` : Run through only the deletion of the database. This will delete the database.

`make kill` : Kill any processes that are currently accessing the database. This is only to be used if an error is raised when running `make clean`

The prisma schema is [here](./prisma/schema.prisma)

**NOTE:** You need `ffmpeg` added to your path to run this. It uses ffmpeg to generate thumbnails for videos uploaded.

**TODO:**
- Add a weekly vacuum cron job

The entire API is now ready. I haven't tested for every scenario, because it would take too much time to do it, so I'll fix any bugs I find in the future when working on the frontend.

I'm very happy with how this turned out. I learned a **lot**.

This is by far my biggest project ever. Yes. Including BruhOS, including the chip8 hardware emulation project I did, even the fullstack Node project I made last year.

This was really fucking fun, especially working with my homie [@PseudoCodeNerd](https://github.com/PseudoCodeNerd).

His support and enthusiasm kept me going all through this. He helped me debug some really tough logic issues too. I wouldn't have made any of this without his help and inspiration.

I literally can't imagine myself building the frontend without him. He's really helpful and knowledgeable.

This (hopefully) will be integrated into a bigger Dwitter project later.
