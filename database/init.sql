CREATE TABLE IF NOT EXISTS flight_turnaround (
  id INTEGER PRIMARY KEY,
  flight_no TEXT,
  aircraft_reg TEXT,
  stand_no TEXT,
  arrival_time TEXT,
  departure_time TEXT,
  turnaround_status TEXT,
  delay_reason TEXT
);

CREATE TABLE IF NOT EXISTS ground_task (
  id INTEGER PRIMARY KEY,
  turnaround_id TEXT,
  task_type TEXT,
  team_id TEXT,
  planned_start TEXT,
  deadline TEXT,
  actual_finish TEXT,
  status TEXT,
  blocker_note TEXT
);

CREATE TABLE IF NOT EXISTS ground_resource (
  id INTEGER PRIMARY KEY,
  resource_code TEXT,
  resource_type TEXT,
  location TEXT,
  availability_status TEXT,
  maintenance_due_at TEXT,
  owner_team TEXT
);

CREATE TABLE IF NOT EXISTS resource_booking (
  id INTEGER PRIMARY KEY,
  resource_id TEXT,
  turnaround_id TEXT,
  task_id TEXT,
  start_time TEXT,
  end_time TEXT,
  booking_status TEXT,
  conflict_reason TEXT
);

CREATE TABLE IF NOT EXISTS delay_event (
  id INTEGER PRIMARY KEY,
  turnaround_id TEXT,
  delay_type TEXT,
  minutes TEXT,
  root_cause TEXT,
  responsibility_team TEXT,
  resolved_at TEXT
);

CREATE TABLE IF NOT EXISTS audit_log (
  id INTEGER PRIMARY KEY,
  actor TEXT,
  action TEXT,
  target_type TEXT,
  target_id TEXT,
  created_at TEXT
);
