import { Entity, Schema} from "redis-om";
import { client } from "./cache.js";

/* our entity */
class Message extends Entity {};

/* create a Schema for Message */
const msgSchema = new Schema(Message, {
  matchId: { type: 'string' },
  username: { type: 'string' },
  text: { type: 'text' },
  time: { type: 'number', sortable: true}
});

/* use the client to create a Repository just for Message */
const msgRepository = client.fetchRepository(msgSchema);

await msgRepository.createIndex();

/* Util function for cache*/
export async function createMsg(matchId, msg) {
  const message = {matchId: matchId, username: msg.username, text: msg.text, time: msg.time};
  await msgRepository.createAndSave(message);
  return true;
}

export async function getMsgInOrder(matchId) {
  return await msgRepository.search().where('matchId').equals(matchId).sortBy('time', 'ASC').return.all();
}

export async function contain(matchId) {
  return await msgRepository.search().where('matchId').equals(matchId).return.count();
}
