run: migrate
	go run .

migrate: clean
	go run github.com/prisma/prisma-client-go db push
	go run github.com/prisma/prisma-client-go generate

clean:
	sudo -H -u postgres bash -c 'psql -U postgres -c "DROP DATABASE dev;"' 
