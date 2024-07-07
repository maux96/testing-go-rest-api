docker image rm postgres_for_go_tests
docker build ./database/ -t postgres_for_go_tests
docker run -p 54321:5432 --rm -e POSTGRES_PASSWORD=Admin123 postgres_for_go_tests:latest 