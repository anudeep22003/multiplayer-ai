import axios from "axios";
import { isSession, readSession, type Session } from "@/features/login/session";
import { VITE_API_URL } from "@/lib/env";

export const api = axios.create({ baseURL: VITE_API_URL });

api.interceptors.request.use((request) => {
  const session = readSession();
  if (session) request.headers.set("X-Dev-User-Id", session.userId);
  return request;
});

interface LoginRequest {
  username: string;
}

export interface LoginResponse extends Session {
  isNew: boolean;
}

export async function login(username: string): Promise<LoginResponse> {
  const request: LoginRequest = { username };
  const { data } = await api.post<unknown>("/login", request);
  if (
    !isSession(data) ||
    !("isNew" in data) ||
    typeof data.isNew !== "boolean"
  ) {
    throw new Error("The login service returned an invalid response.");
  }
  return { userId: data.userId, username: data.username, isNew: data.isNew };
}
