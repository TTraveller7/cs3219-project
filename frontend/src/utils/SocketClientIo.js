export const socketMatchingServiceConfig = {
    cors: {
        origin: "http://localhosts:3000",
        methods: ["GET", "POST"],
    },
    cookie: true,
    withCredentials: true,
    path: "/matching-service/",
    transports: ["websocket"]
};

// SocketIO Events emit
export const MATCHING_USER = "matchingUser";
export const LEAVING_ROOM = "leavingRoom";
export const MATCHING_TIMEOUT = "matchingTimeout";

// SocketIO events listening
// From MATCHING_USER
export const matchSuccess = "matchSuccess";
export const matchPending = "matchPending";
export const idle = 'idle';
// From LEAVING_ROOM
export const leaveSuccess = "leaveSuccess";
// From MATCHING_TIMEOUT
export const matchFail = "matchFail";
export const error = "error";

export const socketChatServiceConfig = {
    cors: {
        origin: "http://localhost:3000",
        methods: ["GET", "POST"],
    },
    withCredentials: true,
    cookie: true,
    path: "/chat-service/",
};

// SocketIO event emits
export const chatMessage = 'chatMessage';

// SocketIO events listening (chat service)
export const updateMessage = 'updateMessage';

export const socketCollabServiceConfig = {
    cors: {
        origin: "http://localhost:3000",
        methods: ["GET", "POST"],
    },
    withCredentials: true,
    cookie: true,
    path: "/collab-service/"
};

// SocketIO events listening (collab service)
export const loadDocument = 'loadDocument';
export const receiveChanges = 'receiveChanges';
export const createSuccess = 'createSuccess';

// SocketIO event emits
export const getDocument = 'getDocument';
export const createNewDocument = 'createNewDocument';
export const sendChanges = 'sendChanges';
export const saveDocument = 'saveDocument';
