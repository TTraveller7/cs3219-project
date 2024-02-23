export function formatMessage(username, text) {
  return {
    username,
    text,
    time: Date.now()
  };
}