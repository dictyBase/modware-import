package arangodb

import (
	"context"
	"fmt"
	"strings"

	"github.com/dictyBase/arangomanager"
	"github.com/dictyBase/modware-import/internal/collection"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// DeleteCmd deletes all data from the given collections
var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete all data from all the collections",
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		session, db, err := arangomanager.NewSessionDb(
			&arangomanager.ConnectParams{
				User:     viper.GetString("arangodb-user"),
				Pass:     viper.GetString("arangodb-pass"),
				Database: viper.GetString("database"),
				Host:     viper.GetString("arangodb-host"),
				Port:     viper.GetInt("arangodb-port"),
				Istls:    viper.GetBool("is-secure"),
			},
		)
		if err != nil {
			return fmt.Errorf("error in connecting to database %s", err)
		}
		registry.SetArangoSession(session)
		registry.SetArangodbConnection(db)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		db := registry.GetArangodbConnection()
		colls, err := db.Handler().Collections(context.Background())
		if err != nil {
			return fmt.Errorf(
				"error in retrieving collections for database %s %s",
				viper.GetString("database"),
				err,
			)
		}
		logger := registry.GetLogger()
		var names []string
		for _, c := range colls {
			props, err := c.Properties(context.Background())
			if err != nil {
				return fmt.Errorf(
					"error in retrieving properties of collection %s %s",
					c.Name(), err,
				)
			}
			if props.IsSystem {
				continue
			}
			names = append(names, c.Name())
		}
		logger.Debugf("received %s collections for deletion", strings.Join(names, " "))
		if len(viper.GetStringSlice("exclude")) > 0 {
			logger.Debugf(
				"excluding the %s collections",
				strings.Join(viper.GetStringSlice("exclude"), ""),
			)
			names = collection.Remove(names, viper.GetStringSlice("exclude")...)
		}
		logger.Debugf("going to delete %s collections", strings.Join(names, " "))
		err = db.Truncate(names...)
		if err != nil {
			return err
		}
		logger.Infof("deleted %s collections", strings.Join(names, " "))
		return nil
	},
}

func init() {
	DeleteCmd.Flags().StringP(
		"database",
		"d",
		"",
		"name of arangodb database whose data from all the non-system collections will be deleted",
	)
	DeleteCmd.Flags().StringSliceP(
		"exclude",
		"e",
		[]string{""},
		"collection to exclude from deletion",
	)
	viper.BindPFlags(DeleteCmd.Flags())
}
