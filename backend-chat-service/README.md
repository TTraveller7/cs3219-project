# Chat Service

## Socket Events

Socket io server is listening at localhost:5300. The server will handle one of the following events, and emit an event to client as response.
- **Event names** are written in bold.
- ARGUMENTS are uppercase. All arguments are strings.

**chatMessage** (USERNAME, TEXT)
- Update user's message to the room where the user joins.
- Response events
    | Event name | Arguments | Return |
    | -- | -- | -- | 
    | **updateMessage** | USERNAME | A message object `{username, text, time}` |
