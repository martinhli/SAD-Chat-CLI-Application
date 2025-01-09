# SAD-Chat-CLI-Application

To use the Chat CLI application, just run the executable chatcli with the command:
`./chatcli nats://localhost:4222 chatroom YourName`, where chatroom can be changed to whatever name of the channel the user want, and YourName can be changed to whatever username.

First the program will connect to the NATS server at nats://localhost:4222, and display all messages from the last hour. Then, the program will prompt the user to input a message:
`Enter an input: `. The user can input a message, and it will show up with their username. To disconnect, the user can disconnect gracefully by typing `Cmd + C`.
