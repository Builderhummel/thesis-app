# Build and run Docker image
## Build
`docker build -t thesis-app -f docker/Dockerfile .`

## Run
    docker run -e DB_IP='127.0.0.1' \
     -e DB_PORT='3306' \
     -e DB_USERNAME='username' \
     -e DB_PASSWORD='password' \
     -e DB_NAME='thesis-app' \
     -e LDAP_URL='ldaps://ldap.example.com:636' \
     -e LDAP_DN='ou=people,dc=tu-example,dc=de' \
     -e JWT_SECRET='e7bb....' \
     thesis-app

Exchange the JWT Secret to a valid one!!!

# docker-compose
## File
