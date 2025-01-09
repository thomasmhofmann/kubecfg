// Copyright 2017 The kubecfg authors
//
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bitnami/kubecfg/pkg/kubecfg"
)

const (
	flagCreate   = "create"
	flagSkipGc   = "skip-gc"
	flagGcTag    = "gc-tag"
	flagDryRun   = "dry-run"
	flagValidate = "validate"
)

func init() {
	RootCmd.AddCommand(updateCmd)
	updateCmd.PersistentFlags().Bool(flagCreate, true, "Create missing resources")
	updateCmd.PersistentFlags().Bool(flagSkipGc, false, "Don't perform garbage collection, even with --"+flagGcTag)
	updateCmd.PersistentFlags().String(flagGcTag, "", "Add this tag to updated objects, and garbage collect existing objects with this tag and not in config")
	updateCmd.PersistentFlags().Bool(flagDryRun, false, "Perform only read-only operations")
	updateCmd.PersistentFlags().Bool(flagValidate, true, "Validate input against server schema")
	updateCmd.PersistentFlags().Bool(flagIgnoreUnknown, false, "Don't fail validation if the schema for a given resource type is not found")
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Kubernetes resources with local config",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := cmd.Flags()
		var err error
		c := kubecfg.UpdateCmd{}

		validate, err := flags.GetBool(flagValidate)
		if err != nil {
			return err
		}

		c.Create, err = flags.GetBool(flagCreate)
		if err != nil {
			return err
		}

		c.GcTag, err = flags.GetString(flagGcTag)
		if err != nil {
			return err
		}

		c.SkipGc, err = flags.GetBool(flagSkipGc)
		if err != nil {
			return err
		}

		c.DryRun, err = flags.GetBool(flagDryRun)
		if err != nil {
			return err
		}

		c.Client, c.Mapper, c.Discovery, err = getDynamicClients(cmd)
		if err != nil {
			return err
		}

		c.DefaultNamespace, err = defaultNamespace(clientConfig)
		if err != nil {
			return err
		}

		objs, err := readObjs(cmd, args)
		if err != nil {
			return err
		}

		if validate {
			v := kubecfg.ValidateCmd{
				Mapper:    c.Mapper,
				Discovery: c.Discovery,
			}

			v.IgnoreUnknown, err = flags.GetBool(flagIgnoreUnknown)
			if err != nil {
				return err
			}

			if err := v.Run(objs, cmd.OutOrStdout()); err != nil {
				return err
			}
		}

		return c.Run(cmd.Context(), objs)
	},
}
