# Matching Service

## Set up instructions

**1.** Clear docker volumes. `docker volume prune`, then `y`

**2.** Two ways to start:
- Under ../postgres directory, enter `docker compose up -d`; then under project root directory, enter `docker compose -f matching-service.yaml -d`. After the images are ready, this should takes around 40s for all components to be ready.
- Or you can start up each component individually:
    - Under ../postgres directory, enter `docker compose up -d`
    - Under redis, enter `docker compose up -d`
    - Under kafka, enter `docker compose up`
    - Wait until kafka pop up all the loggings
    - Under matching-service, enter `go run .`

## HTTP API

Note:
- Url parameters are represented by {param}. For example, **/match/{matchId}** can be **/match/66dd2f02-78f0-4c57-a98c-adf06f3a36fc**

**[Post] /match/create**
- Create a match request. Fail to create the user is not idle (i.e., the user is in a matching queue/has an ongoing match).
- `difficulty` should be one of the following: `easy`, `medium`, `hard`
- Request body example:
    ```json
    {
        "username":"teemo",
        "difficulty":"easy"
    }
    ```
- Response
    | Code | Description | Response body example |
    | -- | -- | -- |
    | 200 | Match request created | `{"message":"Match request sent"}`
    | 400 | Request schema unmatch; user status is not idle; internal error | `{"error":"Fail to create user"}`

**[GET] /match/{matchId}**
- Get match with the match id. If success, the returned json will contain a message and a match object.
    - In the match object, `endedAt` is 
- Response
    | Code | Description | Response body example |
    | -- | -- | -- |
    | 200 | Match exists and successfully fetched | See below |
    | 400 | No match with match id `matchId` exists; | `{"error":"Fail to find match with id s"}` |
- For 200 ok, response body example:
    ```json
        {
            "match": {
                "matchId": "66dd2f02-78f0-4c57-a98c-adf06f3a36fc",
                "difficulty": "medium",
                "usernameA": "teemo16",
                "usernameB": "teemo17",
                "isEnded": false,
                "createdAt": "2022-09-23T12:04:07.474886052Z",
                "endedAt": "0001-01-01T00:00:00Z"
            },
            "message": "Successfully retrieved match"
        }
    ```

## Socket Events

Socket io server is listening at localhost:5200. The server will handle one of the following events, and emit an event to client as response.
- **Event names** are written in bold.
- ARGUMENTS are uppercase. All arguments are strings.

**matchingUser** USERNAME
- Check a user's matching status. 
- USER_STATUS is one of `idle`, `pending_easy`, `pending_medium`, `pending_hard`, `matched`.
- If the user is successfully matched with another user (and the user is not already in a room), this event will trigger the server to join the two users into a room.
- Response events
    | Event name | Argument | Description |
    | -- | -- | -- | 
    | **matchSuccess** | USER_STATUS ANOTHER_USERNAME | Successfully matched with another user |
    | **matchPending** | USER_STATUS | User is idle or pending in a queue |
    | **error** | ERROR_MSG | User did not send matching request to server before (see /match/create API); internal error |

**leavingRoom** USERNAME
- Let a user leave the room.
- Response events
    | Event name | Arguments | Description |
    | -- | -- | -- | 
    | **leaveSuccess** | USERNAME | User successfully leaves room |
    | **error** | ERROR_MSG | User is not in a room; internal error |

**matchingTimeout** USERNAME
- Let a user leaves the matching queue.
- Response events
    | Event name | Arguments | Description |
    | -- | -- | -- | 
    | **matchFail** | USERNAME `Timeout` | User successfully leaves the matching queue |
    | **error** | ERROR_MSG | User is not in a queue; internal error |