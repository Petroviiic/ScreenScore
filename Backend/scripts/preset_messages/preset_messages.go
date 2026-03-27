package scripts

import (
	"context"
	"log"

	"github.com/Petroviiic/ScreenScore/internal/db"
	"github.com/Petroviiic/ScreenScore/internal/env"
	"github.com/Petroviiic/ScreenScore/internal/storage"
)

type PresetMessage struct {
	Text     string
	Price    int
	Rarity   string
	IsActive bool
}

type MessageCategory string

const (
	CategoryCommon    MessageCategory = "common"
	CategoryRare      MessageCategory = "rare"
	CategoryLegendary MessageCategory = "legendary"
)

func Seed() error {
	type dbConfig struct {
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
		dbAddr       string
	}
	cfg := dbConfig{
		maxOpenConns: env.GetInt("DB_MAX_OPEN_CONS", 30),
		maxIdleConns: env.GetInt("DB_MAX_IDLE_CONS", 30),
		maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		dbAddr:       env.GetString("DB_ADDR", "postgresql://user:user123@localhost:5432/screenscore?sslmode=disable"),
	}
	db, err := db.NewDb(cfg.dbAddr, cfg.maxIdleConns, cfg.maxOpenConns, cfg.maxIdleTime)
	if err != nil {
		log.Panic("error connecting to db")
		return err
	}

	storage := storage.NewStorage(db)

	msgs := []PresetMessage{
		{Text: "Phone on charger, brain on vacation! 🔋", Price: 150, Rarity: "common", IsActive: true},
		{Text: "Someone was productive today. Well done! 👏", Price: 180, Rarity: "common", IsActive: true},
		{Text: "Just a reminder that the real world exists. 🌍", Price: 200, Rarity: "common", IsActive: true},
		{Text: "Close your eyes, breathe, stop scrolling. 🧘‍♂️", Price: 220, Rarity: "common", IsActive: true},
		{Text: "Less screen = More life. Simple as that. ✨", Price: 250, Rarity: "common", IsActive: true},
		{Text: "Your focus is your currency. Don't waste it! 💸", Price: 280, Rarity: "common", IsActive: true},
		{Text: "Well-deserved break from notifications. Enjoy! ☕", Price: 300, Rarity: "common", IsActive: true},
		{Text: "Off-screen mode: ON. 🚫📱", Price: 300, Rarity: "common", IsActive: true},
		{Text: "Look at people, not pixels. 👀", Price: 320, Rarity: "common", IsActive: true},
		{Text: "Congrats king, you beat the algorithm! 🏆", Price: 350, Rarity: "common", IsActive: true},

		{Text: "While you scroll, I dominate. Greetings from the top! 🏔️", Price: 550, Rarity: "rare", IsActive: true},
		{Text: "Digital detox level: Expert. 🥋", Price: 650, Rarity: "rare", IsActive: true},
		{Text: "Your screen time is my motivation. ⏱️", Price: 700, Rarity: "rare", IsActive: true},
		{Text: "This user actually has a hobby that isn't TikTok. 🔥", Price: 750, Rarity: "rare", IsActive: true},
		{Text: "Battery 100%, Time saved: Infinite. ⚡", Price: 800, Rarity: "rare", IsActive: true},
		{Text: "Winner of the week in the 'I have a life' category. 🏅", Price: 850, Rarity: "rare", IsActive: true},
		{Text: "Notifications on Mute, Ambitions on Max. 🚀", Price: 900, Rarity: "rare", IsActive: true},
		{Text: "I own my phone, it doesn't own me. 🦾", Price: 1000, Rarity: "rare", IsActive: true},
		{Text: "Happiness is where the Wi-Fi is weak. 🌲", Price: 1100, Rarity: "rare", IsActive: true},
		{Text: "Choosing what you want most over what you want now. 💎", Price: 1200, Rarity: "rare", IsActive: true},

		{Text: "Master of Focus. 👑", Price: 2500, Rarity: "legendary", IsActive: true},
		{Text: "This message is from the future because I don't waste time. 🛸", Price: 3000, Rarity: "legendary", IsActive: true},
		{Text: "MYTHICAL CREATURE: User with <2h screen time. 🦄", Price: 3500, Rarity: "legendary", IsActive: true},
		{Text: "Your Screen Time report is my appetizer. 🍷", Price: 4000, Rarity: "legendary", IsActive: true},
		{Text: "Untouchable in the real world. 🌫️", Price: 4500, Rarity: "legendary", IsActive: true},
		{Text: "I forgot my passcode from all this detoxing. 🧠", Price: 5000, Rarity: "legendary", IsActive: true},
		{Text: "My time is more expensive than your annual data plan. 💰", Price: 5500, Rarity: "legendary", IsActive: true},
		{Text: "Focus level: Buddhist monk on steroids. 🧘‍♂️⚡", Price: 6000, Rarity: "legendary", IsActive: true},
		{Text: "I bought this message with your lost minutes. 😎", Price: 6500, Rarity: "legendary", IsActive: true},
		{Text: "THE G.O.A.T. OF DETOX. 🐐🐐🐐", Price: 8000, Rarity: "legendary", IsActive: true},
	}

	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	for _, msg := range msgs {
		err := storage.MessageStorage.InsertNewPresetMessage(ctx, tx, msg.Text, msg.Price, msg.Rarity, msg.IsActive)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
