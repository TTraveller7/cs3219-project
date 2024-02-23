export const BACKEND_URL = process.env.REACT_APP_BACKEND_URL || "http://localhost:8080"

export const BACKEND_USER_SVC_URL = process.env.URI_USER_SVC || BACKEND_URL;
export const BACKEND_MATCHING_SERVICE_URL = process.env.URI_MATCHING_SVC || BACKEND_URL;
export const BACKEND_QUESTION_SERVICE_URL = process.env.URI_QUESTION_SVC || BACKEND_URL;

const PREFIX_USER_SVC_USER_CREATE = '/user/create'
const PREFIX_USER_SVC_LOGIN = '/login'
const PREFIX_USER_SVC_VERIFICATION = '/validate'
const PREFIX_USER_SVC_LOGOUT = '/logout'
const PREFIX_USER_SVC_DELETE = '/user/delete'
const PREFIX_USER_SVC_USER_CHANGE_PASSWORD = '/user/changepwd'

const PREFIX_MATCHING_SERVICE_USER_MATCHING = "/match/create"
const PREFIX_MATCHING_SVC_GET_USER = "/user/"
const PREFIX_MATCHING_SVC_GET_MATCH = "/match/"

const PREFIX_QUESTION_SVC_QUESTION_START = "/question/start";
const PREFIX_QUESTION_SVC_ANSWER_CREATE = "/answer/create";
const PREFIX_QUESTION_SVC_GET_CURRENT_QUESTION = "/question/curr";

export const URL_USER_SVC_USER_CREATE = BACKEND_USER_SVC_URL + PREFIX_USER_SVC_USER_CREATE
export const URL_USER_SVC_LOGIN = BACKEND_USER_SVC_URL + PREFIX_USER_SVC_LOGIN
export const URL_USER_SVC_VALIDATE = BACKEND_USER_SVC_URL + PREFIX_USER_SVC_VERIFICATION
export const URL_USER_SVC_LOGOUT = BACKEND_USER_SVC_URL + PREFIX_USER_SVC_LOGOUT
export const URL_USER_SVC_USER_DELETE = BACKEND_USER_SVC_URL + PREFIX_USER_SVC_DELETE;
export const URL_USER_SVC_USER_CHANGE_PASSWORD = BACKEND_USER_SVC_URL + PREFIX_USER_SVC_USER_CHANGE_PASSWORD;

export const URL_MATCHING_SVC_MATCHING = BACKEND_MATCHING_SERVICE_URL + PREFIX_MATCHING_SERVICE_USER_MATCHING;
export const URL_MATCHING_SVC_GET_USER = BACKEND_MATCHING_SERVICE_URL + PREFIX_MATCHING_SVC_GET_USER;
export const URL_MATCHING_SVC_GET_MATCH = BACKEND_MATCHING_SERVICE_URL + PREFIX_MATCHING_SVC_GET_MATCH;

export const URL_QUESTION_SVC_QUESTION_START = BACKEND_QUESTION_SERVICE_URL + PREFIX_QUESTION_SVC_QUESTION_START;
export const URL_QUESTION_SVC_ANSWER_CREATE = BACKEND_QUESTION_SERVICE_URL + PREFIX_QUESTION_SVC_ANSWER_CREATE;
export const URL_QUESTION_SVC_GET_CURRENT_QUESTION = BACKEND_QUESTION_SERVICE_URL + PREFIX_QUESTION_SVC_GET_CURRENT_QUESTION;