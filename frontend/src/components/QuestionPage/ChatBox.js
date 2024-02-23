import React, {useEffect, useRef, useState} from "react";
import {ChatBubble, Send} from "@mui/icons-material";
import {Box, colors, Fab, Grow, IconButton, List, Paper, Stack} from "@mui/material";
import "./chatBoxStyle.css";
import ChatInfo, {CHAT_BUBBLE_REMOTE, CHAT_BUBBLE_USER} from "../../model/ChatInfo";
import {formatDateTimestampToMonthDayTime} from "../../utils/FormatDate";
import {io} from "socket.io-client";
import {
	chatMessage,socketChatServiceConfig,
	updateMessage
} from "../../utils/SocketClientIo";
import {useOutletContext} from "react-router-dom";
import {BACKEND_URL} from "../../configs";
let socketIoChatClient;

function ChatBoxSection() {
	const [toggleChat, setToggleChat] = useState(true);
	const [currentUserMessage, setCurrentUserMessage] = useState('');
	const [chats, setChats] = useState([]);
	const [username, match] = useOutletContext();

	// For scrolling to the bottom of the chat box when sending message
	const messagesEndRef = useRef(null)

	function toggleChatBox() {
		setToggleChat(!toggleChat);
	}

	function pressEnterFromChat(e) {
		if (e.key === 'Enter') {
			submitMessage();
		}
	}

	function scrollDownToEndOfMessage() {
		messagesEndRef.current?.scrollIntoView({
			block: 'end',
			inline: 'nearest'
		});
	}

	function addChatToChatsArray(chat) {
		setChats(chats => [...chats, chat]);
	}

	function submitMessage() {
		if (currentUserMessage != '') {
			socketIoChatClient.emit(chatMessage, username, currentUserMessage)
			let dateNow = new Date()
			let dateNowFormatted = formatDateTimestampToMonthDayTime(dateNow)
			addChatToChatsArray(new ChatInfo(username, currentUserMessage, dateNowFormatted, CHAT_BUBBLE_USER));
			setCurrentUserMessage('')
		}
	}

	function receiveMessage(message) {
		const time = formatDateTimestampToMonthDayTime(new Date(message.time))
		if (message.username === username) {
			addChatToChatsArray(new ChatInfo(message.username, message.text, time, CHAT_BUBBLE_USER))
		} else {
			addChatToChatsArray(new ChatInfo(message.username, message.text, time, CHAT_BUBBLE_REMOTE))
		}
	}

	useEffect(() => {
		socketIoChatClient = io.connect(BACKEND_URL, socketChatServiceConfig);
		socketIoChatClient.on(updateMessage, receiveMessage)
		return () => {
			socketIoChatClient.disconnect();
		}
	}, [])

	useEffect(scrollDownToEndOfMessage,[chats])

	const userInput = (
		<Stack
			padding={1}
			direction="row"
		>
			<input
				type="text"
				placeholder="send message..."
				value={currentUserMessage}
				onChange={(e) => setCurrentUserMessage(e.target.value)}
				onKeyDown={pressEnterFromChat}
				className="user-input"
			/>
			<IconButton onClick={submitMessage}>
				<Send/>
			</IconButton>
		</Stack>
	);

	const chatData = (
		<Box
			paddingY={2}
			paddingX={1}
			className="chat-data"
		>
			<Stack
				spacing={2}
			>
				{chats.map(
					(chat, index) => {
						if (chat.type == CHAT_BUBBLE_USER) {
							return (
								<UserChatBubble
									key={index}
									message={chat.message}
									time={chat.time}
									username={chat.username}
								/>
							)
						} else {
							return (
								<RemoteChatBubble
									key={index}
									message={chat.message}
									time={chat.time}
									username={chat.username}
								/>
							)
						}
					})}
				<div ref={messagesEndRef}/>
			</Stack>
		</Box>
	);

	const remoteUserName = (
		<Box
			paddingY={1}
			paddingX={2}
			className='remote-user-name'
		>
			{match.getRemoteUser(username)}
		</Box>
	)

	const chatBox = (
		<Paper elevation={3}>
			<Box>
				<Stack
					padding={0}
					className="chat-box"
				>
					{remoteUserName}
					{chatData}
					{userInput}
				</Stack>
			</Box>
		</Paper>
	);

	return (
		<Box className="section">
			<Grow in={toggleChat}>
				{chatBox}
			</Grow>
			<Fab
				color="primary"
				onClick={toggleChatBox}
			>
				<ChatBubble />
			</Fab>
		</Box>
	)
}

function UserChatBubble({message, time, username}) {
	return (
		<Paper
			elevation={3}
			className='current-user-chat'
		>
			<Stack>
				<Box className='chat-username'>
					{username}
				</Box>
				<Box className='chat-text'>
					{message}
				</Box>
				<Box className='chat-time'>
					{time}
				</Box>
			</Stack>
		</Paper>
	)
}

function RemoteChatBubble({message, time, username}) {
	return (
		<Paper
			elevation={3}
			className='remote-user-chat'
			style={{
				backgroundColor: '#2376d2',
				color: 'white'
			}}
		>
			<Stack>
				<Box className='chat-username'>
					{username}
				</Box>
				<Box className='chat-text'>
					{message}
				</Box>
				<Box className='chat-time'>
					{time}
				</Box>
			</Stack>
		</Paper>
	)
}

export default ChatBoxSection