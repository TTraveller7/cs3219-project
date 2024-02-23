import express from "express";
const app = new express();
import http from "http";
import { Server } from "socket.io";
import cors from "cors";
import { formatMessage } from "./util/messages-format.js";
import {createMsg, getMsgInOrder, contain} from "./util/message.js";


app.use(cors());

const server = http.createServer(app);

const io = new Server(5300, {
  path: "/chat-service/",
});

const botName = 'PeerPrep Bot'

io.on('connection', async (socket) => {
  const matchId = socket.request.headers["x-matchid"];
  const username = socket.request.headers['x-username'];
  console.log(`User Connected to chat: ${username}`);
  console.log(`User is Matched to chat: ${matchId}`);

  // Greet to user
  socket.emit('updateMessage', formatMessage(botName, 'Welcome to Code Takes Two!'));

  socket.join(matchId);

  socket.emit('updateMessage', formatMessage(botName, `${username} has joined`));

  if (await contain(matchId) != 0) {
    const msgs = await getMsgInOrder(matchId);

    msgs.forEach(msg => {
      io.to(matchId).emit('updateMessage', {
        username: msg.username,
        text: msg.text,
        time: msg.time
      });
    });
  }

  // Update message to room
  socket.on('chatMessage', (name, msg) => {
    socket.to(matchId).emit('updateMessage', formatMessage(name, msg));
    createMsg(matchId, formatMessage(name, msg));
  });

  socket.on('disconnect', (reason) => {
    socket.to(matchId).emit('updateMessage', formatMessage(botName, `${username} has left`));
  })
});