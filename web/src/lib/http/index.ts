import { HttpClient } from "./http";

export * from "./http";
export * from "./errors";

export const http = new HttpClient();
export default http;