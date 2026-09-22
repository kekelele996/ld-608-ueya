package services

import "time"

func minutesDuration(n int) time.Duration { return time.Duration(n) * time.Minute }

func resolveNow() time.Time { return time.Now() }

func nowTime() time.Time { return time.Now() }
