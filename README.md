# Story GraphQL TS

![wakatime](https://wakatime.com/badge/user/bc9f5590-5561-4bf9-a568-e58c70266176/project/28842fc7-1047-4b1c-af9d-10867cbe0c1e.svg)


## Main Features

- **Login:** Allows users to authenticate and obtain a user token.
- **Register:** Enables user registration.
- **Get Guest Token:** Returns a UUID token for guest users.
- **Auth and Check Guest Token:** Verifies the validity of a guest token.
- **Get Stories:** Retrieves a list of stories based on the user's role:
  - For users: Returns a list of all stories created by the user.
  - For guest users: Returns a list of all stories scanned by the guest user.
- **Get Subscribed Stories:** Retrieves a list of subscribed stories based on the user's role:
  - For users: Returns a list of all subscribed stories created by the user.
  - For guest users: Returns a list of all subscribed stories scanned by the guest user.
- **Delete Account:** Deactivates a user account and keeps a record of the deletion:
  - This route is available to both user roles.
- **Get Story:** Retrieves the details of a specific story based on the user's role:
  - This route is available to both user roles.
- **Scan Story:** Scans a story for a guest user:
  - This route is only available to guest users.
- **Create New Story:** Creates a new story:
  - This route is only available to users.
- **Edit Story:** Edits an existing story:
  - This route is only available to users.
- **Delete My Story:** Deletes a story created by the user:
  - This route is only available to users.

## Running

**Docker:**

```bash
# In Production, build a multi stage image
docker-compose -f docker-compose-prod.yml up

# In Development, uses air for live reloading and have database ports exposed
docker-compose -f docker-compose-dev.yml up
```

## Swagger documentation

```bash
# Install swag-go generator:
go install github.com/swaggo/swag/cmd/swag@latest

# Generate documents:
swag i -g docs.go -d ./controller/echo/,./DTO
```

### Internal swagger ui:

To enable swagger ui in the project, change the server configuration in `config.yml` to look like this:

```yaml
server:
  listen: ":8080"
  debug: ":8001"  # can be used to expose pprof data
  docs:
    path: "/swagger*" # remove path to disable swagger ui.
    # these are values for swagger ui.
    base_path: "/api/v1/"
    host: "localhost"
    title: "StoryGoApi"
    version: "v1"
    description: "Your description"
```

Then run the project and go to: `localhost:8080/swagger`

To disable swagger ui just remove the `docs` part of config file or remove the `path` from config.

## Ent Orm

Database schema is available at `ent/schema` directory.

To regenerate the orm method after changing `ent/schema` run the following command:

```bash
go generate ./ent
```
