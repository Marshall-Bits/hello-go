# Hello Go
This is a simple Go program that creates a server and listens to the port 8008.

The server has two endpoints:
- `/`: returns a simple message.
- `/all-robots`: connects to a mongo database and returns all the robots in the collection.

## How to run
1. Clone the repository.
2. Install golang if you haven't already. You can download it from [here](https://golang.org/dl/).
3. Run the following command in the root of the project:
```bash
go run main.go
```
4. Open your browser and go to `http://localhost:8008/` to see the message.

*Keep in mind that you need to restart the terminal after installing golang to be able to run the `go` command.

## How to test
Add your mongo connection string in a .env file in the root of the project:
```bash
MONGO_URI=your_connection_string/robots
```

Here's a few robots you can add to your database to test the route:
```json
{
    "name": "R2D2",
    "type": "Astromech droid",
    "mass": 32
}
{
    "name": "C-3PO",
    "type": "Protocol droid",
    "mass": 75
}
{
    "name": "BB-8",
    "type": "Astromech droid",
    "mass": 18
}
```	

Now you can test the `/all-robots` endpoint by going to `http://localhost:8008/all-robots`.

