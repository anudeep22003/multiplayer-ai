import axios from "axios";
import { VITE_API_URL } from "@/lib/env";
import { isSession, readSession } from "@/lib/session";

export const api = axios.create({ baseURL: VITE_API_URL });

api.interceptors.request.use((request) => {
  const session = readSession();
  if (session) request.headers.set("X-Dev-User-Id", session.userId);
  return request;
});

export async function login(username: string) {
  const { data } = await api.post<unknown>("/login", { username });
  if (
    !isSession(data) ||
    !("isNew" in data) ||
    typeof data.isNew !== "boolean"
  ) {
    throw new Error("The login service returned an invalid response.");
  }
  return { userId: data.userId, username: data.username, isNew: data.isNew };
}
