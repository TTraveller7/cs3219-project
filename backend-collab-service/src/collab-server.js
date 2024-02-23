import { Server } from "socket.io";
import { formatUpdate } from "./util/delta-format.js";
import { defaultData, createDoc, saveDoc, getDocsByMatchId } from "./util/document.js";

const io = new Server(5400, {
  path: "/collab-service/",
});

io.on('connection', async socket => {
  const matchId = socket.request.headers["x-matchid"];
  const username = socket.request.headers["x-username"];
  console.log(`User Connected to Collab: ${username}`);
  console.log(`User is Matched to Collab: ${matchId}`);

  socket.join(matchId);

  socket.on("createNewDocument", async (username) => {
    let document = await getDocsByMatchId(matchId);

    if (!document) {
      document = await createDoc(matchId);
      io.to(matchId).emit("createSuccess", "Successfully create a shared document");
    } else {
      document.data = defaultData;
      document.documentUpdated = Date.now();
      await saveDoc(document);
      io.to(matchId).emit("receiveChanges", formatUpdate(username, defaultData))
    }
  });

  socket.on("getDocument", async (username) => {
    /* find or create document in cache */
    let document = await getDocsByMatchId(matchId);

    if (!document) {
      document = await createDoc(matchId);
      io.to(matchId).emit("createSuccess", "Successfully create a shared document");
    }

    socket.emit("loadDocument", document.data);

    socket.on("sendChanges", delta => {
      socket.to(matchId).emit("receiveChanges", formatUpdate(username, delta))
    });

    socket.on("saveDocument", async data => {
      document.data = data;
      document.documentUpdated = Date.now();
      await saveDoc(document);
    });
  })

});