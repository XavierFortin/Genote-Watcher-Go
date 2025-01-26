# GenoteWatcher-Go

Small app that sends notifications to a discord webhook when a new note is added
or changed on genote

The original app is written in typescript but this one has been rewritten in go
to make things faster and easier to share.

## Requirements

- Create a .env file at the root with the following keys:
  - **GENOTE_USER** : Contains your UdS email to login into genote
  - **GENOTE_PASSWORD** : Contains your UdS password to login into genote
  - **DISCORD_WEBHOOK** : Your desired Discord webhook url
  - **TIME_INTERVAL** : A time interval used as a replacement for a cron job.
    - Valid formats are `ms`, `s`, `m` and `h`
    - Example of possible time intervals:
      - `0`: Runs only once
      - `300ms`: Runs every 300ms
      - `2h45m`: Runs every 2h and 45 min
    - Do not put the interval too short. There is a possibility that your IP
      gets blocked if it is spamming too fast.

## Run the standalone executable
- Download the executable
- Create a .env file next to the executable in the same folder with the required informations
- You can specify a port or use the default one which is `4000` 
### Windows
- Run `./genote-watcher.exe --port <PORT>` 
### Linux
- Run `./genote-watcher --port <PORT>`

## Run with Docker

- Get the Image from `docker pull enox/genote-watcher`
- Run the container and make sure that the 3 env variables are set by either:
  - Running `docker run --env-file <env_file_name> enox/genote-watcher:latest --port <PORT>`
  - Running
    `docker run -e <env_name1>=<env_value1> <env_nameX>=<env_valueX> enox/genote-watcher:latest --port <PORT>`
  - Adding the environment variables to a docker-compose
  - If you need to restart the container you can run
    `docker start <name_of_container>`. It is important to start an already
    started container so it can track changes over time. If a new container is
    created, it will not work

## Access the dashboard

- You can go in your web browser with your specified port (default: 4000) to view the dashboard
- ![image](https://github.com/user-attachments/assets/b5302d39-58f8-44c8-914d-a5cb07969642)

## Build the app from scratch

- You will need to have go installed
- Run the build.ps1 script. It will build a windows and a linux/amd64 executable
