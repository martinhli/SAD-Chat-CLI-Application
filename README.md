# SAD-Chat-CLI-Application

To use the Chat CLI application, the user needs to have Go, Docker and NATS installed.

First, the user needs to ensure that a NATS server is running. To run the NATS server, run the following command
```
docker run -d --name nats-server -p 4222:4222 -p 8222:8222 -v "$(pwd)/nats-server.conf:/etc/nats/nats-server.conf" nats:latest -c /etc/nats/nats-server.conf
```

Run the executable chatcli with the following command:\
```
./chatcli nats://localhost:4222 chatroom YourName
```
chatroom can be changed to whatever name of the channel the user want, and YourName can be changed to whatever username.

First the program will connect to the NATS server at nats://localhost:4222, and display all messages from the last hour. Then, the program will prompt the user to input a message:
`Enter an input: `. The user can input a message, and it will show up with their username. The program will fetch and display all messages in the chat, including new messages. After 5 seconds, if there are no new messages, the program will timeout and the user can input a new message. To disconnect, the user can type `Cmd + C` which will shut down the program gracefully.
