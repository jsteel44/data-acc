package actions_impl

import (
	"errors"
	"fmt"
	"github.com/RSE-Cambridge/data-acc/internal/pkg/dacctl"
	parsers2 "github.com/RSE-Cambridge/data-acc/internal/pkg/dacctl/actions_impl/parsers"
	"github.com/RSE-Cambridge/data-acc/internal/pkg/dacctl/workflow_impl"
	"github.com/RSE-Cambridge/data-acc/internal/pkg/datamodel"
	"github.com/RSE-Cambridge/data-acc/internal/pkg/facade"
	"github.com/RSE-Cambridge/data-acc/internal/pkg/fileio"
	"github.com/RSE-Cambridge/data-acc/internal/pkg/store"
	"log"
	"strings"
)

func NewDacctlActions(keystore store.Keystore, disk fileio.Disk) dacctl.DacctlActions {
	return &dacctlActions{
		session: workflow_impl.NewSessionFacade(keystore),
		disk:    disk,
	}
}

type dacctlActions struct {
	session facade.Session
	disk    fileio.Disk
}

func checkRequiredStrings(c dacctl.CliContext, flags ...string) error {
	var errs []string
	for _, flag := range flags {
		if str := c.String(flag); str == "" {
			errs = append(errs, flag)
		}
	}
	if len(errs) > 0 {
		errStr := fmt.Sprintf("Please provide these required parameters: %s", strings.Join(errs, ", "))
		log.Println(errStr)
		return errors.New(errStr)
	}
	return nil
}

func (d *dacctlActions) getSessionName(c dacctl.CliContext) (datamodel.SessionName, error) {
	err := checkRequiredStrings(c, "token")
	if err != nil {
		return "", err
	}

	token := c.String("token")
	if !parsers2.IsValidName(token) {
		return "", fmt.Errorf("badly formatted session name: %s", token)
	}

	return datamodel.SessionName(token), nil
}

func (d *dacctlActions) DeleteBuffer(c dacctl.CliContext) error {
	sessionName, err := d.getSessionName(c)
	if err != nil {
		return err
	}
	hurry := c.Bool("hurry")
	return d.session.DeleteSession(sessionName, hurry)
}

func (d *dacctlActions) DataIn(c dacctl.CliContext) error {
	sessionName, err := d.getSessionName(c)
	if err != nil {
		return err
	}
	return d.session.CopyDataIn(sessionName)
}

func (d *dacctlActions) PreRun(c dacctl.CliContext) error {
	sessionName, err := d.getSessionName(c)
	if err != nil {
		return err
	}
	err = checkRequiredStrings(c, "nodehostnamefile")
	if err != nil {
		return err
	}

	computeHosts, err := parsers2.GetHostnamesFromFile(d.disk, c.String("nodehostnamefile"))
	if err != nil {
		return err
	}
	if len(computeHosts) < 1 {
		return errors.New("unable to mount to zero compute hosts")
	}

	loginNodeFilename := c.String("jobexecutionnodefile")
	var loginNodeHosts []string
	if loginNodeFilename != "" {
		loginNodeHosts, err = parsers2.GetHostnamesFromFile(d.disk, loginNodeFilename)
		if err != nil {
			return err
		}
	}

	return d.session.Mount(sessionName, computeHosts, loginNodeHosts)
}

func (d *dacctlActions) PostRun(c dacctl.CliContext) error {
	sessionName, err := d.getSessionName(c)
	if err != nil {
		return err
	}
	return d.session.Unmount(sessionName)
}

func (d *dacctlActions) DataOut(c dacctl.CliContext) error {
	sessionName, err := d.getSessionName(c)
	if err != nil {
		return err
	}
	return d.session.CopyDataOut(sessionName)
}

func (d *dacctlActions) GenerateAnsible(c dacctl.CliContext) (string, error) {
	sessionName, err := d.getSessionName(c)
	if err != nil {
		return "", err
	}
	return d.session.GenerateAnsible(sessionName)
}
