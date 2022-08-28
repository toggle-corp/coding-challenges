#/bin/bash
set -e

# keep track of the last executed command
trap 'last_command=$current_command; current_command=$BASH_COMMAND' DEBUG
# echo an error message before exiting
trap 'echo "`date`:: \"${last_command}\" command filed with exit code $?." > ~/deploy-error.log' EXIT

cd ~/coding-challenges
git pull
docker-compose stop
docker-compose pull
docker-compose up
