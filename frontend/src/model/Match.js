import axios from "axios";
import {URL_MATCHING_SVC_GET_MATCH, URL_MATCHING_SVC_GET_USER} from "../configs";

export class Match {
    constructor(match) {
        this.matchId = match.matchId;
        this.difficulty = match.difficulty;
        this.usernameA = match.usernameA;
        this.usernameB = match.usernameB;
        this.isEnded = match.isEnded;
        this.createdAt = match.createdAt;
    }

    getRemoteUser(currentUser) {
        return currentUser !== this.usernameA ? this.usernameA : this.usernameB;
    }
}

async function checkIfUserHasMatch(username) {
    const res = await axios.get(URL_MATCHING_SVC_GET_USER + username, {withCredentials: true})
        .catch((err) => {
            return undefined;
        });
    if (res && res.data.user.MatchId) {
        return res.data.user.MatchId
    }
    return undefined
}

async function getMatch(matchId) {
    const res = await axios.get(URL_MATCHING_SVC_GET_MATCH + matchId, {withCredentials: true})
        .catch((err) => {
            return undefined;
        });
    if (res && res.data.match) {
        return new Match(res.data.match)
    }
    return undefined
}

export async function getUserInMatch(username) {
    const matchId = await checkIfUserHasMatch(username)
    if (matchId === undefined) {
        return undefined;
    }
    return getMatch(matchId);
}