const { Server } = require("socket.io");
const axios = require('axios').default;
const { EMBED_USERNAME, EMBED_MATCHID, BACKEND_MATCHING_SERVICE_TO_IDLE, BACKEND_QUESTION_SERVICE_GET_NEXT, MATCHING_SERVICE_HEALTHCHECK } = require('./config');
const { STATUS_CODE_BAD_REQUEST } = require('./constant');
const { isPending, notMe } = require('./utils');

axios.get(MATCHING_SERVICE_HEALTHCHECK)
.catch(err => {
  if (err.response) {
    console.log(err.response.data.error);
  } else {
    console.log('Error', err);
  }
  process.exit();
});

const io = new Server(5200, {
  path: "/matching-service/",
  transports: ["websocket"]
});

io.on("connection", (socket) => {
  console.log(`User Connected: ${socket.id}`);
  const username = socket.request.headers["x-username"];
  console.log(`User Connected: ${username}`); 

  const matchId = socket.request.headers["x-matchid"];
  if ( typeof matchId !== 'undefined' && matchId ) {
    socket.join(matchId);
    console.log(`User is Matched: ${matchId}`);
  }

  socket.on("matchingUser", async (name) => {
    try {
      const res = await axios.get(EMBED_USERNAME(name), {withCredentials: true});
      console.log(res.data.message);
      let user = res.data.user;
      
      if (user.Status === "matched") {
        try {
          const res = await axios.get(EMBED_MATCHID(user.MatchId), {withCredentials: true});
          console.log(res.data.message);
          let match = res.data.match;

          // socket.join(match.MatchId);
          socket.emit("matchSuccess", user.Status, notMe(name, match.usernameA, match.usernameB)); 
        } catch (err) {
          if (error.response) {
            console.log(err.response.data.error);
          } else {
            console.log('Error', err);
          }
          
        }
      } else {
        socket.emit("matchPending", user.Status);
      }
    } catch (err) { 
      if (err.response) {
        if (err.response.status === STATUS_CODE_BAD_REQUEST) {
          console.log(`Fail to retrieve user: ${name}`)
          socket.emit("error", "User did not send matching request to server before");
        }
      } else {
        console.log('Error', err);
      }
    }
  });

  socket.on("leavingRoom", async (name) => {
    try {
      const res = await axios.get(EMBED_USERNAME(name), {withCredentials: true});
      let user = res.data.user;

      if (user.Status === "matched") {
        const res = await axios.get(EMBED_MATCHID(user.MatchId), {withCredentials: true});
        console.log(res.data.message);
        let match = res.data.match;

        socket.to(match.matchId).emit("leaveSuccess", name);
        socket.leave(match.matchId);
        // axios.put(EMBED_MATCHID(match.MatchID) + BACKEND_MATCHING_SERVICE_END)
        //   .catch((err) => {
        //     if (err.response) {
        //       console.log(err.response.data.error);
        //     } else {
        //       console.log('Error', err.message);
        //     }
        //   });
        axios.put(EMBED_USERNAME(name) + BACKEND_MATCHING_SERVICE_TO_IDLE, {}, {withCredentials: true})
          .catch((err) => {
            if (err.response) {
              console.log(err.response.data.error);
            } else {
              console.log('Error', err.message);
            }
          });
        socket.emit("leaveSuccess", name)
      } else {
        console.log(`User: ${name} is not in a room`);
        socket.emit("error", "User is not in a room");
      }
    } catch (err) {
      if (err.response) {
        if (err.response.status === STATUS_CODE_BAD_REQUEST) {
          console.log(`Fail to retrieve user: ${name}`)
          socket.emit("error", "User did not send matching request to server before");
        }
      } else {
        console.log('Error', err);
      }
    }
  });

  socket.on("matchingTimeout", async (name) => {
    try {
      const res = await axios.get(EMBED_USERNAME(name), {withCredentials: true});
      console.log(res.data.message);
      let user = res.data.user;
  
      if (isPending(user.Status)) {
        socket.emit("matchFail", user.Username, "Timeout")
        
        axios.put(EMBED_USERNAME(name) + BACKEND_MATCHING_SERVICE_TO_IDLE, {}, {withCredentials: true})
          .catch((err) => {
            if (err.response) {
              console.log(err.response.data.error);
            } else {
              console.log('Error', err.message);
            }
          })
        
      } else {
        socket.emit("error", "User is not in a queue");
      }
    } catch (err) {
      if (err.response) {
        console.log(err.response.data.error);
      } else {
        console.log('Error', err.message);
      }
      socket.emit("error", "User is not in a queue");
    }
  });

  socket.on("fetchingQuestion", async (name, qid) => {
    try {
      const res = await axios.get(EMBED_USERNAME(name), {withCredentials: true});
      console.log(res.data.message);
      let user = res.data.user;
      if (user.Status == "matched") {
        console.log(BACKEND_QUESTION_SERVICE_GET_NEXT);
        const res = await axios.get(BACKEND_QUESTION_SERVICE_GET_NEXT, {
          params: {currQid: qid },
          headers: {
            "x-matchid": user.MatchId
          } 
        });
        console.log(res.data.message);
        
        io.in(user.MatchId).emit("fetchSuccess", res.data.question);
      } else {
        io.in(user.MatchId).emit("error", "Fail to fetch a question: User is not matched");
      }
    } catch (err) {
      if (err.response) {
        console.log(err.response.data.error);
      } else {
        console.log('Error', err.message);
      }
      io.in(user.MatchId).emit("error", "Fail to fetch a question");
    }
  });

  socket.on("disconnecting", async (reason) => {
    try {
      console.log(`User:${username} is disconnecting`);
      const res = await axios.get(EMBED_USERNAME(username), {withCredentials: true});
      let user = res.data.user;
      
      if (isPending(user.Status)) {
        // Set user status to idle if user is in matching status.
        axios.put(EMBED_USERNAME(name)+BACKEND_MATCHING_SERVICE_TO_IDLE, {}, {withCredentials: true});
      }
    } catch (err) {
      if (err.response) {
        console.log(err.response.data.error);
      } else {
        console.log('Error', err.message);
      }
    }
  });


  socket.on("disconnect", async (reason) => {
    try {
      
    } catch(err) {
      if (err.response) {
        console.log(`Fail to retrieve user: ${username}`);
      } else {
        console.log('Error', err);
      }
    }
  });

});


// io.engine.on("initial_headers", (headers, req) => {
//   headers["set-cookie"] = "MatchId=";
// });

// io.engine.on("headers", (headers, req) => {
//   const name = headers["X-Username"];
//   const res = axios.get(EMBED_USERNAME(name))
//     .catch((err) => {
//       if (err.response) {
//         if (err.response.status === STATUS_CODE_BAD_REQUEST) {
//           console.log(`Fail to retrieve user: ${name}`)
//         }
//       } else {
//         console.log('Error', err);
//       }
//     });
  
//   let user = res.data.user;

//   if (user.Status == "matched") {
//     headers["set-cookie"] = "MatchId="+user.MatchID;
//   }
// });