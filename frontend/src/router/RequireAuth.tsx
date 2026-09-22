import React from "react";
import { Navigate } from "react-router-dom";
import { useAppSelector } from "../stores/hooks";

// RequireAuth is the frontend route guard mirroring backend JWT middleware.
export const RequireAuth: React.FC<React.PropsWithChildren> = ({ children }) => {
  const token = useAppSelector((s) => s.auth.token);
  if (!token) return <Navigate to="/dashboard" replace />;
  return <>{children}</>;
};
