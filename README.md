# CS3219-AY22-23-Project-Skeleton

## Contents

- [To Start](#to-start)
- [Notes](#notes)
- [User Service API](#user-service-api)
- [Matching Service API](#matching-service-api)
- [Question Service API](#question-service-api)
- [Matching Service Socketio Events](#matching-service-socketio-events)
- [Chat Service Socketio Events](#chat-service-socket-events)
- [Frontend](#frontend)

## Relavent links
- Report: [link](./14-FinalReport.pdf)
- Kubernetes config repository: [link](https://github.com/TTraveller7/codetakestwo-k8s)

## To start
- Run `docker compose up` in root directory.
- Wait for all logs to come up.

## Notes
- All backend services uses `localhost:8080` as hostname & port. 
- *Authorization required*: Request should carry a cookie `Authorization` containing JWT. The JWT could be obtained by a successful login.
    - Authorization is automatically done for certain endpoints (i.e. no need to query `/validate` before each request).
- The `error` attribute in response body usually can be directly displayed to user. However, it is not useful for debugging the service. See user service logging for more information about what really happens.
    - Internal error usually refers to an unsuccessful database/cache operation. They are represented by status code 400. In other words, a 400 error code (usually with a generic error message) could potentially means that an internal error occurs in the backend.
- Url parameters are represented by {param}. For example, **/match/{matchId}** can be **/match/66dd2f02-78f0-4c57-a98c-adf06f3a36fc**. 
- Similarly, url query parameters are represented by queryKey={param}. For example, **/question/next?currQid={qid}** can be **/question/next?currQid=1**

## User Service API

**[Post] /user/create**
- Create a new user with username and password. Fail to create the user if the username already exists in the database.
- Request body example:
    ```json
    {
        "username":"teemo",
        "password":"onduty"
    }
    ```
- Response
    | Code | Description | Response body example |
    | -- | -- | -- |
    | 202 | Successfully created | `{"message":"Successfully created"}`
    | 400 | Fail to create user | `{"error":"Fail to create user"}`

**[Post] /login**
- Login to an existing account. If credentials are correct, a new JWT will be returned as cookie.
- Request body example:
    ```json
    {
        "username":"teemo",
        "password":"onduty"
    }
    ```
- Response
    | Code | Description | Response body example |
    | -- | -- | -- |
    | 200 | Successfully login | `{"message":"Successfully generate JWT"}`
    | 400 | Fail to login | `{"error":"Invalid username or password"}`
- Cookie example
    | Key | Value |
    | -- | -- |
    | Authorization | eyJhbGc... (rest of JWT omitted) | 

**[Post] /validate**
- Check if a JWT is valid. For valid JWTs, return the username associated with the JWT.
- Response
    | Code | Description | Response body example |
    | -- | -- | -- |
    | 200 | JWT is valid | `{"message":"Successfully log in", "username":"Felix"}`
    | 401 | No JWT in request; invalid JWT | `{"error":"Unauthorized"}`

**[Post] /logout** *Authorization required*
- Logout a user. Request should carry username in the body, which is checked against JWT payload. After a successful logout, the current JWT will be blacklisted, and the local browser will remove the JWT from cookie.
- Request body example
    ```json
    {
        "username":"teemo"
    }
    ```
- Response
    | Code | Description | Response body example |
    | -- | -- | -- |
    | 200 | JWT is valid | `{"message":"Delete JWT from cookie"}`
    | 400 | Username and JWT payload unmatched | `{"error":"Fail to log out"}`
    | 401 | No JWT in request; invalid JWT; Username and JWT payload unmatched | `{"error":"Unauthorized"}`

**[Post] /user/delete** *Authorization required*
- Delete an existing user. Request should carry username in the body, which is checked against JWT payload.
- Request body example
    ```json
    {
        "username":"teemo"
    }
    ```
- Response
    | Code | Description | Response body example |
    | -- | -- | -- |
    | 202 | JWT is valid | `{"message":"Successfully delete user teemo"}`
    | 400 | Username and JWT payload unmatched | `{"error":"Fail to delete user"}`
    | 401 | No JWT in request; invalid JWT  | `{"error":"Unauthorized"}`

**[Post] /user/changepwd** *Authorization required*
- Change the password for a user. Request should carry both old password and new password in the body. 
- Request body example
    ```json
    {
        "oldPassword":"felix",
        "newPassword":"frank"
    }
    ```
- Response
    | Code | Description | Response body example |
    | -- | -- | -- |
    | 200 | JWT is valid | `{"message":"Successfully changed password"}`
    | 400 | Schema unmatch; Old password is not aligned with database; Old password is the same as new password | `{"error":"Wrong old password"}`
    | 401 | No JWT in request; invalid JWT  | `{"error":"Unauthorized"}`

## Matching Service API

Note:
- All endpoints in matching service requires authorization.
- The user endpoints in match API are to retrieve **match information of a user**, whereas user service enpoints are to retrieve **user credentials**.

**[Post] /match/create** *Authorization required*
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
    | 401 | No JWT in request; invalid JWT  | `{"error":"Unauthorized"}`

**[GET] /match/{matchId}** *Authorization required*
- Get match with the match id. If success, the returned json will contain a message and a match object.
- Response
    | Code | Description | Response body example |
    | -- | -- | -- |
    | 200 | Match exists and successfully fetched | See below |
    | 400 | No match with match id `matchId` exists; | `{"error":"Fail to find match with id s"}` |
    | 401 | No JWT in request; invalid JWT  | `{"error":"Unauthorized"}`
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

**[PUT] /match/{matchId}/end** *Authorization required*
- End an existing match. End an ended match has not effect.
- Response
    | Code | Description | Response body example |
    | -- | -- | -- |
    | 200 | Match is successfully ended; match already ends  | `{"message":"Match already ended"}`
    | 400 | No match with id `matchId`  | `{"error":"Fail to find match with id bf2e19e0-8eeb-46e0-9b26-34c77d111f13"}`
    | 401 | No JWT in request; invalid JWT  | `{"error":"Unauthorized"}`

**[GET] /user/{username}** *Authorization required*
- Get user object with username. The user object in matching service will contain match-related info.
- Response
    | Code | Description | Response body example |
    | -- | -- | -- |
    | 200 | User exists and successfully fetched | See below |
    | 400 | No user with username `username` exists; | `{"error":"Fail to find user with name aNonExistingName"}` |
    | 401 | No JWT in request; invalid JWT  | `{"error":"Unauthorized"}`
- For 200 ok, response body example:
    ```json
    {
        "message": "Successfully retrieved user",
        "user": {
            "Username": "teemo61",
            "Status": "matched",
            "LastMatchRequestId": "fd308efe-59f9-4838-9087-00ed3003c9ad",
            "MatchId": "bf2e19e0-8eeb-46e0-9b26-34c77d111f12"
        }
    }
    ```

**[PUT] /user/{username}/toidle** *Authorization required*
- Set user status to idle. Trying to set an idle user to idle has no effect.
- Response
    | Code | Description | Response body example |
    | -- | -- | -- |
    | 200 | User status successfully set to idle; user is already idle | `{"message": "User status set to idle"}` |
    | 400 | No user with username `username` exists | `{"error":"Fail to find user with name aNonExistingName"}` |
    | 401 | No JWT in request; invalid JWT  | `{"error":"Unauthorized"}`

## Question Service API

**[GET] /question/{questionId}** *Authorization required*
- Get question with questionId.
- Response
    | Code | Description | Response body example |
    | -- | -- | -- |
    | 200 | Question exists and successfully fetched | See below |
    | 400 | No question with id `questionId` exists; | `{"error":"Fail to retrieve question"}` |
    | 401 | No JWT in request; invalid JWT  | `{"error":"Unauthorized"}`
- For 200 ok, response body example:
    ```json
    {
        "message": "Successfully retrieved question",
        "question": {
            "id": 1,
            "difficulty": "easy",
            "name": "Two Sum",
            "description": "Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.\r\n    You may assume that each input would have exactly one solution, and you may not use the same element twice.\r\n    You can return the answer in any order."
        }
    }
    ```

**[GET] /question/start** *Authorization required* *MatchId cookie required*
- Create question list for the match if not exists. Fetch the first question (regardless of what is the latest question).
- For two parallel requests, both are likely to succeed.
    - When the request fails, frontend can try to resend the request.
- This API requires a match id cookie, which can be fetched from `/user/{username}` API for a matched user.
- Response
    | Code | Description | Response body example |
    | -- | -- | -- |
    | 200 | Question successfully fetched | See below |
    | 400 | Fail to fetch lock in 5s; Internal error | `{"error":"Fail to get next question"}` |
    | 401 | No JWT (`Authorization` and `MatchId`) in request; invalid JWT  | `{"error":"Unauthorized"}`
- For 200 ok, response body example:
    ```json
    {
        "message": "Successfully retrieved question",
        "question": {
            "id": 1,
            "difficulty": "easy",
            "name": "Two Sum",
            "description": "Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.\r\n    You may assume that each input would have exactly one solution, and you may not use the same element twice.\r\n    You can return the answer in any order."
        }
    }
    ```

**[GET] /question/next** *Authorization required* *MatchId cookie required*
- Get next question for a match. The API should be called by **socketio** to retrieve questions.
- For multiple parallel requests on the same `currQid`, only one will succeed. The unsuccessful requests will receive a 400 immediately.
- This API requires a match id cookie, which can be fetched from `/user/{username}` API for a matched user.
- Response
    | Code | Description | Response body example |
    | -- | -- | -- |
    | 200 | Question successfully fetched | See below |
    | 400 | The questions in the list are used up; Backend is processing another */question/next* request; internal error | `{"error":"Fail to fetch the next question. Please try again later."}` |
    | 401 | No JWT (`Authorization` and `MatchId`) in request; invalid JWT  | `{"error":"Unauthorized"}`
- For 200 ok, response body example:
    ```json
    {
        "message": "Successfully retrieved question",
        "question": {
            "id": 1,
            "difficulty": "easy",
            "name": "Two Sum",
            "description": "Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.\r\n    You may assume that each input would have exactly one solution, and you may not use the same element twice.\r\n    You can return the answer in any order."
        }
    }
    ```

**[GET] /question/curr** *Authorization required* *MatchId cookie required*
- Get the latest question in a match.
- The API call could fail if the question list is yet to constructed, i.e. /question/start request is still being handled. The frontend can retry a few times.
- This API requires a match id cookie, which can be fetched from `/user/{username}` API for a matched user.
- Response
    | Code | Description | Response body example |
    | -- | -- | -- |
    | 200 | Question successfully fetched | See below |
    | 400 | The question list is yet to be constructed; internal error | `{"error":"Fail to fetch the current question. Please try again later."}` |
    | 401 | No JWT (`Authorization` and `MatchId`) in request; invalid JWT  | `{"error":"Unauthorized"}`
- For 200 ok, response body example:
    ```json
    {
        "message": "Successfully retrieved question",
        "question": {
            "id": 1,
            "difficulty": "easy",
            "name": "Two Sum",
            "description": "Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.\r\n    You may assume that each input would have exactly one solution, and you may not use the same element twice.\r\n    You can return the answer in any order."
        }
    }
    ```

**[GET] /answer?questionId={questionId}&matchId={matchId}** *Authorization required*
- For a certain match, get the latest answer for a question.
- If the request is successful, the response body will contain the answer as well as the associated question.
- Response
    | Code | Description | Response body example |
    | -- | -- | -- |
    | 200 | Answer successfully fetched | See below |
    | 400 | No answer associated with (matchId, questionId) exists. | `{"error":"Fail to get answer"}` |
    | 401 | No JWT in request; invalid JWT  | `{"error":"Unauthorized"}`
- For 200 ok, response body example:
    ```json
    {
        "answer": {
            "ID": 1,
            "CreatedAt": "2022-10-21T11:12:16.664865Z",
            "UpdatedAt": "2022-10-21T11:12:16.664865Z",
            "DeletedAt": null,
            "matchId": "37367a63-b27b-4d11-912e-76f83feb5fa0",
            "code": "print('helloewwww world!')",
            "questionId": 1,
            "question": {
                "ID": 1,
                "CreatedAt": "2022-10-21T11:09:51.342629Z",
                "UpdatedAt": "2022-10-21T11:09:51.342629Z",
                "DeletedAt": null,
                "difficulty": "easy",
                "name": "Two Sum",
                "description": "Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.\n\t\tYou may assume that each input would have exactly one solution, and you may not use the same element twice.\n\t\tYou can return the answer in any order."
            }
        },
        "message": "Successfully retrieved answer"
    }
    ```

**[POST] /answer/create** *Authorizatoin required* *MatchId Cookie required*
- For a certain match, save the an answer associated with the question.
- This API requires a match id cookie, which can be fetched from `/user/{username}` API for a matched user.
- Request body example:
    ```json
    {
        "questionId": 1,
        "code": "print('helloewwww world!')"
    }
    ```
- Response
    | Code | Description | Response body example |
    | -- | -- | -- |
    | 200 | Answer successfully created | `{ "message": "Successfully saved answer" }` |
    | 400 | Missing fields in request body; | `{"error": "Fail to save answer"}` |
    | 401 | No JWT (`Authorization` and `MatchId`) in request; invalid JWT  | `{"error": "Unauthorized"}`

## Matching Service Socketio Events

The server will handle one of the following events, and emit an event to client as response.
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

**fetchingQuestion** USERNAME QID
- Fetch next question by given current question ID
- Response events
    | Event name | Return | Description |
    | -- | -- | -- | 
    | **fetchSuccess** | QUESTION | A question object `{id, difficulty, name, description}` |
    | **error** | ERROR_MSG | User is not matched; internal error |

## Chat Service Socket Events

The server will handle one of the following events, and emit an event to client as response.
- **Event names** are written in bold.
- ARGUMENTS are uppercase. All arguments are strings.

**chatMessage** USERNAME TEXT
- Update user's message to the room where the user joins.
- Response events
    | Event name | Return | Description |
    | -- | -- | -- | 
    | **updateMessage** | USERNAME | A message object `{username, text, time}` |

## Collaboration Service Socket Events

The server will handle one of the following events, and emit an event to client as response.
- **Event names** are written in bold.
- ARGUMENTS are uppercase. All arguments are strings.

**createNewDocument** USERNAME
- Create a record of document in cache with match id as its key.
- Response events
    | Event name | Return | Description |
    | -- | -- | -- | 
    | **createSuccess** | STRING | "Successfully create a shared document" |

**getDocument** USERNAME
- Get document in cache using match id.
- This event should only be emitted once for every connection
- Response events
    | Event name | Return | Description |
    | -- | -- | -- | 
    | **loadDocument** | DATA | The content of document |
    | **createSuccess** | STRING | "Successfully create a shared document" |

**sendChanges** DELTA
- Update the changes to other clients in the room.
- Should be emitted after getDocument event
- Response events
    | Event name | Return | Description |
    | -- | -- | -- | 
    | **receiveChanges** | CHANGE | A object `{username, delta}` |

**saveDocument** DATA
- Save the content of document in server.
- Should be emitted after getDocument event


## Frontend
1. Install npm packages using `npm i`.
2. Run Frontend using `npm start`.

**Note:** After all updates on the frontend, run `npm run build` in the frontend directory before updating the docker file or pushing into the GitHub repo.
