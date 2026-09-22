import React from "react";
import { LoginModal } from "../components/LoginModal";
import { useAppSelector } from "../stores/hooks";

// LoginGate shows the role picker whenever there is no JWT in the store.
export const LoginGate: React.FC<React.PropsWithChildren> = ({ children }) => {
  const token = useAppSelector((s) => s.auth.token);
  return (
    <>
      {children}
      {!token && <LoginModal />}
    </>
  );
};
