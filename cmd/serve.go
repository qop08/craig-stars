package cmd

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/db"
	"github.com/sirgwain/craig-stars/game"
	"github.com/sirgwain/craig-stars/server"

	"github.com/spf13/cobra"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Debug().Msgf("%s took %s", name, elapsed)
}

func newServeCmd() *cobra.Command {
	var generateUniverse bool
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the webserver",
		Long:  `Start a local gin-gonic webserver and serve requests.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			db := db.NewClient()
			cfg := config.GetConfig()
			db.Connect(cfg)
			db.ExecSchema("schema.sql")

			if generateUniverse {
				if err := generateTestGame(db, *cfg); err != nil {
					return err
				}
			}
			server.Start(db, *cfg)
			return nil
		},
	}
	serveCmd.Flags().BoolVar(&generateUniverse, "generate-test-game", false, "Generate a test user and game")

	return serveCmd
}

func generateTestGame(db server.DBClient, config config.Config) error {
	defer timeTrack(time.Now(), "generateTestGame")

	admin, adminRace, err := createTestUser(db, "admin", config.GeneratedUserPassword, game.RoleAdmin)
	if err != nil {
		return err
	}

	user2, user2Race, err := createTestUser(db, "craig", config.GeneratedUserPassword, game.RoleUser)
	if err != nil {
		return err
	}

	// create a game runner to host some games
	gameRunner := server.NewGameRunner(db)

	// admin user will host a game with an ai player
	if _, err := gameRunner.HostGame(admin.ID, game.NewGameSettings().
		// WithSize(game.SizeTiny).
		// WithDensity(game.DensitySparse).
		WithHost(adminRace.ID).
		WithAIPlayer(game.AIDifficultyNormal)); err != nil {
		return err
	}

	// also create a medium size game with 25 turns generated
	mediumGame, err := gameRunner.HostGame(admin.ID, game.NewGameSettings().
		WithName("Medium Game").
		WithSize(game.SizeMedium).
		WithHost(adminRace.ID).
		WithAIPlayer(game.AIDifficultyNormal))
	if err != nil {
		return err
	}
	for i := 0; i < 25; i++ {
		gameRunner.SubmitTurn(mediumGame.ID, mediumGame.HostID)
		if _, err := gameRunner.CheckAndGenerateTurn(mediumGame.ID); err != nil {
			log.Error().Err(err).Msg("check and generate new turn")
		}
	}

	// user2 will also host a game so with an open player slot
	_, err = gameRunner.HostGame(user2.ID, game.NewGameSettings().
		WithName("Joinable Game").
		WithHost(user2Race.ID).
		WithOpenPlayerSlot())
	if err != nil {
		return err
	}

	return nil
}

func createTestUser(db server.DBClient, username string, password string, role game.Role) (*game.User, *game.Race, error) {
	user, err := db.GetUserByUsername(username)
	if err != nil {
		return nil, nil, err
	}

	// default the password to the username if it's empty
	if password == "" {
		password = username
	}

	if user == nil {
		user = game.NewUser(username, password, role)
		err := db.CreateUser(user)
		if err != nil {
			return nil, nil, err
		}

	}

	races, err := db.GetRacesForUser(user.ID)
	if err != nil {
		return nil, nil, err
	}

	var race game.Race
	if len(races) == 0 {
		race = game.Humanoids()
		race.UserID = user.ID

		if err := db.CreateRace(&race); err != nil {
			return nil, nil, err
		}

		// race = game.PPs()
		// race.UserID = user.ID

		// if err := db.CreateRace(race); err != nil {
		// 	return nil, nil, err
		// }
	} else {
		race = races[0]
	}

	return user, &race, nil
}

func init() {
	rootCmd.AddCommand(newServeCmd())
}
