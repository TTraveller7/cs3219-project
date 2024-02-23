import {Client} from "redis-om";

/* pulls the Redis URL from .env */
const url = process.env.COLLAB_REDIS_URL || "redis://localhost:6383";

/* create and open the Redis OM Client */
const client = await new Client().open(url);
export { client };