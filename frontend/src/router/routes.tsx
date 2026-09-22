import { createBrowserRouter, Navigate } from "react-router-dom";
import { AppShell } from "../components/AppShell";
import { RequireAuth } from "./RequireAuth";
import { DashboardPage } from "../pages/DashboardPage";
import { TurnaroundsPage } from "../pages/TurnaroundsPage";
import { TasksPage } from "../pages/TasksPage";
import { ResourcesPage } from "../pages/ResourcesPage";
import { DelaysPage } from "../pages/DelaysPage";
import { LoginGate } from "./LoginGate";

export const router = createBrowserRouter([
  {
    path: "/",
    element: (
      <LoginGate>
        <RequireAuth>
          <AppShell>
            <DashboardPage />
          </AppShell>
        </RequireAuth>
      </LoginGate>
    ),
    children: [],
  },
  {
    path: "/dashboard",
    element: (
      <LoginGate>
        <RequireAuth>
          <AppShell>
            <DashboardPage />
          </AppShell>
        </RequireAuth>
      </LoginGate>
    ),
  },
  {
    path: "/turnarounds",
    element: (
      <LoginGate>
        <RequireAuth>
          <AppShell>
            <TurnaroundsPage />
          </AppShell>
        </RequireAuth>
      </LoginGate>
    ),
  },
  {
    path: "/tasks",
    element: (
      <LoginGate>
        <RequireAuth>
          <AppShell>
            <TasksPage />
          </AppShell>
        </RequireAuth>
      </LoginGate>
    ),
  },
  {
    path: "/resources",
    element: (
      <LoginGate>
        <RequireAuth>
          <AppShell>
            <ResourcesPage />
          </AppShell>
        </RequireAuth>
      </LoginGate>
    ),
  },
  {
    path: "/delays",
    element: (
      <LoginGate>
        <RequireAuth>
          <AppShell>
            <DelaysPage />
          </AppShell>
        </RequireAuth>
      </LoginGate>
    ),
  },
  { path: "*", element: <Navigate to="/dashboard" replace /> },
]);
