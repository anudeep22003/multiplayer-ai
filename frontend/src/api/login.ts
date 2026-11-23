import axios from "axios";
import { z } from "zod";
import { readSession, SessionSchema } from "@/features/login/session";
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

const LoginResponseSchema = SessionSchema.extend({ isNew: z.boolean() });
export type LoginResponse = z.infer<typeof LoginResponseSchema>;

export async function login(username: string): Promise<LoginResponse> {
  const request: LoginRequest = { username };
  const { data } = await api.post<unknown>("/login", request);
  const parsed = LoginResponseSchema.safeParse(data);
  if (!parsed.success) {
    throw new Error("The login service returned an invalid response.");
  }
  return parsed.data;
}
