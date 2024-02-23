class ChatInfo {
    constructor(username, message, time, type) {
        this.username = username;
        this.message = message;
        this.time = time;
        this.type = type;
    }
}
export const CHAT_BUBBLE_USER = "USER";
export const CHAT_BUBBLE_REMOTE  = "REMOTE";
export default ChatInfo;
