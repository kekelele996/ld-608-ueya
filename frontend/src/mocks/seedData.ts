export const mockData = {
  "flightTurnaround": [
    {
      "id": 1,
      "flight_no": "flight no 1",
      "aircraft_reg": "aircraft reg 1",
      "stand_no": "stand no 1",
      "arrival_time": "2026-06-11T09:00:00Z",
      "departure_time": "2026-06-11T09:00:00Z",
      "turnaround_status": "ON_STAND",
      "delay_reason": "delay reason 1"
    },
    {
      "id": 2,
      "flight_no": "flight no 2",
      "aircraft_reg": "aircraft reg 2",
      "stand_no": "stand no 2",
      "arrival_time": "2026-06-12T09:00:00Z",
      "departure_time": "2026-06-12T09:00:00Z",
      "turnaround_status": "IN_SERVICE",
      "delay_reason": "delay reason 2"
    },
    {
      "id": 3,
      "flight_no": "flight no 3",
      "aircraft_reg": "aircraft reg 3",
      "stand_no": "stand no 3",
      "arrival_time": "2026-06-13T09:00:00Z",
      "departure_time": "2026-06-13T09:00:00Z",
      "turnaround_status": "ARRIVING",
      "delay_reason": "delay reason 3"
    }
  ],
  "groundTask": [
    {
      "id": 1,
      "turnaround_id": 1,
      "task_type": "CATERING",
      "team_id": 1,
      "planned_start": "planned start 1",
      "deadline": "deadline 1",
      "actual_finish": "actual finish 1",
      "status": "ON_STAND",
      "blocker_note": "blocker note 1"
    },
    {
      "id": 2,
      "turnaround_id": 2,
      "task_type": "BAGGAGE",
      "team_id": 2,
      "planned_start": "planned start 2",
      "deadline": "deadline 2",
      "actual_finish": "actual finish 2",
      "status": "IN_SERVICE",
      "blocker_note": "blocker note 2"
    },
    {
      "id": 3,
      "turnaround_id": 3,
      "task_type": "REFUEL",
      "team_id": 3,
      "planned_start": "planned start 3",
      "deadline": "deadline 3",
      "actual_finish": "actual finish 3",
      "status": "ARRIVING",
      "blocker_note": "blocker note 3"
    }
  ],
  "groundResource": [
    {
      "id": 1,
      "resource_code": "resource code 1",
      "resource_type": "CATERING",
      "location": "location 1",
      "availability_status": "ON_STAND",
      "maintenance_due_at": "2026-06-11T09:00:00Z",
      "owner_team": "owner team 1"
    },
    {
      "id": 2,
      "resource_code": "resource code 2",
      "resource_type": "BAGGAGE",
      "location": "location 2",
      "availability_status": "IN_SERVICE",
      "maintenance_due_at": "2026-06-12T09:00:00Z",
      "owner_team": "owner team 2"
    },
    {
      "id": 3,
      "resource_code": "resource code 3",
      "resource_type": "REFUEL",
      "location": "location 3",
      "availability_status": "ARRIVING",
      "maintenance_due_at": "2026-06-13T09:00:00Z",
      "owner_team": "owner team 3"
    }
  ],
  "resourceBooking": [
    {
      "id": 1,
      "resource_id": 1,
      "turnaround_id": 1,
      "task_id": 1,
      "start_time": "2026-06-11T09:00:00Z",
      "end_time": "2026-06-11T09:00:00Z",
      "booking_status": "ON_STAND",
      "conflict_reason": "conflict reason 1"
    },
    {
      "id": 2,
      "resource_id": 2,
      "turnaround_id": 2,
      "task_id": 2,
      "start_time": "2026-06-12T09:00:00Z",
      "end_time": "2026-06-12T09:00:00Z",
      "booking_status": "IN_SERVICE",
      "conflict_reason": "conflict reason 2"
    },
    {
      "id": 3,
      "resource_id": 3,
      "turnaround_id": 3,
      "task_id": 3,
      "start_time": "2026-06-13T09:00:00Z",
      "end_time": "2026-06-13T09:00:00Z",
      "booking_status": "ARRIVING",
      "conflict_reason": "conflict reason 3"
    }
  ],
  "delayEvent": [
    {
      "id": 1,
      "turnaround_id": 1,
      "delay_type": "CATERING",
      "minutes": "minutes 1",
      "root_cause": "root cause 1",
      "responsibility_team": "responsibility team 1",
      "resolved_at": "2026-06-11T09:00:00Z"
    },
    {
      "id": 2,
      "turnaround_id": 2,
      "delay_type": "BAGGAGE",
      "minutes": "minutes 2",
      "root_cause": "root cause 2",
      "responsibility_team": "responsibility team 2",
      "resolved_at": "2026-06-12T09:00:00Z"
    },
    {
      "id": 3,
      "turnaround_id": 3,
      "delay_type": "REFUEL",
      "minutes": "minutes 3",
      "root_cause": "root cause 3",
      "responsibility_team": "responsibility team 3",
      "resolved_at": "2026-06-13T09:00:00Z"
    }
  ]
} as const;
