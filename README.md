### REST server written in go
Made for a simple task managing app

## TODO
[x] POST   /task/              :  create a task, returns ID
[x] GET    /task/<taskid>      :  returns a single task by ID
[x] GET    /task/              :  returns all tasks
[x] DELETE /task/<taskid>      :  delete a task by ID
[] GET     /tag/<tagname>      :  returns list of tasks with this tag
[] GET     /due/<yy>/<mm>/<dd> :  returns list of tasks due by this date
[x] POST   /login/             : verifes user and returns a new JWT token
[x] POST   /register/          : creates a new account

[x] TLS connection
[x] JWT implementation
[x] MongoDB
