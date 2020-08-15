builddockerimages:
	docker build -t=donutloop/prices:latest -f=Dockerfile.Prices .
	docker build -t=donutloop/ranks:latest -f=Dockerfile.Ranks .
	docker build -t=donutloop/topcoins:latest -f=Dockerfile.Topcoins .
