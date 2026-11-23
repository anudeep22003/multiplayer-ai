import { isAxiosError } from "axios";
import { useState } from "react";
import { type LoginResponse, login } from "@/api/login";
import {
  clearSession,
  readSession,
  type Session,
  saveSession,
} from "./session";

type LoginState =
  | { status: "idle" }
  | { status: "submitting" }
  | { status: "failed"; error: string }
  | { status: "signedIn"; session: Session; isNew: boolean; error?: string };

function restoreLogin(): LoginState {
  const session = readSession();
  return session
    ? { status: "signedIn", session, isNew: false }
    : { status: "idle" };
}

export function useLogin() {
  const [state, setState] = useState(restoreLogin);

  async function handleLogin(username: string) {
    if (state.status === "submitting" || state.status === "signedIn") return;
    const trimmedUsername = username.trim();
    if (!trimmedUsername) {
      setState({ status: "failed", error: "Enter a username." });
      return;
    }
    setState({ status: "submitting" });
    let result: LoginResponse;
    try {
      result = await login(trimmedUsername);
    } catch (cause) {
      const error =
        isAxiosError(cause) && !cause.response
          ? "Couldn’t reach the login service. Check your connection and try again."
          : "Login failed. Check your username and try again.";
      setState({ status: "failed", error });
      return;
    }
    try {
      saveSession(result);
    } catch {
      setState({
        status: "failed",
        error:
          "Couldn’t save your session. Allow browser storage and try again.",
      });
      return;
    }
    setState({ status: "signedIn", session: result, isNew: result.isNew });
  }

  function handleLogout() {
    if (state.status !== "signedIn") return;
    try {
      clearSession();
      setState({ status: "idle" });
    } catch {
      setState({
        ...state,
        error:
          "Couldn’t clear your saved session. Allow browser storage and try again.",
      });
    }
  }

  function clearLoginError() {
    if (state.status === "failed") setState({ status: "idle" });
  }

  return { state, handleLogin, handleLogout, clearLoginError };
}
