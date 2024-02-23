const axios = require('axios').default;
const { EMBED_USERNAME, EMBED_MATCHID, BACKEND_MATCHING_SERVICE_END, BACKEND_MATCHING_SERVICE_TO_IDLE } = require('../config');
const { STATUS_CODE_BAD_REQUEST } = require('../constant');
const { isPending } = require('../utils');

const MatchingUser = async (socket, name) => {
  try {
    const res = await axios.get(EMBED_USERNAME(name));
    console.log(res.data.message);
    let user = res.data.user;
    
    if (user.Status === "matched") {
      try {
        const res = await axios.get(EMBED_MATCHID(user.matchID));
        console.log(res.data.message);
        let match = res.data.match;

        socket.join(match.MatchID);
        socket.emit("matchSuccess", user.Status, match.UsernameB); 
      } catch (err) {
        console.log(err.response.data.error);
      }
    } else {
      socket.emit("matchPending", user.Status);
    }
  } catch (err) { 
    if (err.response.status === STATUS_CODE_BAD_REQUEST) {
      console.log(`Fail to retrieve user: ${name}`)
      socket.emit("error", "User did not send matching request to server before");
    }
  }
};

const LeavingRoom = async (socket, name) => {
  try {
    const res = await axios.get(EMBED_USERNAME(name));
    console.log(res.data.message);
    let user = res.data.user;

    if (user.Status === "matched") {
      const res = await axios.get(EMBED_MATCHID(user.MatchID));
      console.log(res.data.message);
      let match = res.data.match;

      socket.to(match.MatchID).emit("leaveSuccess", match.MatchID);
      socket.leave(match.MatchID);
      axios.put(EMBED_MATCHID(match.MatchID) + BACKEND_MATCHING_SERVICE_END)
        .catch((err) => {
          console.log(err.response.error);
        });
      axios.put(EMBED_USERNAME(name) + BACKEND_MATCHING_SERVICE_TO_IDLE)
        .catch((err) => {
          console.log(err.response.error);
        });

    } else {
      console.log(`User: ${name} is not in a room`);
      socket.emit("error", "User is not in a room");
    }
  } catch (err) {
    if (err.response.status === STATUS_CODE_BAD_REQUEST) {
      console.log(`Fail to retrieve user: ${name}`);
      socket.emit("error", "User did not send matching request to server before");
    }
  }
};

const MatchingTimeout = async (socket, name) => {
  try {
    const res = await axios.get(EMBED_USERNAME(name));
    console.log(res.data.message);
    let user = res.data.user;

    if (isPending(user.Status)) {
      socket.emit("matchFail", user.Username, "Timeout")
      
      axios.put(EMBED_USERNAME(name) + BACKEND_MATCHING_SERVICE_TO_IDLE)
        .catch((err) => {
          console.log(err.response.error);
        })
			
    } else {
      console.log(err.response.data.error);
      socket.emit("error", "User is not in a queue");
    }
  } catch (err) {
    console.log(err.response.data.error);
    socket.emit("error", "User is not in a queue");
  }
};

const Disconnecting = (socket, reason) => {
  console.log(socket.id);
  if (socket.adpter.rooms !== null) {
    for (const room of socket.rooms) {
      socket.to(room).emit("leaveSuccess", user.Username);
      axios.put(EMBED_MATCHID(room) + BACKEND_MATCHING_SERVICE_END)
      .catch((err) => {
        console.log(err.response.error);
      });
    axios.put(EMBED_USERNAME(name) + BACKEND_MATCHING_SERVICE_TO_IDLE)
      .catch((err) => {
        console.log(err.response.error);
      });
    }
  }
};

module.exports = {
  MatchingUser,
  LeavingRoom,
  MatchingTimeout,
  Disconnecting
};