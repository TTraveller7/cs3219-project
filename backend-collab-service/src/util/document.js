import { Entity, Schema} from "redis-om";
import { client } from "./cache.js";

const defaultData = "";

/* our entity */
class Document extends Entity {};

/* create a Schema for Document */
const docSchema = new Schema(Document, {
  matchId: { type: 'string' },
  documentUpdated: { type: 'date' },
  data: { type: 'text' }
});

/* use the client to create a Repository just for Document */
const docRepository = client.fetchRepository(docSchema);

await docRepository.createIndex();

/* Util function for cache*/
async function createDoc(matchId) {
  let document = {matchId: matchId, documentUpdated: Date.now(), data: defaultData};
  document = await docRepository.createAndSave(document);
  return document;
}

async function saveDoc(document) {
  await docRepository.save(document);
  return true;
}

async function getDoc(matchId) {
  return await docRepository.fetch(id);
}

async function getDocsByMatchId(matchId) {
  return await docRepository.search().where('matchId').equals(matchId).return.first();
}

export {
  createDoc,
  saveDoc,
  getDoc,
  getDocsByMatchId,
  defaultData,
};