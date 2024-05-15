# Todo

A simple backend for to-do list

# Run
1. Download project
2. Edit `config.yml` file
3. Run `go run main.go`

# Routes
- `POST /api/register` - User registration
- `POST /api/login` - User authentication
- `GET /api/auth` - Route used to verify user authentication
- `POST /api/logout` - Route for user logout (session termination)
- `POST /api/settings/password` - Route to change the password
- `POST /api/todo` - Create a new task by the user
- `POST /api/todo/{id}` - Route to completion of a specific task
- `GET /api/todo?limit=10&page=1` - Route to get a list of all user tasks  with a page limit
- `PUT /api/todo` - Edit an existing task
- `DELETE /api/todo/{id}` - Deleting a specific task