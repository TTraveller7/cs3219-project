import { Client } from 'redis-om';

/* pulls the Redis URL from .env */
const url = process.env.CHAT_REDIS_URL || "redis://localhost:6384";

/* create and open the Redis OM Client */
export const client = await new Client().open(url);
