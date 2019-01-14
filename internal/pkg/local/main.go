package local

//
//import (
//	"errors"
//	"fmt"
//	"github.com/fatih/color"
//	"github.com/olekukonko/tablewriter"
//	"github.com/pgmtc/orchard-cli/internal/pkg/common"
//	"os"
//	"strconv"
//	"strings"
//	"time"
//)
//
//func Parse(args []string) error {
//
//	actions := common.MakeActions()
//	actions["status"] = status
//	actions["stop"] = common.ComponentActionHandler(stopContainer, config)
//	actions["start"] = common.ComponentActionHandler(startContainer, config)
//	actions["remove"] = common.ComponentActionHandler(removeContainer, config)
//	actions["create"] = common.ComponentActionHandler(createContainer, config)
//	actions["pull"] = common.ComponentActionHandler(pullImage, config)
//	actions["logs"] = logsHandler(dockerPrintLogs, false)
//	actions["watch"] = logsHandler(dockerPrintLogs, true)
//	return common.ParseParams(actions, args)
//}
//
//func logsHandler(handler func(component common.Component, follow bool) error, follow bool) func(args []string) error {
//	return func(args []string) error {
//		if len(args) == 0 {
//			return errors.New(fmt.Sprintf("Missing component Name. Available components = %s", config.CurrentPro.ComponentNames()))
//		}
//		componentId := args[0]
//		componentMap := common.ComponentMap()
//		if component, ok := componentMap[componentId]; ok {
//			return handler(component, follow)
//		}
//		return errors.New(fmt.Sprintf("Cannot find component '%s'. Available components = %s", componentId, common.ComponentNames()))
//	}
//}
