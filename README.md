# Dwitter Go API

I'm trying to make a Twitter clone using Go, GraphQL, and Prisma with a Postgresql backend.

I use `make` to make my workflow easier (no pun intended), the commands are:

`make run` : Run only the API. Don't migrate databases.

`make migrate` : Run through only migration of the database. If DB cannot be migrated, you will need to delete it.

`make clean` : Run through only the deletion of the database. This will delete the database.

`make kill` : Kill any processes that are currently accessing the database. This is only to be used if an error is raised when running `make clean`

The prisma schema is [here](./prisma/schema.prisma)

**NOTE:** You need `ffmpeg` added to your path to run this. It uses ffmpeg to generate thumbnails for videos uploaded.

> .env contains ACCESS_SECRET, REFRESH_SECRET, DISCORD_CLIENT_ID, DISCORD_CLIENT_SECRET

> cdn_key.json is the key to Google Firebase

**TODO:**
- Add more OAuth providers
- If username/email is already taken, allow them to set a new one
- If user has already signed up, log them in when using OAuth flow
- Advanced Search for Tweets
- Frontend in Vue

**TODO: (Advanced)**
- Add support for becoming an OAuth identity provider

The entire API is now ready. I haven't tested for every scenario, because it would take too much time to do it, so I'll fix any bugs I find in the future when working on the frontend.

I'm very happy with how this turned out. I learned a **lot**.

This is by far my biggest project ever. Yes. Including BruhOS, including the chip8 hardware emulation project I did, even the fullstack Node project I made last year.

This was really fucking fun, especially working with my homie [@PseudoCodeNerd / @Madhav Sharma / @mdvsh](https://github.com/mdvsh).

His support and enthusiasm kept me going all through this. He helped me debug some really tough logic issues too. I wouldn't have made any of this without his help and inspiration.

I literally can't imagine myself building the frontend without him. He's really helpful and knowledgeable.

This (hopefully) will be integrated into a bigger Dwitter project later.

## Update

So um... College/Uni happened. I am pretty busy now, but I really want to complete this project some other time.

I think its more worthwhile to work on smaller projects and keep my practice. I'm already starting to get a bit rusty on the OAuth flow, and the Access/Refresh Token flow after a 5 month break from programming.

I am working on another project for now, but the progress has been incredibly slow due to me being busy.

It makes me sad how I don't really have time to do all the stuff I enjoy so dearly anymore.

Or well, I'll have to learn to make time, and it might take me anywhere from a week to a month of getting used to my new schedule.

Things are changing. Things are changing fast, and I'm starting to feel old for some reason. I'm 18. I have no reason to feel this way.

PseudoCodeNerd is now mdvsh, I don't get to talk to him as much anymore, both of us are stuck in our own work, I'm doing new stuff now, I switched back to Windows from Linux for school, I'm done with Grade 12, I need to keep a good GPA, I need to learn all kinds of new stuff.

Man.

Future me, if you're reading this, don't let time pass you by like I did. 5 months down the drain and no code at all. Wasn't even a voluntary break.

I hope I find time.

I hope I can start working on the Vue frontend for this soon. <3

Sorry if I went a bit off-topic on this README, but I felt like it was necessary. My code usually contains bits about my life anyways. See you hopefully in the next commit (Hopefully at most 1 month from now).
