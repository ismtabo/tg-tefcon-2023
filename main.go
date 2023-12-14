package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"time"

	"github.com/Telefonica/tg-tefcon-2023/assets"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/samber/lo"
)

var client = http.DefaultClient

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, handler),
		bot.WithMessageTextHandler("/help", bot.MatchTypeExact, helpHandler),
		bot.WithMessageTextHandler("/map", bot.MatchTypeExact, mapHandler),
		bot.WithMessageTextHandler("/rooms", bot.MatchTypeExact, roomsHandler),
		bot.WithMessageTextHandler("/current_events", bot.MatchTypeExact, currentEventsHandler),
		bot.WithMessageTextHandler("/next_events", bot.MatchTypeExact, nextEventsHandler),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Hola, bienvenido a @tefconbot, tu asistente en la TefCON 2023. Puedes usar el comando /help para ver los comandos disponibles.",
	})
}

func helpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Comandos disponibles:\n/rooms\n/map\n/current_events/\n/next_events",
	})
}

func mapHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendPhoto(ctx, &bot.SendPhotoParams{
		ChatID: update.Message.Chat.ID,
		Photo: &models.InputFileUpload{
			Filename: "map.jpg",
			Data:     bytes.NewReader(assets.Map),
		},
	})
}

func roomsHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	rooms, err := getRooms()
	if err != nil {
		slog.Error("Error getting rooms: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Ups, algo falló. ¡Qué chopecha!",
		})
		return
	}
	txt := ""
	for _, room := range rooms {
		txt += room.Name
		if room.Occupancy > 0 {
			txt += fmt.Sprintf(" (%d %% of %d)", room.Occupancy, room.Capacity)
		}
		txt += "\n"
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   txt,
	})
}

const baseurl = "https://tefcon.tid.es/api/v1"

type Rooms []Room

type Room struct {
	ID              int64  `json:"id"`
	ShortName       string `json:"short_name"`
	Name            string `json:"name"`
	Location        string `json:"location"`
	Capacity        int64  `json:"capacity"`
	Color           string `json:"color"`
	Occupancy       int64  `json:"occupancy"`
	ShowOccupancy   bool   `json:"show_occupancy"`
	OccupancyEditor int64  `json:"occupancy_editor"`
}

func getRooms() ([]Room, error) {
	url := baseurl + "/rooms"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	rooms := []Room{}
	if err := json.NewDecoder(resp.Body).Decode(&rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}

type BasicInfo []BasicInfoElement

type BasicInfoElement struct {
	ID            int64        `json:"id"`
	IsActive      bool         `json:"is_active"`
	StartDateTime string       `json:"start_date_time"`
	EndDateTime   string       `json:"end_date_time"`
	Event         Event        `json:"event"`
	MeetingRoom   *MeetingRoom `json:"meeting_room"`
}

type Event struct {
	ID                 int64     `json:"id"`
	Name               string    `json:"name"`
	DisplayDescription bool      `json:"display_description"`
	EventType          EventType `json:"event_type"`
	DisplayOwner       *string   `json:"display_owner"`
}

type MeetingRoom struct {
	ID              int64           `json:"id"`
	ShortName       ShortName       `json:"short_name"`
	Name            string          `json:"name"`
	Location        Location        `json:"location"`
	Capacity        int64           `json:"capacity"`
	Color           Color           `json:"color"`
	Occupancy       int64           `json:"occupancy"`
	ShowOccupancy   bool            `json:"show_occupancy"`
	OccupancyEditor OccupancyEditor `json:"occupancy_editor"`
}

type OccupancyEditor struct {
	ID              int64         `json:"id"`
	Password        string        `json:"password"`
	LastLogin       interface{}   `json:"last_login"`
	IsSuperuser     bool          `json:"is_superuser"`
	Username        Username      `json:"username"`
	FirstName       string        `json:"first_name"`
	LastName        string        `json:"last_name"`
	Email           string        `json:"email"`
	IsStaff         bool          `json:"is_staff"`
	IsActive        bool          `json:"is_active"`
	DateJoined      string        `json:"date_joined"`
	Groups          []interface{} `json:"groups"`
	UserPermissions []interface{} `json:"user_permissions"`
}

type EventType string

const (
	Other  EventType = "OTHER"
	Speech EventType = "SPEECH"
)

type Color string

const (
	C466Ef    Color = "#c466ef"
	E66C64    Color = "#E66C64"
	Eac344    Color = "#eac344"
	The59C2C9 Color = "#59C2C9"
	The64566A Color = "#64566A"
)

type Location string

const (
	CentroDeFormación Location = "Centro de Formación"
	EdificioCentral   Location = "Edificio Central"
)

type Username string

const (
	Sala1 Username = "sala1"
	Sala2 Username = "sala2"
	Sala3 Username = "sala3"
	Sala4 Username = "sala4"
	Sala5 Username = "sala5"
)

type ShortName string

const (
	Auditorio1         ShortName = "AUDITORIO 1"
	Auditorio2         ShortName = "AUDITORIO 2"
	CentroFormaciónA11 ShortName = "CENTRO FORMACIÓN A11"
	CentroFormaciónJ21 ShortName = "CENTRO FORMACIÓN J21"
	SalaCiria          ShortName = "SALA CIRIA"
)

func getBasicInfo() ([]BasicInfoElement, error) {
	url := baseurl + "/events/basicInfo/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	basicInfo := []BasicInfoElement{}
	if err := json.NewDecoder(resp.Body).Decode(&basicInfo); err != nil {
		return nil, err
	}
	return basicInfo, nil
}

func (a BasicInfo) Len() int           { return len(a) }
func (a BasicInfo) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BasicInfo) Less(i, j int) bool { return a[i].StartDateTime < a[j].StartDateTime }

func getCurrentEvents() ([]BasicInfoElement, error) {
	basicInfo, err := getBasicInfo()
	if err != nil {
		return nil, err
	}
	currentEvents := lo.Filter(basicInfo, func(info BasicInfoElement, _index int) bool {
		return info.IsActive
	})
	return currentEvents, nil
}

func currentEventsHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	currentEvents, err := getCurrentEvents()
	if err != nil {
		slog.Error("Error getting current events: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Ups, algo falló. ¡Qué chopecha!",
		})
		return
	}
	txt := "Eventos en curso:\n"
	for _, info := range currentEvents {
		txt += fmt.Sprintf("- %s", info.Event.Name)
		if info.MeetingRoom != nil {
			txt += fmt.Sprintf(" en %s", info.MeetingRoom.Name)
			txt += fmt.Sprintf(" (%d %% of %d)", info.MeetingRoom.Occupancy, info.MeetingRoom.Capacity)
		}
		txt += "\n"
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   txt,
	})
}

func getNextEvents() ([]BasicInfoElement, error) {
	basicInfo, err := getBasicInfo()
	if err != nil {
		return nil, err
	}
	eventsBySlot := lo.GroupBy[BasicInfoElement, string](basicInfo, func(info BasicInfoElement) string {
		return info.StartDateTime
	})
	slots := lo.Map(basicInfo, func(info BasicInfoElement, _index int) string {
		return info.StartDateTime
	})
	nextSlots := lo.Filter(slots, func(slot string, _index int) bool {
		ts, err := time.Parse(time.RFC3339, slot)
		if err != nil {
			slog.Error("Error parsing time: %v", err)
			return false
		}
		return time.Now().Before(ts)
	})
	sort.Strings(nextSlots)
	nextSlot := nextSlots[0]
	nextEvents := eventsBySlot[nextSlot]
	return nextEvents, nil
}

func nextEventsHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	nextEvents, err := getNextEvents()
	if err != nil {
		slog.Error("Error getting next events: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Ups, algo falló. ¡Qué chopecha!",
		})
		return
	}
	txt := "Próximos eventos:\n"
	for _, info := range nextEvents {
		ts, err := time.Parse(time.RFC3339, info.StartDateTime)
		if err != nil {
			slog.Error("Error parsing time: %v", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Ups, algo falló. ¡Qué chopecha!",
			})
			return
		}
		txt += fmt.Sprintf("- %s", info.Event.Name)
		txt += fmt.Sprintf(" (%s)", ts.Format("15:04"))
		if info.MeetingRoom != nil {
			txt += fmt.Sprintf(" en %s", info.MeetingRoom.Name)
		}
		txt += "\n"
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   txt,
	})
}
