package models

type Slot struct {
    ID        int    `json:"id"`
    StartTime string `json:"start_time"`
    EndTime   string `json:"end_time"`
}